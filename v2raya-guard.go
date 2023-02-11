package v2rayaguard

import (
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

func Run(_username, _password, _serverBaseUrl, cronExp string) {
	username = _username
	password = _password
	serverBase = _serverBaseUrl
	// fmt.Printf("username: %+v, password: %+v, serverBase: %+v cronExp: %+v \n", username, password, serverBase, cronExp)
	fmt.Printf("v2rayaguard running: %s\n", " comming soon.")
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
