---
date: 2026-09-22
---
## Certbot

```shell
sudo certbot certonly --manual --preferred-challenges dns -d "DOMAIN" --key-type rsa
```

## SSH Key

```shell
ssh-keygen -R [HOST_IP_ADDRESS]
```

## 圧縮＆解凍

```shell
# Compress
tar -zcvf xxxx.tar.gz directory
# Uncompress
tar -zxvf xxxx.tar.gz
```

## サービス一覧

```shell
systemctl list-unit-files -t service | grep enabled

systemctl list-unit-files -t service | grep disabled
```

## ファイル転送

```shell
# from local to remote
scp /file/to/send username@remote:/where/to/put

# from remote to local
scp username@remote:/file/to/send /where/to/put

# from remote to remote
scp username@remote_1:/file/to/send username@remote_2:/where/to/put
```
