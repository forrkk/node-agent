package main

import (
	"crypto/rand"
)

func NewRnd(n int, s string) []byte {
	const stdChrs = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	if n <= 0 {
		return nil
	}
	var b = make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	if s == "" {
		s = stdChrs
	}
	for i := range b {
		b[i] = s[b[i]%byte(len(s))]
	}
	return b
}
