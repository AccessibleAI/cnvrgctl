set -e
downloadTools(){
  echo "[$(hostname -f)] downloading tools"

  rkeBinDst=/usr/local/bin/rke
  if [ -f $rkeBinDst ]; then
    echo "[$(hostname -f)] $rkeBinDst present, skipping"
  else
    echo "[$(hostname -f)] downloading rke..."
    curl -Lso $rkeBinDst https://github.com/rancher/rke/releases/download/v1.2.7/rke_linux-amd64
    chmod 0755 $rkeBinDst
  fi

  k9sBinDst=/usr/local/bin/k9s
  if [ -f $k9sBinDst ]; then
    echo "[$(hostname -f)] $k9sBinDst present, skipping"
  else
    echo "[$(hostname -f)] downloading k9s..."
    mkdir -p tmp \
     && cd tmp \
     && curl -Lso k9s.tar.gz https://github.com/derailed/k9s/releases/download/v0.24.7/k9s_Linux_x86_64.tar.gz \
     && tar zxvf k9s.tar.gz \
     && cp ./k9s $k9sBinDst \
     && cd ../ \
     && rm -fr tmp
    chmod 0755 $k9sBinDst
  fi

  kubectlBinDst=/usr/local/bin/kubectl
  if [ -f $kubectlBinDst ]; then
    echo "[$(hostname -f)] $kubectlBinDst present, skipping"
  else
    echo "[$(hostname -f)] downloading kubectl..."
    curl -Lso $kubectlBinDst https://dl.k8s.io/release/v1.20.5/bin/linux/amd64/kubectl
    chmod 0755 $kubectlBinDst
  fi

  helmBinDst=/usr/local/bin/helm
  if [ -f $helmBinDst ]; then
    echo "[$(hostname -f)] $helmBinDst present, skipping"
  else
    echo "[$(hostname -f)] downloading helm"
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    chmod 0755 $helmBinDst
    helm repo add cnvrg https://charts.cnvrg.io
  fi
}

hasSudo() {
    local prompt
    prompt=$(sudo -nv 2>&1)
    if [ $? -eq 0 ]; then
    echo "has_sudo__pass_set"
    elif echo $prompt | grep -q '^sudo:'; then
    echo "has_sudo__needs_pass"
    else
    echo "no_sudo"
    fi
}

patchSshUser(){
  userSudo=$(hasSudo)
  if [ $userSudo == "has_sudo__pass_set" ]; then

    echo "[$(hostname -f)] user has sudo access and password is set, no need to patch"

    cnvrgSudoersGroupExists=$(cat /etc/group | grep cnvrg-sudoers | wc -l)
    if [ $cnvrgSudoersGroupExists -eq 0 ]; then
      sudo groupadd cnvrg-sudoers
    else
      echo "[$(hostname -f)] cnvrg-sudoers group already exists"
    fi

    sudo su root -c 'echo "%cnvrg-sudoers ALL=(ALL:ALL) NOPASSWD:ALL" > /etc/sudoers.d/cnvrg-sudoers'

  elif [ $userSudo == "has_sudo__needs_pass" ]; then

    cnvrgSudoersGroupExists=$(cat /etc/group | grep cnvrg-sudoers | wc -l)
    if [ $cnvrgSudoersGroupExists -eq 0 ]; then
      echo $PASSWD | 2>&1 sudo -S groupadd cnvrg-sudoers
    else
      echo "cnvrg-sudoers group already exists"
    fi
    echo $PASSWD | 2>&1 sudo -S su root -c 'echo "%cnvrg-sudoers ALL=(ALL:ALL) NOPASSWD:ALL" > /etc/sudoers.d/cnvrg-sudoers'
    echo $PASSWD | 2>&1 sudo -S usermod -a -G cnvrg-sudoers {{ .Data.SshUser }}

  else

    >&2 echo "user does not have sudo access, unable proceed with deployment"
    exit 1

  fi
}

createUser(){
  userExists=$(cat /etc/passwd | grep {{ .Data.CnvrgUser }} | wc -l)
  if [ $userExists -eq 0 ]; then
    echo "[$(hostname -f)] creating user cnvrg"
    useradd -m -d /home/{{ .Data.CnvrgUser }} -s /bin/bash -p paMfuNMgwFAX2 --groups sudo {{ .Data.CnvrgUser }}
  else
    echo "[$(hostname -f)] user for cnvrg already exists, skipping user creation"
  fi
}

workdirs(){
  mkdir -p /home/{{ .Data.CnvrgUser }}/.ssh
  mkdir -p /home/{{ .Data.CnvrgUser }}/.kube
  mkdir -p /home/{{ .Data.CnvrgUser }}/rke-cluster
  chown -R {{ .Data.CnvrgUser }}:{{ .Data.CnvrgUser }} /home/{{ .Data.CnvrgUser }}
}

addUserToGroups(){
  usermod -a -G sudo,docker,cnvrg-sudoers {{ .Data.CnvrgUser }}
}

installDocker(){
  2>&1 apt update -y
  2>&1 apt install docker.io -y
  2>&1 systemctl enable docker
}

generateSSHKeys(){
  if [ -f ~/.ssh/id_rsa ]; then
    echo "[$(hostname -f)] ssh keys exists, skipping"
  else
    echo "[$(hostname -f)] generating ssh keys"
    ssh-keygen -b 2048 -t rsa -f ~/.ssh/id_rsa -q -N ""
    cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys
  fi
}

getMainIp(){
  iface=$(cat /proc/net/route | head -n2 | tail -n1 | awk '{print $1}')
  echo $(ip -4 addr show $iface | grep -oP '(?<=inet\s)\d+(\.\d+){3}')
  sleep 1 # to make sure stdout stream reached the client
}

removeRke(){
  userExists=$(cat /etc/passwd | grep {{ .Data.CnvrgUser }} | wc -l)
  if [ $userExists -eq 1 ]; then
    cd /home/{{ .Data.CnvrgUser }}/rke-cluster && rke -d remove --force && rm -fr  ~/.kube/config
  else
    echo "[$(hostname -f)] K8s already removed"
  fi
}

delUser() {
  userExists=$(cat /etc/passwd | grep {{ .Data.CnvrgUser }} | wc -l)
  if [ $userExists -eq 1 ]; then
    2>&1 killall -u {{ .Data.CnvrgUser }}
    2>&1 userdel -fr {{ .Data.CnvrgUser }}
  else
    echo "[$(hostname -f)] user {{ .Data.CnvrgUser }} already removed"
  fi
}


actions="downloadTools|createUser|installDocker|generateSSHKeys|addUserToGroups|patchSshUser|getMainIp|workdirs"
if [ "$#" -ne 1 ]; then
    echo "[$(hostname -f)] missing action parameter, provide one of the following: $actions"
    exit 1
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
"addUserToGroups")
  addUserToGroups
  ;;
"patchSshUser")
  patchSshUser
  ;;
"getMainIp")
  getMainIp
  ;;
"removeRke")
  removeRke
  ;;
"delUser")
  delUser
  ;;
"workdirs")
  workdirs
  ;;
*)
  echo "[$(hostname -f)] ERROR: acceptable values for action: $actions"
  exit 1
  ;;
esac
