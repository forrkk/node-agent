package main
import (
	"os"
	"io/ioutil"
	"os/exec"
)

const (
	wodbyNS = `{
"apiVersion":"v1beta3",
"kind": "Namespace",
"metadata": {"name": "wodby"},
"status": {"phase": "Active"},
"labels": {"name": "wodby"}
}`
	wodbyETCDsvc = `{
        "kind": "Service",
        "apiVersion": "v1beta3",
        "metadata": {
          "name": "etcd",
          "labels": {
            "name": "etcd"
          }
        },
        "spec": {
          "ports": [
            {
              "name": "etcd",
              "protocol": "TCP",
              "port": 4001,
              "targetPort": "etcd"
            }
          ]
        }
      }`
	wodbySvc = `{
        "kind": "Service",
        "apiVersion": "v1beta3",
        "metadata": {
          "name": "wodby-svc",
          "labels": {
            "name": "wodby-svc"
          }
        },
        "spec": {
          "selector": {
            "name": "wodby-svc"
          },
          "ports": [
            {
              "name": "agent",
              "protocol": "TCP",
              "port": 8125,
              "targetPort": "agent"
            },
            {
              "name": "dns-udp",
              "protocol": "UDP",
              "port": 53,
              "targetPort": "dns-udp"
            },
            {
              "name": "dns-tcp",
              "protocol": "TCP",
              "port": 53,
              "targetPort": "dns-tcp"
            }
          ]
        }
      }`
)

func initServices() error {
	_ = os.MkdirAll("/opt/wodby/etc", 0644)
	err := ioutil.WriteFile("/opt/wodby/etc/wodby_ns.json", []byte(wodbyNS), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "create", "-f", "/opt/wodby/etc/wodby_ns.json").Output()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/opt/wodby/etc/etcd_svc.json", []byte(wodbyETCDsvc), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "--namespace=wodby", "create", "-f", "/opt/wodby/etc/etcd_svc.json").Output()
	if err != nil {
		return err
	}
	return nil
	err = ioutil.WriteFile("/opt/wodby/etc/wodby_svc.json", []byte(wodbySvc), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "--namespace=wodby", "create", "-f", "/opt/wodby/etc/wodby_svc.json").Output()
	if err != nil {
		return err
	}
	return nil
}