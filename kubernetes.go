package main
import (
	"os"
	"os/exec"
	"net/http"
	"io"
	"io/ioutil"
	"text/template"
)

const (
	kubeApiUpstartScript = `description "Kubernetes API Server"
	start on started etcd
	stop on stopping etcd
	respawn
	kill timeout 5
	script
		exec /opt/kubernetes/bin/kube-apiserver \
		--bind-address=127.0.0.1 \
        --insecure_bind_address=127.0.0.1 \
        --allow_privileged=true \
        --insecure_port=8080 \
        --kubelet_https=true \
        --secure_port=6443 \
        --portal_net=172.20.0.0/14 \
        --etcd_servers=http://127.0.0.1:4001 \
        --logtostderr=true \
        --profiling=false \
        --authorization-mode=ABAC \
        --authorization-policy-file=/opt/kubernetes/etc/policy.json \
        --token-auth-file=/opt/kubernetes/etc/tokens.csv
	end script`
	kubeControllerUpstartScript = `description "Kubernetes Controller Manager"
	start on runlevel [2345]
	stop on runlevel [!2345]
	respawn
	kill timeout 5
	script
		exec /opt/kubernetes/bin/kube-controller-manager \
        --machines=127.0.0.1 \
        --master=127.0.0.1:8080 \
        --concurrent-endpoint-syncs=10 \
        --node-monitor-grace-period=3m \
        --node-monitor-period=60s \
        --node-startup-grace-period=5m \
        --node-memory=1Gi \
        --resource-quota-sync-period=1m \
        --logtostderr=true
	end script`
	kubeSchedulerUpstartScript = `description "Kubernetes Scheduler Service"
	start on runlevel [2345]
	stop on runlevel [!2345]
	respawn
	kill timeout 5
	script
		exec /opt/kubernetes/bin/kube-scheduler --master=127.0.0.1:8080 --profiling=false
	end script`
	kubeProxyUpstartScript = `description "Kubernetes Proxy Service"
	start on runlevel [2345]
	stop on runlevel [!2345]
	limit nofile 65536 65536
	respawn
	kill timeout 5
	script
		exec /opt/kubernetes/bin/kube-proxy \
		--master=127.0.0.1:8080 \
		--logtostderr=true
	end script`
	kubeKubeletUpstartScript = `description "Kubernetes Kubelet"
	start on runlevel [2345]
	stop on runlevel [!2345]
	respawn
	kill timeout 5
	script
		exec /opt/kubernetes/bin/kubelet \
        --address=127.0.0.1 \
        --port=10250 \
        --hostname_override=127.0.0.1 \
        --api_servers=127.0.0.1:8080 \
        --allow_privileged=true \
        --logtostderr=true \
        --healthz_bind_address=127.0.0.1 \
        --healthz_port=10248 \
        --maximum-dead-containers=10 \
        --maximum-dead-containers-per-container=0 \
        --enable-debugging-handlers=false \
        --global_housekeeping_interval=10m \
        --housekeeping_interval=1m \
        --max_housekeeping_interval=1h \
        --node-status-update-frequency=60s \
        --sync-frequency=15s \
        --cluster-domain=wodby.local \
        --cluster-dns={{.DNSIP}}
	end script`
)

func downloadKubernetes() error {
	_ = os.MkdirAll("/opt/kubernetes/bin", 0755)
	_ = os.MkdirAll("/opt/kubernetes/etc", 0644)
	files := []string{"kube-apiserver",
		"kubectl",
		"kube-controller-manager",
		"kube-scheduler",
		"kube-proxy",
		"kubelet"}
	for i := range files {
		file, err := os.OpenFile("/opt/kubernetes/bin/"+files[i], os.O_CREATE|os.O_RDWR, 0755)
		defer file.Close()
		if err != nil {
			return err
		}
		resp, err := http.Get("http://sfo.registry.wodby.com:81/releases/kubernetes/v0.16.0/bin/linux/amd64/"+files[i])
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile("/opt/kubernetes/etc/policy.json", []byte(`{"user":"wodby-agent"}`), 0600)
	if err != nil {
		return err
	}
	config.KubeToken = string(NewRnd(64, ""))
	err = ioutil.WriteFile("/opt/kubernetes/etc/tokens.csv", []byte(config.KubeToken + ",wodby-agent,1000"), 0600)
	if err != nil {
		return err
	}
	return nil
}

func installKubernetes() error {
	err := downloadKubernetes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/init/kube-apiserver.conf", []byte(kubeApiUpstartScript), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/init/kube-controller.conf", []byte(kubeControllerUpstartScript), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/init/kube-scheduler.conf", []byte(kubeSchedulerUpstartScript), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/init/kube-proxy.conf", []byte(kubeProxyUpstartScript), 0644)
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-apiserver", "start").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-controller", "start").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-scheduler", "start").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-proxy", "start").Output()
	if err != nil {
		return err
	}
	return nil
}

func installKubelet() error {
	f, err := os.OpenFile("/etc/init/kube-kubelet.conf", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = template.Must(template.New("kubelet").Parse(kubeKubeletUpstartScript)).Execute(f, config)
	if err != nil {
		return err
	}
	f.Close()
	_, err = exec.Command("service", "kube-kubelet", "start").Output()
	if err != nil {
		return err
	}
	return nil
}

func UninstallKubernetes() error {
	var err error
	_, err = exec.Command("service", "kube-kubelet", "stop").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-proxy", "stop").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-scheduler", "stop").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-controller", "stop").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("service", "kube-apiserver", "stop").Output()
	if err != nil {
		return err
	}
    err = os.RemoveAll("/etc/init/kube-kubelet.conf")
	if err != nil {
		return err
	}
    err = os.RemoveAll("/etc/init/kube-proxy.conf")
	if err != nil {
		return err
	}
    err = os.RemoveAll("/etc/init/kube-scheduler.conf")
	if err != nil {
		return err
	}
    err = os.RemoveAll("/etc/init/kube-controller.conf")
	if err != nil {
		return err
	}
    err = os.RemoveAll("/etc/init/kube-apiserver.conf")
	if err != nil {
		return err
	}
    err = os.RemoveAll("/opt/kubernetes")
	if err != nil {
		return err
	}
    return nil
}