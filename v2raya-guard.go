package v2rayaguard

import "fmt"

// 支持启动时显示构建日期和构建版本
// 需要通过命令 ` go build -ldflags "-X main.build=`git rev-parse HEAD`" ` 打包
var build = "not set"


func Run() {
	fmt.Printf("Build: %s\n", build)
}