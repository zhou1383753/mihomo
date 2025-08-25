package main

import (
	"clash-admin/pkg/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func parseXdaili(input string) [][]string {
	input = strings.Trim(input, "[]")
	input = strings.ReplaceAll(input, "\"", "")

	entries := strings.Split(input, ",")

	var result [][]string
	for _, entry := range entries {
		fields := strings.Split(entry, "|")
		result = append(result, fields)
	}

	return result
}

func updateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	xdaili := r.Form["xdaili"]
	proxyType := r.FormValue("type")
	xdailiParsed := parseXdaili(xdaili[0])
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

	// sk5/http代理格式为1.1.1.1|端口|账号|密码
	// ss代理格式为1.1.1.1|端口|加密方式|密码<br>
	// vm代理格式为1.1.1.1|端口|uuid|alterId
	switch proxyType {
	case "socks5":
		for _, entry := range xdailiParsed {
			if len(entry) < 5 {
				http.Error(w, "填写的格式错误", http.StatusBadRequest)
				return
			}
			name, server, port, username, password := entry[0], entry[1], entry[2], entry[3], entry[4]

			for i, proxy := range config.Proxies {
				if proxy.Name == name {
					porti, err := strconv.Atoi(port)
					if err != nil || !util.CheckPort(porti) {
						http.Error(w, "无效的端口号", http.StatusBadRequest)
						return
					}
					config.Proxies[i].Server = server
					config.Proxies[i].Port = porti
					config.Proxies[i].Username = username
					config.Proxies[i].Password = password
					config.Proxies[i].Type = proxyType
					break
				}
			}
		}

	case "ss":
		for _, entry := range xdailiParsed {
			if len(entry) < 5 {
				http.Error(w, "填写的格式错误", http.StatusBadRequest)
				return
			}
			name, server, port, cipher, password := entry[0], entry[1], entry[2], entry[3], entry[4]

			switch cipher {
			case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm", "chacha20-ietf-poly1305", "xchacha20-ietf-poly1305":
			case "aes-128-cfb", "aes-192-cfb", "aes-256-cfb", "rc4-md5", "chacha20-ietf", "xchacha20":
			case "aes-128-ctr", "aes-192-ctr", "aes-256-ctr":
			default:
				http.Error(w, "无效的加密方式", http.StatusBadRequest)
				return
			}

			for i, proxy := range config.Proxies {
				if proxy.Name == name {
					porti, err := strconv.Atoi(port)
					if err != nil || !util.CheckPort(porti) {
						http.Error(w, "无效的端口号", http.StatusBadRequest)
						return
					}
					config.Proxies[i].Server = server
					config.Proxies[i].Port = porti
					config.Proxies[i].Cipher = cipher
					config.Proxies[i].Password = password
					config.Proxies[i].Type = proxyType
					break
				}
			}
		}

	case "vmess":
		for _, entry := range xdailiParsed {
			if len(entry) < 6 {
				http.Error(w, "填写的格式错误", http.StatusBadRequest)
				return
			}
			name, server, port, uuid, alterId, cipher := entry[0], entry[1], entry[2], entry[3], entry[4], entry[5]

			if !util.CheckUUID(uuid) {
				http.Error(w, "无效的uuid", http.StatusBadRequest)
				return
			}
			if !util.CheckAlterID(alterId) {
				http.Error(w, "无效的alterId", http.StatusBadRequest)
				return
			}

			switch cipher {
			case "auto", "aes-128-gcm", "chacha20-poly1305", "none":
			default:
				http.Error(w, "无效的加密方式", http.StatusBadRequest)
				return
			}

			for i, proxy := range config.Proxies {
				if proxy.Name == name {
					porti, err := strconv.Atoi(port)
					if err != nil || !util.CheckPort(porti) {
						http.Error(w, "无效的端口号", http.StatusBadRequest)
						return
					}
					config.Proxies[i].Server = server
					config.Proxies[i].Port = porti
					config.Proxies[i].UUID = uuid
					config.Proxies[i].AlterID = alterId
					config.Proxies[i].Type = proxyType
					config.Proxies[i].Cipher = cipher
					break
				}
			}
		}

	case "http":
		for _, entry := range xdailiParsed {
			if len(entry) < 5 {
				http.Error(w, "填写的格式错误", http.StatusBadRequest)
				return
			}
			name, server, port, username, password := entry[0], entry[1], entry[2], entry[3], entry[4]

			for i, proxy := range config.Proxies {
				if proxy.Name == name {
					porti, err := strconv.Atoi(port)
					if err != nil || !util.CheckPort(porti) {
						http.Error(w, "无效的端口号", http.StatusBadRequest)
						return
					}
					config.Proxies[i].Server = server
					config.Proxies[i].Port = porti
					config.Proxies[i].Username = username
					config.Proxies[i].Password = password
					config.Proxies[i].Type = proxyType
					break
				}
			}
		}

	default:
		http.Error(w, "Invalid proxy type", http.StatusBadRequest)
		return
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

	// Convert to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
