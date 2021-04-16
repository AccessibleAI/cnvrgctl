cluster_name: cnvrg
ingress:
  provider: none
monitoring:
  provider: none
nodes:
  - address: {{ .Data.Server }}
    user: {{ .Data.User }}
    ssh_key_path: {{ .Data.SshPrivateKey }}
    role:
      - controlplane
      - etcd
      - worker
authentication:
  strategy: x509
  sans:
    - {{ .Data.ExternalIp }}