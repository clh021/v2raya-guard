package v2rayaguard

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func Run(_username, _password, _serverBaseUrl, cronExp string) {
	username = _username
	password = _password
	serverBase = _serverBaseUrl
	// log.Printf("username: %+v, password: %+v, serverBase: %+v cronExp: %+v \n", username, password, serverBase, cronExp)
	log.Printf("v2rayaguard running: %s\n", " comming soon.")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	cron := cron.New()
	log.Println("add task.name:restart,cron:" + cronExp)
	cron.AddFunc(cronExp, restart)
	log.Println("add task.name:running check,cron:0 * * * *")
	cron.AddFunc("*/1 * * * *", func() {
		if !isRunning() {
			log.Println("v2ray is not running. try start")
			restart()
		}
	})
	cron.Start()
	wg.Wait()
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
