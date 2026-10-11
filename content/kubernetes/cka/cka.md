---
date: 2026-09-22
---
##### clusterinfo 取得

```shell
kubectl cluster-info

<<output
Kubernetes control plane is running at https://10.0.1.10:6443
CoreDNS is running at https://10.0.1.10:6443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

output
```

##### cluster config 確認

```shell
kubectl config view

<<output
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://10.0.1.10:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: DATA+OMITTED
    client-key-data: DATA+OMITTED

output
```

##### cluster event

```shell
kubectl get events
kubectl get events --watch
kubectl get events -w & # & for keep watch in backround, fg to bring it back, ctrl + c to kill the progress.
```

##### get all api resources and version

```shell
kubectl api-resources

<<output
NAME                              SHORTNAMES   APIVERSION                             NAMESPACED   KIND
bindings                                       v1                                     true         Binding
componentstatuses                 cs           v1                                     false        ComponentStatus
configmaps                        cm           v1                                     true         ConfigMap
endpoints                         ep           v1                                     true         Endpoints
events                            ev           v1                                     true         Event
limitranges                       limits       v1                                     true         LimitRange
namespaces                        ns           v1                                     false        Namespace
nodes                             no           v1                                     false        Node
persistentvolumeclaims            pvc          v1                                     true         PersistentVolumeClaim
persistentvolumes                 pv           v1                                     false        PersistentVolume
pods                              po           v1                                     true         Pod
...

output
```

##### get all resources

```shell
kubectl get all --all-namespaces
kubectl get all -A

<<output
NAMESPACE      NAME                                           READY   STATUS    RESTARTS      AGE
kube-flannel   pod/kube-flannel-ds-47mcp                      1/1     Running   5 (28h ago)   2d22h
kube-flannel   pod/kube-flannel-ds-8mjm9                      1/1     Running   4 (28h ago)   2d22h
kube-flannel   pod/kube-flannel-ds-x5cls                      1/1     Running   5 (28h ago)   2d22h
kube-system    pod/coredns-5dd5756b68-8xrqg                   1/1     Running   4 (28h ago)   2d22h
kube-system    pod/coredns-5dd5756b68-lcwfq                   1/1     Running   4 (28h ago)   2d22h
kube-system    pod/etcd-k8s-master.local                      1/1     Running   4 (28h ago)   2d22h
kube-system    pod/kube-apiserver-k8s-master.local            1/1     Running   4 (28h ago)   2d22h
kube-system    pod/kube-controller-manager-k8s-master.local   1/1     Running   4 (28h ago)   2d22h
kube-system    pod/kube-proxy-4czdv                           1/1     Running   4 (28h ago)   2d22h
kube-system    pod/kube-proxy-cz6rw                           1/1     Running   4 (28h ago)   2d22h
kube-system    pod/kube-proxy-kf5m6                           1/1     Running   4 (28h ago)   2d22h
kube-system    pod/kube-scheduler-k8s-master.local            1/1     Running   4 (28h ago)   2d22h

NAMESPACE     NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
default       service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP                  2d22h
kube-system   service/kube-dns     ClusterIP   10.96.0.10   <none>        53/UDP,53/TCP,9153/TCP   2d22h

NAMESPACE      NAME                             DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
kube-flannel   daemonset.apps/kube-flannel-ds   3         3         3       3            3           <none>                   2d22h
kube-system    daemonset.apps/kube-proxy        3         3         3       3            3           kubernetes.io/os=linux   2d22h

NAMESPACE     NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
kube-system   deployment.apps/coredns   2/2     2            2           2d22h

NAMESPACE     NAME                                 DESIRED   CURRENT   READY   AGE
kube-system   replicaset.apps/coredns-5dd5756b68   2         2         2       2d22h

output
```

##### context

```shell
kubectl config get-contexts
kubectl config use-context {$CONTEXT_NAME}
kubectl config set-context {$CONTEXT_NAME}
kubectl config set-context --current --namespace={$NAMESPACE_NAME}
kubectl config delete-context {$CONTEXT_NAME}

kubectl --context {$CONTEXT_NAME} get pod

<<output
CURRENT   NAME                          CLUSTER      AUTHINFO           NAMESPACE
*         kubernetes-admin@kubernetes   kubernetes   kubernetes-admin

output
```

##### ./kube/config → コピーしてローカルからサーバ操作する

```
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0F...
    server: https://10.0.1.10:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
- context:
    cluster: ""
    user: ""
  name: test-context
current-context: kubernetes-admin@kubernetes
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJ...
    client-key-data: LS0tLS1CRUdJTiBS...
```

##### kubectl explain

1. https://kubernetes.io/ja/docs/home/
2. kubectl explain

```shell
kubectl explain pod
kubectl explain pod.metadata
kubectl explain pod.spec
kubectl explain pod.spec.containers
```

3. kubectl --dry-run=[server | client]

```shell
kubectl run web --image=nginx --dry-run=client -o yaml > ngxin.yaml
```

4. kubectl diff -f nginx.yaml

##### Pod 操作

```shell
kubectl get pod
kubectl get pod {$POD_NAME}
kubectl get pod {$POD_NAME} --show-labels
kubectl get pod {$POD_NAME} -o yaml
kubectl get pod {$POD_NAME} -o jsonpath='{.metadata.name}'
kubectl delete pod {$POD_NAME}
kubectl describe pod {$POD_NAME}
kubectl apply -f pod-manifest.yaml

<<output
Name:             web
Namespace:        default
Priority:         0
Service Account:  default
Node:             k8s-worker-02.local/10.0.1.12
Start Time:       Wed, 18 Oct 2023 05:52:07 +0000
Labels:           run=web
Annotations:      <none>
Status:           Running
IP:               10.244.2.3
IPs:
  IP:  10.244.2.3
Containers:
  web:
    Container ID:   containerd://a5be9044b40442c3ad9b6e3e96ab9546e4bf8129085a4a1ef49ca83e44e2ce0d
    Image:          nginx
    Image ID:       docker.io/library/nginx@sha256:b4af4f8b6470febf45dc10f564551af682a802eda1743055a7dfc8332dffa595
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Wed, 18 Oct 2023 05:52:15 +0000
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-4qtb8 (ro)
Conditions:
  Type              Status
  Initialized       True
  Ready             True
  ContainersReady   True
  PodScheduled      True
Volumes:
  kube-api-access-4qtb8:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  23m   default-scheduler  Successfully assigned default/web to k8s-worker-02.local
  Normal  Pulling    23m   kubelet            Pulling image "nginx"
  Normal  Pulled     23m   kubelet            Successfully pulled image "nginx" in 7.168s (7.168s including waiting)
  Normal  Created    23m   kubelet            Created container web
  Normal  Started    23m   kubelet            Started container web

output
```

##### Pod 内コマンド実行

```shell
kubectl exec {$POD_NAME} -- [COMMAND]
kubectl exec {$POD_NAME} -it -- sh # login to container
kubectl exec {$POD_NAME} -c {$CONTAINER_NAME} -- [COMMAND] # for multi-container pod, specific container name
kubectl exec {$POD_NAME} -c {$CONTAINER_NAME} -- sh
```

##### kubectl < > api verbose log

```shell
kubectl get pods -v 6/7/8/9

<<output
I1018 06:31:48.673747  474232 loader.go:395] Config loaded from file:  /home/admin/.kube/config
I1018 06:31:48.692066  474232 round_trippers.go:553] GET https://10.0.1.10:6443/api/v1/namespaces/default/pods?limit=500 200 OK in 10 milliseconds
NAME     READY   STATUS    RESTARTS   AGE
client   1/1     Running   0          12m
my-pod   2/2     Running   0          10m

output

<<output
I1018 06:32:05.305545  474300 loader.go:395] Config loaded from file:  /home/admin/.kube/config
I1018 06:32:05.311464  474300 round_trippers.go:463] GET https://10.0.1.10:6443/api/v1/namespaces/default/pods?limit=500
I1018 06:32:05.311489  474300 round_trippers.go:469] Request Headers:
I1018 06:32:05.311501  474300 round_trippers.go:473]     Accept: application/json;as=Table;v=v1;g=meta.k8s.io,application/json;as=Table;v=v1beta1;g=meta.k8s.io,application/json
I1018 06:32:05.311520  474300 round_trippers.go:473]     User-Agent: kubectl/v1.28.2 (linux/amd64) kubernetes/89a4ea3
I1018 06:32:05.325688  474300 round_trippers.go:574] Response Status: 200 OK in 14 milliseconds
NAME     READY   STATUS    RESTARTS   AGE
client   1/1     Running   0          12m
my-pod   2/2     Running   0          10m

output

<<output
I1018 06:32:17.823658  474369 loader.go:395] Config loaded from file:  /home/admin/.kube/config
I1018 06:32:17.831210  474369 round_trippers.go:463] GET https://10.0.1.10:6443/api/v1/namespaces/default/pods?limit=500
I1018 06:32:17.831235  474369 round_trippers.go:469] Request Headers:
I1018 06:32:17.831263  474369 round_trippers.go:473]     Accept: application/json;as=Table;v=v1;g=meta.k8s.io,application/json;as=Table;v=v1beta1;g=meta.k8s.io,application/json
I1018 06:32:17.831282  474369 round_trippers.go:473]     User-Agent: kubectl/v1.28.2 (linux/amd64) kubernetes/89a4ea3
I1018 06:32:17.842015  474369 round_trippers.go:574] Response Status: 200 OK in 10 milliseconds
I1018 06:32:17.842104  474369 round_trippers.go:577] Response Headers:
I1018 06:32:17.842123  474369 round_trippers.go:580]     X-Kubernetes-Pf-Flowschema-Uid: 04f2ada1-b960-498e-8efb-75bcbc42e8c4
I1018 06:32:17.842139  474369 round_trippers.go:580]     X-Kubernetes-Pf-Prioritylevel-Uid: 1c9d6433-fde5-40cf-930d-3f206b1df684
I1018 06:32:17.842156  474369 round_trippers.go:580]     Date: Wed, 18 Oct 2023 06:32:17 GMT
I1018 06:32:17.842171  474369 round_trippers.go:580]     Audit-Id: 1482f003-33a3-4bc3-ab0f-4a8db90faf40
I1018 06:32:17.842185  474369 round_trippers.go:580]     Cache-Control: no-cache, private
I1018 06:32:17.842222  474369 round_trippers.go:580]     Content-Type: application/json
I1018 06:32:17.842399  474369 request.go:1212] Response Body: {"kind":"Table","apiVersion":"meta.k8s.io/v1","metadata":{"resourceVersion":"230550"},"columnDefinitions":[{"name":"Name","type":"string","format":"name","description":"Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names","priority":0},{"name":"Ready","type":"string","format":"","description":"The aggregate readiness state of this pod for accepting traffic.","priority":0},{"name":"Status","type":"string","format":"","description":"The aggregate status of the containers in this pod.","priority":0},{"name":"Restarts","type":"string","format":"","description":"The number of times the containers in this pod have been restarted and when the last container in this pod has restarted.","priority":0},{ [truncated 5680 chars]
NAME     READY   STATUS    RESTARTS   AGE
client   1/1     Running   0          12m
my-pod   2/2     Running   0          11m

output
```

##### kubectl proxy 使って認証なしでアクセス

```shell
kubectl proxy &
curl http://127.0.0.1:8001/api/v1/namespaces/default/pods?limit=500
```

##### namespace 操作

```shell
kubectl get namespaces
kubectl create namespace {$NAMESPACE_NAME}
kubectl delete namespace {$NAMESPACE_NAME}
kubectl describe namespaces {$NAMESPACE_NAME}

kubectl config set-context --current --namespace={$NAMESPACE_NAME} # change context default namespace
```

##### static Pods

```shell
/etc/kubernetes/manifests/ 配下にpodのyamlを配置する
# APIサーバ通さずに作成される、削除の時はyamlを削除する

# kubectlのコンフィグファイルからmanifestの場所を特定する（staticPodPath）
sudo cat /var/lib/kubelet/config.yaml

# kubectlのconfigの場所を特定する
ps -ef | grep kubelet
```

##### init containers

```yaml
# run containers in order

apiVersion: v1
kind: Pod
metadata:
  name: pod-with-init-containers
spec:
  initContainers:
    - name: init-service
      image: busybox
      command: ["sh", "-c", "echo waiting for sercice; sleep 4"]
    - name: init-database
      image: busybox
      command: ["sh", "-c", "echo waiting for database; sleep 4"]
  containers:
    - name: app-container
      image: nginx
```

##### Pod lifecycle

```shell
# https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/

# container state
waiting > running > termination

# pod state
pending > runnning > succeeded > failed or unknow

# Pod削除：まずAPI ServerからPodへSIGTERMを送信、30秒（デフォルト）後終了しなかったら、SIGKILLを送る、etcd更新
# Pod削除：GracePeriodSecondsを設定
<<yaml
apiVersion: v1
kind: Pod
metadata:
  name: web
spec:
  terminationGracePeriodSeconds: 10
  containers:
  - image: nginx
    name: web
yaml

# Pod削除：強制削除
kubectl delete pod {$POD_NAME} --grace-period=0 --force
```

##### container restart policy (default Always, OnFailure, Never)

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: web
  name: web
spec:
  containers:
    - image: nginx
      name: web
      resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}
```

##### controller manager

```
1. kube controller
 - replica set
 - deployment
 - daemon set
 - stateful set
 - job
 - cron job

2. cloud controller (loadbalancer)
```

##### labels 操作

```shell
kubectl label pod {pod_name} {key=value}
kubectl label pod {pod_name} {key=value} --overwrite
kubectl label pod {pod_name} {key}-  #remove label

kubectl get pods --show-labels
kubectl get pods --selector {key=value}
kubectl get pods -l 'tier in (PROD, TEST)'  #notin

<<yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod-1
  labels:
    app: v1
    tier: PROD
spec:
  containers:
  - name: nginx
    image: nginx

yaml
```

##### deployment

```shell
# pod controlled by replicaset, replicaset controlled by deployment

kubectl get deployment
kubectl get replicaset

kubectl create deployment {$DEPLOYMENT_NAME} --image=nginx --replicas=3 --dry-run=client -o yaml > deploy.yaml
kubectl delete deployment {$DEPLOYMENT_NAME}

kubectl edit deployment {$DEPLOYMENT_NAME}
kubectl scale deployment {$DEPLOYMENT_NAME}--replicas=3
```

##### deployment update strategy

```shell
# 1. StrategyType: RollingUpdate | RollingUpdateStrategy:  25% max unavailable, 25% max surge
# 2. Recreate | scale down to 0 and recreate new pod

kubectl set image deployment/{$DEPLOYMENT_NAME} nginx=nginx:1.14.2
kubectl rollout history deployment {$DEPLOYMENT_NAME}
kubectl rollout history deployment {$DEPLOYMENT_NAME} --revision 1
kubectl rollout undo deployment {$DEPLOYMENT_NAME} --to-revision 1
kubectl rollout restart deployment {$DEPLOYMENT_NAME}  # create a new replicaset and use same image to create new pod
```

##### daemon set

```shell
kubectl get daemonsets -o wide
kubectl delete daemonset {ds_name}

# why daemon set can run on master node? tolerations!
<<yaml
tolerations:
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: NoSchedule

yaml
```

##### daemon set update strategy

```
1. RollingUpdate (like deployment)
2. OnDelete (remove old pod then update)
```

##### job & cronjob

```shell
# https://kubernetes.io/docs/concepts/workloads/controllers/job/

kubectl get job -o wide
kubectl describe job {$JOB_NAME}
kubectl delete job {$JOB_NAME}

kubectl create job my-job --image=busybox -- sh -c "sleep 50"

<<yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: my-job
spec:
  template:
    spec:
      containers:
      - name: my-job
        image: busybox
        command: ["sh",  "-c", "sleep 50"]
      restartPolicy: Never

yaml
```

```shell
# kubectl get cronjob
# https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/

<<yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox:1.28
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure

yaml
```

##### metric server

```shell
kubectl top node
kubectl top pod

# Install metric server
# https://github.com/kubernetes-sigs/metrics-server
# add - --kubelet-insecure-tls for development environment
<<yaml
spec:
      containers:
      - args:
        - --cert-dir=/tmp
        - --secure-port=4443
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --kubelet-use-node-status-port
        - --metric-resolution=30s
        - --kubelet-insecure-tls

yaml
```

##### etcd database

```shell
# install etcdctl tools https://github.com/etcd-io/etcd/releases/
# https://etcd.io/docs/v3.5/install/
sudo cp etcdctl /usr/local/bin
etcdctl version

kubectl describe pod etcd-k8s-master.local -n kube-system
sudo ls /var/lib/etcd/member/snap

# backup
sudo ETCDCTL_API=3 etcdctl --endpoints="127.0.0.1:2379" --cacert="/etc/kubernetes/pki/etcd/ca.crt" --cert="/etc/kubernetes/pki/etcd/server.crt" --key="/etc/kubernetes/pki/etcd/server.key" snapshot save /var/lib/etcd/member/snap/{$BACKUP_FILE_NAME}.db
sudo ls /var/lib/etcd/member/snap

# check backup status
sudo ETCDCTL_API=3 etcdctl --write-out=table snapshot status /var/lib/etcd/member/snap/{$BACKUP_FILE_NAME}.db

# restore
sudo ETCDCTL_API=3 etcdctl snapshot restore /var/lib/etcd/member/snap/{$BACKUP_FILE_NAME}.db

```

##### cluster upgrade

```shell
# upgrade master node

apt update
apt-cache madison kubeadm

apt-mark unhold kubeadm && apt-get update && apt-get install -y kubeadm='1.2x.x-*' && apt-mark hold kubeadm
kubeadm version

kubeadm upgrade plan
kubectl drain <node-to-drain> --ignore-daemonsets

sudo kubeadm upgrade apply v1.2x.x

apt-mark unhold kubelet kubectl && \
apt-get update && apt-get install -y kubelet='1.2x.x-*' kubectl='1.2x.x-*' && \
apt-mark hold kubelet kubectl

sudo systemctl daemon-reload
sudo systemctl restart kubelet

kubectl uncordon <node-to-uncordon>

# upgrade worker node
# run drain command from master node
kubectl drain <node-to-drain> --ignore-daemonsets

apt update
apt-cache madison kubeadm

apt-mark unhold kubeadm && apt-get update && apt-get install -y kubeadm='1.2x.x-*' && apt-mark hold kubeadm
kubeadm version

kubeadm upgrade plan

sudo kubeadm upgrade node v1.2x.x

apt-mark unhold kubelet && \
apt-get update && apt-get install -y kubelet='1.2x.x-*' && \
apt-mark hold kubelet

sudo systemctl daemon-reload
sudo systemctl restart kubelet

kubectl uncordon <node-to-uncordon>
```

##### node selector

```yaml
# key=value一致すれば対象ノードで作成される
# すでに作成されたPodに影響しない

apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
    - image: nginx
      name: nginx
  nodeSelector:
    node: worker2
  restartPolicy: Always
```

##### affinity and anti-affinity

```shell
# https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector

# node affinity: Podが条件一致するNodeを選択する、Podに定義する
# pod affinity: 他のPodと条件一致する場合展開 or 展開しない
```

##### taints and tolerations

```shell
# Taints Nodeが特定のPodを排除する、Nodeに定義する
# Tolerations がTaintsとマッチする場合、Pod作成可能、Podに定義する

kubectl taint nodes k8s-worker-01.local key1=value1:NoSchedule  # add
kubectl taint nodes k8s-worker-01.local key1=value1:NoSchedule-  # delete

# tolerations
<<yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: web
  name: web
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  strategy: {}
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - image: nginx
        name: nginx
      tolerations:
      - key: "color"
        operator: "Equal"
        value: "red"
        effect: "NoSchedule"

yaml
```

##### cordon, drain, uncordon

```shell
# cordon: Podの作成を禁止できる、Nodeのメンテナンス時に利用. すでに実行中のPodは影響受けない
# drain: gracefullyでPodを停止する
# uncordon: Nodeをscheduleableに戻す

kubectl cordon {$NODE_NAME}
kubectl drain {$NODE_NAME} --ingore-daemonsets
kubectl uncordon {$NODE_NAME}

# manual scheduling
# taints, cordon関係なくPodの配布Nodeを指定できる

<<yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx
  nodeName: kube-01

yaml
```

##### volumes

volumes: pod と同じファイフサイクル、Pod 内のコンテナ間共有
persistent volumes: 独立のストレージ、k8s cluster の外

```yaml
# emptyDir: Pod内のディレクトリでデータを保存する、Pod削除後データも削除される

apiVersion: v1
kind: Pod
metadata:
  name: multicontainer-pod
spec:
  containers:
    - name: producer
      image: busybox
      command: ["sh", "-c", "while true; do echo $(hostname) $(date) >> /var/log/index.html; sleep 10; done"]
      volumeMounts:
        - name: webcontent
          mountPath: /var/log
    - name: consumer
      image: nginx
      ports:
        - containerPort: 80
      volumeMounts:
        - name: webcontent
          mountPath: /usr/share/nginx/html
  volumes:
    - name: webcontent
      emptyDir: {}
```

```yaml
# hostPath: ホストのディレクトリ上データを保存する、Pod削除してもデータはホスト上に保存される（worker）
# PersistentVolumesではない、毎回異なるWorker上で起動される可能性があるため。

apiVersion: v1
kind: Pod
metadata:
  name: multicontainer-pod
spec:
  containers:
    - name: producer
      image: busybox
      command: ["sh", "-c", "while true; do echo $(hostname) $(date) >> /var/log/index.html; sleep 10; done"]
      volumeMounts:
        - name: webcontent
          mountPath: /var/log
    - name: consumer
      image: nginx
      ports:
        - containerPort: 80
      volumeMounts:
        - name: webcontent
          mountPath: /usr/share/nginx/html
  volumes:
    - name: webcontent
      hostPath:
        path: /tmp
        type: Directory
```

```yaml
# nfs server
# https://learn-k8s-from-scratch.readthedocs.io/en/latest/storage/nfs.html
```

##### persistent volumes and persistent volumes claims

```yaml
# 一つのPVは同時に複数のPVCとBoundできない、複数とBoundする場合はDynamic persistent volumes参照
# PersistentVolume
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-nfs
spec:
  capacity:
    storage: 4Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  nfs:
    server: 192.168.56.20
    path: "/export/volumes/pod"

# PersistentVolumeClaim
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-nfs
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi

# use persistentvolume in pod
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      volumes:
      - name: webcontent
        persistentVolumeClaim:
          claimName: pvc-nfs
      containers:
      - image: nginx
        name: nginx
        ports:
        - containerPort: 80
        volumeMounts:
        - name: webcontent
          mountPath: "/usr/share/nginx/html/web-app"
```

##### reclaim policy

```
Retain: PVC削除後データ残る、PVの状態はRelease、他のPVCとBoundできない
Delete: PVC削除後データも削除される、PVの状態はAvailable、他のPVCとBoundできる
Recycle: 非推奨
```

##### environment variables

1. command line arguments
2. environment variables
3. configmaps

environment variables can not be updated once the pod is created

##### secret

1. updatable
2. immutable (immutable: true)

```shell
# base64 encoded
# encryption can be configured
# stored in etcd
# namespaced based
# pod can not start if reference unavailable secrets

kubectl get secrets
kubectl get secrets {$SECRET_NAME} -o yaml
kubectl edit secrets {$SECRET_NAME}

kubectl create secret generic {$SECRET_NAME} --from-literal=key1=supersecret --from-literal=key2=topsecret
```

create secret from yaml

```yaml
# base64 encode secret value
# echo xxx | base64
# echo xxx | base64 --decode

apiVersion: v1
kind: Secret
metadata:
  name: mysql-secret
type: Opaque
data:
  MYSQL_PASSWORD: ZGVtbw==
  MYSQL_ROOT_PASSWORD: cm9vdA==
  MYSQL_USER: ZGVtbw==
```

create secret from config file

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
stringData:
  config.yaml: |
    MYSQL_ROOT_PASSWORD: root
    MYSQL_PASSWORD: demo
    MYSQL_USER: demo
```

use secret in pod

```yaml
# secretKeyRef

apiVersion: v1
kind: Pod
metadata:
  name: mysql
spec:
  containers:
    - name: mysql
      image: mysql:8.0
      env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: MYSQL_ROOT_PASSWORD
        - name: MYSQL_USER
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: MYSQL_USER
        - name: MYSQL_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: MYSQL_PASSWORD
```

```yaml
# volumeMounts

apiVersion: v1
kind: Pod
metadata:
  name: pod-env
spec:
  volumes:
    - name: appconfig
      secret:
        secretName: mysql-secret
  containers:
    - name: busybox
      image: busybox
      command: ["sh", "-c", "while true; do echo $(date) >> /tmp/index.html; sleep 10; done"]
      volumeMounts:
        - name: appconfig
          mountPath: "/etc/appconfig"
```

##### configMap

secret とほぼ同じ、暗号化必要ないデータを保存する
env との違いは後からも変更可能、Pod 作り直す必要ない

from command line

```shell
kubectl create configmap {$CONFIG_NAME} --from-file=key1=/path/to/bar/file1.txt --from-file=key2=/path/to/bar/file2.txt
```

from file

```shell
kubectl create configmap {$CONFIG_NAME} --from-file=key1=/path/to/bar/file1.txt --from-file=key2=/path/to/bar/file2.txt
```

from yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: appconfig
data:
  MYSQL_ROOT_PASSWORD: root
  MYSQL_USER: demo
  MYSQL_PASSWORD: demo
```

use configmap

from env

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mysql
spec:
  containers:
    - name: mysql
      image: mysql:8.0
      env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            configMapKeyRef:
              name: mysql-secret
              key: MYSQL_ROOT_PASSWORD
        - name: MYSQL_USER
          valueFrom:
            configMapKeyRef:
              name: mysql-secret
              key: MYSQL_USER
        - name: MYSQL_PASSWORD
          valueFrom:
            configMapKeyRef:
              name: mysql-secret
              key: MYSQL_PASSWORD
```

from volume mount

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-env
spec:
  volumes:
    - name: appconfig
      configMap:
        name: appconfig
  containers:
    - name: busybox
      image: busybox
      command: ["sh", "-c", "while true; do echo $(date) >> /tmp/index.html; sleep 10; done"]
      volumeMounts:
        - name: appconfig
          mountPath: "/etc/appconfig"
```

##### cluster dns

coredns configmap

```
Name:         coredns
Namespace:    kube-system
Labels:       <none>
Annotations:  <none>

Data
====
Corefile:
----
.:53 {
    errors
    health {
       lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf {
       max_concurrent 1000
    }
    cache 30
    loop
    reload
    loadbalance
}


BinaryData
====

Events:  <none>
```

##### service

1. clusterip

```shell
kubectl expose deployment {$DEPLOYMENT_NAME} --port=80 --target-port=8080  # target port===container port
```

2. nodeport

```shell
kubectl expose deployment {$DEPLOYMENT_NAME} --target-port=8080 --type=NodePort
```

3. loadbalancer

```shell
kubectl expose deployment {$DEPLOYMENT_NAME} --port=80 --target-port=8080 --type=LoadBalancer
```

##### Create new kubeconfig context

1. Create Certificate

```shell
# create a private key
openssl genrsa -out demouser.key 2048

# generate CSR
# CN(common name) is your username, o(Organization) is the Group
openssl req -new -key demouser.key -out demouser.csr -subj "/CN=demouser"

# the CertificateSigning Request needs to be base64 encoded
cat demouser.csr | base64 | tr -d "\n"
```

2. Create CertificateSigning Requst object

```yaml
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: demouser
spec:
  groups:
    - system:authenticated
  request: put the base64 encoded csr here
  signerName: kubernetes.io/kube-apiserver-client
  usages:
    - client auth
```

3. Approve CertificateSigning Requst object

```shell
# approe the CSR
kubectl certificate approve demouser

# retrieve the certificate from the CSR object, and decode it from base64
kubectl get certificatesigningrequests demouser -o jsonpath='{.status.certificate}' | base64 --decode > demouser.crt

# check certificate
openssl x509 -in demouser.crt -text
```

4. Update kubeconfig file

```shell
# add demo user
kubectl config set-credentials demouser --client-key=demouser.key --client-certificate=demouser.crt --embed-certs=true

# check
kubectl config get-users

# create context
kubectl config get-contexts
kubectl config set-context demo --user=demouser --cluster=kubernetes

# check
kubectl config get-contexts

# change context
kubectl config use-context demo
kubectl config get-contexts

kubectl get nodes -v 6

```

##### role based access control

1. create roles (resource level)

```shell
kubectl create role demorole --verb=get,list --resource=pods --namespace ns1
kubectl create role demorole --verb=* --resource=pods --namespace ns1
```

2. clusterrole: node, persistent volume level (not exist in namespace)

```shell
kubectl create clusterrole democlusterrole --verb=get,list --resource=nodes
```

3. binding roles (rolebinding)

```shell
kubectl crate rolebinding demorolebinding --role=demorole --user=demouser --namespace ns1
kubectl get rolebinding
```

4. binding roles (clusterrolebinding)

```shell
kubectl create clusterrolebindding democlusterrolebinding --clusterrole=democlusterrole --user=demouser
kubectl get clusterrolebinding
```

```shell
# test resource
kubectl config use-context kubernetes-admin@kubernetes
kubectl create namespace ns1
kubectl create deployment web1 --namespace=ns1 --image=gcr.io/google-samples/hello-app:1.0 --port=8080 --replicas=2

kubectl auth can-i list pod # yes
kubectl auth can-i list pod --as demouser # no

# create role and role binding
kubectl create role demorole --verb=get,list --resource=pods --namespace ns1
kubectl create rolebinding demorolebinding --role=demorole --user=demouser --namespace ns1

# check
kubectl auth can-i list pod --as demouser # no
kubectl auth can-i list pod --as demouser --namespace ns1 # yes

kubectl get pods --namespace ns1 --as demouser

kubectl auth can-i delete pod --as demouser --namespace ns1 # no

kubectl auth can-i list node --as demouser # no

kubectl auth can-i list deployment --as demouser --namespace ns1 # no

# create clusterrole and cluster role binding
kubectl create clusterrole democlusterrole --verb=list --resource=node
kubectl create clusterrolebinding democlusteerrolebinding --clusterrole=democlusterrole --user=demouser

kubectl auth can-i list node # yes
kubectl auth can-i list node --as demouser # yes
```

##### service account

```shell
kubectl create serviceaccount demosa
kubectl describe serviceaccount demosa
```

```yaml
# use service account in pod

apiVersion: v1
kind: Pod
metadata:
  name: client
spec:
  serviceAccount: demosa
  containers:
    - name: client
      image: xiaopeng163/net-box:latest
      command:
        - sh
        - -c
        - "sleep 1000000"
```

```shell
# ServiceAccount Authentication
kubectl exec -it client -- sh

cd /var/run/secrets/kubernetes.io/serviceaccount
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
CACERT=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt

# without token
curl --cacert $CACERT -X GET https://kubernetes.default.svc.cluster.local/api

# with token
curl --cacert $CACERT --header "Authorization: Bearer $TOKEN" -X GET https://kubernetes.default.svc.cluster.local/api

# service account no api access authorization
curl --cacert $CACERT --header "Authorization: Bearer $TOKEN" -X GET https://kubernetes.default.svc.cluster.local/api/v1/namespaces/default/pods?limit=500

# ServiceAccount Authorization
kubectl auth can-i list pods --as=system:serviceaccount:default:demosa
kubectl get pods -v 6 --as=system:serviceaccount:default:demosa # api error

# rbac
kubectl create role demorole --verb=get,list --resource=pods
kubectl create rolebinding demorolebinding --role=demorole --serviceaccount=default:demosa
kubectl auth can-i list pods --as=system:serviceaccount:default:demosa

curl --cacert $CACERT --header "Authorization: Bearer $TOKEN" -X GET https://kubernetes.default.svc.cluster.local/api/v1/namespaces/default/pods?limit=500
```

##### logging and monitor

1. Logs

```shell
# Pod
kubectl logs -n $NAMESPACE $POD_NAME
kubectl logs -n $NAMESPACE $POD_NAME -f # keep listening
kubectl logs -n $NAMESPACE $POD_NAME --tail $NUM # last $NUM logs

# Container
kubectl logs $POD_NAME -c $CONTAINER_NAME

# Use container runtime list all container on node
sudo crictl --runtime-endpoint unix:///run/containerd/containerd.sock ps
sudo crictl --runtime-endpoint unix:///run/containerd/containerd.sock image
sudo crictl --runtime-endpoint unix:///run/containerd/containerd.sock logs $CONTAINER_ID

# Check log file at /var/log
cd /var/log/containers

# Check systemd service
systemctl status kubelet
sudo journalctl -u kubelet | grep -i ERROR
sudo journalctl -u kubelet --since today --no-pager
```

2. Events

```shell
kubectl get events
kubectl get events --field-selector type=Warning,reason=Failed
kubectl get events --watch & # fg and ctrl+c to break

kubectl describe pod $POD_NAME
```

3. JsonPath

```shell
kubectl get pods -o json
kubectl get pods -o jsonpath='{.items[*].metadata.name}'

kubectl get pods --all-namespaces -o jsonpath='{.items[*].spec.containers[*].image}{"\n"}'
kubectl get pods -A -o jsonpath='{.items[*].metadata.name}{"\n"}' --sort-by=.metadata.name

kubectl get nodes -o jsonpath="{.items[*].status.addresses[?(@.type=='InternalIP')].address}"
```

4. Monitor

```shell
kubectl top pods
kubectl top nodes

# Install metric server
# https://github.com/kubernetes-sigs/metrics-server
# add - --kubelet-insecure-tls for development environment
<<yaml
spec:
      containers:
      - args:
        - --cert-dir=/tmp
        - --secure-port=4443
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --kubelet-use-node-status-port
        - --metric-resolution=30s
        - --kubelet-insecure-tls

yaml
```

##### Troubleshooting

1. Application Failure

- kubectl logs
- kubectl event
- systemctl
- journalctl
- system logs

```shell
kubectl edit pod $POD_NAME
kubectl edit deployment $DEPLOYMENT_NAME
kubectl edit service $SERVICE_NAME

kubectl get pod $POD_NAME -o yaml > pod.yaml
kubectl get deployment $DEPLOYMENT_NAME -o yaml > deployment.yaml
kubectl get service $SERVICE_NAME -o yaml > service.yaml
```

2. Control Plane Failure

- Server online
- Network reachability
- systemd
- container runtime
- kubelet
- static pod manifest

```shell
# static pod config
cat /var/lib/kubelet/config.yaml # staticPodPath
sudo ls -l /etc/kubernetes/manifests
```

3. Worker Node Failure

- Server online
- Network reachability
- systemd
- container runtime
- kubelet
- kube-proxy

```shell
systemctl status kubelet.service
systemctl enable kubelet.service
systemctl start kubelet.service

# journallog
sudo journalctl -u kubelet.service --no-pager

# systemd service unit config
cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

# kubelet config
cat /var/lib/kubelet/config.yaml
```
