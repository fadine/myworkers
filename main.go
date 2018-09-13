/**
* @author Nguyen Quang Huy
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/fadine/myworkers/global"
	"github.com/fadine/myworkers/queue"
	"github.com/fadine/myworkers/workers"

	"github.com/olebedev/config"
	"github.com/urfave/cli"
)

var ()

func init() {
	//do something
}

func main() {

	_ = loadConfig("./conf/app.json")
	app_name, _ := global.Cfg.String("app_name")
	app_version, _ := global.Cfg.String("app_version")
	app_usage, _ := global.Cfg.String("app_usage")

	app := cli.NewApp()

	app.Name = app_name
	app.Version = app_version
	app.Usage = app_usage

	app.Commands = []cli.Command{
		{
			Name:   "email-engine",
			Usage:  "Worker to send mail",
			Action: actionHandler,
		},
	}

	app.Action = actionHandler

	fmt.Println(global.WorkerId)
	app.Run(os.Args)
}

func actionHandler(c *cli.Context) error {

	if len(c.Command.Name) == 0 {
		//panic("Command name not found")
		c.Command.Name = "email-engine"
	}

	queueName, _ := global.Cfg.String("queue-" + c.Command.Name)

	var onServing = false
	if len(queueName) == 0 {
		onServing = true
	}

	// Serving worker, running forever - no depend
	if onServing {

		msg := queue.ServingMessage{Action: c.Command.Name}
		wg := sync.WaitGroup{}
		for {
			select {
			case <-global.GracefulStop:
				return nil
			default:
				wg.Add(1)
				executeHandler(msg, &wg)
			}
		}

		wg.Wait()
		global.MustStop <- true
	} else {

		//Queue worker
		queueService := queue.GetService()

		go func() {

			wg := sync.WaitGroup{}
			for msg := range queueService.GetMessageFromQueue(queueName) {
				wg.Add(1)
				go executeHandler(msg, &wg)
			}

			wg.Wait()
			global.MustStop <- true
			close(global.MustStop)
		}()
	}

	<-global.MustStop
	return nil
}

func executeHandler(msg queue.IQueueMessage, group *sync.WaitGroup) {

	defer func() {
		if err := recover(); err != nil {
			group.Done()
		}
	}()

	handler := workers.GetHandler(msg.GetAction())
	handler.Process(msg, group)
}

func loadConfig(cf string) error {
	// read config file
	var raw []byte
	var err error

	//raw, err = ioutil.ReadFile("config/dev.json")
	raw, err = ioutil.ReadFile(cf)
	if err != nil {
		return err
	}

	jsonString := string(raw)

	cO, _ := config.ParseJson(jsonString)
	global.Cfg = *cO

	return nil
}
