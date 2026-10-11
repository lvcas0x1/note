---
date: 2026-09-22
---
# Debian 12 with Containerd runtime

### MASTER NODE

##### ensure all the system are up to date

```shell
sudo apt update
sudo apt upgrade -y
```

##### add hostname to hosts file (set all hostname and private to all nodes)

```shell
sudo hostnamectl set-hostname "k8s-master.local"
echo "PRIVATEIP $(hostname)" | sudo tee -a /etc/hosts
```

##### allow the port through the firewall

```shell
# master node
sudo ufw allow 6443/tcp
sudo ufw allow 2379/tcp
sudo ufw allow 2380/tcp
sudo ufw allow 10250/tcp
sudo ufw allow 10251/tcp
sudo ufw allow 10252/tcp
sudo ufw allow 10255/tcp
sudo ufw reload

# for calico network plugin
sudo ufw allow 179/tcp
sudo ufw allow 4789/udp
sudo ufw allow 51820/udp
sudo ufw allow 51821/udp
```

##### change iptables to legacy

```shell
# choose iptables-legacy
sudo update-alternatives --config iptables
```

##### disable swap

```shell
sudo swapoff -a
```

##### option: permanently disable swap in /etc/fstab

```shell
sudo vim /etc/fstab
```

##### enable required kernel modules

```shell
sudo tee /etc/modules-load.d/k8s.conf <<EOF
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter
```

##### modify the below file as showen

```shell
sudo tee /etc/sysctl.d/kubernetes.conf<<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
```

##### reload the settings

```shell
sudo sysctl --system
```

##### ensure the modules are correctly loaded

```shell
lsmod | grep br_netfilter
```

##### reboot os

```shell
sudo reboot
```

##### install requirements for container runtime

```shell
sudo apt install -y curl gnupg2 software-properties-common apt-transport-https ca-certificates
```

##### add docker repo to the system

```shell
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/debian.gpg
sudo add-apt-repository "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
```

##### install containerd

```shell
sudo apt update
sudo apt install -y containerd.io
```

##### configure containerd

```shell
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
```

##### use cgroup driver

```shell
sudo vi /etc/containerd/config.toml

# modify [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options] SystemdCgroup = true
```

```
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
  ...
  SystemdCgroup = true
```

##### restart the service

```shell
sudo systemctl restart containerd
sudo systemctl enable containerd
systemctl status containerd
```

##### install kubeadm bootstrapping tool

```shell
sudo apt install curl apt-transport-https gnupg -y
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmour -o /etc/apt/trusted.gpg.d/cgoogle.gpg
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
```

##### update apt package index

```shell
sudo apt update
```

##### install k8s tools

```shell
sudo apt install kubelet kubeadm kubectl -y
sudo apt-mark hold kubelet kubeadm kubectl
```

##### verify installation

```shell
kubelet --version && kubectl version --client && kubeadm version
```

##### ensure kubelet has been installed

```shell
sudo systemctl enable kubelet
```

##### pull all the required container images on the master node

```shell
sudo kubeadm config images pull
```

##### option: specify the desired runtimes

```shell
# CRI-O
sudo kubeadm config images pull --cri-socket /var/run/crio/crio.sock

# Containerd
sudo kubeadm config images pull --cri-socket /run/containerd/containerd.sock

# Docker
sudo kubeadm config images pull --cri-socket /run/cri-dockerd.sock
```

##### bootstrap master node

```shell
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address PRIVATEIP --apiserver-cert-extra-sans GLOBALIP
```

##### option: specify your -cri-socket for runtimes

| runtime    | socket                          |
| ---------- | ------------------------------- |
| docker     | /run/cri-dockerd.sock           |
| containerd | /run/containerd/containerd.sock |
| cri-o      | /var/run/crio/crio.sock         |

##### ensure bootstrap is successed

```
Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 10.0.1.10:6443 --token da4pbt.q5nja0b883wf6rx3 \
	--discovery-token-ca-cert-hash sha256:e2fdc62f105a53062f90d50cbdf7abd475542447473ed729902d39989d6b2c43
```

##### option: join-command recreate

```shell
kubeadm token create --print-join-command
```

##### To start using your cluster, you need to run the following as a regular user:

```shell
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

##### check the cluster information

```shell
kubectl cluster-info
kubectl get nodes -o wide # not ready
kubectl get pods -A -o wide # ready except coredns
```

##### install and configure the network plugin (flannel)

```shell
wget https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
```

##### modify flannel network configuration if you are using a custom pod cidr and not the default 10.244.0.0/16

```shell
vi kube-flannel.yml
```

```json
net-conf.json: |
  {
    "Network": "192.168.0.0/16",
    "Backend": {
      "Type": "vxlan"
    }
  }
```

##### apply the flannel manifest

```shell
kubectl apply -f kube-flannel.yml
```

##### verify if the node and pods is ready

```shell
kubectl get nodes -o wide
kubectl get pods -A -o wide
```

### WORKER NODE

##### ensure all the system are up to date

```shell
sudo apt update
sudo apt upgrade -y
```

##### add hostname to hosts file (set all hostname and private to all nodes)

```shell
sudo hostnamectl set-hostname "k8s-worker-XX.local"
echo "PRIVATEIP $(hostname)" | sudo tee -a /etc/hosts
```

##### allow the port through the firewall

```shell
# worker node
sudo ufw allow 10250/tcp
sudo ufw allow 30000:32767/tcp
sudo ufw reload

# for calico network plugin
sudo ufw allow 179/tcp
sudo ufw allow 4789/udp
sudo ufw allow 51820/udp
sudo ufw allow 51821/udp
```

##### change iptables to legacy

```shell
# choose iptables-legacy
sudo update-alternatives --config iptables
```

##### disable swap

```shell
sudo swapoff -a
```

##### option: permanently disable swap in /etc/fstab

```shell
sudo vim /etc/fstab
```

##### enable required kernel modules

```shell
sudo tee /etc/modules-load.d/k8s.conf <<EOF
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter
```

##### modify the below file as showen

```shell
sudo tee /etc/sysctl.d/kubernetes.conf<<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
```

##### reload the settings

```shell
sudo sysctl --system
```

##### ensure the modules are correctly loaded

```shell
lsmod | grep br_netfilter
```

##### reboot os

```shell
sudo reboot
```

##### install requirements for container runtime

```shell
sudo apt install -y curl gnupg2 software-properties-common apt-transport-https ca-certificates
```

##### add docker repo to the system

```shell
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/debian.gpg
sudo add-apt-repository "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
```

##### install containerd

```shell
sudo apt update
sudo apt install -y containerd.io
```

##### configure containerd

```shell
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
```

##### use cgroup driver

```shell
sudo vi /etc/containerd/config.toml

# modify [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options] SystemdCgroup = true
```

```
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
  ...
  SystemdCgroup = true
```

##### restart the service

```shell
sudo systemctl restart containerd
sudo systemctl enable containerd
systemctl status containerd
```

##### install kubeadm bootstrapping tool

```shell
sudo apt install curl apt-transport-https gnupg -y
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmour -o /etc/apt/trusted.gpg.d/cgoogle.gpg
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
```

##### update apt package index

```shell
sudo apt update
```

##### install k8s tools

```shell
sudo apt install kubelet kubeadm kubectl -y
sudo apt-mark hold kubelet kubeadm kubectl
```

##### verify installation

```shell
kubelet --version && kubectl version --client && kubeadm version
```

##### ensure kubelet has been installed

```shell
sudo systemctl enable kubelet
```

##### join worker node to master node

```shell
sudo kubeadm join 10.0.1.10:6443 --token da4pbt.q5nja0b883wf6rx3 --discovery-token-ca-cert-hash sha256:e2fdc62f105a53062f90d50cbdf7abd475542447473ed729902d39989d6b2c43

```

##### ensure if worker node is ready (on master node side)

```shell
kubectl get nodes -o wide
```
