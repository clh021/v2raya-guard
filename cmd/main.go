package main

import (
	"fmt"
	"os"

	v "github.com/clh021/v2raya-guard"
)

// 支持启动时显示构建日期和构建版本
// 需要通过命令 ` go build -ldflags "-X main.build=`git rev-parse HEAD`" ` 打包
var build = "not set"

func main() {
	fmt.Printf("Build: %s\n", build)
	v.Run(
		os.Getenv("username"),
		os.Getenv("password"),
		os.Getenv("serverbaseurl"),
		os.Getenv("cronExp"),
	)
}
