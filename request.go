package v2rayaguard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
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

func request(action, method string, body io.Reader) (*Response, error) {
	url := serverBase + action
	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if strings.TrimSpace(authorization) != "" {
		req.Header.Set("Authorization", authorization)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 401 {
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
		authorization = fmt.Sprintf("%s", token)
		log.Println("login success")
		return nil
	}
	return fmt.Errorf("login failed. response is:%v", response)
}
func postJson(action string, body map[string]interface{}) (*Response, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return request(action, "post", bytes.NewReader(buf))
}

func restart() {
	res, err := request("touch", "GET", nil)
	if err != nil {
		return
	}
	log.Println("get server list success")
	touch, _ := res.Data["touch"].(map[string]interface{})
	params := make([]map[string]interface{}, 0)
	serverList, _ := touch["servers"].([]interface{})
	for _, s := range serverList {
		params = append(params, map[string]interface{}{
			"id":    s.(map[string]interface{})["id"],
			"_type": "server",
		})
	}
	subscriptions, _ := touch["subscriptions"].([]interface{})
	for index, sub := range subscriptions {
		for _, s := range sub.(map[string]interface{})["servers"].([]interface{}) {
			params = append(params, map[string]interface{}{
				"id":    s.(map[string]interface{})["id"],
				"_type": "subscriptionServer",
				"sub":   index,
			})
		}
	}
	if len(params) == 0 {
		return
	}
	res, err = request("v2ray", "delete", nil)
	if err != nil || res.isFailed() {
		log.Printf("stop v2ray failed.%s\n", res.Data)
	} else {
		log.Println("stop v2ray success")
	}
	paramStr, _ := json.Marshal(params)
	action := "httpLatency?whiches=" + string(paramStr)
	res, err = request(action, "get", nil)
	if err != nil || res.isFailed() {
		log.Printf("ping servers failed.%s\n", res.Data)
	} else {
		log.Println("ping servers success")
	}
	whiches := res.Data["whiches"].([]interface{})
	pings := make(servers, 0)
	for _, which := range whiches {
		d := which.(map[string]interface{})
		p, e := time.ParseDuration(d["pingLatency"].(string))
		if e != nil {
			continue
		}
		pings = append(pings, &server{
			id:          int(d["id"].(float64)),
			_type:       d["_type"].(string),
			sub:         int(d["sub"].(float64)),
			pingLatency: p,
		})
	}
	sort.Sort(pings)
	connectedServers := touch["connectedServer"].([]interface{})
	for _, cs := range connectedServers {
		m := cs.(map[string]interface{})
		param := fmt.Sprintf(`{"id":%v,"_type":"%v","sub":%v,"outbound":"proxy"}`, m["id"], m["_type"], m["sub"])
		res, err = request("connection", "delete", strings.NewReader(param))
		if err != nil || res.isFailed() {
			log.Printf("unselect %s failed\n", param)
		} else {
			log.Printf("unselect %s success\n", param)
		}
	}
	for index, s := range pings {
		if index > 5 {
			break
		}
		param := fmt.Sprintf(`{"id":%d,"_type":"%s","sub":%d,"outbound":"proxy"}`, s.id, s._type, s.sub)
		res, err = request("connection", "post", strings.NewReader(param))
		if err != nil || res.isFailed() {
			log.Printf("select %s failed\n", param)
		} else {
			log.Printf("select %s success\n", param)

		}
	}
	res, err = request("v2ray", "post", nil)
	if err != nil || res.isFailed() {
		log.Printf("start v2ray failed.%s\n", res.Data)
	} else {
		log.Printf("start v2ray success.%s\n", res.Data)
	}
}

func isRunning() bool {
	res, err := request("touch", "GET", nil)
	if err != nil || res.isFailed() {
		return false
	}
	data := res.Data
	return data["running"].(bool)
}
