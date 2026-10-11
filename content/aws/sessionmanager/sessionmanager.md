---
date: 2026-09-22
---
##### セッションマネージャーをセットアップ

##### Session Manager を利用する二つの方法

1. Internet 通信あり：NATGateway 作成
2. Internet 通信なし：VPC エンドポイント作成

##### 共通設定事項：

1. Install SSM Agent
2. Create policy - json(S3 Arn 及び KMS Key は事前に作成する必要がある)

```
参考：https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/getting-started-create-iam-instance-profile.html
```

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
        "ssmmessages:OpenDataChannel",
        "ssm:UpdateInstanceInformation"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["logs:CreateLogStream", "logs:PutLogEvents", "logs:DescribeLogGroups", "logs:DescribeLogStreams"],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject"],
      "Resource": "arn:aws:s3:::DOC-EXAMPLE-BUCKET/s3-bucket-prefix/*"
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetEncryptionConfiguration"],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["kms:Decrypt"],
      "Resource": "key-name"
    },
    {
      "Effect": "Allow",
      "Action": "kms:GenerateDataKey",
      "Resource": "*"
    }
  ]
}
```

3. Create role with policy
4. Attach Role to EC2
5. Restart EC2 instance

##### VPC Endpoint の設定

```
参考：https://repost.aws/ja/knowledge-center/ec2-systems-manager-vpc-endpoints
```

以下の Endpoint を作成：

```
com.amazonaws.[region].ssm
com.amazonaws.[region].ec2messages
com.amazonaws.[region].ssmmessages
com.amazonaws.[region].kms
com.amazonaws.[region].s3
com.amazonaws.[region].logs
```

##### 接続ユーザーの指定：

1. ssm-user を利用して Session Manager より EC2 にログイン
2. 新規ユーザー作成

```shell
sudo useradd -m USERNAME
sudo passwd USERNAME
```

3. Session Manager > Linux インスタンスの Run As サポートを有効にする
4. Operating system user name: デフォルトユーザー設定
5. IAM コンソールより付与対象の IAM ユーザーのタグを以下のように追加

```
key: SSMSessionRunAs
value: USERNAME (os username)
```

6. Session Manager より EC2 にログイン
7. whoami でログインユーザー確認
