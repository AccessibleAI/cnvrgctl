cluster_name: cnvrg
nodes:
  - address: {{ .Data.Server }}
    user: {{ .Data.User }}
    ssh_key_path: {{ .Data.SshPrivateKey }}
    role:
      - controlplane
      - etcd
      - worker