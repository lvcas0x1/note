---
date: 2026-09-22
---
##### resize os disk when root is 100% (Azure)

```sh
1. Run df -h to verify your root partition is full (100%)
2. Run lsblk and then lsblk -f to get block device details
3. sudo mount -o size=10M,rw,nodev,nosuid -t tmpfs tmpfs /tmp
4. sudo growpart /dev/DEVICE_ID PARTITION_NUMBER
5. lsblk to verify partition has expanded
6. sudo resize2fs /dev/DEVICE_IDPARTITION_NUMBER
7. Run df -h to verify your resized disk
8. sudo umount /tmp
```

##### remove host key

```
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The fingerprint for the ED25519 key sent by the remote host is
SHA256:lUkXlp8WdJK361NzQ99kLlFep1d0LjfgY5dJ7MmfSno.
Please contact your system administrator.
Add correct host key in /Users/azaria/.ssh/known_hosts to get rid of this message.
Offending ECDSA key in /Users/azaria/.ssh/known_hosts:36
Host key for 3.112.186.239 has changed and you have requested strict checking.
Host key verification failed.
```

```shell
ssh-keygen -R <hostname or IP address>
```
