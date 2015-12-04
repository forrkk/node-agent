package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	RegToken    string `toml:"reg_token"`
	NodeUUID    string `toml:"node_uuid"`
	AuthKey     string `toml:"access_key"`
	HeartBeat   uint16 `toml:"update_interval"`
	Initialised bool   `toml:"initialised"`
	ReqPorts    []int  `toml:"requred_ports"`
	DNSIP       string `toml:"dns_svc_ip"`
	ETCDIP      string `toml:"etcd_svc_ip"`
	KubeToken   string `toml:"kube_token"`
	EnvName     string `toml:"env_name"`
	EnvType     string `toml:"env_type"`
}

var (
	config     Config
	configFile string
)

const (
	defaultConfigDir  = "/opt/wodby/etc/"
	defaultConfigFile = defaultConfigDir + "config.toml"
)

func initConfig() {
	if configFile == "" {
		if _, err := os.Stat(defaultConfigDir); os.IsNotExist(err) {
			if err = os.MkdirAll(defaultConfigDir, 0640); err != nil {
				log.Fatalln(err)
			}
			if _, err = os.OpenFile(defaultConfigFile, os.O_CREATE, 0640); err != nil {
				log.Fatalln(err)
			}
		}
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		} else {
			_, err = os.OpenFile(defaultConfigFile, os.O_CREATE, 0640)
			if err != nil {
				log.Fatalln(err)
			}
			configFile = defaultConfigFile
		}
	}
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatalln(err)
	}
}

func updateConfig() (bool, error) {
	buf, err := os.OpenFile(configFile, os.O_RDWR, 0640)
	if err != nil {
		return false, err
	}
	defer buf.Close()
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return false, err
	}
	return true, nil
}
