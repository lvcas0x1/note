---
date: 2026-10-10
---
# Openshift

## Installation

- Download and Install Openshift at https://developers.redhat.com/products/openshift-local

## Create first instance

- crc version

```bash
CRC version: 2.64.0+
OpenShift version: 4.22.14
MicroShift version: 4.22.13
```

- crc setup

```bash
CRC is constantly improving and we would like to know more about usage (more details at https://developers.redhat.com/article/tool-data-collection)
Your preference can be changed manually if desired using 'crc config set consent-telemetry <yes/no>'
Would you like to contribute anonymous usage statistics? [y/N]: N
No worry, you can still enable telemetry manually with the command 'crc config set consent-telemetry yes'.
INFO Using bundle path /Users/lvcas0x1/.crc/cache/crc_vfkit_4.22.14_arm64.crcbundle
INFO Checking if running macOS version >= 13.x
INFO Checking if running as non-root
INFO Checking if crc-admin-helper executable is cached
INFO Checking if running on a supported CPU architecture
INFO Checking if crc executable symlink exists
INFO Creating symlink for crc executable
INFO Checking minimum RAM requirements
INFO Checking if running emulated on Apple silicon
INFO Checking if vfkit is installed
INFO Checking if CRC bundle is extracted in '$HOME/.crc'
INFO Checking if /Users/lvcas0x1/.crc/cache/crc_vfkit_4.22.14_arm64.crcbundle exists
INFO Getting bundle for the CRC executable
INFO Downloading bundle: /Users/lvcas0x1/.crc/cache/crc_vfkit_4.22.14_arm64.crcbundle...
5.65 GiB / 5.65 GiB [----------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 76.75 MiB/s
INFO Uncompressing /Users/lvcas0x1/.crc/cache/crc_vfkit_4.22.14_arm64.crcbundle
crc.img:  31.00 GiB / 31.00 GiB [----------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00%
oc:  96.50 MiB / 96.50 MiB [---------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00%
INFO Checking if old launchd config for tray and/or daemon exists
INFO Checking REST API token
INFO Creating REST API token
INFO Checking if crc daemon plist file is present and loaded
INFO Adding crc daemon plist file and loading it
INFO Checking SSH port availability
Your system is correctly setup for using CRC. Use 'crc start' to start the instance
```

- crc config

```
Modifies crc configuration properties.
Properties:

* bundle                               Bundle path/URI - absolute or local path, http, https or docker URI (string, like 'https://foo.com/crc_vfkit_4.22.14_arm64.crcbundle', 'docker://quay.io/myorg/crc_vfkit_4.22.14_arm64.crcbundle:2.64.0' default '/Users/lvcas0x1/.crc/cache/crc_vfkit_4.22.14_arm64.crcbundle' )
* consent-telemetry                    Consent to collection of anonymous usage data (yes/no)
* cpus                                 Number of CPU cores (must be greater than or equal to '4')
* developer-password                   User defined developer password
* disable-update-check                 Disable update check (true/false, default: false)
* disk-size                            Total size in GiB of the disk (must be greater than or equal to '31')
* enable-bundle-quay-fallback          If bundle download from the default location fails, fallback to quay.io (true/false, default: false)
* enable-cluster-monitoring            Enable cluster monitoring Operator (true/false, default: false)
* enable-emergency-login               Enable emergency login for 'core' user. Password is randomly generated. (true/false, default: false)
* enable-experimental-features         Enable experimental features (true/false, default: false)
* enable-shared-dirs                   Mounts the host's home directory into the CRC VM (true/false, default: true)
* host-network-access                  Allow TCP/IP connections from the CRC VM to services running on the host (true/false, default: false)
* http-proxy                           HTTP proxy URL (string, like 'http://my-proxy.com:8443')
* https-proxy                          HTTPS proxy URL (string, like 'https://my-proxy.com:8443')
* ingress-http-port                    HTTP port to use for OpenShift ingress/routes on the host (1024-65535, default: 80)
* ingress-https-port                   HTTPS port to use for OpenShift ingress/routes on the host (1024-65535, default: 443)
* kubeadmin-password                   User defined kubeadmin password
* memory                               Memory size in MiB (must be greater than or equal to '10752')
* modify-hosts-file                    Allow CRC to modify the system hosts file (true/false, default: true)
* nameserver                           IPv4 address of nameserver (string, like '1.1.1.1 or 8.8.8.8')
* no-proxy                             Hosts, ipv4 addresses or CIDR which do not use a proxy (string, comma-separated list such as '127.0.0.1,192.168.100.1/24')
* persistent-volume-size               Total size in GiB of the persistent volume used by the CSI driver for microshift preset (must be greater than or equal to '15')
* preset                               Virtual machine preset (valid values are: [openshift okd microshift])
* proxy-ca-file                        Path to an HTTPS proxy certificate authority (CA)
* pull-secret-file                     Path of image pull secret (download from https://console.redhat.com/openshift/create/local)
* skip-check-admin-helper-cached       Skip preflight check (true/false, default: false)
* skip-check-bundle-extracted          Skip preflight check (true/false, default: false)
* skip-check-crc-symlink               Skip preflight check (true/false, default: false)
* skip-check-daemon-launchd-plist      Skip preflight check (true/false, default: false)
* skip-check-m1-cpu                    Skip preflight check (true/false, default: false)
* skip-check-mac-version               Skip preflight check (true/false, default: false)
* skip-check-old-autostart             Skip preflight check (true/false, default: false)
* skip-check-ram                       Skip preflight check (true/false, default: false)
* skip-check-resolver-file-permissions Skip preflight check (true/false, default: false)
* skip-check-rest-api-token            Skip preflight check (true/false, default: false)
* skip-check-root-user                 Skip preflight check (true/false, default: false)
* skip-check-ssh-port                  Skip preflight check (true/false, default: false)
* skip-check-supported-cpu-arch        Skip preflight check (true/false, default: false)
* skip-check-vfkit-installed           Skip preflight check (true/false, default: false)

Usage:
  crc config SUBCOMMAND [flags]
  crc config [command]

Available Commands:
  get         Get a crc configuration property
  set         Set a crc configuration property
  unset       Unset a crc configuration property
  view        Display all assigned crc configuration properties

Flags:
  -h, --help   help for config

Global Flags:
      --log-level string   log level (e.g. "debug | info | warn | error") (default "info")

Use "crc config [command] --help" for more information about a command.
```

- crc start -p pull-secret.txt

```
INFO Using bundle path /Users/lvcas0x1/.crc/cache/crc_vfkit_4.22.14_arm64.crcbundle
INFO Checking if running macOS version >= 13.x
INFO Checking if running as non-root
INFO Checking if crc-admin-helper executable is cached
INFO Checking if running on a supported CPU architecture
INFO Checking if crc executable symlink exists
INFO Checking minimum RAM requirements
INFO Checking if running emulated on Apple silicon
INFO Checking if vfkit is installed
INFO Checking if old launchd config for tray and/or daemon exists
INFO Checking REST API token
INFO Checking if crc daemon plist file is present and loaded
INFO Checking SSH port availability
INFO Loading bundle: crc_vfkit_4.22.14_arm64...
INFO Creating CRC VM for OpenShift 4.22.14...
INFO Generating new SSH key pair...
INFO Generating new password for the kubeadmin user
INFO Starting CRC VM for openshift 4.22.14...
INFO CRC instance is running with IP 127.0.0.1
INFO CRC VM is running
INFO Updating authorized keys...
INFO Configuring shared directories
INFO Check internal and public DNS query...
INFO Check DNS query from host...
INFO Verifying validity of the kubelet certificates...
INFO Starting kubelet service
INFO Waiting for kube-apiserver availability... [takes around 2min]
INFO Adding user's pull secret to the cluster...
INFO Updating SSH key to machine config resource...
INFO Overriding password for developer user
INFO Changing the password for the users
INFO Updating cluster ID...
INFO Updating root CA cert to admin-kubeconfig-client-ca configmap...
INFO Starting openshift instance... [waiting for the cluster to stabilize]
INFO 2 operators are progressing: authentication, image-registry
INFO Operator authentication is progressing
INFO All operators are available. Ensuring stability...
INFO Operators are stable (2/3)...
INFO Operators are stable (3/3)...
INFO Waiting until the user's pull secret is written to the instance disk...
INFO Adding crc-admin and crc-developer contexts to kubeconfig...
Started the OpenShift cluster.

The server is accessible via web console at:
  https://console-openshift-console.apps-crc.testing

Log in as administrator:
  Username: kubeadmin
  Password: xHD8L-v4M5G-Mkv59-L7Axy

Log in as user:
  Username: developer
  Password: developer

Use the 'oc' command line interface:
  $ eval $(crc oc-env)
  $ oc login -u developer https://api.crc.testing:6443
```

- crc status

```
CRC VM:          Running
OpenShift:       Running (v4.22.14)
RAM Usage:       6.905GB of 10.93GB
Disk Usage:      25.48GB of 32.68GB (Inside the CRC VM)
Cache Usage:     39.45GB
Cache Directory: /Users/lvcas0x1/.crc/cache
```

- crc config set cpus 8
- crc config set memory 16384
- crc config set disk-size 64
- crc config view

- crc oc-env

```
export PATH="/Users/lvcas0x1/.crc/bin/oc:$PATH"
# Run this command to configure your shell:
# eval $(crc oc-env)
```

- oc version

```
Client Version: 4.22.14
Kustomize Version: v5.7.1
Server Version: 4.22.14
Kubernetes Version: v1.35.6
```

- oc get nodes

```
NAME   STATUS   ROLES                         AGE   VERSION
crc    Ready    control-plane,master,worker   17d   v1.35.6
```

- oc logout

```
Logged "kubeadmin" out on "https://api.crc.testing:6443"

# run oc command after logout openshift cluster
Error from server (Forbidden): nodes is forbidden: User "system:anonymous" cannot list resource "nodes" in API group "" at the cluster scope
```

- oc login

```
oc login --token=sha256~crscFdl7NSeS6-YpucAGtIAzNWMmDL19rqJmsaFGVhk --server=https://api.crc.testing:6443
```

- oc project

```
Using project "default" on server "https://api.crc.testing:6443".
```

- oc get project

```
~/Documents/github/openshift ❯ oc get projects                                                                                               ○ default/api-crc-testing:6443/kubeadmin 08:53:16
NAME                                               DISPLAY NAME   STATUS
default                                                           Active
hostpath-provisioner                                              Active
kube-node-lease                                                   Active
kube-public                                                       Active
kube-system                                                       Active
openshift                                                         Active
openshift-apiserver                                               Active
openshift-apiserver-operator                                      Active
openshift-authentication                                          Active
openshift-authentication-operator                                 Active
openshift-cloud-network-config-controller                         Active
openshift-cloud-platform-infra                                    Active
openshift-cluster-machine-approver                                Active
openshift-cluster-samples-operator                                Active
openshift-cluster-storage-operator                                Active
openshift-cluster-version                                         Active
openshift-config                                                  Active
openshift-config-managed                                          Active
openshift-config-operator                                         Active
openshift-console                                                 Active
openshift-console-operator                                        Active
openshift-console-user-settings                                   Active
openshift-controller-manager                                      Active
openshift-controller-manager-operator                             Active
openshift-dns                                                     Active
openshift-dns-operator                                            Active
openshift-etcd                                                    Active
openshift-etcd-operator                                           Active
openshift-host-network                                            Active
openshift-image-registry                                          Active
openshift-infra                                                   Active
openshift-ingress                                                 Active
openshift-ingress-canary                                          Active
openshift-ingress-operator                                        Active
openshift-kni-infra                                               Active
openshift-kube-apiserver                                          Active
openshift-kube-apiserver-operator                                 Active
openshift-kube-controller-manager                                 Active
openshift-kube-controller-manager-operator                        Active
openshift-kube-scheduler                                          Active
openshift-kube-scheduler-operator                                 Active
openshift-kube-storage-version-migrator                           Active
openshift-kube-storage-version-migrator-operator                  Active
openshift-machine-api                                             Active
openshift-machine-config-operator                                 Active
openshift-marketplace                                             Active
openshift-monitoring                                              Active
openshift-multus                                                  Active
openshift-network-console                                         Active
openshift-network-diagnostics                                     Active
openshift-network-node-identity                                   Active
openshift-network-operator                                        Active
openshift-node                                                    Active
openshift-nutanix-infra                                           Active
openshift-oauth-apiserver                                         Active
openshift-openstack-infra                                         Active
openshift-operator-lifecycle-manager                              Active
openshift-operators                                               Active
openshift-ovirt-infra                                             Active
openshift-ovn-kubernetes                                          Active
openshift-route-controller-manager                                Active
openshift-service-ca                                              Active
openshift-service-ca-operator                                     Active
openshift-user-workload-monitoring                                Active
openshift-vsphere-infra                                           Active
```

- oc project $project_name

```
Now using project "openshift" on server "https://api.crc.testing:6443".
```

- oc get pod

```
NAME                                   READY   STATUS      RESTARTS   AGE
openshift-from-zero-1-build            0/1     Completed   0          3d3h
openshift-from-zero-799d6f96db-nqj7q   1/1     Running     1          3d3h
```

- oc delete pod $pod_name

```
pod "openshift-from-zero-1-build" deleted from rabbit-dev namespace
```
