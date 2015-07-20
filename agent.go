// agent.go
package main

import (
	"log"
	"runtime"
	"os"
	"time"
)

func init() {
	if len(os.Args) == 2 {
		config.RegToken = string(os.Args[1])
	}
}

func main() {
	switch runtime.GOOS {
	case "linux": break
	default: log.Fatalln("not implemented")
	}
	if _, err := GetOsInfo(); err != nil {
		log.Fatalln(err)
	}
	if !IsRoot() {
		log.Fatalln("must be root")
	}
	initConfig()
	if !config.Initialised {
		if config.RegToken != "" {
			sys, err := GetOsInfo()
			if err != nil {
				log.Fatalln(err)
			}
			switch sys["id"] {
			case "ubuntu": break
			default: log.Fatalln("not implemented")
			}
			if config.ReqPorts == nil {
				config.ReqPorts = append(config.ReqPorts, 22, 80, 443)
			}
			ps := GetPortAvailability(config.ReqPorts)
			for k, v := range ps {
				if !v {
					log.Fatalln("port %s isnot free, but necessary", k)
				}
			}
			_, err = GetDockerStatus()
			if err != nil {
				if err = installDocker(); err != nil {
					log.Fatalln(err)
				}
			}
			err = installEtcd()
			if err != nil {
				log.Fatalln(err)
			}
			err = installKubernetes()
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(5 * time.Second)
			err = initServices()
			if err != nil {
				log.Fatalln(err)
			}
			err = installKubelet()
			if err != nil {
				log.Fatalln(err)
			}
			resp, err := registerNode()
			if err != nil {
				log.Fatalln(err)
			}
			if resp.Error.Code != 0 {
				log.Fatalln(resp.Error.Code)
			}
			config.AuthKey = resp.Result.AuthKey
			config.NodeUUID = resp.Result.NodeUUID
			config.Initialised = true
			err = installRC()
			if err != nil {
				log.Fatalln(err)
			}
			updateConfig()
			err = SelfInstall()
			if err != nil {
				log.Fatalln(err)
			}
			err = SelfStart()
			if err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		} else {
			log.Fatalln("the node isn't initialised and token wasn't provided")
		}
	}
go checkVersion()
select{}
}
