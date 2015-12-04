package main

import (
	"encoding/json"
	"time"
)

const Version = "0.0.1"

type Ver struct {
	Version string `json:"result"`
}

func getVersion() (string, error) {
	var j Ver
	url := "https://api.wodby.com/api/v1/nodes/" + config.NodeUUID + "/version"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.AuthKey,
	}
	resp, err := SendReq("GET", url, nil, headers)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(resp, &j); err != nil {
		return "", err
	}
	return j.Version, nil
}

func checkVersion() {
	var ver string
	for {
		time.Sleep(60 * time.Second)
		ver, _ = getVersion()
		if ver != "" && ver != Version {
			selfUpgrade(ver)
		}
	}
}
