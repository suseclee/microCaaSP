# microCaaSP
The microCaaSP is small, simple, and light way of kubernetes. The microCaaSP has only one master node, which is the stripped down version of full CaaSP.
The purpose of the microCaaSP is for the people of marketing, sales, and container application developers for CaaSP.
The microCaaSP deployment downloads two files: microCaaSP.xml and microCaaSP.qcow2. microCaaSP.qcow2 from http://10.84.128.39. `http://10.84.128.39` is behind SUSE vpn.
Therefore, in order to use the microCaaSP, you should be connected via SUSE vpn. The microCaaSP.qcow2 you need for deployment is size of 3.9G. Once the image is downloaded, the microCaaSP deployment will continue to create the master node. The image will be saved in cache.

After you deploy the microCaaSP, within a mintute the node and Kubernetes in the node are ready to use by`microCaaSP login`. The microCaaSP includes the clusterctl binary. By `clusterctl init`, users will be also ready to deploy full spec of Kubernetes in any platform that cluster-api supports.
More information about cluster-API is shown https://github.com/kubernetes-sigs/cluster-api. Now you are ready to deploy full stack of Kuberntes any of platforms.  

If you done using the microCaaSP, you can simple destroy by `microCaaSP destroy`


# Prerequisite
### 1. VPN to SUSE
In order for microCaaSP to download image and xml files, you need to run via connection through SUSE VPN

### 2. kvm related tools
These are essential and related packages for kvm virtualization
```
sudo zypper in  qemu-kvm guestfs-tools libvirt libvirt-daemon-qemu virt-manager bridge-utils
sudo systemctl start libvirtd  && sudo systemctl enable libvirtd
sudo usermod -aG libvirt,kvm $USER
```

### 3. Go (for building/installing)
Users need to install GO

# Build/Install microCaaSP
You can download microCaasP binary from the release tap instead of building/installing microCaaSP. 
This insruction shows how to build and install microCaaSP instead of downloading microCaaSP binary.
For building microCaaSP, create a folder under ~/go/source path and clone the repo.
### 1. download source files
```
go get git@github.com:suseclee/microCaaSP.git
cd ~/go/src/github.com/suseclee/microCaaSP
```
##### or
```
mkdir -p ~/go/src/github.com/suseclee
git clone git@github.com:suseclee/microCaaSP.git
cd microCaaSP
```
### 2. create microCaaSP binary
The microCaaSP binary will be available in GO bin path.
```
make
```
##### or
The microCaaSP binary will be avialble in the project root folder.
```
make build
```


# Usage
```
microCaaSP
Usage:
   [command]

Available Commands:
  deploy      deploy microCaaSP cluster
  destroy     destroy microCaaSP cluster
  help        Help about any command
  login       login microCaaSP cluster

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```
