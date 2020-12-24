package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CnvrgAppSpec struct {
	Debug         string `json:"debug"`
	DryRun        string `json:"dryRun"`
	DumpDir       string `json:"dumpDir"`
	ClusterDomain string `json:"clusterDomain"`
	Orchestrator  string `json:"orchestrator"`
	PrivilegedSa  string `json:"privilegedSa"`
	SecurityMode  string `json:"securityMode"`
	IngressType   string `json:"ingressType"`
	Tenancy       struct {
		Enabled        string `json:"enabled"`
		DedicatedNodes string `json:"dedicatedNodes"`
		Cnvrg          struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"cnvrg"`
	} `json:"tenancy"`
	AppConfigs struct {
		CnvrgStorageUseIamRole string `json:"cnvrgStorageUseIamRole"`
		FeatureFlags           string `json:"featureFlags"`
		SMTP                   struct {
			Domain   string      `json:"domain"`
			Password string      `json:"password"`
			Port     interface{} `json:"port"`
			Server   string      `json:"server"`
			Username string      `json:"username"`
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
	CnvrgApp struct {
		CPU                       string      `json:"cpu"`
		CustomAgentTag            string      `json:"customAgentTag"`
		Edition                   string      `json:"edition"`
		Enabled                   string      `json:"enabled"`
		Image                     string      `json:"image"`
		Intercom                  string      `json:"intercom"`
		Memory                    string      `json:"memory"`
		NodePort                  interface{} `json:"nodePort"`
		Port                      interface{} `json:"port"`
		Replicas                  interface{} `json:"replicas"`
		SidekiqCPU                string      `json:"sidekiqCpu"`
		SidekiqMemory             string      `json:"sidekiqMemory"`
		SidekiqReplicas           interface{} `json:"sidekiqReplicas"`
		SidekiqSearchkickCPU      string      `json:"sidekiqSearchkickCpu"`
		SidekiqSearchkickMemory   string      `json:"sidekiqSearchkickMemory"`
		SidekiqSearchkickReplicas interface{} `json:"sidekiqSearchkickReplicas"`
		PassengerMaxPoolSize      interface{} `json:"passengerMaxPoolSize"`
		SvcName                   string      `json:"svcName"`
		ResourcesRequestEnabled   string      `json:"resourcesRequestEnabled"`
		EnableReadinessProbe      string      `json:"enableReadinessProbe"`
		Fixpg                     string      `json:"fixpg"`
		sidekiqPrestopHook        struct {
			Enabled     string      `json:"enabled"`
			KillTimeout interface{} `json:"killTimeout"`
		}
	} `json:"cnvrgApp"`
	CnvrgRouter struct {
		Enabled  string      `json:"enabled"`
		Image    string      `json:"image"`
		NodePort interface{} `json:"nodePort"`
		Port     interface{} `json:"port"`
		SvcName  string      `json:"svcName"`
	} `json:"cnvrgRouter"`
	Conf struct {
		Enabled             string `json:"enabled"`
		GcpKeyfileMountPath string `json:"gcpKeyfileMountPath"`
		GcpKeyfileName      string `json:"gcpKeyfileName"`
		GcpStorageSecret    string `json:"gcpStorageSecret"`
	} `json:"conf"`
	Es struct {
		CPULimit      string      `json:"cpuLimit"`
		CPURequest    string      `json:"cpuRequest"`
		Enabled       string      `json:"enabled"`
		FsGroup       string      `json:"fsGroup"`
		Image         string      `json:"image"`
		JavaOpts      string      `json:"javaOpts"`
		MaxMapImage   string      `json:"maxMapImage"`
		MemoryLimit   string      `json:"memoryLimit"`
		MemoryRequest string      `json:"memoryRequest"`
		NodePort      interface{} `json:"nodePort"`
		PatchEsNodes  string      `json:"patchEsNodes"`
		Port          interface{} `json:"port"`
		RunAsGroup    string      `json:"runAsGroup"`
		RunAsUser     string      `json:"runAsUser"`
		StorageClass  string      `json:"storageClass"`
		StorageSize   string      `json:"storageSize"`
		SvcName       string      `json:"svcName"`
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
		Image    string      `json:"image"`
		NodePort interface{} `json:"nodePort"`
		Port     interface{} `json:"port"`
		SvcName  string      `json:"svcName"`
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
		Enabled         string      `json:"enabled"`
		PerTryTimeout   string      `json:"perTryTimeout"`
		RetriesAttempts interface{} `json:"retriesAttempts"`
	} `json:"ingress"`
	Istio struct {
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
		CPULimit      string      `json:"cpuLimit"`
		CPURequest    string      `json:"cpuRequest"`
		Enabled       string      `json:"enabled"`
		Image         string      `json:"image"`
		MemoryLimit   string      `json:"memoryLimit"`
		MemoryRequest string      `json:"memoryRequest"`
		NodePort      interface{} `json:"nodePort"`
		Port          interface{} `json:"port"`
		SvcName       string      `json:"svcName"`
	} `json:"kibana"`
	Minio struct {
		Enabled       string      `json:"enabled"`
		Image         string      `json:"image"`
		MemoryRequest string      `json:"memoryRequest"`
		NodePort      interface{} `json:"nodePort"`
		Port          interface{} `json:"port"`
		Replicas      interface{} `json:"replicas"`
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
	Pg struct {
		CPURequest    string      `json:"cpuRequest"`
		Dbname        string      `json:"dbname"`
		Enabled       string      `json:"enabled"`
		FsGroup       string      `json:"fsGroup"`
		Image         string      `json:"image"`
		MemoryRequest string      `json:"memoryRequest"`
		Pass          string      `json:"pass"`
		Port          interface{} `json:"port"`
		RunAsGroup    string      `json:"runAsGroup"`
		RunAsUser     string      `json:"runAsUser"`
		StorageClass  string      `json:"storageClass"`
		StorageSize   string      `json:"storageSize"`
		SvcName       string      `json:"svcName"`
		User          string      `json:"user"`
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
		Port     interface{} `json:"port"`
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
	Seeder struct {
		Image   string `json:"image"`
		SeedCmd string `json:"seedCmd"`
	} `json:"seeder"`
	UseHTTPS string `json:"useHttps"`

	Monitoring struct {
		DcgmExporter struct {
			Enabled string      `yaml:"enabled"`
			Image   string      `yaml:"image"`
			Port    interface{} `yaml:"port"`
		} `yaml:"dcgmExporter"`
		DefaultServiceMonitors struct {
			Enabled string `yaml:"enabled"`
		} `yaml:"defaultServiceMonitors"`
		Enabled string `yaml:"enabled"`
		Grafana struct {
			Enabled  string      `yaml:"enabled"`
			Image    string      `yaml:"image"`
			NodePort interface{} `yaml:"nodePort"`
			Port     interface{} `yaml:"port"`
			SvcName  string      `yaml:"svcName"`
		} `yaml:"grafana"`
		IdleMetricsExporter struct {
			Enabled string `yaml:"enabled"`
		} `yaml:"idleMetricsExporter"`
		KubeStateMetrics struct {
			Enabled string `yaml:"enabled"`
			Image   string `yaml:"image"`
		} `yaml:"kubeStateMetrics"`
		MetricsServer struct {
			Enabled string `yaml:"enabled"`
			Image   string `yaml:"image"`
		} `yaml:"metricsServer"`
		MinioExporter struct {
			Enabled string `yaml:"enabled"`
			Image   string `yaml:"image"`
		} `yaml:"minioExporter"`
		NodeExporter struct {
			Enabled string      `yaml:"enabled"`
			Image   string      `yaml:"image"`
			Port    interface{} `yaml:"port"`
		} `yaml:"nodeExporter"`
		Prometheus struct {
			CPURequest    string      `yaml:"cpuRequest"`
			Enabled       string      `yaml:"enabled"`
			Image         string      `yaml:"image"`
			MemoryRequest string      `yaml:"memoryRequest"`
			NodePort      interface{} `yaml:"nodePort"`
			Port          interface{} `yaml:"port"`
			StorageClass  string      `yaml:"storageClass"`
			StorageSize   string      `yaml:"storageSize"`
			SvcName       string      `yaml:"svcName"`
		} `yaml:"prometheus"`
		PrometheusOperator struct {
			Enabled string `yaml:"enabled"`
			Images  struct {
				ConfigReloaderImage           string `yaml:"configReloaderImage"`
				KubeRbacProxyImage            string `yaml:"kubeRbacProxyImage"`
				OperatorImage                 string `yaml:"operatorImage"`
				PrometheusConfigReloaderImage string `yaml:"prometheusConfigReloaderImage"`
			} `yaml:"images"`
		} `yaml:"prometheusOperator"`
		SidekiqExporter struct {
			Enabled string `yaml:"enabled"`
			Image   string `yaml:"image"`
		} `yaml:"sidekiqExporter"`
	} `yaml:"monitoring"`
}

type CnvrgAppStatus struct{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CnvrgApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CnvrgAppSpec   `json:"spec,omitempty"`
	Status CnvrgAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CnvrgAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CnvrgApp `json:"items"`
}
