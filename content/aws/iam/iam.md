---
date: 2026-09-22
---
##### 認証方法及び優先順位

```text
1. コード内、またはCLI
2. 環境変数
・AWS_ACCESS_KEY_ID
・AWS_SECRET_ACCESS_KEY
3. プロファイル ~/.aws/credentials
4. インスタンスプロファイル（IAMロール）
```

##### ポリシー

1. 何ベースのポリシーか？

リソースベースポリシー: プリンシパル

```text
リソースベースポリシーは、リソース自体（例：ファイル、データベース、サービスなど）に対するアクセス制御ルールを定義します。
このタイプのポリシーは、特定のリソースがどのようにアクセスされるべきかを定義し、それに基づいてアクセス権を付与または拒否します。
リソースベースポリシーは通常、リソース自体の属性や、リソースに対する操作（読み取り、書き込み、削除など）を基にアクセス制御を行います。
```

```json
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Principal": { "AWS": "arn:aws:iam::777788889999:user/bob" },
    "Action": ["s3:PutObject", "s3:PutObjectAcl"],
    "Resource": "arn:aws:s3:::example-bucket/*"
  }
}
```

アイデンティティベースポリシー

```text
アイデンティティベースポリシーは、ユーザーやグループ、または他のアイデンティティ（例：ロールやサービスアカウント）に基づいてアクセス制御ルールを定義します。
このタイプのポリシーは、特定のアイデンティティがどのリソースにアクセスできるか、またはどの操作を実行できるかを定義します。
アイデンティティベースポリシーは通常、ユーザーのロールやグループメンバーシップ、または他のアイデンティティ属性を基にアクセス制御を行います。
```

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:Get*", "s3:List*"],
      "Resource": "*"
    }
  ]
}
```

2. ポリシーの適用先

```
(a)IAM User
(b)IAM Group
(c)IAM Role
```

3. ポリシーの種類

```
管理ポリシー：「ポリシー」単体として存在できて、複数のプリンシパルエンティティにつけたり外したりできる。
  ・AWS管理ポリシー：予めAWSによって用意されている管理ポリシー。大まかに各AWSサービスごとにフル権限と読み取り権限が用意されている。
  ・カスタマー管理ポリシー：利用者が独自に作成する管理ポリシー。ジェネレータで作成するか手書きで作成する。
インラインポリシー：「ポリシー」単体では存在できず、プリンシパルエンティティと1対1の関係にある。
```

4. ポリシーの評価基準

```
(a)デフォルトでは拒否
(b)デフォルトの拒否より明示的な許可が強い
(c)明示的な許可より明示的な拒否が強い
```

##### アクセス許可の境界

```
IAMの拡張機能です。IAMエンティティに付与できるアクセス許可の最大を設定するのに使用される。
アクセス許可の境界が設定されたエンティティでは、アイデンティティベースのポリシーとアクセス許可の境界の両方で許可されたアクションのみ実行できる。
```

署名 v4 認証(SigV4)??

##### Python で IAM 操作

```shell
aws configure

region=$(curl http://169.254.169.254/latest/meta-data/placement/region -s)

# aws cloud9が使用するリクエストの認証情報を確認
aws sts get-caller-identity

# s3 - list all
aws s3 ls

# s3 - get bucket name with deleteme
bucketToDelete=$(aws s3api list-buckets --output text --query 'Buckets[?contains(Name, `deleteme`) == `true`] | [0].Name')

# s3 - remove bucket
aws s3 rb s3://$bucketToDelete

# --debug
# 認証情報の検出、指定されたパラメータの解析、AWSに送信されるリクエストの構築、内容、レスポンスの内容、フォーマット済み出力
aws s3 rb s3://$bucketToDelete --debug

# policy - get arn
policyArn=$(aws iam list-policies --output text --query 'Policies[?PolicyName == `S3-Delete-Bucket-Policy`].Arn')

# policy - get content
aws iam get-policy-version --policy-arn $policyArn --version-id v1

# policy - attach
aws iam attach-role-policy --policy-arn $policyArn --role-name notes-application-role

# role - list attached policy
aws iam list-attached-role-policies --role-name notes-application-role

```
