---
date: 2026-09-22
---
## Linux 共通手順実施

### システムを最新にアップデートする

```bash
dnf update -y
```

Output

```bash
Upgraded:
  NetworkManager-1:1.40.16-15.0.1.al8.x86_64                                       NetworkManager-initscripts-updown-1:1.40.16-15.0.1.al8.noarch                NetworkManager-libnm-1:1.40.16-15.0.1.al8.x86_64
  NetworkManager-team-1:1.40.16-15.0.1.al8.x86_64                                  NetworkManager-tui-1:1.40.16-15.0.1.al8.x86_64                               alinux-release-3.9.1-1.al8.x86_64
  alinux-repos-3.9.1-1.al8.x86_64                                                  audit-3.0.7-5.0.1.al8.x86_64                                                 audit-libs-3.0.7-5.0.1.al8.x86_64
  authselect-1.2.6-2.al8.x86_64                                                    authselect-compat-1.2.6-2.al8.x86_64                                         authselect-libs-1.2.6-2.al8.x86_64
  automake-1.16.2-8.0.1.al8.noarch                                                 avahi-libs-0.7-21.0.1.al8.1.x86_64                                           bind-export-libs-32:9.11.36-11.0.1.al8.x86_64
  binutils-2.35-12.3.al8.x86_64                                                    bpftool-5.10.134-16.3.al8.x86_64                                             chkconfig-1.19.2-1.al8.x86_64
  cloud-init-23.2.2-1.0.1.al8.noarch                                               container-selinux-2:2.229.0-1.al8.noarch                                     coreutils-8.30-15.0.3.al8.x86_64
  coreutils-common-8.30-15.0.3.al8.x86_64                                          cpp-10.2.1-3.8.al8.x86_64                                                    crypto-policies-20230731-1.git3177e06.al8.noarch
  crypto-policies-scripts-20230731-1.git3177e06.al8.noarch                         cryptsetup-libs-2.3.7-7.0.1.al8.x86_64                                       cups-client-1:2.2.6-54.0.1.al8.x86_64
  cups-libs-1:2.2.6-54.0.1.al8.x86_64                                              curl-7.61.1-34.0.1.al8.x86_64                                                dbus-1:1.12.8-26.0.1.al8.x86_64
  dbus-common-1:1.12.8-26.0.1.al8.noarch                                           dbus-daemon-1:1.12.8-26.0.1.al8.x86_64                                       dbus-libs-1:1.12.8-26.0.1.al8.x86_64
  dbus-tools-1:1.12.8-26.0.1.al8.x86_64                                            device-mapper-8:1.02.181-13.al8.0.1.al8.x86_64                               device-mapper-event-8:1.02.181-13.al8.0.1.al8.x86_64
  device-mapper-event-libs-8:1.02.181-13.al8.0.1.al8.x86_64                        device-mapper-libs-8:1.02.181-13.al8.0.1.al8.x86_64                          dmidecode-1:3.3-5.0.6.al8.x86_64
  dnf-4.7.0-19.0.1.al8.noarch                                                      dnf-data-4.7.0-19.0.1.al8.noarch                                             dnf-plugins-core-4.0.21-23.al8.noarch
  dracut-049-228.git20230802.0.1.al8.x86_64                                        dracut-config-rescue-049-228.git20230802.0.1.al8.x86_64                      dracut-network-049-228.git20230802.0.1.al8.x86_64
  dracut-squash-049-228.git20230802.0.1.al8.x86_64                                 elfutils-debuginfod-client-0.189-3.al8.x86_64                                elfutils-default-yama-scope-0.189-3.al8.noarch
  elfutils-libelf-0.189-3.al8.x86_64                                               elfutils-libs-0.189-3.al8.x86_64                                             emacs-filesystem-1:27.2-9.0.1.al8.noarch
  ethtool-2:6.6-1.al8.x86_64                                                       file-5.33-25.al8.x86_64                                                      file-libs-5.33-25.al8.x86_64
  findutils-1:4.6.0-21.0.1.al8.x86_64                                              fuse-2.9.7-17.al8.x86_64                                                     fuse-common-3.3.0-17.al8.x86_64
  fuse-libs-2.9.7-17.al8.x86_64                                                    fwupd-1.7.8-2.0.1.al8.x86_64                                                 gcc-10.2.1-3.8.al8.x86_64
  glib2-2.68.4-11.al8.x86_64                                                       glibc-2.32-1.16.al8.x86_64                                                   glibc-all-langpacks-2.32-1.16.al8.x86_64
  glibc-common-2.32-1.16.al8.x86_64                                                glibc-devel-2.32-1.16.al8.x86_64                                             glibc-headers-x86-2.32-1.16.al8.noarch
  gmp-1:6.2.0-13.0.1.al8.x86_64                                                    gnutls-3.6.16-8.0.2.al8.x86_64                                               grub2-common-1:2.02-150.0.1.al8.noarch
  grub2-efi-x64-1:2.02-150.0.1.al8.x86_64                                          grub2-pc-1:2.02-150.0.1.al8.x86_64                                           grub2-pc-modules-1:2.02-150.0.1.al8.noarch
  grub2-tools-1:2.02-150.0.1.al8.x86_64                                            grub2-tools-efi-1:2.02-150.0.1.al8.x86_64                                    grub2-tools-extra-1:2.02-150.0.1.al8.x86_64
  grub2-tools-minimal-1:2.02-150.0.1.al8.x86_64                                    grubby-8.40-48.0.1.al8.x86_64                                                gzip-1.9-14.al8.x86_64
  hwdata-0.314-8.19.0.2.al8.noarch                                                 iproute-6.2.0-5.0.1.al8.x86_64                                               iptables-1.8.5-9.0.1.al8.x86_64
  iptables-ebtables-1.8.5-9.0.1.al8.x86_64                                         iptables-libs-1.8.5-9.0.1.al8.x86_64                                         iputils-20180629-11.0.1.al8.x86_64
  irqbalance-2:1.9.2-1.0.1.al8.x86_64                                              kbd-2.0.4-11.0.1.al8.x86_64                                                  kbd-legacy-2.0.4-11.0.1.al8.noarch
  kbd-misc-2.0.4-11.0.1.al8.noarch                                                 keentune-target-2.4.0-1.al8.noarch                                           keentuned-2.4.0-1.al8.x86_64
  kernel-headers-5.10.134-16.3.al8.x86_64                                          kernel-tools-5.10.134-16.3.al8.x86_64                                        kernel-tools-libs-5.10.134-16.3.al8.x86_64
  kexec-tools-2.0.26-8.0.3.1.al8.x86_64                                            kpartx-0.8.4-39.0.2.al8.x86_64                                               krb5-libs-1.18.2-26.0.1.al8.x86_64
  libasan-10.2.1-3.8.al8.x86_64                                                    libatomic-10.2.1-3.8.al8.x86_64                                              libblkid-2.32.1-45.0.1.1.al8.1.x86_64
  libblockdev-2.28-4.al8.x86_64                                                    libblockdev-crypto-2.28-4.al8.x86_64                                         libblockdev-fs-2.28-4.al8.x86_64
  libblockdev-loop-2.28-4.al8.x86_64                                               libblockdev-mdraid-2.28-4.al8.x86_64                                         libblockdev-part-2.28-4.al8.x86_64
  libblockdev-swap-2.28-4.al8.x86_64                                               libblockdev-utils-2.28-4.al8.x86_64                                          libcap-2.48-6.0.1.al8.x86_64
  libcurl-7.61.1-34.0.1.al8.x86_64                                                 libdnf-0.63.0-17.0.1.1.al8.x86_64                                            libfastjson-0.99.9-5.al8.x86_64
  libfdisk-2.32.1-45.0.1.1.al8.1.x86_64                                            libgcc-10.2.1-3.8.al8.x86_64                                                 libgomp-10.2.1-3.8.al8.x86_64
  libibverbs-46.0-1.0.3.al8.1.x86_64                                               libmount-2.32.1-45.0.1.1.al8.1.x86_64                                        libnftnl-1.2.2-3.0.1.al8.x86_64
  libnghttp2-1.33.0-5.al8.x86_64                                                   libsmartcols-2.32.1-45.0.1.1.al8.1.x86_64                                    libsolv-0.7.20-6.al8.x86_64
  libssh-0.9.6-11.al8.x86_64                                                       libssh-config-0.9.6-11.al8.noarch                                            libsss_idmap-2.9.1-4.al8.5.x86_64
  libsss_nss_idmap-2.9.1-4.al8.5.x86_64                                            libstdc++-10.2.1-3.8.al8.x86_64                                              libubsan-10.2.1-3.8.al8.x86_64
  libudisks2-2.9.0-16.0.1.al8.x86_64                                               libuuid-2.32.1-45.0.1.1.al8.1.x86_64                                         libxml2-2.9.7-18.0.2.al8.x86_64
  lvm2-8:2.03.14-13.al8.0.1.al8.x86_64                                             lvm2-libs-8:2.03.14-13.al8.0.1.al8.x86_64                                    mdadm-4.2-8.0.1.al8.x86_64
  memstrack-0.2.5-2.0.1.al8.x86_64                                                 nftables-1:1.0.4-3.al8.x86_64                                                nscd-2.32-1.16.al8.x86_64
  numactl-libs-2.0.16-1.0.1.al8.x86_64                                             openssh-8.0p1-20.0.1.al8.x86_64                                              openssh-clients-8.0p1-20.0.1.al8.x86_64
  openssh-server-8.0p1-20.0.1.al8.x86_64                                           openssl-1:1.1.1k-12.0.1.al8.x86_64                                           openssl-libs-1:1.1.1k-12.0.1.al8.x86_64
  pam-1.3.1-27.al8.x86_64                                                          perf-5.10.134-16.3.al8.x86_64                                                perl-HTTP-Tiny-0.074-2.0.1.al8.1.noarch
  platform-python-3.6.8-56.0.1.2.al8.1.x86_64                                      platform-python-devel-3.6.8-56.0.1.2.al8.1.x86_64                            platform-python-pip-9.0.3-23.0.1.al8.1.noarch
  postfix-2:3.5.8-7.al8.x86_64                                                     procps-ng-3.3.15-14.0.1.al8.x86_64                                           python2-setuptools-39.0.1-12.2.module+al8+9+77049424.noarch
  python2-setuptools-wheel-39.0.1-12.2.module+al8+9+77049424.noarch                python3-audit-3.0.7-5.0.1.al8.x86_64                                         python3-cffi-1.11.5-6.0.1.al8.x86_64
  python3-cryptography-3.2.1-7.al8.x86_64                                          python3-dnf-4.7.0-19.0.1.al8.noarch                                          python3-dnf-plugin-post-transaction-actions-4.0.21-23.al8.noarch
  python3-dnf-plugins-core-4.0.21-23.al8.noarch                                    python3-hawkey-0.63.0-17.0.1.1.al8.x86_64                                    python3-libdnf-0.63.0-17.0.1.1.al8.x86_64
  python3-libs-3.6.8-56.0.1.2.al8.1.x86_64                                         python3-libxml2-2.9.7-18.0.2.al8.x86_64                                      python3-nftables-1:1.0.4-3.al8.x86_64
  python3-perf-5.10.134-16.3.al8.x86_64                                            python3-pip-9.0.3-23.0.1.al8.1.noarch                                        python3-pip-wheel-9.0.3-23.0.1.al8.1.noarch
  python3-pytz-2017.2-11.al8.noarch                                                python3-requests-2.20.0-4.al8.noarch                                         python3-rpm-4.14.3-27.0.5.2.al8.x86_64
  python3-setools-4.3.0-5.al8.x86_64                                               python3-unbound-1.16.2-6.al8.x86_64                                          rng-tools-6.16-1.0.1.al8.x86_64
  rpm-4.14.3-27.0.5.2.al8.x86_64                                                   rpm-build-libs-4.14.3-27.0.5.2.al8.x86_64                                    rpm-libs-4.14.3-27.0.5.2.al8.x86_64
  rpm-plugin-selinux-4.14.3-27.0.5.2.al8.x86_64                                    rpm-plugin-systemd-inhibit-4.14.3-27.0.5.2.al8.x86_64                        rsyslog-8.2102.0-15.1.al8.x86_64
  selinux-policy-3.14.3-128.0.1.al8.1.noarch                                       selinux-policy-targeted-3.14.3-128.0.1.al8.1.noarch                          shadow-utils-2:4.6-19.0.1.al8.x86_64
  sqlite-3.26.0-19.al8.x86_64                                                      sqlite-libs-3.26.0-19.al8.x86_64                                             sssd-client-2.9.1-4.al8.5.x86_64
  sudo-1.9.5p2-1.0.1.al8.x86_64                                                    sysak-2.4.0-1.x86_64                                                         sysstat-11.7.3-11.0.1.al8.x86_64
  systemd-239-78.0.4.1.al8.x86_64                                                  systemd-libs-239-78.0.4.1.al8.x86_64                                         systemd-pam-239-78.0.4.1.al8.x86_64
  systemd-udev-239-78.0.4.1.al8.x86_64                                             tcpdump-14:4.9.3-4.0.1.al8.x86_64                                            tpm2-tss-2.3.2-5.0.2.al8.x86_64
  tuned-2.21.0-1.0.1.2.al8.noarch                                                  tzdata-2024a-1.0.1.1.al8.noarch                                              udisks2-2.9.0-16.0.1.al8.x86_64
  unbound-libs-1.16.2-6.al8.x86_64                                                 util-linux-2.32.1-45.0.1.1.al8.1.x86_64                                      util-linux-user-2.32.1-45.0.1.1.al8.1.x86_64
  virt-what-1.25-4.al8.x86_64                                                      volume_key-libs-0.3.11-6.al8.x86_64                                          which-2.21-20.0.1.al8.x86_64
  xfsprogs-5.0.0-12.0.1.al8.x86_64                                                 yum-4.7.0-19.0.1.al8.noarch                                                  yum-utils-4.0.21-23.al8.noarch
Installed:
  kernel-5.10.134-16.3.al8.x86_64                       kernel-core-5.10.134-16.3.al8.x86_64       kernel-devel-5.10.134-16.3.al8.x86_64      kernel-modules-5.10.134-16.3.al8.x86_64      kernel-modules-extra-5.10.134-16.3.al8.x86_64
  kernel-modules-internal-5.10.134-16.3.al8.x86_64      python3-pyserial-3.1.1-8.1.al8.noarch

Complete!
```

### 日時アップデート設定

```bash
echo -e  "#\!/usr/bin/env bash\nset -e\n/usr/bin/dnf -y -t update" > /etc/cron.daily/update.sh
chmod a+x /etc/cron.daily/update.sh

ls -ltr /etc/cron.daily/update.sh
```

Output

```bash
-rwxr-xr-x 1 root root 54 May 29 19:30 /etc/cron.daily/update.sh
```

```bash
cat /etc/cron.daily/update.sh
```

Output

```bash
#\!/usr/bin/env bash
set -e
/usr/bin/dnf -y -t update
```

### Timezone 変更

```bash
timedatectl set-timezone Asia/Tokyo
timedatectl
```

Output

```bash
Local time: Tue 2024-05-28 10:39:12 JST
Universal time: Tue 2024-05-28 01:39:12 UTC
RTC time: Tue 2024-05-28 01:39:11
Time zone: Asia/Tokyo (JST, +0900)
System clock synchronized: yes
NTP service: active
RTC in local TZ: no
```

### Log rotation

```bash
sed -i "s/rotate 4/rotate 13/" /etc/logrotate.conf
sed -i "s/#compress/compress/" /etc/logrotate.conf

cat /etc/logrotate.conf | grep ^rotate
```

Output

```bash
rotate 13
```

```bash
cat /etc/logrotate.conf | grep compress
```

Output

```bash
# uncomment this if you want your log files compressed
compress
```

### cloud monitor

https://www.alibabacloud.com/help/ja/doc-detail/183482.htm#task-1950491

```bash
ps aux | grep argusagent | grep -v grep
```

Output

```bash
root       46574  0.0  0.1 101596  7028 ?        Ss   14:41   0:00 /usr/local/cloudmonitor/bin/argusagent -d
root       46576  0.2  0.8 1467948 29744 ?       Sl   14:41   0:00 /usr/local/cloudmonitor/bin/argusagent
```

### 一般ユーザ作成

```bash
useradd -m -k /root/ -g wheel ecs-user
echo "ecs-user:`curl -s http://100.100.100.200/latest/meta-data/instance-id`" | chpasswd
sed -i 's/PermitRootLogin yes/PermitRootLogin no/g' /etc/ssh/sshd_config
sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/g' /etc/ssh/sshd_config
cat /dev/null > /root/.ssh/authorized_keys
systemctl restart sshd
systemctl status sshd
```

Output

```bash
● sshd.service - OpenSSH server daemon
   Loaded: loaded (/usr/lib/systemd/system/sshd.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2024-05-29 20:32:07 JST; 6s ago
     Docs: man:sshd(8)
           man:sshd_config(5)
 Main PID: 46625 (sshd)
    Tasks: 1 (limit: 22672)
   Memory: 948.0K
   CGroup: /system.slice/sshd.service
           └─46625 /usr/sbin/sshd -D -oCiphers=aes256-gcm@openssh.com,chacha20-poly1305@openssh.com,aes256-ctr,aes256-cbc,aes128-gcm@openssh.com,aes128-ctr,aes128-cbc -oMACs=hmac-sha2-256-etm@openssh.com,hmac-sha1-etm@openssh.com,umac-128->

May 29 20:32:07 iZ6we3soupdad1mrn2mlxuZ systemd[1]: Starting OpenSSH server daemon...
May 29 20:32:07 iZ6we3soupdad1mrn2mlxuZ sshd[46625]: Server listening on 0.0.0.0 port 22.
May 29 20:32:07 iZ6we3soupdad1mrn2mlxuZ systemd[1]: Started OpenSSH server daemon.
```

### Logout root user and Login with ecs-user

```bash
ssh ecs-user@{PUBLIC_IP_ADDRESS} -i {CERTIFICATION.pem}
```

Output

```bash
Welcome to Alibaba Cloud Elastic Compute Service !

Updates Information Summary: available
    1 Security notice(s)
        1 Important Security notice(s)
Security: kernel-core-5.10.134-16.3.al8.x86_64 is an installed security update
Security: kernel-core-5.10.134-16.1.al8.x86_64 is the currently running version
Run "dnf upgrade-minimal --security" to apply all updates.More details please refer to:
https://help.aliyun.com/document_detail/416274.html
```

### sudo su -

```bash
sudo su - {INSTANCE_ID}
```

### Cloud init script

```bash
vi /root/init.sh

#!/bin/sh

#ssh key copy
curl -s http://100.100.100.200/latest/meta-data/public-keys/0/openssh-key > /home/ecs-user/.ssh/authorized_keys

#change permition
chown ecs-user:wheel /home/ecs-user/.ssh/authorized_keys

#change ecs-user password
echo "ecs-user:`curl -s http://100.100.100.200/latest/meta-data/instance-id`" | chpasswd
```

### Cloud init file

```bash
vi /etc/cloud/cloud.cfg.d/useradd.cfg

#cloud-config
runcmd:
  - 'sh /root/init.sh'
```

## Install Guacamole

### Install epel

```
AlibabaCloudLinux3では実施しない：sudo dnf install -y epel-release
```

### Rpmfusion for FFMPEG

```bash
sudo dnf localinstall --nogpgcheck https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-8.noarch.rpm
```

Output

```bash
Installed:
  rpmfusion-free-release-8-0.1.noarch

Complete!
```

### Dependencies

```bash
sudo dnf install cairo-devel libjpeg-turbo-devel libpng-devel libtool libuuid-devel
```

Output

```bash
Installed:
  brotli-devel-1.0.6-3.1.al8.x86_64                 bzip2-devel-1.0.6-26.2.al8.x86_64            cairo-1.17.4-7.al8.x86_64                       cairo-devel-1.17.4-7.al8.x86_64                 cmake-filesystem-3.26.5-1.0.2.al8.x86_64
  expat-devel-2.2.5-11.al8.x86_64                   fontconfig-2.13.1-4.al8.x86_64               fontconfig-devel-2.13.1-4.al8.x86_64            freetype-devel-2.10.4-9.al8.x86_64              glib2-devel-2.68.4-11.al8.x86_64
  graphite2-devel-1.3.14-9.0.1.al8.x86_64           harfbuzz-devel-2.7.4-8.0.1.al8.x86_64        harfbuzz-icu-2.7.4-8.0.1.al8.x86_64             libX11-1.7.0-8.al8.x86_64                       libX11-common-1.7.0-8.al8.noarch
  libX11-devel-1.7.0-8.al8.x86_64                   libX11-xcb-1.7.0-8.al8.x86_64                libXau-1.0.9-8.al8.x86_64                       libXau-devel-1.0.9-8.al8.x86_64                 libXext-1.3.4-8.al8.x86_64
  libXext-devel-1.3.4-8.al8.x86_64                  libXrender-0.9.10-7.2.al8.x86_64             libXrender-devel-0.9.10-7.2.al8.x86_64          libblkid-devel-2.32.1-45.0.1.1.al8.1.x86_64     libffi-devel-3.1-24.0.1.al8.x86_64
  libicu-devel-60.3-2.1.al8.x86_64                  libjpeg-turbo-2.0.90-6.0.1.al8.x86_64        libjpeg-turbo-devel-2.0.90-6.0.1.al8.x86_64     libmount-devel-2.32.1-45.0.1.1.al8.1.x86_64     libpng-devel-2:1.6.34-5.2.al8.x86_64
  libselinux-devel-2.9-8.2.al8.x86_64               libsepol-devel-2.9-3.0.1.al8.x86_64          libtool-2.4.6-25.3.al8.x86_64                   libuuid-devel-2.32.1-45.0.1.1.al8.1.x86_64      libxcb-1.13.1-1.3.al8.x86_64
  libxcb-devel-1.13.1-1.3.al8.x86_64                pcre-cpp-8.42-6.0.1.al8.x86_64               pcre-devel-8.42-6.0.1.al8.x86_64                pcre-utf16-8.42-6.0.1.al8.x86_64                pcre-utf32-8.42-6.0.1.al8.x86_64
  pcre2-devel-10.32-3.0.1.al8.x86_64                pcre2-utf16-10.32-3.0.1.al8.x86_64           pcre2-utf32-10.32-3.0.1.al8.x86_64              pixman-0.40.0-6.al8.x86_64                      pixman-devel-0.40.0-6.al8.x86_64
  sysprof-capture-devel-3.40.1-3.0.1.al8.x86_64     xorg-x11-proto-devel-2022.2-1.al8.noarch     zlib-devel-1.2.11-20.9.al8.alnx.x86_64

Complete!
```

```bash
sudo dnf module enable libuv
```

Output

```bash
Complete!
```

```bash
sudo dnf install ffmpeg ffmpeg-devel freerdp-devel pango-devel libssh2-devel libtelnet-devel libvncserver-devel libwebsockets-devel pulseaudio-libs-devel openssl-devel libvorbis-devel libwebp-devel libgcrypt-devel uuid-devel
```

Output

```bash
Installed:
  SDL2-2.0.10-2.1.al8.x86_64                       alsa-lib-1.2.9-1.al8.x86_64                  cairo-gobject-1.17.4-7.al8.x86_64               cmake-3.26.5-1.0.2.al8.x86_64            cmake-data-3.26.5-1.0.2.al8.noarch
  cmake-rpm-macros-3.26.5-1.0.2.al8.noarch         ffmpeg-4.2.9-1.el8.x86_64                    ffmpeg-devel-4.2.9-1.el8.x86_64                 ffmpeg-libs-4.2.9-1.el8.x86_64           flac-libs-1.3.3-11.al8.x86_64
  freerdp-devel-2:2.2.0-10.0.1.al8.x86_64          freerdp-libs-2:2.2.0-10.0.1.al8.x86_64       fribidi-1.0.4-9.0.1.al8.x86_64                  fribidi-devel-1.0.4-9.0.1.al8.x86_64     gdk-pixbuf2-2.42.6-3.0.1.al8.x86_64
  glibmm24-2.56.0-2.1.al8.x86_64                   gsm-1.0.19-6.al8.x86_64                      jack-audio-connection-kit-1.9.14-4.el8.x86_64   keyutils-libs-devel-1.5.10-9.al8.x86_64  krb5-devel-1.18.2-26.0.1.al8.x86_64
  lame-libs-3.100-6.2.al8.x86_64                   libXfixes-5.0.3-7.2.al8.x86_64               libXft-2.3.3-8.al8.x86_64                       libXft-devel-2.3.3-8.al8.x86_64          libXxf86vm-1.1.4-9.2.al8.x86_64
  libaom-3.6.1-1.el8.x86_64                        libass-0.17.1-1.el8.x86_64                   libasyncns-0.8-14.2.al8.x86_64                  libavdevice-4.2.9-1.el8.x86_64           libbluray-1.0.2-3.2.al8.x86_64
  libcdio-2.1.0-6.al8.x86_64                       libcdio-paranoia-10.2+2.0.1-6.0.4.al8.x86_64 libcom_err-devel-1.46.0-1.0.1.al8.x86_64        libconfig-1.5-9.2.al8.x86_64             libdatrie-0.2.13-4.0.1.al8.x86_64
  libdatrie-devel-0.2.13-4.0.1.al8.x86_64          libdav1d-0.5.2-2.el8.x86_64                  libdrm-2.4.115-1.al8.x86_64                     libev-devel-4.33-5.0.1.al8.x86_64        libffado-2.4.7-1.el8.x86_64
  libgcrypt-devel-1.8.5-7.0.1.al8.x86_64           libglvnd-1:1.3.4-1.0.1.al8.x86_64            libglvnd-glx-1:1.3.4-1.0.1.al8.x86_64           libgpg-error-devel-1.42-5.0.1.al8.x86_64 libiec61883-1.2.0-18.2.al8.x86_64
  libkadm5-1.18.2-26.0.1.al8.x86_64                libmfx-1.25-4.el8.x86_64                     libmodplug-1:0.8.9.0-9.el8.x86_64               libogg-2:1.3.4-6.al8.x86_64              libogg-devel-2:1.3.4-6.al8.x86_64
  libpciaccess-0.14-1.2.al8.x86_64                 libraw1394-2.1.2-5.2.al8.x86_64              librsvg2-2.50.7-2.0.2.al8.x86_64                libsamplerate-0.1.9-1.2.al8.x86_64       libsigc++20-2.10.0-6.1.al8.x86_64
  libsndfile-1.0.28-13.0.1.al8.x86_64              libssh2-1.10.0-1.el8.x86_64                  libssh2-devel-1.10.0-1.el8.x86_64               libtelnet-0.23-1.el8.x86_64              libtelnet-devel-0.23-1.el8.x86_64
  libthai-0.1.28-8.al8.x86_64                      libthai-devel-0.1.28-8.al8.x86_64            libtheora-1:1.1.1-21.2.al8.x86_64               libuv-1:1.42.0-1.al8.x86_64              libuv-devel-1:1.23.1-2.module_el8+8698+518af944.x86_64
  libv4l-1.14.2-3.2.al8.x86_64                     libva-2.13.0-2.0.1.al8.x86_64                libvdpau-1.4-2.1.al8.x86_64                     libverto-devel-0.3.2-2.al8.x86_64        libvmaf-1.3.15-2.el8.x86_64
  libvncserver-0.9.11-17.1.al8.x86_64              libvncserver-devel-0.9.11-17.1.al8.x86_64    libvorbis-1:1.3.7-5.0.1.al8.x86_64              libvorbis-devel-1:1.3.7-5.0.1.al8.x86_64 libvpx-1.7.0-10.0.1.al8.x86_64
  libwayland-client-1.21.0-1.al8.x86_64            libwayland-cursor-1.21.0-1.al8.x86_64        libwayland-egl-1.21.0-1.al8.x86_64              libwayland-server-1.21.0-1.al8.x86_64    libwebp-1.2.0-7.0.1.al8.x86_64
  libwebp-devel-1.2.0-7.0.1.al8.x86_64             libwebsockets-4.0.1-2.el8.x86_64             libwebsockets-devel-4.0.1-2.el8.x86_64          libwinpr-2:2.2.0-10.0.1.al8.x86_64       libwinpr-devel-2:2.2.0-10.0.1.al8.x86_64
  libxkbcommon-devel-1.0.3-4.al8.x86_64            libxkbfile-1.1.0-1.1.al8.x86_64              libxml++-2.40.1-10.el8.x86_64                   libxml2-devel-2.9.7-18.0.2.al8.x86_64    libxshmfence-1.3-2.2.al8.x86_64
  lzo-minilzo-2.08-14.2.al8.x86_64                 mesa-filesystem-23.1.4-1.al8.x86_64          mesa-libGL-23.1.4-1.al8.x86_64                  mesa-libglapi-23.1.4-1.al8.x86_64        ocl-icd-2.2.13-4.0.1.al8.x86_64
  openal-soft-1.18.2-7.2.al8.x86_64                opencore-amr-0.1.5-7.el8.x86_64              openjpeg2-2.4.0-5.al8.x86_64                    openssl-devel-1:1.1.1k-12.0.1.al8.x86_64 opus-1.3-0.4.beta.2.al8.x86_64
  pango-1.48.7-3.0.1.al8.x86_64                    pango-devel-1.48.7-3.0.1.al8.x86_64          pugixml-1.13-1.el8.x86_64                       pulseaudio-libs-15.0-2.0.1.al8.x86_64    pulseaudio-libs-devel-15.0-2.0.1.al8.x86_64
  pulseaudio-libs-glib2-15.0-2.0.1.al8.x86_64      soxr-0.1.3-4.el8.x86_64                      speex-1.2.0-1.2.al8.x86_64                      srt-libs-1.4.1-3.el8.x86_64              uuid-1.6.2-43.1.al8.x86_64
  uuid-devel-1.6.2-43.1.al8.x86_64                 vapoursynth-libs-51-1.el8.x86_64             vid.stab-1.1.0-12.20190213gitaeabc8d.el8.x86_64 vo-amrwbenc-0.1.3-8.el8.x86_64           wayland-devel-1.21.0-1.al8.x86_64
  x264-libs-0.157-12.20190717git34c06d1.el8.x86_64 x265-libs-3.1.2-1.el8.x86_64                 xvidcore-1.3.7-9.el8.x86_64                     xz-devel-5.4.4-2.al8.x86_64              zimg-3.0.5-1.el8.x86_64
  zvbi-0.2.35-9.el8.x86_64

Complete!
```

```bash
sudo dnf install rpmdevtools doxygen systemd-devel gnutls-devel
```

Output

```bash
Installed:
  annobin-11.13-2.0.1.1.al8.x86_64       boost-regex-1.66.0-13.0.1.al8.x86_64   ctags-5.8-23.0.1.al8.x86_64               doxygen-1:1.8.14-12.1.al8.x86_64       dwz-0.14-3.al8.x86_64                  efi-srpm-macros-3-3.1.al8.noarch
  elfutils-0.189-3.al8.x86_64            gc-8.0.4-7.0.1.al8.x86_64              gdb-headless-9.2-7.1.0.4.al8.x86_64       ghc-srpm-macros-1.4.2-7.1.al8.noarch   gmp-c++-1:6.2.0-13.0.1.al8.x86_64      gmp-devel-1:6.2.0-13.0.1.al8.x86_64
  gnutls-c++-3.6.16-8.0.2.al8.x86_64     gnutls-dane-3.6.16-8.0.2.al8.x86_64    gnutls-devel-3.6.16-8.0.2.al8.x86_64      go-srpm-macros-2-17.1.al8.noarch       guile-5:2.0.14-7.2.al8.x86_64          libbabeltrace-1.5.8-10.0.1.al8.x86_64
  libidn2-devel-2.2.0-1.2.al8.x86_64     libipt-1.6.1-8.2.al8.x86_64            libtasn1-devel-4.13-4.0.1.al8.x86_64      nettle-devel-3.4.1-7.0.1.al8.x86_64    ocaml-srpm-macros-5-4.1.al8.noarch     openblas-srpm-macros-2-2.1.al8.noarch
  p11-kit-devel-0.23.22-1.1.al8.x86_64   perl-srpm-macros-1-25.1.al8.noarch     qt5-srpm-macros-5.15.3-1.0.1.al8.noarch   redhat-rpm-config-125-1.4.al8.noarch   rpm-build-4.14.3-27.0.5.2.al8.x86_64   rpmdevtools-8.10-8.1.al8.noarch
  rust-srpm-macros-5-2.1.al8.noarch      source-highlight-3.1.9-11.al8.x86_64   systemd-devel-239-78.0.4.1.al8.x86_64     unzip-6.0-46.0.1.al8.x86_64            zip-3.0-23.2.al8.x86_64                zstd-1.5.1-2.0.2.al8.x86_64

Complete!
```

```bash
sudo dnf install ghostscript ghostscript-devel
```

Output

```bash
Installed:
  adobe-mappings-cmap-20171205-12.al8.noarch           adobe-mappings-cmap-deprecated-20171205-12.al8.noarch             adobe-mappings-pdf-20180407-10.al8.noarch                  ghostscript-9.54.0-14.al8.x86_64
  ghostscript-tools-fonts-9.54.0-14.al8.x86_64         ghostscript-tools-printing-9.54.0-14.al8.x86_64                   google-droid-sans-fonts-20120715-13.1.al8.noarch           jbig2dec-libs-0.19-7.0.1.al8.x86_64
  jbigkit-libs-2.1-14.2.al8.x86_64                     lcms2-2.9-2.2.al8.x86_64                                          libgs-9.54.0-14.al8.x86_64                                 libgs-devel-9.54.0-14.al8.x86_64
  libijs-0.35-5.2.al8.x86_64                           libpaper-1.1.28-4.al8.x86_64                                      libtiff-4.4.0-10.0.1.al8.x86_64                            urw-base35-bookman-fonts-20200910-6.al8.noarch
  urw-base35-c059-fonts-20200910-6.al8.noarch          urw-base35-d050000l-fonts-20200910-6.al8.noarch                   urw-base35-fonts-20200910-6.al8.noarch                     urw-base35-fonts-common-20200910-6.al8.noarch
  urw-base35-gothic-fonts-20200910-6.al8.noarch        urw-base35-nimbus-mono-ps-fonts-20200910-6.al8.noarch             urw-base35-nimbus-roman-fonts-20200910-6.al8.noarch        urw-base35-nimbus-sans-fonts-20200910-6.al8.noarch
  urw-base35-p052-fonts-20200910-6.al8.noarch          urw-base35-standard-symbols-ps-fonts-20200910-6.al8.noarch        urw-base35-z003-fonts-20200910-6.al8.noarch

Complete!
```

### Install Guacamole-Server

```bash
wget https://downloads.apache.org/guacamole/1.5.4/source/guacamole-server-1.5.4.tar.gz
tar -zxvf guacamole-server-1.5.4.tar.gz
cd guacamole-server-1.5.4/
./configure --with-init-dir=/etc/init.d
```

Output

```bash
------------------------------------------------
guacamole-server version 1.5.4
------------------------------------------------

   Library status:

     freerdp2 ............ yes
     pango ............... yes
     libavcodec .......... yes
     libavformat.......... yes
     libavutil ........... yes
     libssh2 ............. yes
     libssl .............. yes
     libswscale .......... yes
     libtelnet ........... yes
     libVNCServer ........ yes
     libvorbis ........... yes
     libpulse ............ yes
     libwebsockets ....... yes
     libwebp ............. yes
     wsock32 ............. no

   Protocol support:

      Kubernetes .... yes
      RDP ........... yes
      SSH ........... yes
      Telnet ........ yes
      VNC ........... yes

   Services / tools:

      guacd ...... yes
      guacenc .... yes
      guaclog .... yes

   FreeRDP plugins: /usr/lib64/freerdp2
   Init scripts: /etc/init.d
   Systemd units: no

Type "make" to compile guacamole-server.
```

```bash
make
make install
ldconfig
```

### Start guacd

```bash
sudo systemctl daemon-reload
sudo systemctl start guacd
sudo systemctl enable guacd
sudo systemctl status guacd
```

Output

```bash
● guacd.service - LSB: Guacamole proxy daemon
   Loaded: loaded (/etc/rc.d/init.d/guacd; generated)
   Active: active (running) since Wed 2024-05-29 14:50:21 JST; 8s ago
     Docs: man:systemd-sysv-generator(8)
    Tasks: 1 (limit: 22672)
   Memory: 9.9M
   CGroup: /system.slice/guacd.service
           └─89688 /usr/local/sbin/guacd -p /var/run/guacd.pid

May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ systemd[1]: Starting LSB: Guacamole proxy daemon...
May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ guacd[89685]: Guacamole proxy daemon (guacd) version 1.5.4 started
May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ guacd[89683]: Starting guacd:
May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ guacd[89685]: guacd[89685]: INFO:        Guacamole proxy daemon (guacd) version 1.5.4 started
May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ guacd[89683]: SUCCESS
May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ systemd[1]: Started LSB: Guacamole proxy daemon.
May 29 14:50:21 iZ6weh01ci1sdsul7688tvZ guacd[89688]: Listening on host ::1, port 4822
```

```bash
ss -tap |grep 4822
```

Output

```bash
LISTEN    0      5              [::1]:4822               [::]:*        users:(("guacd",pid=89688,fd=4))
```

### Install guacamole-client

```bash
sudo dnf install java-latest-openjdk-devel.x86_64
```

Output

```bash
Installed:
  ModemManager-glib-1.20.2-1.al8.x86_64            abattis-cantarell-fonts-0.301-4.al8.noarch                adobe-source-code-pro-fonts-2.030.1.050-12.al8.1.noarch         adwaita-cursor-theme-40.1.1-3.al8.noarch
  adwaita-icon-theme-40.1.1-3.al8.noarch           at-spi2-atk-2.38.0-4.al8.x86_64                           at-spi2-core-2.40.3-1.al8.x86_64                                atk-2.36.0-5.0.1.al8.x86_64
  avahi-glib-0.7-21.0.1.al8.1.x86_64               bluez-libs-5.63-1.0.1.al8.x86_64                          colord-libs-1.4.5-4.0.1.al8.x86_64                              copy-jdk-configs-4.0-3.al8.noarch
  dconf-0.28.0-4.1.al8.x86_64                      exempi-2.6.0-0.2.20211007gite23c213.al8.x86_64            exiv2-0.27.5-2.al8.x86_64                                       exiv2-libs-0.27.5-2.al8.x86_64
  fdk-aac-free-2.0.0-8.al8.x86_64                  flatpak-1.10.8-1.0.1.al8.x86_64                           flatpak-selinux-1.10.8-1.0.1.al8.noarch                         flatpak-session-helper-1.10.8-1.0.1.al8.x86_64
  gdk-pixbuf2-modules-2.42.6-3.0.1.al8.x86_64      geoclue2-2.6.0-7.al8.x86_64                               geoclue2-libs-2.6.0-7.al8.x86_64                                giflib-5.2.1-9.al8.x86_64
  glib-networking-2.56.1-1.1.2.al8.x86_64          graphene-1.10.6-4.al8.x86_64                              gsettings-desktop-schemas-40.0-6.0.1.al8.x86_64                 gstreamer1-1.22.1-2.al8.x86_64
  gstreamer1-plugins-base-1.22.1-1.al8.x86_64      gtk-update-icon-cache-3.24.31-2.0.2.1.al8.x86_64          gtk3-3.24.31-2.0.2.1.al8.x86_64                                 hicolor-icon-theme-0.17-2.1.al8.noarch
  iso-codes-4.6.0-3.0.1.al8.noarch                 java-latest-openjdk-1:22.0.1.0.8-1.rolling.el8.x86_64     java-latest-openjdk-devel-1:22.0.1.0.8-1.rolling.el8.x86_64     java-latest-openjdk-headless-1:22.0.1.0.8-1.rolling.el8.x86_64
  javapackages-filesystem-5.3.1-7.3.al8.noarch     libXcomposite-0.4.5-7.al8.x86_64                          libXcursor-1.2.0-7.al8.x86_64                                   libXdamage-1.1.5-7.al8.x86_64
  libXi-1.7.10-1.1.al8.x86_64                      libXinerama-1.1.4-1.2.al8.x86_64                          libXrandr-1.5.2-8.al8.x86_64                                    libXtst-1.2.3-7.2.al8.x86_64
  libXv-1.0.11-7.2.al8.x86_64                      libappstream-glib-0.7.18-4.0.1.al8.x86_64                 libcanberra-0.30-18.1.al8.x86_64                                libcanberra-gtk3-0.30-18.1.al8.x86_64
  libepoxy-1.5.8-1.0.2.al8.x86_64                  libexif-0.6.22-6.al8.x86_64                               libfontenc-1.1.3-8.2.al8.x86_64                                 libgexiv2-0.12.3-1.0.1.al8.x86_64
  libglvnd-egl-1:1.3.4-1.0.1.al8.x86_64            libgsf-1.14.47-5.0.1.al8.x86_64                           libgxps-0.3.2-3.al8.x86_64                                      libimobiledevice-1.3.0-5.0.2.al8.x86_64
  libiptcdata-1.0.5-9.al8.x86_64                   libldac-2.0.2.3-10.al8.x86_64                             libmodman-2.0.1-17.2.al8.x86_64                                 libnotify-0.7.9-8.al8.x86_64
  libosinfo-1.9.0-3.al8.x86_64                     libplist-2.2.0-5.0.1.al8.x86_64                           libproxy-0.4.15-5.2.2.al8.x86_64                                libsbc-1.4-9.al8.x86_64
  libsoup-2.62.3-4.0.1.al8.x86_64                  libstemmer-0-10.585svn.2.al8.x86_64                       libtdb-1.4.8-3.0.1.al8.x86_64                                   libtracker-sparql-3.1.2-3.0.1.al8.x86_64
  libusbmuxd-2.0.2-5.0.1.al8.x86_64                libvisual-1:0.4.0-25.1.al8.x86_64                         libxslt-1.1.32-6.1.al8.x86_64                                   lksctp-tools-1.0.18-3.2.al8.x86_64
  llvm-libs-15.0.7-1.0.4.al8.x86_64                lua-5.3.4-12.0.2.al8.x86_64                               lua-posix-35.0-8.al8.x86_64                                     mesa-libEGL-23.1.4-1.al8.x86_64
  mesa-libgbm-23.1.4-1.al8.x86_64                  mesa-vulkan-drivers-23.1.4-1.al8.x86_64                   orc-0.4.28-3.1.al8.x86_64                                       osinfo-db-20230518-1.al8.noarch
  osinfo-db-tools-1.9.0-1.1.al8.x86_64             ostree-libs-2022.2-8.al8.x86_64                           p11-kit-server-0.23.22-1.1.al8.x86_64                           pipewire-0.3.67-2.0.1.al8.x86_64
  pipewire-libs-0.3.67-2.0.1.al8.x86_64            poppler-20.11.0-10.0.1.al8.x86_64                         poppler-data-0.4.9-1.1.al8.noarch                               poppler-glib-20.11.0-10.0.1.al8.x86_64
  rtkit-0.11-19.2.al8.x86_64                       sound-theme-freedesktop-0.8-9.1.al8.noarch                totem-pl-parser-3.26.6-2.0.1.al8.x86_64                         tracker-3.1.2-3.0.1.al8.x86_64
  tracker-miners-3.1.2-4.0.1.al8.x86_64            ttmkfdir-3.0.9-54.2.al8.x86_64                            tzdata-java-2024a-1.0.1.1.al8.noarch                            upower-0.99.14-1.0.2.al8.x86_64
  vulkan-loader-1.3.250.1-1.al8.x86_64             webrtc-audio-processing-0.3.1-8.al8.x86_64                xdg-desktop-portal-1.8.1-1.al8.x86_64                           xdg-desktop-portal-gtk-1.12.0-3.al8.x86_64
  xml-common-0.6.3-50.1.al8.noarch                 xorg-x11-font-utils-1:7.5-41.0.2.al8.x86_64               xorg-x11-fonts-Type1-7.5-19.1.al8.noarch

Complete!
```

```bash
java -version
```

Output

```bash
openjdk version "22.0.1" 2024-04-16
OpenJDK Runtime Environment (Red_Hat-22.0.1.0.8-1) (build 22.0.1+8)
OpenJDK 64-Bit Server VM (Red_Hat-22.0.1.0.8-1) (build 22.0.1+8, mixed mode, sharing)
```

### Install tomcat

```bash
wget https://downloads.apache.org/tomcat/tomcat-9/v9.0.89/bin/apache-tomcat-9.0.89.tar.gz
tar -zxvf apache-tomcat-9.0.89.tar.gz -C /opt
sudo ln -s apache-tomcat-9.0.89 /opt/tomcat
ls -la /opt/
```

Output

```bash
total 12
drwxr-xr-x.  3 root root 4096 May 29 18:28 .
dr-xr-xr-x. 18 root root 4096 May 29 17:54 ..
drwxr-xr-x   9 root root 4096 May 29 18:27 apache-tomcat-9.0.89
lrwxrwxrwx   1 root root   20 May 29 18:28 tomcat -> apache-tomcat-9.0.89
```

### Create tomcat user

```bash
sudo groupadd -g 53 tomcat
sudo useradd -M -d /opt/tomcat -s /sbin/nologin -c "Apache Tomcat" -g 53 -G 53 -u 53 tomcat
sudo id tomcat
sudo cat /etc/passwd | grep tomcat
```

Output

```bash
tomcat:x:53:53:Apache Tomcat:/opt/tomcat:/sbin/nologin
```

### Change tomcat directory owner

```bash
sudo chown -R tomcat:tomcat /opt/apache-tomcat-9.0.89
sudo chown -R tomcat:tomcat /opt/tomcat
sudo ls -la /opt
```

Output

```bash
total 12
drwxr-xr-x.  3 root   root   4096 May 29 14:58 .
dr-xr-xr-x. 18 root   root   4096 May 29 14:32 ..
drwxr-xr-x   9 tomcat tomcat 4096 May 29 14:58 apache-tomcat-9.0.89
lrwxrwxrwx   1 tomcat tomcat   20 May 29 14:58 tomcat -> apache-tomcat-9.0.89
```

### Set remote ip value

```bash
vi /opt/tomcat/conf/server.xml

      <Host name="localhost"  appBase="webapps"
            unpackWARs="true" autoDeploy="false">

        <Valve className="org.apache.catalina.valves.AccessLogValve" directory="logs"
               prefix="localhost_access_log" suffix=".txt"
               maxDays="30" requestAttributesEnabled="true" pattern="combined" />

        <Valve className="org.apache.catalina.valves.RemoteIpValve"
               internalProxies="10\.\d+\.\d+\.\d+|192\.168\.\d+\.\d+|169\.254\.\d+\.\d+|127\.\d+\.\d+\.\d+|172\.(1[6-9]|2[0-9]|3[0-1])\.\d+\.\d+|100\.\d+\.\d+\.\d+|0:0:0:0:0:0:0:1"
               remoteIpHeader="x-forwarded-for"
               proxiesHeader="x-forwarded-by"
               protocolHeader="x-forwarded-proto" />

        <Valve className="org.apache.catalina.valves.ErrorReportValve" showReport="false" showServerInfo="false"/>
```

### Create systemd unit file

```bash
sudo vi /etc/systemd/system/tomcat.service

[Unit]
Description=Apache Tomcat 9
After=syslog.target network.target

[Service]
User=tomcat
Group=tomcat
Type=forking
RemainAfterExit=yes
SuccessExitStatus=143

EnvironmentFile=/etc/tomcat/tomcat.conf
Environment=CATALINA_PID=/opt/tomcat/temp/tomcat.pid
Environment=JAVA_HOME=/usr/lib/jvm/jre
Environment=CATALINA_HOME=/opt/tomcat
Environment=CATALINA_BASE=/opt/tomcat

ExecStart=/opt/tomcat/bin/startup.sh
ExecStop=/opt/tomcat/bin/shutdown.sh

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl status tomcat -l
```

Output

```bash
● tomcat.service - Apache Tomcat 9
   Loaded: loaded (/etc/systemd/system/tomcat.service; disabled; vendor preset: enabled)
   Active: inactive (dead)
```

```bash
sudo mkdir /etc/tomcat
sudo chown tomcat:tomcat /etc/tomcat
ls -la /etc/tomcat
```

Output

```bash
total 16
drwxr-xr-x    2 tomcat tomcat  4096 May 29 15:01 .
drwxr-xr-x. 118 root   root   12288 May 29 15:01 ..
```

```bash
sudo vi /etc/tomcat/tomcat.conf

CATALINA_OPTS="-Xms2048M -Xmx2048M -XX:MaxMetaspaceSize=256M -XX:NewRatio=3 -server -XX:+UseConcMarkSweepGC -XX:+UseParNewGC -Djava.awt.headless=true"
```

```bash
cd /opt/tomcat/bin
sudo touch setenv.sh
sudo chmod +x setenv.sh
sudo vi setenv.sh

export CATALINA_OPTS="-XX:+UseG1GC"
```

```bash
sudo systemctl start tomcat
sudo systemctl enable tomcat
sudo systemctl status tomcat -l
```

Output

```bash
● tomcat.service - Apache Tomcat 9
   Loaded: loaded (/etc/systemd/system/tomcat.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2024-05-29 15:03:42 JST; 6s ago
 Main PID: 91420 (java)
    Tasks: 35 (limit: 22672)
   Memory: 109.1M
   CGroup: /system.slice/tomcat.service
           └─91420 /usr/lib/jvm/jre/bin/java -Djava.util.logging.config.file=/opt/tomcat/conf/logging.properties -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager -Djdk.tls.ephemeralDHKeySize=2048 -Djava.protocol.handler.pk>

May 29 15:03:42 iZ6weh01ci1sdsul7688tvZ systemd[1]: Starting Apache Tomcat 9...
May 29 15:03:42 iZ6weh01ci1sdsul7688tvZ startup.sh[91413]: Tomcat started.
May 29 15:03:42 iZ6weh01ci1sdsul7688tvZ systemd[1]: Started Apache Tomcat 9.
```

### Check heap size

```bash
jstat -gccapacity "MainPID"
# ENGCMX と OGCMX の合計が2048MBであること。（数値はKB表示）
```

Output

```bash
   NGCMN        NGCMX         NGC          S0C     S1C              EC         OGCMN        OGCMX         OGC           OC         MCMN       MCMX        MC       CCSMN     CCSMX     CCSC     YGC    FGC   CGC
         0.0     917504.0      33792.0         0.0      4096.0      29696.0          0.0     917504.0      25600.0      25600.0        0.0  1114112.0    13248.0       0.0 1048576.0    1536.0      4     0     0
```

### Download Guacamole-Client&Plugins

```bash
## Guacamole-Client
wget https://downloads.apache.org/guacamole/1.5.4/binary/guacamole-1.5.4.war

### Guacamole DB Plugin
wget https://downloads.apache.org/guacamole/1.5.4/binary/guacamole-auth-jdbc-1.5.4.tar.gz

### TOTP Plugin
wget https://downloads.apache.org/guacamole/1.5.4/binary/guacamole-auth-totp-1.5.4.tar.gz

### Logo Plugin
wget https://bucket.oss/remote-login-logo-1.3.0.jar
```

### 不要なモジュールの削除

```bash
ls -ltr /opt/tomcat/webapps/
```

Output

```bash
total 20
drwxr-x---  3 tomcat tomcat 4096 May 29 14:58 ROOT
drwxr-x--- 16 tomcat tomcat 4096 May 29 14:58 docs
drwxr-x---  7 tomcat tomcat 4096 May 29 14:58 examples
drwxr-x---  6 tomcat tomcat 4096 May 29 14:58 host-manager
drwxr-x---  6 tomcat tomcat 4096 May 29 14:58 manager
```

```bash
sudo rm -rf /opt/tomcat/webapps/*
ls -ltr /opt/tomcat/webapps/
```

Output

```bash
total 0
```

### war ファイルのオーナー変更

```bash
sudo chown tomcat:tomcat guacamole-1.5.4.war
```

### war ファイルのデプロイ

```bash
sudo cp -rp guacamole-1.5.4.war /opt/tomcat/webapps/guacamole.war
ls -ltr /opt/tomcat/webapps/
```

Output

```bash
total 16996
-rw-r--r-- 1 tomcat tomcat 17401039 Mar 30 13:40 guacamole.war
```

### Access Check

```
http://{PUBLIC_IP_ADDRESS}:8080/guacamole
```

### 拡張機能、設定ファイルの配置ディレクトリ作成

```bash
sudo mkdir -pv /etc/guacamole/{extensions,lib}
```

Output

```bash
mkdir: created directory '/etc/guacamole'
mkdir: created directory '/etc/guacamole/extensions'
mkdir: created directory '/etc/guacamole/lib'
```

```bash
ls -ltr /etc/guacamole/
```

Output

```bash
total 8
drwxr-xr-x 2 root root 4096 May 29 20:46 lib
drwxr-xr-x 2 root root 4096 May 29 20:46 extensions
```

### プラグインのデプロイ

```bash
# JDBCプラグイン
tar -zxvf guacamole-auth-jdbc-1.5.4.tar.gz
sudo cp guacamole-auth-jdbc-1.5.4/mysql/guacamole-auth-jdbc-mysql-1.5.4.jar /etc/guacamole/extensions/
# TOTPプラグイン
tar -zxvf guacamole-auth-totp-1.5.4.tar.gz
sudo cp guacamole-auth-totp-1.5.4/guacamole-auth-totp-1.5.4.jar /etc/guacamole/extensions/
# Logoプラグイン
sudo cp remote-login-logo-1.3.0.jar /etc/guacamole/extensions/

ls -ltr /etc/guacamole/extensions/
```

Output

```bash
total 15656
-rw-r--r-- 1 root root 11036568 May 29 15:24 guacamole-auth-jdbc-mysql-1.5.4.jar
-rw-r--r-- 1 root root  4918007 May 29 15:25 guacamole-auth-totp-1.5.4.jar
-rw-r--r-- 1 root root    72417 May 29 15:26 remote-login-logo-1.3.0.jar
```

## Install MySQL

```bash
sudo dnf localinstall https://dev.mysql.com/get/mysql84-community-release-el8-1.noarch.rpm
```

Output

```bash
Installed:
  mysql84-community-release-el8-1.noarch

Complete!
```

```bash
sudo dnf repolist enabled | grep "mysql.*-community.*"
```

Output

```bash
mysql-8.4-lts-community       MySQL 8.4 LTS Community Server
mysql-connectors-community    MySQL Connectors Community
mysql-tools-8.4-lts-community MySQL Tools 8.4 LTS Community
```

```bash
sudo dnf repolist allupdate rpmfusion-free-updatesinfo mysql-connector-javaupdate
rpm -ivh https://dev.mysql.com/get/mysql84-community-release-el8-1.noarch.rpm
```

Output

```bash
Retrieving https://dev.mysql.com/get/mysql84-community-release-el8-1.noarch.rpm
warning: /var/tmp/rpm-tmp.sya9iU: Header V4 RSA/SHA256 Signature, key ID a8d3785c: NOKEY
Verifying...                          ################################# [100%]
Preparing...                          ################################# [100%]
        package mysql84-community-release-el8-1.noarch is already installed
```

```bash
sudo dnf info mysql-community-server
```

Output

```bash
MySQL 8.4 LTS Community Server                                                                                                                                                                                   12 MB/s | 215 kB     00:00
MySQL Connectors Community                                                                                                                                                                                       14 MB/s | 128 kB     00:00
MySQL Tools 8.4 LTS Community                                                                                                                                                                                   194 kB/s |  97 kB     00:00
Available Packages
Name         : mysql-community-server
Version      : 8.4.0
Release      : 1.el8
Architecture : x86_64
Size         : 61 M
Source       : mysql-community-8.4.0-1.el8.src.rpm
Repository   : mysql-8.4-lts-community
Summary      : A very fast and reliable SQL database server
URL          : http://www.mysql.com/
License      : Copyright (c) 2000, 2024, Oracle and/or its affiliates. Under GPLv2 license as shown in the Description field.
Description  : The MySQL(TM) software delivers a very fast, multi-threaded, multi-user,
             : and robust SQL (Structured Query Language) database server. MySQL Server
             : is intended for mission-critical, heavy-load production systems as well
             : as for embedding into mass-deployed software. MySQL is a trademark of
             : Oracle and/or its affiliates
             :
             : The MySQL software has Dual Licensing, which means you can use the MySQL
             : software free of charge under the GNU General Public License
             : (http://www.gnu.org/licenses/). You can also purchase commercial MySQL
             : licenses from Oracle and/or its affiliates if you do not wish to be bound by the terms of
             : the GPL. See the chapter "Licensing and Support" in the manual for
             : further info.
             :
             : The MySQL web site (http://www.mysql.com/) provides the latest news and
             : information about the MySQL software.  Also please see the documentation
             : and the manual for more information.
             :
             : This package includes the MySQL server binary as well as related utilities
             : to run and administer a MySQL server.
```

```bash
sudo dnf install mysql-community-server
```

Output

```bash
Installed:
  mysql-community-client-8.4.0-1.el8.x86_64    mysql-community-client-plugins-8.4.0-1.el8.x86_64    mysql-community-common-8.4.0-1.el8.x86_64    mysql-community-icu-data-files-8.4.0-1.el8.x86_64    mysql-community-libs-8.4.0-1.el8.x86_64
  mysql-community-server-8.4.0-1.el8.x86_64

Complete!
```

```bash
mysqld --version
```

Output

```bash
/usr/sbin/mysqld  Ver 8.4.0 for Linux on x86_64 (MySQL Community Server - GPL)
```

```bash
sudo dnf info mysql-connector-java
```

Output

```bash
Last metadata expiration check: 0:01:43 ago on Wed 29 May 2024 03:29:58 PM JST.
Available Packages
Name         : mysql-connector-java
Epoch        : 1
Version      : 8.0.30
Release      : 1.el8
Architecture : noarch
Size         : 2.4 M
Source       : mysql-connector-java-8.0.30-1.el8.src.rpm
Repository   : mysql-connectors-community
Summary      : Standardized MySQL database driver for Java
URL          : http://dev.mysql.com/downloads/connector/j/
License      : GPLv2
Description  : MySQL provides connectivity for client applications developed in the
             : Java programming language with MySQL Connector/J, a driver that
             : implements the [Java Database Connectivity (JDBC) API]
             : (http://www.oracle.com/technetwork/java/javase/jdbc/).
             :
             : MySQL Connector/J 8.0 is a JDBC Type 4 driver that is compatible with
             : the [JDBC 4.2](http://docs.oracle.com/javase/8/docs/technotes/guides/jdbc/)
             : specification. The Type 4 designation means that the driver is a pure
             : Java implementation of the MySQL protocol and does not rely on the
             : MySQL client libraries.
             :
             : For detailed information please visit the official
             : [MySQL Connector/J documentation]
             : (http://dev.mysql.com/doc/connector-j/en/).
```

```bash
sudo dnf install mysql-connector-java
```

Output

```bash
Installed:
  java-1.8.0-alibaba-dragonwell-headless-3:8.15.16.372-2.al8.x86_64                                                                    mysql-connector-j-1:8.4.0-1.el8.noarch

Complete!
```

```bash
sudo ln -s /usr/share/java/mysql-connector-java.jar /etc/guacamole/lib/mysql-connector-java.jar
ls -la /etc/guacamole/lib/mysql-connector-java.jar
```

Output

```bash
lrwxrwxrwx 1 root root 40 May 29 15:32 /etc/guacamole/lib/mysql-connector-java.jar -> /usr/share/java/mysql-connector-java.jar
```

```bash
sudo systemctl start mysqld
sudo systemctl enable mysqld
sudo systemctl status mysqld -l
```

Output

```bash
● mysqld.service - MySQL Server
   Loaded: loaded (/usr/lib/systemd/system/mysqld.service; enabled; vendor preset: enabled)
   Active: active (running) since Wed 2024-05-29 20:50:24 JST; 14s ago
     Docs: man:mysqld(8)
           http://dev.mysql.com/doc/refman/en/using-systemd.html
 Main PID: 92301 (mysqld)
   Status: "Server is operational"
    Tasks: 36 (limit: 22672)
   Memory: 450.8M
   CGroup: /system.slice/mysqld.service
           └─92301 /usr/sbin/mysqld

May 29 20:50:16 iZ6we3soupdad1mrn2mlxuZ systemd[1]: Starting MySQL Server...
May 29 20:50:24 iZ6we3soupdad1mrn2mlxuZ systemd[1]: Started MySQL Server.
```

```bash
grep 'temporary password' /var/log/mysqld.log
```

Output

```bash
2024-05-29T06:33:00.514333Z 6 [Note] [MY-010454] [Server] A temporary password is generated for root@localhost: 7__9odRIOdNS
```

```bash
mysql_secure_installation
```

Output

```bash
The existing password for the user account root has expired. Please set a new password.

New password: ${PASSWORD}

Re-enter new password:
The 'validate_password' component is installed on the server.
The subsequent steps will run with the existing configuration
of the component.
Using existing password for root.

Estimated strength of the password: 100
Change the password for root ? ((Press y|Y for Yes, any other key for No) : y

New password:

Re-enter new password:

Estimated strength of the password: 100
Do you wish to continue with the password provided?(Press y|Y for Yes, any other key for No) : y
By default, a MySQL installation has an anonymous user,
allowing anyone to log into MySQL without having to have
a user account created for them. This is intended only for
testing, and to make the installation go a bit smoother.
You should remove them before moving into a production
environment.

Remove anonymous users? (Press y|Y for Yes, any other key for No) : y
Success.


Normally, root should only be allowed to connect from
'localhost'. This ensures that someone cannot guess at
the root password from the network.

Disallow root login remotely? (Press y|Y for Yes, any other key for No) : y
Success.

By default, MySQL comes with a database named 'test' that
anyone can access. This is also intended only for testing,
and should be removed before moving into a production
environment.


Remove test database and access to it? (Press y|Y for Yes, any other key for No) : y
 - Dropping test database...
Success.

 - Removing privileges on test database...
Success.

Reloading the privilege tables will ensure that all changes
made so far will take effect immediately.

Reload privilege tables now? (Press y|Y for Yes, any other key for No) : y
Success.

All done!
```

### MySQL ログイン(Remote Gateway DB/USER 作成)

```bash
mysql -u root -p
CREATE DATABASE guacamole_rg_db;
SHOW DATABASES;
```

Output

```bash
+--------------------+
| Database           |
+--------------------+
| guacamole_rg_db    |
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
5 rows in set (0.00 sec)
```

#### DB ユーザー作成

```bash
CREATE USER 'guacamole_rg_user'@'localhost' IDENTIFIED BY '${PASSWORD}';
SELECT user, host FROM mysql.user;
```

Output

```bash
+-------------------+-----------+
| user              | host      |
+-------------------+-----------+
| guacamole_rg_user | localhost |
| mysql.infoschema  | localhost |
| mysql.session     | localhost |
| mysql.sys         | localhost |
| root              | localhost |
+-------------------+-----------+
5 rows in set (0.00 sec)
```

#### 権限付与

```bash
GRANT SELECT,INSERT,UPDATE,DELETE ON guacamole_rg_db.* TO 'guacamole_rg_user'@' localhost';
FLUSH PRIVILEGES;
SHOW GRANTS for guacamole_rg_user@localhost;
```

Output

```bash
+------------------------------------------------------------------------------------------------+
| Grants for guacamole_rg_user@localhost                                                         |
+------------------------------------------------------------------------------------------------+
| GRANT USAGE ON *.* TO `guacamole_rg_user`@`localhost`                                          |
| GRANT SELECT, INSERT, UPDATE, DELETE ON `guacamole_rg_db`.* TO `guacamole_rg_user`@`localhost` |
+------------------------------------------------------------------------------------------------+
2 rows in set (0.00 sec)
```

```bash
quit
```

### 初期データ作成

```bash
ls -lt guacamole-auth-jdbc-1.5.4/mysql/schema
```

Output

```bash
total 28
-rw-r--r-- 1 ecs-user 1000 20174 Jul 21  2021 001-create-schema.sql
-rw-r--r-- 1 ecs-user 1000  2876 Jul 21  2021 002-create-admin-user.sql
drwxr-xr-x 2 ecs-user 1000  4096 Jul 21  2021 upgrade
```

```bash
cat guacamole-auth-jdbc-1.5.4/mysql/schema/001-create-schema.sql | mysql -u root -p guacamole_rg_db
cat guacamole-auth-jdbc-1.5.4/mysql/schema/002-create-admin-user.sql | mysql -u root -p guacamole_rg_db
```

### インポートデータの確認

```bash
mysql -u guacamole_rg_user -p
SHOW DATABASES;
USE guacamole_rg_db;
SELECT * FROM guacamole_entity;
```

Output

```bash
+-----------+-----------+------+
| entity_id | name      | type |
+-----------+-----------+------+
|         1 | guacadmin | USER |
+-----------+-----------+------+
1 row in set (0.00 sec)
```

```bash
quit
```

### Guacamole 設定ファイルの作成

```bash
sudo vi /etc/guacamole/guacamole.properties

# Guacamole Server (guacd)の指定
guacd-hostname: localhost
guacd-port: 4822

# MySQL properties
mysql-hostname: localhost
mysql-port: 3306
mysql-database: guacamole_rg_db
mysql-username: guacamole_rg_user
mysql-password: ${PASSWORD}

# Password complexity
mysql-user-password-min-length： 8
mysql-user-password-require-multiple-case：true
mysql-user-password-require-symbol：true
mysql-user-password-require-digit：true
mysql-user-password-prohibit-username：true
mysql-user-password-history-size: 6

# Restricting authentication to database users only
mysql-user-required: true

# Auto-creating database users
mysql-auto-create-accounts: true

# TOTP
totp-issuer: SB Cloud RG
```

```bash
ls -ltr /etc/guacamole/
```

Output

```bash
total 12
drwxr-xr-x 2 root root 4096 May 29 15:26 extensions
drwxr-xr-x 2 root root 4096 May 29 15:32 lib
-rw-r--r-- 1 root root  686 May 29 15:44 guacamole.properties
```

### 各種ディレクトリ作成

#### guac_record ディレクトリ作成

```bash
sudo mkdir -pv /opt/guac_record/{screenrecording,typescript}
```

Output

```bash
mkdir: ディレクトリ '/opt/guac_record' を作成しました
mkdir: ディレクトリ '/opt/guac_record/screenrecording' を作成しました
mkdir: ディレクトリ '/opt/guac_record/typescript' を作成しました
```

```bash
ls -ltr /opt/guac_record/
```

Output

```bash
total 8
drwxr-xr-x 2 root root 4096 May 29 20:54 typescript
drwxr-xr-x 2 root root 4096 May 29 20:54 screenrecording
```

#### guac_record ディレクトリオーナー変更

```bash
# sudo chown -R guacd:guacd /opt/guac_record/
# ls -ltr /opt/guac_record/
合計 8
drwxr-xr-x 2 guacd guacd 4096  7月  8 20:38 typescript
drwxr-xr-x 2 guacd guacd 4096  7月  8 20:38 screenrecording
```

#### guac_drive ディレクトリ作成

```bash
sudo mkdir -pv /opt/guac_drive/
```

Output

```bash
mkdir: ディレクトリ '/opt/guac_drive/' を作成しました
```

```bash
ls -ltr /opt/
```

Output

```bash
total 12
drwxr-xr-x 9 tomcat tomcat 4096 May 29 20:41 apache-tomcat-9.0.89
lrwxrwxrwx 1 tomcat tomcat   20 May 29 20:41 tomcat -> apache-tomcat-9.0.89
drwxr-xr-x 4 root   root   4096 May 29 20:54 guac_record
drwxr-xr-x 2 root   root   4096 May 29 20:54 guac_drive
```

#### guac_drive ディレクトリオーナー変更

```bash
# sudo chown -R guacd:guacd /opt/guac_drive/
# ls -la /opt/guac_drive
合計 8
drwxr-xr-x  2 guacd guacd 4096  7月  8 20:50 .
drwxr-xr-x. 5 root  root  4096  7月  8 20:50 ..
```

#### 自動削除設定(7 日で消す)

```bash
sudo vi /etc/tmpfiles.d/guac_record.conf

d /opt/guac_record/screenrecording 1755 root root 7d
d /opt/guac_record/typescript 1755 root root 7d
```

### Bug-Fix_1 Modify Hosts

```bash
sudo vi /etc/hosts

#::1    localhost       localhost.localdomain   localhost6      localhost6.localdomain6
127.0.0.1       localhost       localhost.localdomain   localhost4      localhost4.localdomain4
```

### Restart VM

### Guacamole ログイン

```
http://PUBLIC_IP_ADDREESS:8080/guacamoleへアクセスし、画面が表示されること。

以下でログインできること
ID:guacadmin
Password:guacadmin
```

## Image

### runcmd log 削除

```bash
ls -ltr /var/lib/cloud/instance/scripts/runcmd
ls -ltr /var/lib/cloud/instance/user-data.txt

# ファイルが存在しているかつ、ファイルサイズが0以上の場合以下を実施
# cloud-initを利用したimage作成の場合実施
cat /dev/null > /var/lib/cloud/instance/scripts/runcmd
cat /dev/null > /var/lib/cloud/instance/user-data.txt
```

### log の強制 rotate と不要 log 削除

```bash
/usr/sbin/logrotate -f /etc/logrotate.conf
ls -ltr /var/log/*$(date '+%Y%m%d').gz

rm -rf /var/log/*$(date '+%Y%m%d').gz
```

### ssh 不要情報削除

```bash
shred -u /etc/ssh/*_key /etc/ssh/*_key.pub

# 該当ファイルがなくなっていること
ls -ltr /etc/ssh/*_key
ls -ltr /etc/ssh/*_key.pub
```

### cron ジョブ確認

```bash
ls -ltaR /var/spool/cron/
ls -ltaR /etc/cron*
```

### root コマンドログ削除

```bash
cat /dev/null > /root/.bash_history
```

### root ターミナルプロセス kill

```bash
kill -9 $$
```

### ecs-user コマンドログ削除

```bash
cat /dev/null > ~/.bash_history
```

### ecs-user ターミナルプロセス kill

```bash
kill -9 $$
```
