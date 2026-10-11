---
date: 2026-09-22
---
##### ECS(Fargate)のコンテナに Exec でログインする

```
# exec設定確認ツール
https://github.com/aws-containers/amazon-ecs-exec-checker
```

1. 以下の Session Manager 権限必要のため、ポリシーを作成

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssmmessages:CreateControlChannel",
        "ssmmessages:CreateDataChannel",
        "ssmmessages:OpenControlChannel",
        "ssmmessages:OpenDataChannel"
      ],
      "Resource": "*"
    }
  ]
}
```

2. 上記ポリシーを ecsTaskExecutionRole にアタッチ

3. 既存 ECS Service のコマンド実行を true にする

```shell
aws ecs update-service --region ap-northeast-1 --cluster ${CLUSTER_NAME} --service ${SERVICE_NAME} --enable-execute-command
```

4. 必要に応じてタスクを再起動する

5. 確認コマンド

```shell
aws ecs describe-tasks --cluster ${CLUSTER_NAME} --task ${TASK_ID}
```

6. ログインする

```shell
aws ecs execute-command --region ap-northeast-1 --cluster ${CLUSTER_NAME} --task ${TASK_ID} --container ${CONTAINER_NAME} --interactive --command "bin/sh"
```
