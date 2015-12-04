package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

const (
	wodbyRC = `{
        "apiVersion": "v1",
        "kind": "ReplicationController",
        "metadata": {
          "labels": {
            "name": "wodby-svc"
          },
          "name": "wodby-svc"
        },
        "spec": {
          "replicas": 1,
          "selector": {
            "name": "wodby-svc"
          },
          "template": {
            "metadata": {
              "labels": {
                "name": "wodby-svc"
              }
            },
            "spec": {
              "containers": [
                {
                  "name": "edge",
                  "image": "wodby/edge:0.1",
                  "imagePullPolicy": "Always",
                  "securityContext": {
                    "privileged": true
                  },
				  "volumeMounts": [
					{
					  "name": "backups",
					  "mountPath": "/usr/share/nginx/html/backups"
					}
				  ],
                  "env": [
                    {
                      "name": "WODBY_NODE_UUID",
                      "value": "{{.NodeUUID}}"
                    }
                  ],
                  "ports": [
                    {
                      "containerPort": 80,
                      "hostPort": 80,
                      "protocol":"TCP",
                      "name": "http"
                    },
                    {
                      "containerPort": 443,
                      "hostPort": 443,
                      "protocol":"TCP",
                      "name": "https"
                    },
                    {
                      "containerPort": 22,
                      "hostPort": 2222,
                      "protocol":"TCP",
                      "name": "ssh"
                    }
                  ]
                },
                {
                  "name": "agent",
                  "image": "wodby/agent:0.1",
                  "imagePullPolicy": "Always",
                  "env": [
                    {
                      "name": "WODBY_NODE_UUID",
                      "value": "{{.NodeUUID}}"
                    },
                    {
                      "name": "WODBY_NAMESPACE",
                      "value": "services"
                    },
                    {
                      "name": "WODBY_ENVIRONMENT_TYPE",
                      "value": "{{.EnvType}}"
                    },
                    {
                      "name": "WODBY_ENVIRONMENT_NAME",
                      "value": "{{.EnvName}}"
                    },
                    {
                      "name": "WODBY_NODE_ACCESS_KEY",
                      "value": "{{.AuthKey}}"
                    },
                    {
                      "name": "WODBY_KUBE_TOKEN",
                      "value": "{{.KubeToken}}"
                    }
                  ],
                  "ports": [
                    {
                      "containerPort": 8125,
                      "protocol":"TCP",
                      "name": "agent"
                    }
                  ]
                },
                {
                  "name": "skydns",
                  "image": "wodby/skydns:0.1",
                  "imagePullPolicy": "Always",
                  "ports": [
                    {
                      "containerPort": 53,
                      "protocol":"TCP",
                      "name": "dns-tcp"
                    },
                    {
                      "containerPort": 53,
                      "protocol":"UDP",
                      "name": "dns-udp"
                    }
                  ]
                },
                {
                  "name": "mta",
                  "image": "wodby/mta-alpine",
                  "imagePullPolicy": "Always",
                  "volumeMounts": [
                    {
                      "name": "mta",
                      "mountPath": "/var/spool/smtpd"
                    }
                  ],
                  "ports": [
                    {
                      "name": "smtp",
                      "containerPort": 25,
                      "protocol": "TCP"
                    }
                  ]
                }
              ],
              "volumes": [
                {
                  "name": "es",
                  "hostPath": {
                    "path": "/srv/wodby/services/es"
                  }
                },
                {
                  "name": "redis",
                  "hostPath": {
                    "path": "/srv/wodby/services/redis"
                  }
                },
                {
                  "name": "mta",
                  "hostPath": {
                    "path": "/srv/wodby/services/mta"
                  }
                },
				{
				  "name": "backups",
				  "hostPath": {
					"path": "/srv/wodby/backups"
				  }
				}
              ]
            }
          }
        }
      }`
	wodbyK2Src = `{
	"apiVersion": "v1",
	"kind": "ReplicationController",
	"metadata": {
		"labels": {
			"name": "kube2sky"
		},
		"name": "kube2sky"
	},
	"spec": {
		"replicas": 1,
		"selector": {
			"name": "kube2sky"
		},
		"template": {
			"metadata": {
				"labels": {
					"name": "kube2sky"
				}
			},
			"spec": {
				"hostNetwork": true,
				"containers": [{
					"name": "kube2sky",
					"image": "wodby/kube2sky",
					"imagePullPolicy": "Always",
					"args": [
						"-domain=wodby.local",
						"-v=2"
					]
				}]
			}
		}
	}
}`
)

func installRC() error {
	f, err := os.OpenFile("/opt/wodby/etc/wodby_rc.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = template.Must(template.New("wodbyrc").Parse(wodbyRC)).Execute(f, config)
	if err != nil {
		return err
	}
	f.Close()
	err = ioutil.WriteFile("/opt/wodby/etc/kube2sky_rc.json", []byte(wodbyK2Src), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "--namespace=wodby", "create", "-f", "/opt/wodby/etc/wodby_rc.json").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("/opt/kubernetes/bin/kubectl", "--namespace=wodby", "create", "-f", "/opt/wodby/etc/kube2sky_rc.json").Output()
	if err != nil {
		return err
	}
	return nil
}
