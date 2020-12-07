package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:generate controller-gen object paths=$GOFILE

type CnvrgAppSpec struct {
	AppConfigs struct {
		CnvrgStorageUseIamRole string `json:"cnvrgStorageUseIamRole"`
		FeatureFlags           string `json:"featureFlags"`
		SMTP                   struct {
			Domain   string `json:"domain"`
			Password string `json:"password"`
			Port     string `json:"port"`
			Server   string `json:"server"`
			Username string `json:"username"`
		} `json:"smtp"`
	} `json:"appConfigs"`
	AppSecrets struct {
		CheckJobExpiration           string `json:"checkJobExpiration"`
		CnvrgStorageAccessKey        string `json:"cnvrgStorageAccessKey"`
		CnvrgStorageAzureAccessKey   string `json:"cnvrgStorageAzureAccessKey"`
		CnvrgStorageAzureAccountName string `json:"cnvrgStorageAzureAccountName"`
		CnvrgStorageAzureContainer   string `json:"cnvrgStorageAzureContainer"`
		CnvrgStorageBucket           string `json:"cnvrgStorageBucket"`
		CnvrgStorageProject          string `json:"cnvrgStorageProject"`
		CnvrgStorageRegion           string `json:"cnvrgStorageRegion"`
		CnvrgStorageSecretKey        string `json:"cnvrgStorageSecretKey"`
		CnvrgStorageType             string `json:"cnvrgStorageType"`
		DefaultComputeConfig         string `json:"defaultComputeConfig"`
		DefaultComputeName           string `json:"defaultComputeName"`
		ExtractTagsFromCmd           string `json:"extractTagsFromCmd"`
		MinioSseMasterKey            string `json:"minioSseMasterKey"`
		PassengerAppEnv              string `json:"passengerAppEnv"`
		RailsEnv                     string `json:"railsEnv"`
		RunJobsOnSelfCluster         string `json:"runJobsOnSelfCluster"`
		SecretKeyBase                string `json:"secretKeyBase"`
		SentryURL                    string `json:"sentryUrl"`
		StsIv                        string `json:"stsIv"`
		StsKey                       string `json:"stsKey"`
		UseStdout                    string `json:"useStdout"`
	} `json:"appSecrets"`
	Autoscaler struct {
		Enabled string `json:"enabled"`
	} `json:"autoscaler"`
	ClusterDomain string `json:"clusterDomain"`
	CnvrgApp      struct {
		CPU                       string `json:"cpu"`
		CustomAgentTag            string `json:"customAgentTag"`
		Edition                   string `json:"edition"`
		Enabled                   string `json:"enabled"`
		Image                     string `json:"image"`
		Intercom                  string `json:"intercom"`
		Memory                    string `json:"memory"`
		NodePort                  string `json:"nodePort"`
		Port                      string `json:"port"`
		Replicas                  int    `json:"replicas"`
		SidekiqCPU                string `json:"sidekiqCpu"`
		SidekiqMemory             string `json:"sidekiqMemory"`
		SidekiqReplicas           int    `json:"sidekiqReplicas"`
		SidekiqSearchkickCPU      string `json:"sidekiqSearchkickCpu"`
		SidekiqSearchkickMemory   string `json:"sidekiqSearchkickMemory"`
		SidekiqSearchkickReplicas int    `json:"sidekiqSearchkickReplicas"`
		SvcName                   string `json:"svcName"`
	} `json:"cnvrgApp"`
	CnvrgRouter struct {
		Enabled  string `json:"enabled"`
		Image    string `json:"image"`
		NodePort string `json:"nodePort"`
		Port     string `json:"port"`
		SvcName  string `json:"svcName"`
	} `json:"cnvrgRouter"`
	Conf struct {
		Enabled             string `json:"enabled"`
		GcpKeyfileMountPath string `json:"gcpKeyfileMountPath"`
		GcpKeyfileName      string `json:"gcpKeyfileName"`
		GcpStorageSecret    string `json:"gcpStorageSecret"`
	} `json:"conf"`
	Debug   string `json:"debug"`
	DryRun  string `json:"dryRun"`
	DumpDir string `json:"dumpDir"`
	Es      struct {
		CPULimit      string `json:"cpuLimit"`
		CPURequest    string `json:"cpuRequest"`
		Enabled       string `json:"enabled"`
		FsGroup       string `json:"fsGroup"`
		Image         string `json:"image"`
		JavaOpts      string `json:"javaOpts"`
		MaxMapImage   string `json:"maxMapImage"`
		MemoryLimit   string `json:"memoryLimit"`
		MemoryRequest string `json:"memoryRequest"`
		NodePort      string `json:"nodePort"`
		PatchEsNodes  string `json:"patchEsNodes"`
		Port          string `json:"port"`
		RunAsGroup    string `json:"runAsGroup"`
		RunAsUser     string `json:"runAsUser"`
		StorageClass  string `json:"storageClass"`
		StorageSize   string `json:"storageSize"`
		SvcName       string `json:"svcName"`
	} `json:"es"`
	Fluentd struct {
		ContainersPath string `json:"containersPath"`
		CPURequest     string `json:"cpuRequest"`
		Enabled        string `json:"enabled"`
		Image          string `json:"image"`
		JournalPath    string `json:"journalPath"`
		Journald       string `json:"journald"`
		MemoryLimit    string `json:"memoryLimit"`
		MemoryRequest  string `json:"memoryRequest"`
	} `json:"fluentd"`
	Grafana struct {
		Image    string `json:"image"`
		NodePort string `json:"nodePort"`
		Port     string `json:"port"`
		SvcName  string `json:"svcName"`
	} `json:"grafana"`
	Hostpath struct {
		CPULimit         string `json:"cpuLimit"`
		CPURequest       string `json:"cpuRequest"`
		Enabled          string `json:"enabled"`
		HostPath         string `json:"hostPath"`
		Image            string `json:"image"`
		MemoryLimit      string `json:"memoryLimit"`
		MemoryRequest    string `json:"memoryRequest"`
		NodeName         string `json:"nodeName"`
		StorageClassName string `json:"storageClassName"`
	} `json:"hostpath"`
	HTTPS struct {
		Cert                   string `json:"cert"`
		CertSecret             string `json:"certSecret"`
		Enabled                string `json:"enabled"`
		Key                    string `json:"key"`
		UseWildcardCertificate string `json:"useWildcardCertificate"`
	} `json:"https"`
	Ingress struct {
		Enabled string `json:"enabled"`
	} `json:"ingress"`
	IngressType string `json:"ingressType"`
	Istio       struct {
		Enabled               string `json:"enabled"`
		ExternalIP            string `json:"externalIp"`
		GwName                string `json:"gwName"`
		Hub                   string `json:"hub"`
		IngressSvcAnnotations string `json:"ingressSvcAnnotations"`
		MixerImage            string `json:"mixerImage"`
		OperatorImage         string `json:"operatorImage"`
		PilotImage            string `json:"pilotImage"`
		ProxyImage            string `json:"proxyImage"`
		Tag                   string `json:"tag"`
	} `json:"istio"`
	Kibana struct {
		CPULimit      string `json:"cpuLimit"`
		CPURequest    string `json:"cpuRequest"`
		Enabled       string `json:"enabled"`
		Image         string `json:"image"`
		MemoryLimit   string `json:"memoryLimit"`
		MemoryRequest string `json:"memoryRequest"`
		NodePort      string `json:"nodePort"`
		Port          string `json:"port"`
		SvcName       string `json:"svcName"`
	} `json:"kibana"`
	Minio struct {
		Enabled       string `json:"enabled"`
		Image         string `json:"image"`
		MemoryRequest string `json:"memoryRequest"`
		NodePort      string `json:"nodePort"`
		Port          string `json:"port"`
		Replicas      string `json:"replicas"`
		SharedStorage struct {
			Enabled          string `json:"enabled"`
			NfsServer        string `json:"nfsServer"`
			Path             string `json:"path"`
			StorageClassName string `json:"storageClassName"`
		} `json:"sharedStorage"`
		StorageClass string `json:"storageClass"`
		StorageSize  string `json:"storageSize"`
		SvcName      string `json:"svcName"`
	} `json:"minio"`
	Mpi struct {
		Enabled string `json:"enabled"`
	} `json:"mpi"`
	Nfs struct {
		CPULimit         string `json:"cpuLimit"`
		CPURequest       string `json:"cpuRequest"`
		Enabled          string `json:"enabled"`
		Image            string `json:"image"`
		MemoryLimit      string `json:"memoryLimit"`
		MemoryRequest    string `json:"memoryRequest"`
		Path             string `json:"path"`
		Provisioner      string `json:"provisioner"`
		Server           string `json:"server"`
		StorageClassName string `json:"storageClassName"`
	} `json:"nfs"`
	Nvidiadp struct {
		Enabled      string `json:"enabled"`
		Image        string `json:"image"`
		NodeSelector struct {
			Enabled string `json:"enabled"`
			Key     string `json:"key"`
			Value   string `json:"value"`
		} `json:"nodeSelector"`
	} `json:"nvidiadp"`
	Orchestrator string `json:"orchestrator"`
	Pg           struct {
		CPURequest    string `json:"cpuRequest"`
		Dbname        string `json:"dbname"`
		Enabled       string `json:"enabled"`
		FsGroup       string `json:"fsGroup"`
		Image         string `json:"image"`
		MemoryRequest string `json:"memoryRequest"`
		Pass          string `json:"pass"`
		Port          string `json:"port"`
		RunAsGroup    string `json:"runAsGroup"`
		RunAsUser     string `json:"runAsUser"`
		StorageClass  string `json:"storageClass"`
		StorageSize   string `json:"storageSize"`
		SvcName       string `json:"svcName"`
		User          string `json:"user"`
	} `json:"pg"`
	PgBackup struct {
		CronTime     string `json:"cronTime"`
		Enabled      string `json:"enabled"`
		Name         string `json:"name"`
		Path         string `json:"path"`
		ScriptPath   string `json:"scriptPath"`
		StorageClass string `json:"storageClass"`
		StorageSize  string `json:"storageSize"`
	} `json:"pgBackup"`
	PrivilegedSa string `json:"privilegedSa"`
	Prometheus   struct {
		AdapterImage          string `json:"adapterImage"`
		AlertManagerImage     string `json:"alertManagerImage"`
		ConfigReloaderImage   string `json:"configReloaderImage"`
		Enabled               string `json:"enabled"`
		Image                 string `json:"image"`
		KubeRbacProxyImage    string `json:"kubeRbacProxyImage"`
		KubeStateMetricsImage string `json:"kubeStateMetricsImage"`
		KubeletMetrics        struct {
			Port   string `json:"port"`
			Schema string `json:"schema"`
		} `json:"kubeletMetrics"`
		NodeExporterImage             string `json:"nodeExporterImage"`
		NodePort                      string `json:"nodePort"`
		NvidiaExporterImage           string `json:"nvidiaExporterImage"`
		OperatorImage                 string `json:"operatorImage"`
		Port                          string `json:"port"`
		PrometheusConfigReloaderImage string `json:"prometheusConfigReloaderImage"`
		StorageClass                  string `json:"storageClass"`
		StorageSize                   string `json:"storageSize"`
		SvcName                       string `json:"svcName"`
	} `json:"prometheus"`
	Rbac struct {
		Role               string `json:"role"`
		RoleBindingName    string `json:"roleBindingName"`
		ServiceAccountName string `json:"serviceAccountName"`
	} `json:"rbac"`
	Redis struct {
		Enabled string `json:"enabled"`
		Image   string `json:"image"`
		Limits  struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"limits"`
		Port     string `json:"port"`
		Requests struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"requests"`
		SvcName string `json:"svcName"`
	} `json:"redis"`
	Registry struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		URL      string `json:"url"`
		User     string `json:"user"`
	} `json:"registry"`
	SecurityMode string `json:"securityMode"`
	Seeder       struct {
		Image   string `json:"image"`
		SeedCmd string `json:"seedCmd"`
	} `json:"seeder"`
	UseHTTPS string `json:"useHttps"`
}

type CnvrgAppStatus struct{}

type CnvrgApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppSpec   `json:"spec,omitempty"`
	Status CnvrgAppStatus `json:"status,omitempty"`
}

type CnvrgAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgApp `json:"items"`
}
