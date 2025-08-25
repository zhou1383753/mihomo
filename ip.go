package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type IPData struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ContinentCN    string `json:"continentCN"`
		CountryCN      string `json:"countryCN"`
		ZoneCN         string `json:"zoneCN"`
		ProvinceCN     string `json:"provinceCN"`
		CityCN         string `json:"cityCN"`
		CountyCN       string `json:"countyCN"`
		TownCN         string `json:"townCN"`
		IspCN          string `json:"ispCN"`
		ContinentID    int    `json:"continentID"`
		CountryID      int    `json:"countryID"`
		ZoneID         int    `json:"zoneID"`
		ProvinceID     int    `json:"provinceID"`
		CityID         int    `json:"cityID"`
		CountyID       int    `json:"countyID"`
		IspID          int    `json:"ispID"`
		TownID         int    `json:"townID"`
		Latitude       string `json:"latitude"`
		Longitude      string `json:"longitude"`
		OverseasRegion bool   `json:"overseasRegion"`
	} `json:"data"`
}

func fetchIPInfo(ip string) (string, error) {
	url := fmt.Sprintf("https://mesh.if.iqiyi.com/aid/ip/info?version=1.1.1&ip=%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch IP info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var ipdata IPData
	_ = json.Unmarshal(body, &ipdata)

	ipInfo := ipdata.Data.ContinentCN + "-" + ipdata.Data.CountryCN + "-" + ipdata.Data.ProvinceCN + "-" + ipdata.Data.CityCN + "-" + ipdata.Data.IspCN

	return ipInfo, nil
}

// func checkIP(ip string) string {

// 	info, err := qqdb.Find(ip)
// 	if err != nil {
// 		log.Fatalf("Failed to find IP info: %v", err)
// 	}
// 	return info.String()
// }

func handleCheckIP(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	ydaili := r.Form["ydaili"]
	var proxies []Proxy
	json.Unmarshal([]byte(ydaili[0]), &proxies)

	// 读取并解析config.yaml文件
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

	// 更新config.yaml中的相应字段
	for _, entry := range proxies {

		for i, proxy := range config.Proxies {
			if proxy.Name == entry.Name {
				info, _ := fetchIPInfo(entry.Server)
				config.Proxies[i].Address = info
				break
			}
		}
	}

	// 将更新后的配置写回到config.yaml文件
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

	response := Response{
		Code: 0,
		Msg:  "success",
	}

	// Convert to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
