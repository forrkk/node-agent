package main

import (
	"net"
	"strconv"
)

func GetPortAvailability(ports []int) (map[string]bool) {
	res := make(map[string]bool)
	for _, p := range ports {
		if ls, err := net.Listen("tcp4","0.0.0.0:"+strconv.Itoa(p)); err == nil {
			ls.Close()
			res[strconv.Itoa(p)] = true
		} else {
			res[strconv.Itoa(p)] = false
		}
	}
	return res
}