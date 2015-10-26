package main

import (
	"io/ioutil"
	"os/exec"
    "os"
)

const (
	etcdUpstartScript = `description "ECTD key/value storage"
	start on started docker
	stop on stopping docker
	respawn
	kill timeout 10
	pre-start script
		docker pull quay.io/coreos/etcd
		docker rm -f etcd || true
	end script
	script
		exec docker run --rm -v /etc/ssl/certs/:/etc/ssl/certs \
		-v /opt/etcd/data:/default.etcd \
		-p 127.0.0.1:4001:2379 -p 127.0.0.1:2379:2379 \
		--hostname etcd --name etcd quay.io/coreos/etcd \
		--listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
		--advertise-client-urls http://127.0.0.1:2379,http://127.0.0.1:4001,http://etcd.wodby.wodby.local:4001,http://etcd.wodby.wodby.local:2379 \
		--initial-cluster-token wodby --initial-cluster-state new
	end script`
)

func installEtcd() error {
	err := ioutil.WriteFile("/etc/init/etcd.conf", []byte(etcdUpstartScript), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "etcd", "start").Output()
	if err != nil {
		return err
	}
	return nil
}

func UninstallETCD() error {
    err := os.RemoveAll("/etc/init/etcd.conf")
	if err != nil {
		return err
	}
	return nil
}