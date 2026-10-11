package main

import (
	"fmt"
	"os"
	"runtime"
)

// version 由构建参数 -X main.version=... 写入，本地直接运行时是 dev。
var version = "dev"

func main() {
	// version 子命令用来核对下载到的是哪个版本、哪个系统。
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("hello-cli %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		return
	}
	fmt.Println("hello")
}
