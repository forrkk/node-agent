package main

import (
	"runtime"
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"errors"
)

func getLinuxInfo () (map[string]string, error) {
	m := make(map[string]string)
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			a := strings.Split(s.Text(),"=")
			m[strings.ToLower(a[0])] = a[1]
		}
		if err := s.Err(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	if v, err := ioutil.ReadFile("/proc/version"); err == nil {
		m["kernel_ver"] = string(v)
	}
	return m, nil
}

func getInitInfo () (map[string]string, error) {
	m := make(map[string]string)
	if v, err := ioutil.ReadFile("/proc/1/comm"); err == nil {
		if strings.Contains(string(v), "systemd") {
			m["init1"] = "systemd"
		} else {
			m["init1"] = "upstart"
		}
	} else {
		return nil, err
	}
	return m, nil
}

func GetOsInfo() (map[string]string, error) {
	r := make(map[string]string)
	if runtime.GOOS != "linux" {
		err := errors.New("is not implemented")
		return nil, err
	}
	if m, err := getLinuxInfo(); err == nil {
		r = m
	} else {
		return nil, err
	}
	if m, err := getInitInfo(); err == nil {
		for k, v := range m {
			r[k] = v
		}
	}
	r["type"] = runtime.GOOS
	r["arch"] = runtime.GOARCH
	return r, nil
}