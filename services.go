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
    wodbyETCDenp = `{
        "kind": "Endpoints",
        "apiVersion": "v1beta3",
        "metadata": {
          "name": "etcd"
        },
        "subsets": [
          {
            "addresses": [
              { "IP": "127.0.0.1" }
            ],
            "ports": [
              {
                "name": "etcd",
                "protocol": "TCP",
                "port": 4001
              }
            ]
          }
        ]
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
    wodbyGetDNSSvcIP = `#!/bin/sh
      while [ -z ${dnsSvcIp} ];do
        sleep 1
        dnsSvcIp=$(/opt/kubernetes/bin/kubectl --namespace=wodby get svc wodby-svc | grep wodby-svc | awk '{print $4}')
      done
      echo "${dnsSvcIp}" > /opt/wodby/etc/dns_svc_ip`
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
	err = ioutil.WriteFile("/opt/wodby/etc/etcd_enp.json", []byte(wodbyETCDenp), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "--namespace=wodby", "create", "-f", "/opt/wodby/etc/etcd_enp.json").Output()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/opt/wodby/etc/wodby_svc.json", []byte(wodbySvc), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "--namespace=wodby", "create", "-f", "/opt/wodby/etc/wodby_svc.json").Output()
	if err != nil {
		return err
	}
    err = ioutil.WriteFile("/opt/kubernetes/bin/getdnssvcip", []byte(wodbyGetDNSSvcIP), 0755)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/getdnssvcip").Output()
	if err != nil {
		return err
	}
    config.DNSIP = "34.34.35.35"
	return nil
}