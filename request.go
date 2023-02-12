package v2rayaguard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	ActionLogin = "login"
)

var (
	authorization string
	serverBase    = "http://localhost:2017/api/"
	username      string
	password      string
)

// postJson
func postJson(action string, body map[string]interface{}) (*Response, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return request(action, "post", bytes.NewReader(buf))
}

// 登陆
// 一般用于在执行某个接口请求时发现未登陆或登陆过期时自动完成认证
func login() error {
	response, err := postJson(ActionLogin, map[string]interface{}{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}
	if response.isFailed() {
		panic(fmt.Errorf("login failed. response is:%v", response))
	}
	token, ok := response.Data["token"]
	if ok {
		authorization = strings.TrimSpace(fmt.Sprintf("%s", token))
		log.Printf("login success! Get Authorization:%s \n", authorization)
		return nil
	}
	return fmt.Errorf("login failed. response is:%v", response)
}

// 发送请求的封装
// action: 具体路径
// method: 请求方式:GET|POST|PUT 等
// body:   请求发送的数据
func request(action, method string, body io.Reader) (*Response, error) {
	url := serverBase + action
	method = strings.ToUpper(method)
	if strings.Contains(url, "/api/touch") {
		fmt.Printf(".")
	} else {
		log.Printf("[%s]%s", method, url)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 401 {
		log.Println("服务器提示未授权的请求,将自动登陆获得新的授权码")
		login()
		return request(action, method, body)
	}
	if res.StatusCode > 399 {
		return nil, fmt.Errorf("server internal error. code: %d", res.StatusCode)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := &Response{}
	err = json.Unmarshal(bytes, response)
	return response, err
}

func isRunning() bool {
	res, err := request("touch", "GET", nil)
	if err != nil || res.isFailed() {
		log.Println("RunningCheckErr:", err)
		return false
	}
	data := res.Data
	return data["running"].(bool)
}
