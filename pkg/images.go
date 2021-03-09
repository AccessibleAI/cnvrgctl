package pkg

import (
	"github.com/briandowns/spinner"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func LoadCnvrgImages() []string {
	if viper.GetString("image") != "" {
		return []string{
			viper.GetString("image"),
		}
	}
	return []string{
		"docker.io/rancher/coreos-etcd:v3.4.14-rancher1",
		"docker.io/rancher/rke-tools:v0.1.72",
		"docker.io/rancher/hyperkube:v1.20.4-rancher1",
		"docker.io/rancher/pause:3.2",
		"docker.io/rancher/metrics-server:v0.4.1",
		"docker.io/calico/cni:v3.18.0",
		"docker.io/calico/kube-controllers:v3.18.0",
		"docker.io/calico/node:v3.18.0",
		"docker.io/calico/typha:v3.18.0",
		"docker.io/calico/pod2daemon-flexvol:v3.18.0",
		"docker.io/coredns/coredns:1.8.0",
		"docker.io/mpioperator/mpi-operator:v0.2.3",
		"docker.io/cnvrg/cnvrg-redis:v3.0.5.c2",
		"docker.io/istio/operator:1.8.1",
		"docker.io/istio/proxyv2:1.8.1",
		"docker.io/istio/pilot:1.8.1",
		"docker.io/minio/minio:RELEASE.2020-09-17T04-49-20Z",
		"docker.io/centos/postgresql-12-centos7",
		"docker.io/cnvrg/hyper-server:latest",
		"docker.io/cnvrg/cnvrg-boot:v0.25",
		"docker.io/cnvrg/cnvrg-boot:v0.24",
		"docker.io/grafana/grafana:7.2.0",
		"docker.io/strech/sidekiq-prometheus-exporter:0.1.13",
		"docker.io/nvidia/dcgm-exporter:1.7.2",
		"docker.io/nvidia/k8s-device-plugin:v0.7.0",
		"docker.io/cnvrg/cnvrg-es:v7.8.1",
		"docker.io/cnvrg/cnvrg-tools:v0.3",
		"docker.io/bitsensor/elastalert:3.0.0-beta.1",
		"docker.io/fluent/fluentd-kubernetes-daemonset:v1.11-debian-elasticsearch7-1",
		"docker.io/jimmidyson/configmap-reload:v0.3.0",
		"docker.io/cnvrg/app:master-5190-encode",
		"docker.io/cnvrg/cnvrg_cli:latest",
		"docker.io/cnvrg/cnvrg:v5.0",
		"docker.io/cnvrg/cnvrg_gpu:nvidia-tf-19.10",
		"docker.io/cnvrg/cnvrg-onprem",
		"docker.io/cnvrg/cnvrg-operator:2.13.0",
		"k8s.gcr.io/autoscaling/vpa-admission-controller:0.9.0",
		"k8s.gcr.io/autoscaling/vpa-recommender:0.9.0",
		"k8s.gcr.io/autoscaling/vpa-updater:0.9.0",
		"k8s.gcr.io/metrics-server/metrics-server:v0.3.7",
		"quay.io/tigera/operator:v1.15.0",
		"quay.io/tigera/key-cert-provisioner:release-v1.0",
		"quay.io/kubevirt/hostpath-provisioner",
		"quay.io/external_storage/nfs-client-provisioner:latest",
		"quay.io/coreos/prometheus-operator:v0.40.0",
		"quay.io/coreos/prometheus-config-reloader:v0.40.0",
		"quay.io/coreos/kube-rbac-proxy:v0.4.1",
		"quay.io/prometheus/prometheus:v2.22.2",
		"quay.io/prometheus/node-exporter:v0.18.1",
		"quay.io/coreos/kube-state-metrics:v1.9.5",
		"quay.io/kubernetes_incubator/nfs-provisioner:v2.3.0",
		"docker.elastic.co/kibana/kibana-oss:7.8.1",
	}
}

func ListAppImages(username string, password string) (images []string) {
	imagesLength := 10
	s := spinner.New(spinner.CharSets[27], 50*time.Millisecond)
	go StartSpinner(s, "fetching images list...", nil)
	url := "https://registry-1.docker.io/"
	hub, err := registry.New(url, username, password)
	if err != nil {
		logrus.Debug(err.Error())
		logrus.Fatal("error fetching images from docker hub")
	}
	tags, err := hub.Tags("cnvrg/app")
	tagRegex, _ := regexp.Compile("^master-\\d*-encode$")
	var filteredTags []int
	for _, tag := range tags {
		if tagRegex.MatchString(tag) {
			tagNumber, _ := strconv.Atoi(strings.Split(tag, "-")[1])
			filteredTags = append(filteredTags, tagNumber)
		}
	}
	logrus.Info(len(filteredTags))
	if len(filteredTags) == 0 {
		logrus.Fatal("no images available for upgrade")
	}
	sort.Sort(sort.Reverse(sort.IntSlice(filteredTags)))
	for i := 0; i < imagesLength; i++ {
		images = append(images, "docker.io/cnvrg/app:master-"+strconv.Itoa(filteredTags[i])+"-encode")
	}
	s.Stop()
	return
}
