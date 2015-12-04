package main

import (
	"os/user"
)

func IsRoot() bool {
	if u, _ := user.Current(); u.Uid != "0" {
		return false
	}
	return true
}
