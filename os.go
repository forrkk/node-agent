package main

import (
	"runtime"
//	"log"
	"bufio"
	"os"
//	"fmt"
	"strings"
	"errors"
)

func getLinuxInfo () (map[string]string, error) {
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		m := make(map[string]string)
		for s.Scan() {
			a := strings.Split(s.Text(),"=")
			m[strings.ToLower(a[0])] = a[1]
		}
		if err := s.Err(); err != nil {
			return nil, err
		}
		return m, nil
	} else {
		return nil, err
	}
}

func GetOsInfo() (map[string]string, error) {
	if runtime.GOOS != "linux" {
		err := errors.New("is not implemented")
		return nil, err
	}
	if m, err := getLinuxInfo(); err == nil {
		return m, nil
	} else {
		return nil, err
	}
}

//const (
//	osType = runtime.GOOS
//	osArch = runtime.GOARCH
//)

///*var (
//	osName string
//	osVersion string
//)*/

//func GetOsType() string {
//	return osType
//}

//func GetOsArch() string {
//	return osArch
//}

//func detect_os() {
//	file, err := os.Open("/etc/os-release")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()

//	scanner := bufio.NewScanner(file)
//	os_arg := make(map[string]string)
//	for scanner.Scan() {
//		s := strings.Split(scanner.Text(),"=")
//		os_arg[strings.ToLower(s[0])] = s[1]
////    	if strings.HasPrefix(scanner.Text(), "NAME") {
////			fmt.Println(strings.ToLower(strings.TrimPrefix(scanner.Text(),"NAME=")))
////		}
////    	if strings.HasPrefix(scanner.Text(), "ID_LIKE") {
////			fmt.Println(strings.ToLower(strings.TrimPrefix(scanner.Text(),"ID_LIKE=")))
////		}
////    	if strings.HasPrefix(scanner.Text(), "ID") {
////			fmt.Println(strings.ToLower(strings.TrimPrefix(scanner.Text(),"ID=")))
////		}
////		if strings.HasPrefix(scanner.Text(), "VERSION_ID") {
////			fmt.Println(strings.ToLower(strings.TrimPrefix(scanner.Text(),"VERSION_ID=")))
////		}
//	}
//	fmt.Println(os_arg)

//	if err := scanner.Err(); err != nil {
//    log.Fatal(err)
//	}
//}
