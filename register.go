package main

import (
	"encoding/json"
)

type regResponse struct {
	Result struct {
		NodeUUID string `json:"node_uuid"`
		AuthKey  string `json:"access_key"`
	} `json:"result"`
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func registerNode() (*regResponse, error) {
	var j regResponse
	url := "https://api.wodby.com/api/v1/nodes/register"
	data := []byte(`{"token": "` + config.RegToken + `"}`)
	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := SendReq("POST", url, data, headers)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resp, &j); err != nil {
		return nil, err
	}
	return &j, nil
}
