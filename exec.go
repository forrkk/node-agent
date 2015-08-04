package main

import (
	"os/exec"
)

func IsCommandExists(c string) ([]byte, error) {
	err := exec.Command("command").Run()
	if err != nil {
		return nil, err
	}
	p, err := exec.Command("command", "-v", c).Output()
	if err != nil {
		return nil, err
	}
	return []byte(p), nil
}