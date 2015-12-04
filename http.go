package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func SendReq(method, url string, data []byte, headers map[string]string) ([]byte, error) {
	var body io.Reader = bytes.NewReader(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200, 201, 202:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		break
	}
	panic("unreachable")
}
