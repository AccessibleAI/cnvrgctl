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
addons_include:
  - ./cnvrg-crds.yaml
  - ./cnvrg-operator.yaml