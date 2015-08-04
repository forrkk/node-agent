package main

import (
	"os/exec"
)

func IsCommandExists(c string) ([]byte, error) {
	err := exec.Command("command").Run()
	if err != nil {
		return "", err
	}
	p, err := exec.Command("command", "-v", c)
	if err != nil {
		return "", err
	}
	return []byte(p), nil
}