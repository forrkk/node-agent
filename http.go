package main
import (
	"net/http"
	"io"
	"bytes"
	"fmt"
	"io/ioutil"
)

func SendReq(method, url string, data []byte, headers map[string]string) ([]byte, error) {
	var body io.Reader = bytes.NewReader(data)
	req, err := http.NewRequest(method, url, body)
	if err !=nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	fmt.Println("URL ", url)
	fmt.Println("Body ", string(data))
	fmt.Println("H ", req.Header)
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
		fmt.Println(string(b))
		return b, nil
	default:
		break
	}
	panic("unreachable")
}

func GetKey() {
	var data io.Reader
	data = bytes.NewReader([]byte(`{"token":"goUVPEJzYozhnXM4aJNG6kzS6YuKRUs8DLorouxxCmSb4hgB8ji6XEoMrnc22FjP"}`))
	method := "POST"
	url := "https://api.wodby.com/api/v1/nodes/register"
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, data)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	key := "GbpcvmP8uMH6YITTWgBVNvnRDtckVC8UcjY14fTuWbY"
	url = "https://api.wodby.com/api/v1/nodes/3653f2fc-2461-4f27-b5e1-5fc2c2d92a0a/version"
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+key)
	resp, _ = client.Do(req)
	fmt.Println(resp)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}