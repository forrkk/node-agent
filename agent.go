// agent.go
package main

import (
	"fmt"
)

func main() {
	if m, err := GetOsInfo(); err == nil {
		fmt.Println(m["type"])
		fmt.Println(m["arch"])
		fmt.Println(m["pretty_name"])
		fmt.Println(m["id"])
		fmt.Println(m["version_id"])
		fmt.Println(m["id_like"])
		fmt.Println(m["kernel_ver"])
		fmt.Println(m["init1"])
	}
	if m, err := GetDockerStatus(); err == nil {
		fmt.Println(m["docker_path"])
		fmt.Println(m["docker_version"])
		fmt.Println(m["docker_running"])
	}
	p := []int{22,23}
	fmt.Println(IsRoot())
	fmt.Println(GetPortAvailability(p))
}
