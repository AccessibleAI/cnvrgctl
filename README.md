### cnvrgctl - command line tool for managing cnvrg stack

### Download and install `cnvrgctl`
* mac: [cnvrgctl-darwin-x86_64](https://whitening-pn38xqkin816s3fk.s3-us-west-2.amazonaws.com/cnvrgctl-darwin-x86_64)
  ```shell
  curl -#o /usr/local/bin/cnvrgctl \
    https://cnvrg-public-images.s3-us-west-2.amazonaws.com/cnvrgctl-darwin-x86_64 \
  && chmod +x /usr/local/bin/cnvrgctl \
  && cnvrgctl completion bash > /usr/local/etc/bash_completion.d/cnvrgctl
  ```
* linux: [cnvrgctl-linux-x86_64](https://whitening-pn38xqkin816s3fk.s3-us-west-2.amazonaws.com/cnvrgctl-linux-x86_64)
  ```shell
  curl -#o /usr/local/bin/cnvrgctl \
    https://cnvrg-public-images.s3-us-west-2.amazonaws.com/cnvrgctl-linux-x86_64 \
  && chmod +x /usr/local/bin/cnvrgctl \
  && cnvrgctl completion bash > /etc/bash_completion.d/cnvrgctl
  ```
### Usage 
1. [Deploy all-in-one single node K8s cluster for cnvrg](https://github.com/AccessibleAI/cnvrgctl#deploy-all-in-one-single-node-k8s-cluster-for-cnvrg) 
2. [Import images](https://github.com/AccessibleAI/cnvrgctl#importing-images-for-air-gap-setups)

#### Deploy all-in-one single node K8s cluster for cnvrg 
Prerequisite
1. VM or Bare metal Ubuntu 20.04 server with 32 CPUs, 64GB memory, 500GB storage
2. In case of VM, [bridged network (preferred) or nat network](https://superuser.com/questions/227505/what-is-the-difference-between-nat-bridged-host-only-networking) between VM and the host
3. either root user or regular user with sudo access    
4. SSH access to the server either by ssh key or password

Deploy single node K8s cluster for cnvrg deployment 
```shell
# access the server with ssh password  
cnvrgctl cluster up --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=<SSH-PASS>
# access the server with ssh key  
cnvrgctl cluster up --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=</path/to/private/key>
```

Cleanup 
```shell
# access the server with ssh password  
cnvrgctl cluster destroy --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=<SSH-PASS>
# access the server with ssh key  
cnvrgctl cluster destroy --host=<SERVER-IP> --ssh-user=<SSH-USER> --ssh-pass=</path/to/private/key>
```

#### importing images for air gap setups
There are two ways for importing cnvrg stack images into your internal air gap docker registry
 1. Using `cnvrgctl` 
 2. Using bash scripts generated by `cnvrgctl`

To see all the available options, run `cnvrgctl -h`, for example, for listing all the required images for whitening, run `cnvrgctl images dump --list`

##### Using `cnvrgctl` (*preferred method!*)
   
1. Pull all the images on the linux machine connected to internet
    ```shell
    # the cnvrg pull user/pass should be provided as a part of your license subscription
    cnvrgctl images pull --registry-user=<CNVRG-PULL-USER> --registry-pass=<CNVRG-PULL-KEY>
    ```
2. Save all images to local disk as an images archives  
    ```shell
    cnvrgctl images save
    ```
3. Run whitening process, after that copy all the images archives to internal linux machine, once copied, run
    ```shell
    cnvrgctl images load
    ```
4. Tag images with your internal docker registry and docker repo 
    ```shell
    cnvrgctl images tag \
     --registry=<INTERNAL-PRIVATE-REGISTY> \
     --registry-repo=<INTERNAL-REPO-FOR-CNVRG-IMAGES> 
    ```
5. Push tagged images into your internal docker registry
    ```shell
    cnvrgctl images push \
     --registry=<INTERNAL-PRIVATE-REGISTY> \
     --registry-repo=<INTERNAL-REPO-FOR-CNVRG-IMAGES> \
     --registry-user=<INTERNAL-REGISTY-USER> \
     --registry-pass=<INTERNAL-REGISTY-PASSWORD> 
    ```

#### Using bash scripts generated by `cnvrgctl`
1. Generate and execute bash script for pulling all the images on the linux machine connected to internet
    ```shell
    cnvrgctl images dump --pull
    ```
2. Generate and execute bash script for saving all images to local disk as an images archives
    ```shell
    cnvrgctl images dump --save
    ```
3. Run whitening process, after that, copy all the images archives to internal linux machine, once copied, generate bash script for image load 
    ```shell
    cnvrgctl images dump --load
    ```
4. Generate tag script and run it to load images into yours internal registry 
    ```shell
    cnvrgctl images dump --tag \
     --registry=<INTERNAL-PRIVATE-REGISTY> \
     --registry-repo=<INTERNAL-REPO-FOR-CNVRG-IMAGES> 
    ```
5. Generate push script to push tagged images into your internal docker registry
    ```shell
    cnvrgctl images dump --push \
     --registry=<INTERNAL-PRIVATE-REGISTY> \
     --registry-repo=<INTERNAL-REPO-FOR-CNVRG-IMAGES>
    ```