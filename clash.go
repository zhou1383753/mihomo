package main

import (
	"clash-admin/internal/buildflag"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"time"
)

func restartShell() {
	var command string = "/etc/init.d/mymihomo"
	if buildflag.Kernel == "clash" {
		command = "/etc/clash/start.sh"
	}
	go func() {
		time.Sleep(2 * time.Second)
		cmd := exec.Command(command, "restart")
		if err := cmd.Run(); err != nil {
			log.Printf("Command execution failed: %v", err)
		}
	}()

}

func restart() (string, error) {
	baseURL := fmt.Sprintf("%s/restart", CLASH_HOST)
	data := url.Values{}
	data.Set("key", "value")

	resp, err := http.PostForm(baseURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type DelayResponse struct {
	Delay int `json:"delay"`
}

func delay(proxy, testURL string, timeout int) (string, error) {
	var delayResp DelayResponse

	baseURL := fmt.Sprintf("%s/proxies/%s/delay", CLASH_HOST, proxy)
	reqURL := fmt.Sprintf("%s?url=%s&timeout=%d", baseURL, url.QueryEscape(testURL), timeout)

	resp, err := http.Get(reqURL)
	if err != nil {
		return "0", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", err
	}

	err = json.Unmarshal(body, &delayResp)
	if err != nil {
		return "0", err
	}
	dd := strconv.Itoa(delayResp.Delay)

	return dd, nil
}
