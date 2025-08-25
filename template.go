package main

import (
	"bufio"
	"bytes"
	"clash-admin/internal/buildflag"
	"clash-admin/pkg/util"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"golang.org/x/crypto/ssh"
)

//go:embed templates/*.html
var tmplFS embed.FS

var (
	ROOT_PASS_HASH = "$1$B2nm3Tmt$L2x9wwmd6VmRpzE5b57bZ0:19951:0:99999:7:::"
)

func getRootPassword() string {
	b, err := os.ReadFile("/etc/shadow")
	if err != nil {
		return ""
	}
	reader := bytes.NewReader(b)
	bReader := bufio.NewScanner(reader)
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		text := bReader.Text()
		if text[:5] == "root:" {
			return text[5:]
		}
	}
	return ""
}

func sshDial(user, password, addr string) bool {
	if buildflag.NOSSH == "true" {
		return true
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         1 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return false
	}
	defer client.Close()
	return true
}

func templates() {
	// 解析嵌入的模板文件
	tmpl, err := template.ParseFS(tmplFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// 定义HTTP处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		activateFile, err := os.ReadFile(".activation")
		deviceID := util.GetDeviceID()

		if deviceID == "" || err != nil || !util.Decrypt([]byte(deviceID), string(activateFile)) ||
			(runtime.GOARCH != "arm" && runtime.GOARCH != "arm64" && util.SSH_AUTH_PASS != "1" && !sshDial("root", "zhou9921036", "10.0.0.254:22")) {
			// 如果ACTIVE没有值，跳转到/show页面
			http.Redirect(w, r, "/show", http.StatusFound)
			return
		}
		// 定义数据
		data := struct {
			Title   string
			Heading string
			Content string
		}{
			Title:   "My Page Title",
			Heading: "Welcome to My Page",
			Content: "This is a sample content.",
		}

		// 执行index.html模板并将结果写入HTTP响应
		err = tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) {
		// 定义数据
		data := struct {
			MaccEncrypted string
		}{
			MaccEncrypted: "EncryptedValueHere",
		}
		if err == nil {
			data.MaccEncrypted = util.GetDeviceID()
		} else {
			data.MaccEncrypted = err.Error()
		}
		// 执行show.html模板并将结果写入HTTP响应
		err = tmpl.ExecuteTemplate(w, "show.html", data)
		if err != nil {
			log.Printf("Error executing template show.html: %v", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

}
