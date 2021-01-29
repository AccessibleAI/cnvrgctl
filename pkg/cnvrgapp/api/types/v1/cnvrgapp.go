package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
// controller-gen object paths=./...
type CnvrgAppSpec struct {
	CnvrgApp struct {
		Mpi          struct {
			Enabled              string `json:"enabled"`
			Image                string `json:"image"`
			KubectlDeliveryImage string `json:"kubectlDeliveryImage"`
			Registry             struct {
				Name     string `json:"name"`
				URL      string `json:"url"`
				User     string `json:"user"`
				Password string `json:"password"`
			} `json:"registry"`
		} `json:"mpi"`
		Redis struct {
			Enabled string `json:"enabled"`
			Image   string `json:"image"`
			SvcName string `json:"svcName"`
			Port    int    `json:"port"`
			Limits  struct {
				CPU    int    `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"limits"`
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests"`
		} `json:"redis"`
		Networking struct {
			Enabled     string `json:"enabled"`
			IngressType string `json:"ingressType"`
			HTTPS       struct {
				Enabled    string `json:"enabled"`
				Cert       string `json:"cert"`
				Key        string `json:"key"`
				CertSecret string `json:"certSecret"`
			} `json:"https"`
			Istio struct {
				Enabled                  string `json:"enabled"`
				OperatorImage            string `json:"operatorImage"`
				Hub                      string `json:"hub"`
				Tag                      string `json:"tag"`
				ProxyImage               string `json:"proxyImage"`
				MixerImage               string `json:"mixerImage"`
				PilotImage               string `json:"pilotImage"`
				GwName                   string `json:"gwName"`
				ExternalIP               string `json:"externalIp"`
				IngressSvcAnnotations    string `json:"ingressSvcAnnotations"`
				IngressSvcExtraPorts     string `json:"ingressSvcExtraPorts"`
				LoadBalancerSourceRanges string `json:"loadBalancerSourceRanges"`
			} `json:"istio"`
			Ingress struct {
				Enabled         string `json:"enabled"`
				Timeout         string `json:"timeout"`
				RetriesAttempts int    `json:"retriesAttempts"`
				PerTryTimeout   string `json:"perTryTimeout"`
			} `json:"ingress"`
		} `json:"networking"`
		Vpa struct {
			Enabled string `json:"enabled"`
			Images  struct {
				AdmissionImage   string `json:"admissionImage"`
				RecommenderImage string `json:"recommenderImage"`
				UpdaterImage     string `json:"updaterImage"`
			} `json:"images"`
		} `json:"vpa"`
		Minio struct {
			Enabled       string `json:"enabled"`
			Replicas      int    `json:"replicas"`
			Image         string `json:"image"`
			Port          int    `json:"port"`
			StorageSize   string `json:"storageSize"`
			SvcName       string `json:"svcName"`
			NodePort      int    `json:"nodePort"`
			StorageClass  string `json:"storageClass"`
			CPURequest    int    `json:"cpuRequest"`
			MemoryRequest string `json:"memoryRequest"`
			SharedStorage struct {
				Enabled        string `json:"enabled"`
				ConsistentHash struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"consistentHash"`
			} `json:"sharedStorage"`
		} `json:"minio"`
		Storage struct {
			Enabled         string `json:"enabled"`
			CcpStorageClass string `json:"ccpStorageClass"`
			Hostpath        struct {
				Enabled          string `json:"enabled"`
				Image            string `json:"image"`
				HostPath         string `json:"hostPath"`
				StorageClassName string `json:"storageClassName"`
				NodeName         string `json:"nodeName"`
				CPURequest       string `json:"cpuRequest"`
				MemoryRequest    string `json:"memoryRequest"`
				CPULimit         string `json:"cpuLimit"`
				MemoryLimit      string `json:"memoryLimit"`
				ReclaimPolicy    string `json:"reclaimPolicy"`
				DefaultSc        string `json:"defaultSc"`
			} `json:"hostpath"`
			Nfs struct {
				Enabled          string `json:"enabled"`
				Image            string `json:"image"`
				Provisioner      string `json:"provisioner"`
				StorageClassName string `json:"storageClassName"`
				Server           string `json:"server"`
				Path             string `json:"path"`
				CPURequest       string `json:"cpuRequest"`
				MemoryRequest    string `json:"memoryRequest"`
				CPULimit         string `json:"cpuLimit"`
				MemoryLimit      string `json:"memoryLimit"`
				ReclaimPolicy    string `json:"reclaimPolicy"`
				DefaultSc        string `json:"defaultSc"`
			} `json:"nfs"`
		} `json:"storage"`
		Pg struct {
			Enabled        string `json:"enabled"`
			SecretName     string `json:"secretName"`
			Image          string `json:"image"`
			Port           int    `json:"port"`
			StorageSize    string `json:"storageSize"`
			SvcName        string `json:"svcName"`
			Dbname         string `json:"dbname"`
			Pass           string `json:"pass"`
			User           string `json:"user"`
			RunAsUser      int    `json:"runAsUser"`
			FsGroup        int    `json:"fsGroup"`
			StorageClass   string `json:"storageClass"`
			CPURequest     int    `json:"cpuRequest"`
			MemoryRequest  string `json:"memoryRequest"`
			MaxConnections int    `json:"maxConnections"`
			SharedBuffers  string `json:"sharedBuffers"`
		} `json:"pg"`
		PgBackup struct {
			StorageSize  string `json:"storageSize"`
			Enabled      string `json:"enabled"`
			Name         string `json:"name"`
			Path         string `json:"path"`
			ScriptPath   string `json:"script_path"`
			StorageClass string `json:"storageClass"`
			CronTime     string `json:"cronTime"`
		} `json:"pg_backup"`
		CnvrgApp struct {
			Replicas                int    `json:"replicas"`
			Enabled                 string `json:"enabled"`
			Image                   string `json:"image"`
			Port                    int    `json:"port"`
			CPU                     int    `json:"cpu"`
			Memory                  string `json:"memory"`
			SvcName                 string `json:"svcName"`
			Fixpg                   string `json:"fixpg"`
			NodePort                int    `json:"nodePort"`
			PassengerMaxPoolSize    int    `json:"passengerMaxPoolSize"`
			EnableReadinessProbe    string `json:"enableReadinessProbe"`
			InitialDelaySeconds     int    `json:"initialDelaySeconds"`
			ReadinessPeriodSeconds  int    `json:"readinessPeriodSeconds"`
			ReadinessTimeoutSeconds int    `json:"readinessTimeoutSeconds"`
			FailureThreshold        int    `json:"failureThreshold"`
			ResourcesRequestEnabled string `json:"resourcesRequestEnabled"`
			Sidekiq                 struct {
				Enabled  string `json:"enabled"`
				Split    string `json:"split"`
				CPU      string `json:"cpu"`
				Memory   string `json:"memory"`
				Replicas int    `json:"replicas"`
			} `json:"sidekiq"`
			Searchkiq struct {
				Enabled  string `json:"enabled"`
				CPU      string `json:"cpu"`
				Memory   string `json:"memory"`
				Replicas int    `json:"replicas"`
			} `json:"searchkiq"`
			Systemkiq struct {
				Enabled  string `json:"enabled"`
				CPU      string `json:"cpu"`
				Memory   string `json:"memory"`
				Replicas int    `json:"replicas"`
			} `json:"systemkiq"`
			KiqPrestopHook struct {
				Enabled     string `json:"enabled"`
				KillTimeout int    `json:"killTimeout"`
			} `json:"kiqPrestopHook"`
			Hyper struct {
				Enabled                 string `json:"enabled"`
				Image                   string `json:"image"`
				Port                    int    `json:"port"`
				Replicas                int    `json:"replicas"`
				NodePort                int    `json:"nodePort"`
				SvcName                 string `json:"svcName"`
				Token                   string `json:"token"`
				CPURequest              string `json:"cpuRequest"`
				MemoryRequest           string `json:"memoryRequest"`
				CPULimit                int    `json:"cpuLimit"`
				MemoryLimit             string `json:"memoryLimit"`
				EnableReadinessProbe    string `json:"enableReadinessProbe"`
				ReadinessPeriodSeconds  int    `json:"readinessPeriodSeconds"`
				ReadinessTimeoutSeconds int    `json:"readinessTimeoutSeconds"`
			} `json:"hyper"`
			Seeder struct {
				Image           string `json:"image"`
				SeedCmd         string `json:"seedCmd"`
				CreateBucketCmd string `json:"createBucketCmd"`
			} `json:"seeder"`
			Conf struct {
				GcpStorageSecret             string `json:"gcpStorageSecret"`
				GcpKeyfileMountPath          string `json:"gcpKeyfileMountPath"`
				GcpKeyfileName               string `json:"gcpKeyfileName"`
				JobsStorageClass             string `json:"jobsStorageClass"`
				FeatureFlags                 string `json:"featureFlags"`
				SentryURL                    string `json:"sentryUrl"`
				SecretKeyBase                string `json:"secretKeyBase"`
				StsIv                        string `json:"stsIv"`
				StsKey                       string `json:"stsKey"`
				RedisURL                     string `json:"redisUrl"`
				PassengerAppEnv              string `json:"passengerAppEnv"`
				RailsEnv                     string `json:"railsEnv"`
				RunJobsOnSelfCluster         string `json:"runJobsOnSelfCluster"`
				DefaultComputeConfig         string `json:"defaultComputeConfig"`
				DefaultComputeName           string `json:"defaultComputeName"`
				UseStdout                    string `json:"useStdout"`
				ExtractTagsFromCmd           string `json:"extractTagsFromCmd"`
				CheckJobExpiration           string `json:"checkJobExpiration"`
				CnvrgStorageType             string `json:"cnvrgStorageType"`
				CnvrgStorageBucket           string `json:"cnvrgStorageBucket"`
				CnvrgStorageAccessKey        string `json:"cnvrgStorageAccessKey"`
				CnvrgStorageSecretKey        string `json:"cnvrgStorageSecretKey"`
				CnvrgStorageEndpoint         string `json:"cnvrgStorageEndpoint"`
				MinioSseMasterKey            string `json:"minioSseMasterKey"`
				CnvrgStorageAzureAccessKey   string `json:"cnvrgStorageAzureAccessKey"`
				CnvrgStorageAzureAccountName string `json:"cnvrgStorageAzureAccountName"`
				CnvrgStorageAzureContainer   string `json:"cnvrgStorageAzureContainer"`
				CnvrgStorageRegion           string `json:"cnvrgStorageRegion"`
				CnvrgStorageProject          string `json:"cnvrgStorageProject"`
				CustomAgentTag               string `json:"customAgentTag"`
				Intercom                     string `json:"intercom"`
				CnvrgJobUID                  string `json:"cnvrgJobUid"`
				Ldap                         struct {
					Enabled       string `json:"enabled"`
					Host          string `json:"host"`
					Port          string `json:"port"`
					Account       string `json:"account"`
					Base          string `json:"base"`
					AdminUser     string `json:"adminUser"`
					AdminPassword string `json:"adminPassword"`
					Ssl           string `json:"ssl"`
				} `json:"ldap"`
				Registry struct {
					Name     string `json:"name"`
					URL      string `json:"url"`
					User     string `json:"user"`
					Password string `json:"password"`
				} `json:"registry"`
				Rbac struct {
					Role               string `json:"role"`
					ServiceAccountName string `json:"serviceAccountName"`
					RoleBindingName    string `json:"roleBindingName"`
				} `json:"rbac"`
				SMTP struct {
					Server   string `json:"server"`
					Port     string `json:"port"`
					Username string `json:"username"`
					Password string `json:"password"`
					Domain   string `json:"domain"`
				} `json:"smtp"`
			} `json:"conf"`
			CnvrgRouter struct {
				Enabled  string `json:"enabled"`
				Image    string `json:"image"`
				SvcName  string `json:"svcName"`
				NodePort int    `json:"nodePort"`
				Port     int    `json:"port"`
			} `json:"cnvrgRouter"`
		} `json:"cnvrgApp"`
		Monitoring struct {
			Enabled            string `json:"enabled"`
			PrometheusOperator struct {
				Enabled string `json:"enabled"`
				Images  struct {
					OperatorImage                 string `json:"operatorImage"`
					ConfigReloaderImage           string `json:"configReloaderImage"`
					PrometheusConfigReloaderImage string `json:"prometheusConfigReloaderImage"`
					KubeRbacProxyImage            string `json:"kubeRbacProxyImage"`
				} `json:"images"`
			} `json:"prometheusOperator"`
			Prometheus struct {
				Enabled       string `json:"enabled"`
				Image         string `json:"image"`
				CPURequest    int    `json:"cpuRequest"`
				MemoryRequest string `json:"memoryRequest"`
				SvcName       string `json:"svcName"`
				Port          int    `json:"port"`
				NodePort      int    `json:"nodePort"`
				StorageSize   string `json:"storageSize"`
				StorageClass  string `json:"storageClass"`
			} `json:"prometheus"`
			NodeExporter struct {
				Enabled string `json:"enabled"`
				Port    int    `json:"port"`
				Image   string `json:"image"`
			} `json:"nodeExporter"`
			KubeStateMetrics struct {
				Enabled string `json:"enabled"`
				Image   string `json:"image"`
			} `json:"kubeStateMetrics"`
			Grafana struct {
				Enabled  string `json:"enabled"`
				Image    string `json:"image"`
				SvcName  string `json:"svcName"`
				Port     int    `json:"port"`
				NodePort int    `json:"nodePort"`
			} `json:"grafana"`
			DefaultServiceMonitors struct {
				Enabled string `json:"enabled"`
			} `json:"defaultServiceMonitors"`
			SidekiqExporter struct {
				Enabled string `json:"enabled"`
				Image   string `json:"image"`
			} `json:"sidekiqExporter"`
			MinioExporter struct {
				Enabled string `json:"enabled"`
				Image   string `json:"image"`
			} `json:"minioExporter"`
			DcgmExporter struct {
				Enabled string `json:"enabled"`
				Image   string `json:"image"`
				Port    int    `json:"port"`
			} `json:"dcgmExporter"`
			IdleMetricsExporter struct {
				Enabled string `json:"enabled"`
			} `json:"idleMetricsExporter"`
			MetricsServer struct {
				Enabled string `json:"enabled"`
				Image   string `json:"image"`
			} `json:"metricsServer"`
		} `json:"monitoring"`
		Nvidiadp struct {
			Enabled      string `json:"enabled"`
			Image        string `json:"image"`
			NodeSelector struct {
				Enabled string `json:"enabled"`
				Key     string `json:"key"`
				Value   string `json:"value"`
			} `json:"nodeSelector"`
		} `json:"nvidiadp"`
		Logging struct {
			Enabled string `json:"enabled"`
			Es      struct {
				Enabled       string `json:"enabled"`
				Image         string `json:"image"`
				MaxMapImage   string `json:"maxMapImage"`
				Port          string `json:"port"`
				StorageSize   string `json:"storageSize"`
				SvcName       string `json:"svcName"`
				RunAsUser     int    `json:"runAsUser"`
				FsGroup       int    `json:"fsGroup"`
				PatchEsNodes  string `json:"patchEsNodes"`
				NodePort      int    `json:"nodePort"`
				StorageClass  string `json:"storageClass"`
				CPURequest    int    `json:"cpuRequest"`
				MemoryRequest string `json:"memoryRequest"`
				CPULimit      int    `json:"cpuLimit"`
				MemoryLimit   string `json:"memoryLimit"`
				JavaOpts      string `json:"javaOpts"`
			} `json:"es"`
			Elastalert struct {
				Enabled       string `json:"enabled"`
				Image         string `json:"image"`
				Port          string `json:"port"`
				NodePort      int    `json:"nodePort"`
				ContainerPort string `json:"containerPort"`
				StorageSize   string `json:"storageSize"`
				SvcName       string `json:"svcName"`
				StorageClass  string `json:"storageClass"`
				CPURequest    string `json:"cpuRequest"`
				MemoryRequest string `json:"memoryRequest"`
				CPULimit      string `json:"cpuLimit"`
				MemoryLimit   string `json:"memoryLimit"`
				RunAsUser     int    `json:"runAsUser"`
				FsGroup       int    `json:"fsGroup"`
			} `json:"elastalert"`
			Fluentd struct {
				Enabled        string `json:"enabled"`
				Image          string `json:"image"`
				JournalPath    string `json:"journalPath"`
				ContainersPath string `json:"containersPath"`
				Journald       string `json:"journald"`
				CPURequest     string `json:"cpuRequest"`
				MemoryRequest  string `json:"memoryRequest"`
				MemoryLimit    string `json:"memoryLimit"`
			} `json:"fluentd"`
			Kibana struct {
				Enabled       string `json:"enabled"`
				SvcName       string `json:"svcName"`
				Port          int    `json:"port"`
				Image         string `json:"image"`
				NodePort      int    `json:"nodePort"`
				CPURequest    string `json:"cpuRequest"`
				MemoryRequest string `json:"memoryRequest"`
				CPULimit      int    `json:"cpuLimit"`
				MemoryLimit   string `json:"memoryLimit"`
			} `json:"kibana"`
		} `json:"logging"`
		Debug         string `json:"debug"`
		DumpDir       string `json:"dumpDir"`
		DryRun        string `json:"dryRun"`
		ClusterDomain string `json:"clusterDomain"`
		HTTPScheme    string `json:"httpScheme"`
		Otags         string `json:"otags"`
		Tenancy       struct {
			Enabled        string `json:"enabled"`
			DedicatedNodes string `json:"dedicatedNodes"`
			Cnvrg          struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"cnvrg"`
		} `json:"tenancy"`
	} `json:"cnvrgApp"`
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

