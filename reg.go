package main
import (
	"net/http"
	"io"
	"bytes"
	"fmt"
	"io/ioutil"
)


func GetKey() {
	var data io.Reader
	data = bytes.NewReader([]byte(`{"token":"dv4L56c3QTuZKFCpk2BL44QNrL2cLALLynTXvZfwnKGQrLgtVvQpNPEQewQjG3sa"}`))
	method := "POST"
	url := "https://api.wodby.com/api/v1/nodes/register"
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, data)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}