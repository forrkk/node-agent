// agent.go
package main

import (
	"log"
	"runtime"
	"os"
//	"os/signal"
//	"syscall"
	"time"
)

func init() {
//	if configFile == ""
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
		} else {
			log.Fatalln("the node isn't initialised and token wasn't provided")
		}
	}
//	data := []byte(`{"token":"goUVPEJzYozhnXM4aJNG6kzS6YuKRUs8DLorouxxCmSb4hgB8ji6XEoMrnc22FjP"}`)
//	b, err := SendReq("POST", "https://api.wodby.com/api/v1/nodes/register", data, nil)
//	fmt.Println(string(b), err)
	//err := GetPortAvailability([]int{-1, 22, 80, 443, 70000, 8080, 0})
	//fmt.Println(err)
	//GetKey()

	WriteConfig()
	updateConfig()
//	go checkVersion()
//	select{}

//	if m, err := GetOsInfo(); err == nil {
//		fmt.Println(m["type"])
//		fmt.Println(m["arch"])
//		fmt.Println(m["pretty_name"])
//		fmt.Println(m["id"])
//		fmt.Println(m["version_id"])
//		fmt.Println(m["id_like"])
//		fmt.Println(m["kernel_ver"])
//		fmt.Println(m["init1"])
//	}
//	if m, err := GetDockerStatus(); err == nil {
//		fmt.Println(m["docker_path"])
//		fmt.Println(m["docker_version"])
//		fmt.Println(m["docker_running"])
//	}
//	p := []int{22,23}
//	fmt.Println(IsRoot())
//	fmt.Println(GetPortAvailability(p))
}
