package main

import (
	"clash-admin/internal/buildflag"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var CONFIG_FILE = "/etc/mihomo/config.yaml"

func init() {
	if buildflag.Kernel == "clash" {
		CONFIG_FILE = "/etc/clash/config.yaml"
	}
}

const CLASH_HOST = "http://127.0.0.1:9999"

func isValidDomainName(domain string) bool {
	// 域名的正则表达式
	var domainRegex = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)

	return domainRegex.MatchString(domain)
}

func setWriteAndBlack(w http.ResponseWriter, r *http.Request, name string) {
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var ALL_RULES, DIRECT_OR_BLACK_RULES []string
	for _, rule_str := range config.Rules {
		rule, err := parseRule(rule_str)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rule.Type == "DOMAIN-SUFFIX" && rule.Name == name {
			DIRECT_OR_BLACK_RULES = append(DIRECT_OR_BLACK_RULES, rule_str)
		} else {
			ALL_RULES = append(ALL_RULES, rule_str)
		}
	}
	if r.Method == "GET" {
		var resDATA = make([]string, 0)
		for _, rule := range DIRECT_OR_BLACK_RULES {
			resDATA = append(resDATA, strings.Split(rule, ",")[1])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resDATA)
	} else if r.Method == "POST" {
		var NEW_DIRECT_DOMAIN, NEW_DIRECT_OR_BLACK_RULES []string
		// 解析请求体中的 JSON 数据
		if err := json.NewDecoder(r.Body).Decode(&NEW_DIRECT_DOMAIN); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, domain := range NEW_DIRECT_DOMAIN {
			if !isValidDomainName(domain) {
				http.Error(w, "填写的域名无效", http.StatusBadRequest)
				return
			}
			domain = strings.Trim(domain, " ")
			NEW_DIRECT_OR_BLACK_RULES = append(NEW_DIRECT_OR_BLACK_RULES, fmt.Sprintf("%s,%s,%s", "DOMAIN-SUFFIX", domain, name))
		}
		config.Rules = ALL_RULES
		var NewRules = make([]string, len(ALL_RULES)+len(NEW_DIRECT_OR_BLACK_RULES))
		copy(NewRules, NEW_DIRECT_OR_BLACK_RULES)
		copy(NewRules[len(NEW_DIRECT_OR_BLACK_RULES):], config.Rules)
		config.Rules = NewRules
		updatedData, err := yaml.Marshal(&config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		preData, err := readUntilProxies(CONFIG_FILE)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		preData = append(preData, updatedData...)
		err = os.WriteFile(CONFIG_FILE, preData, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		restartShell()
	}
}

func handleSetWrite(w http.ResponseWriter, r *http.Request) {
	setWriteAndBlack(w, r, "DIRECT")
}

func handleSetBlack(w http.ResponseWriter, r *http.Request) {
	setWriteAndBlack(w, r, "REJECT")
}

func main() {
	templates()
	static()
	http.HandleFunc("/listdata", listDataHandler)
	http.HandleFunc("/update-config", updateConfigHandler)
	http.HandleFunc("/delayDetection", delayHandler)
	http.HandleFunc("/check", handleCheckIP)
	http.HandleFunc("/yaml", yamlReset)
	http.HandleFunc("/activate", activate)
	http.HandleFunc("/write-to-config", writeToFile)
	http.HandleFunc("/setWrite", handleSetWrite)
	http.HandleFunc("/setBlack", handleSetBlack)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)

}
