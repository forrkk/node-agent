// agent.go
package main

import (
	"log"
	"runtime"
	"os"
	"time"
	"fmt"
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
	fmt.Print("Fetching system information: ")
	if _, err := GetOsInfo(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("OK")
	fmt.Print("Checking is root: ")
	if !IsRoot() {
		log.Fatalln("must be root")
	}
	fmt.Println("OK")
	initConfig()
	if !config.Initialised {
		if config.RegToken != "" {
			sys, err := GetOsInfo()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Print("Checking is system compatible: ")
			switch sys["id"] {
			case "ubuntu": break
			default: log.Fatalln("not implemented")
			}
			fmt.Println("OK")
			fmt.Print("Registering the server on Wodby: ")
			if config.NodeUUID == "" {
				resp, err := registerNode()
				if err != nil {
					log.Fatalln(err)
				}

				if resp.Error.Code != 0 {
					log.Fatalln(resp.Error.Message)
				}
				config.AuthKey = resp.Result.AuthKey
				config.NodeUUID = resp.Result.NodeUUID
				updateConfig()
				fmt.Println("OK")
			}
			fmt.Print("Checking network ports: ")
			if config.ReqPorts == nil {
				config.ReqPorts = append(config.ReqPorts, 2222, 80, 443, 4001, 8080, 6443, 10248)
			}
			ps := GetPortAvailability(config.ReqPorts)
			for k, v := range ps {
				if !v {
					log.Fatalln("port " + k + " isnot free, but necessary")
				}
			}
			fmt.Println("OK")
			fmt.Print("Checking docker status: ")
			_, err = GetDockerStatus()
			if err != nil {
				if err = installDocker(); err != nil {
					log.Fatalln(err)
				}
			}
			fmt.Println("OK")
			fmt.Print("Installing ETCD: ")
			err = installEtcd()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("OK")
			fmt.Print("Installing Google Kubernetes: ")
			err = installKubernetes()
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(5 * time.Second)
			fmt.Println("OK")
			fmt.Print("Installing Wodby Service: ")
			err = initServices()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("OK")
			err = installKubelet()
			if err != nil {
				log.Fatalln(err)
			}
			err = installRC()
			if err != nil {
				log.Fatalln(err)
			}
			err = addSSHKey()
			if err != nil {
				log.Fatalln(err)
			}
			err = SelfInstall()
			if err != nil {
				log.Fatalln(err)
			}
			err = SelfStart()
			if err != nil {
				log.Fatalln(err)
			}
			config.Initialised = true
			updateConfig()
			fmt.Println(
			`All required software has been installed.
			Now we are connecting this server to Wodby platform.
			Please proceed to the dashboard to see the progress.`)
			os.Exit(0)
		} else {
			log.Fatalln("the server isn't initialised and token wasn't provided")
		}
	}
go checkVersion()
select{}
}
