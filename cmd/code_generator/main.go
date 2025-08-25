package main

import (
	"clash-admin/pkg/util"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {

		code, err := util.Encrypt([]byte(os.Args[1]))
		if err != nil {
			fmt.Println("生成id出错", err)
			return
		}
		fmt.Println(code)
	} else {
		fmt.Println("请传入设备id")
	}
}
