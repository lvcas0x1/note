---
date: 2026-09-22
---
<!-- panels:start -->
<!-- div:title-panel -->
### docker commands

<!-- div:left-panel -->
##### docker login

<!-- div:right-panel -->
```shell
vi ghcr.txt
cat ~/ghcr.txt | docker login ghcr.io -u <github-user-name> --password-stdin
```

<!-- div:left-panel -->
##### login to docker container

<!-- div:right-panel -->
```shell
docker container exec -it コンテナ名 bash
```

<!-- div:left-panel -->
##### docker run with .env file

<!-- div:right-panel -->
```shell
docker run -p 3000:3000 --env-file .env.local image:tag
```

<!-- div:left-panel -->
##### docker run with volume mount

<!-- div:right-panel -->
```shell
docker run -v /home/directory:/container/directory -p 3000:3000 image:tag
```

<!-- div:left-panel -->
##### docker run with memory limit

<!-- div:right-panel -->
```shell
docker run --memory=2g -p 3000:3000 image:tag
```

<!-- div:left-panel -->
##### docker build

<!-- div:right-panel -->
```shell
docker build -t image:latest .
```

<!-- div:left-panel -->
##### docker tag

<!-- div:right-panel -->
```shell
docker tag imageID REPOSITORY:TAG
```

<!-- div:left-panel -->
##### docker stats

<!-- div:right-panel -->
```shell
docker stats
```

<!-- div:left-panel -->
##### docker cache clear

<!-- div:right-panel -->
```shell
docker system df
docker builder prune
```

<!-- div:left-panel -->
##### docker compose

<!-- div:right-panel -->
```shell
docker compose -p <PROJECT_NAME> -f <DOCKER_COMPOSE_FILE_NAME> up -d
docker compose -p <PROJECT_NAME> -f <DOCKER_COMPOSE_FILE_NAME> down
```

<!-- div:left-panel -->
##### dagger

<!-- div:right-panel -->
```
daggerビルドの時、.envなどの環境変数も一緒にイメージにビルドしてしまう　※要注意
```
<!-- panels:end -->