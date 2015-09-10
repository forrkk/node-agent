package main

import (
    "os/exec"
    "strings"
)

func isUFW() bool {
    if !isUFWinstalled() || !isUFWActive() {
        return false
    }
    return true
}

func isUFWinstalled() bool {
    if err := exec.Command("ufw", "version").Run(); err == nil {
        return true
    }
    return false
}

func isUFWActive() bool {
    out, err := exec.Command("ufw", "status").Output()
    if err != nil || strings.Contains(strings.ToLower(out), "inactive") {
        return false
    }
    return true
}