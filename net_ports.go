package main

import (
	"net"
	"strconv"
)

func GetPortAvailability(ports []int) (map[string]bool) {
	res := make(map[string]bool)
	for _, p := range ports {
		if p <= 65535 && p > 0 {
			if ls, err := net.Listen("tcp4", "0.0.0.0:"+strconv.Itoa(p)); err != nil {
				res[strconv.Itoa(p)] = false
			} else {
				ls.Close()
				res[strconv.Itoa(p)] = true
			}
		}
	}
	return res
}