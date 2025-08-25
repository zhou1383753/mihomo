package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func parseRule(ruleStr string) (Rule, error) {
	parts := strings.Split(ruleStr, ",")
	if len(parts) != 3 {
		return Rule{}, fmt.Errorf("invalid rule format: %s", ruleStr)
	}
	return Rule{
		Type: parts[0],
		CIDR: parts[1],
		Name: parts[2],
	}, nil
}

func listDataHandler(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Read the YAML file
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the YAML file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error unmarshaling YAML: %v", err)
		return
	}

	// Parse rules
	var rules []Rule
	for _, ruleStr := range config.Rules {
		rule, err := parseRule(ruleStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing rule: %v", err)
			return
		}
		rules = append(rules, rule)
	}

	// Match CIDR with Proxies
	var responseData []ResponseData
	for _, proxy := range config.Proxies {
		cidr := ""
		for _, rule := range rules {
			if rule.Name == proxy.Name {
				cidr = rule.CIDR
				break
			}
		}
		responseData = append(responseData, ResponseData{
			Name:     proxy.Name,
			Type:     proxy.Type,
			Server:   proxy.Server,
			Port:     proxy.Port,
			Username: proxy.Username,
			Password: proxy.Password,
			UUID:     proxy.UUID,
			Cipher:   proxy.Cipher,
			Address:  proxy.Address,
			Delay:    proxy.Delay,
			CIDR:     cidr,
			URL:      proxy.URL,
		})
	}

	response := Response{
		Code:  0,
		Msg:   "",
		Count: len(responseData),
		Data:  responseData,
	}

	// Convert to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
