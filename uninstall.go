package main

import "os"

func Uninstall() error {
	var err error
	err = SelfUninstall()
	if err != nil {
		return err
	}
	err = UninstallETCD()
	if err != nil {
		return err
	}
	err = os.RemoveAll("/opt/etcd")
	if err != nil {
		return err
	}
	err = UninstallKubernetes()
	if err != nil {
		return err
	}
	err = UninstallDocker()
	if err != nil {
		return err
	}
	err = os.RemoveAll("/opt/wodby")
	if err != nil {
		return err
	}
	return nil
}
