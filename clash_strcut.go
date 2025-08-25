package main

type Response struct {
	Code  int            `json:"code"`
	Msg   string         `json:"msg"`
	Count int            `json:"count"`
	Data  []ResponseData `json:"data"`
}

type Config struct {
	Proxies []Proxy  `yaml:"proxies"`
	Rules   []string `yaml:"rules"`
}

type Proxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	UUID     string `yaml:"uuid"`
	Cipher   string `yaml:"cipher"`
	Address  string `yaml:"address"`
	Delay    string `yaml:"delay"`
	UDP      bool   `yaml:"udp"`
	AlterID  string `yaml:"alterId"`
	URL      string `yaml:"url"`
}

func (p Proxy) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{
		"name":     p.Name,
		"type":     p.Type,
		"server":   p.Server,
		"port":     p.Port,
		"username": p.Username,
		"password": p.Password,
		"cipher":   p.Cipher,
		"uuid":     p.UUID,
		"address":  p.Address,
		"delay":    p.Delay,
		"udp":      p.UDP,
		"alterId":  p.AlterID,
		"url":      p.URL,
	}, nil
}

type Rule struct {
	Type string
	CIDR string
	Name string
}

type ResponseData struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Server   string `json:"Server"`
	Port     int    `json:"Port"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	UUID     string `json:"UUID"`
	Cipher   string `json:"Cipher"`
	Address  string `json:"Address"`
	Delay    string `json:"Delay"`
	CIDR     string `json:"CIDR"`
	URL      string `json:"URL"`
}

type DefaultYaml struct {
	Proxies []Proxy  `yaml:"proxies"`
	Rules   []string `yaml:"rules"`
}
