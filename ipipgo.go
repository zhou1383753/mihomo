package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type writeToFileParam struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
	URL  string `json:"URL"`
}

func parseIpipgoUrl(rawUrl string) []string {
	var result []string
	rawUrl = strings.ReplaceAll(rawUrl, "\n", "%0A")
	resp, err := http.Get(rawUrl)

	if err != nil {
		fmt.Println("Error making GET request:", err)
		return result
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading response body:", err)
	}

	return result
}

func writeToFile(w http.ResponseWriter, r *http.Request) {
	var urls []writeToFileParam
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

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

	var ips []string
	index := 0

	for _, urlInfo := range urls {
		if ips == nil {
			ips = parseIpipgoUrl(urlInfo.URL)
		}
		for i, proxy := range config.Proxies {
			if proxy.Name == urlInfo.Name {
				ip := ips[index]
				config.Proxies[i].Server = strings.Split(ip, ":")[0]
				porti, err := strconv.Atoi(strings.Split(ip, ":")[1])
				if err != nil {
					fmt.Println("Error converting string to int:", err)
					return
				}
				config.Proxies[i].Port = porti
				config.Proxies[i].Type = urlInfo.Type
				config.Proxies[i].URL = urlInfo.URL
				index++
				break
			}
		}
	}

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
	response := Response{
		Code: 0,
		Msg:  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
