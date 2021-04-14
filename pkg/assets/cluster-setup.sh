downloadTools(){
  echo "downloading tools"

  rkeBinDst=/usr/local/bin/rke
  if [[ -f $rkeBinDst ]]; then
    echo "$rkeBinDst present, skipping"
  else
    echo "downloading rke..."
    curl -Lso /usr/local/bin/rke https://github.com/rancher/rke/releases/download/v1.2.7/rke_linux-amd64
  fi

  k9sBinDst=/usr/local/bin/k9s
  if [[ -f $k9sBinDst ]]; then
    echo "$k9sBinDst present, skipping"
  else
    echo "downloading k9s..."
    mkdir -p tmp \
     && cd tmp \
     && curl -Lso k9s.tar.gz https://github.com/derailed/k9s/releases/download/v0.24.7/k9s_Linux_x86_64.tar.gz \
     && tar zxvf k9s.tar.gz \
     && cp ./k9s /usr/local/bin/k9s \
     && cd ../ \
     && rm -fr tmp
  fi

  kubectlBinDst=/usr/local/bin/kubectl
  if [[ -f $kubectlBinDst ]]; then
    echo "$kubectlBinDst present, skipping"
  else
    echo "downloading kubectl..."
    curl -Lso /usr/local/bin/kubectl https://dl.k8s.io/release/v1.20.5/bin/linux/amd64/kubectl
  fi

  helmBinDst=/usr/local/bin/helm
  if [[ -f $helmBinDst ]]; then
    echo "$helmBinDst present, skipping"
  else
    echo "downloading helm"
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    helm repo add cnvrg https://charts.cnvrg.io
  fi
}

createUser(){

  userExists=$(cat /etc/passwd | grep {{ .Data.CnvrgUser }} | wc -l)
  if [[ $userExists -eq 0 ]]; then
    echo "creating user for cnvrg"
    useradd -m -d /home/{{ .Data.CnvrgUser }} -s /bin/bash -p paMfuNMgwFAX2 --groups docker,sudo {{ .Data.CnvrgUser }}
  else
    echo "user for cnvrg already exists, skipping user creation"
  fi
}

installDocker(){
  apt update -y
  apt install docker.io=19.03.8-0ubuntu1.20.04.2 -y
}

generateSSHKeys(){
  if [[ -f ~/.ssh/id_rsa ]]; then
    echo "ssh keys exists, skipping"
  else
    echo "generating ssh keys"
    mkdir -p ~/.ssh
    ssh-keygen -b 2048 -t rsa -f .ssh/id_rsa -q -N ""
    cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys
  fi
}

actions="downloadTools|createUser|installDocker|generateSSHKeys"
if [ "$#" -ne 1 ]; then
    echo "missing action parameter, provide one of the following: $actions"
fi

case $1 in
"downloadTools")
  downloadTools
  ;;
"createUser")
  createUser
  ;;
"installDocker")
  installDocker
  ;;
"generateSSHKeys")
  generateSSHKeys
  ;;
*)
  echo "ERROR: acceptable values for action: $actions"
  exit 1
  ;;
esac
