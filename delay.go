package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

func delayHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	targetURL := r.Form.Get("TargetURL")
	timeoutStr := r.Form.Get("timeout")

	go func(targetURL, timeoutStr string) {
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			http.Error(w, "Invalid timeout parameter", http.StatusBadRequest)
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

		for i := range config.Proxies {
			d, _ := delay(config.Proxies[i].Name, targetURL, timeout)
			config.Proxies[i].Delay = d
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

	}(targetURL, timeoutStr)

	response := Response{
		Code: 0,
		Msg:  "success",
	}

	// Convert to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
