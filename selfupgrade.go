package main

import (
	"os"
	"net/http"
	"io"
	"io/ioutil"
	"os/exec"
)

const (
	wodbyNAUpstartScript = `description "Wodby Node Agent Service"
	start on runlevel [2345]
	stop on runlevel [!2345]
	respawn
	kill timeout 5
	script
		exec /opt/wodby/bin/node-agent
	end script`
)

func selfUpgrade(ver string) error {
	var url string = "https://github.com/Wodby/node-agent/releases/download/v"+ver+"/waiter"
	f, err := os.OpenFile("./node-agent", os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}
	d, err := http.Get(url)
	defer d.Body.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(f, d.Body)
	if err != nil {
		return err
	}
	return nil
}

func SelfInstall() error {
	err := ioutil.WriteFile("/etc/init/wodby-agent.conf", []byte(wodbyNAUpstartScript), 0644)
	if err != nil {
		return err
	}
	return nil
}

func SelfStart() error {
	_, err := exec.Command("service", "wodby-agent", "start").Output()
	if err != nil {
		return err
	}
	return nil
}