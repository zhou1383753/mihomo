package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed default.yaml
var defaultYaml []byte

func yamlReset(w http.ResponseWriter, r *http.Request) {

	var dYaml DefaultYaml
	yaml.Unmarshal(defaultYaml, &dYaml)

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

	// config.Proxies = dYaml.Proxies
	config.Proxies = nil
	// config.Rules = dYaml.Rules

	for i := 1; i <= len(config.Rules); i++ {
		config.Proxies = append(config.Proxies,
			Proxy{Name: fmt.Sprintf("Name%d", i), Type: "socks5", Server: "1.1.1.1", Port: 1080, Username: "user", Password: "pass", UDP: true})
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
