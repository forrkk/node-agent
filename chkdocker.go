package main

import (
	"runtime"
	"errors"
	"os/exec"
	"strings"
)

func checkBinary() (map[string]string, error) {
	if err := exec.Command("command").Run(); err != nil {
		return nil, err
	}
	m := make(map[string]string)
	if res, err := exec.Command("command", "-v", "docker").Output(); err == nil {
		m["docker_path"] = strings.TrimSpace(string(res))
		if v, err := exec.Command(m["docker_path"], "-v").Output(); err == nil {
			m["docker_version"] = string(v)
		} else {
			m["docker_version"] = "unknown"
		}
	}
	return m, nil
}

func isRunning(p string) bool {
	if err := exec.Command(p, "info").Run(); err == nil {
		return true
	}
	return false
}

func GetDockerStatus() (map[string]string, error) {
	m := make(map[string]string)
	if runtime.GOOS != "linux" {
		return nil, errors.New("not implemented")
	}
	if v, err := checkBinary(); err == nil {
		m = v
		if isRunning(m["docker_path"]) {
			m["docker_running"] = "yes"
		} else {
			m["docker_running"] = "no"
		}
	} else {
		return nil, err
	}
	return m, nil
}