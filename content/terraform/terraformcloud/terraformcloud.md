---
date: 2026-09-22
---
# Try Terraform cloud with Github

## 目次

1. アーキテクチャー
2. Slack 通知アプリ設定
3. TerraformCloud 設定
4. 実行および確認
5. Slack 通知確認

### 1. アーキテクチャー

#### 1.1 以下のようなアーキテクチャーを作成していく

![image](./img/architecture-2023-03-07.svg)

#### 1.2 重要機能

Terraform Cloud を利用する上で、以下いくつかの用語を理解する必要がある  
■ Organization  
■ WorkSpace  
■ Variable  
■ Variable Set  
■ Run Trigger

### 2. Slack 通知アプリ設定

以下の Url より Slack App 作成ページにアクセスする  
https://api.slack.com/

Create an app をクリックする
![image](./img/slack_2022-12-06_12-46-31.png)

From Scratch をクリックする
![image](./img/slack_2022-12-06_12-46-55.png)

App Name を入力し、Create App をクリックする
![image](./img/slack_2022-12-06_12-47-29.png)

Incoming WebHooks をクリックし、On に変更する
![image](./img/slack_2022-12-06_12-47-58.png)

Add New Webhook to Workspace をクリックする
![image](./img/slack_2022-12-06_12-49-19.png)

事前に作成した通知用チャンネルを選択し、Allow をクリックする
![image](./img/slack_2022-12-06_12-49-40.png)

生成された WebHook Url をコピーする
![image](./img/slack_2022-12-06_12-50-54.png)

### 3. Terraform Cloud 設定

以下の Url より Terraform Cloud のアカウントを作成する
https://cloud.hashicorp.com/products/terraform

作成後、以下の画面が表示され、Start from scratch をクリックする
![image](./img/terraform_2022-12-02_17-18-18.png)

新しい Organization を作成する
![image](./img/terraform_2022-12-02_17-19-41.png)

今回は Github と連携のため、VersionControlWorkflow を選択する
![image](./img/terraform_2022-12-02_17-21-33.png)

Github.com を選択する
![image](./img/terraform_2022-12-02_17-21-51.png)

認証を許可する
![image](./img/terraform_2022-12-02_17-22-08.png)

選択されたリポジトリのみ許可する
![image](./img/terraform_2022-12-02_17-23-51.png)

二段階認証がある場合は認証を行う
![image](./img/terraform_2022-12-02_17-24-01.png)

連携するリポジトリを選択する
![image](./img/terraform_2022-12-02_17-24-47.png)

新しい Work Space を作成する
![image](./img/terraform_2022-12-02_17-26-20.png)

WorkingDirectory を設定する（他の設定値は任意）
![image](./img/terraform_2022-12-02_17-26-36.png)

Continue to Workspace overview をクリックする
![image](./img/terraform_2022-12-02_17-26-53.png)

Configure Variables をクリックする
![image](./img/terraform_2022-12-02_17-28-42.png)

Terraform デプロイ必要なキーやシークレットなどを定義する
![image](./img/terraform_2022-12-06_19-19-13.png)

Settings ＞ Notifications を開き、新規通知設定を行う
![image](./img/terraform_2022-12-06_19-21-39.png)

Slack を選択し、WebHook Url を入力する
![image](./img/terraform_2022-12-06_19-22-23.png)

以上で Terraform Cloud の設定が完了

### 4. 実行および確認

Start New run をクリックする
![image](./img/terraform_2022-12-06_19-25-22.png)

Plan Only を選択し、Start Run をクリックする
![image](./img/terraform_2022-12-06_19-25-43.png)

Plan が正しく実行され、デプロイ対象リソース一覧が表示される
![image](./img/terraform_2022-12-06_19-28-20.png)

Comment を記載し、Apply を実行する
![image](./img/terraform_2022-12-06_19-28-46.png)

Apply が正しく実行されていることを確認する
![image](./img/terraform_2022-12-06_19-30-15.png)

Terraform State ファイルは左ペインの States より確認可能
![image](./img/terraform_2022-12-06_19-30-58.png)

また Overview より、デプロイされたリソース数など確認可能
![image](./img/terraform_2022-12-06_19-31-33.png)

Destroy の場合は Settings ＞ Destruction and Deletion を選択し、Queue Destroy plan をクリックする
![image](./img/terraform_2022-12-06_19-32-15.png)

Workspace 名を入力し、削除を実行する
![image](./img/terraform_2022-12-06_19-32-36.png)

### 5. Slack 通知確認

以上で Terraform Cloud を使って一連のデプロイができた  
Slack 通知もこんな感じで送信される

![image](./img/slack_2022-12-06_19-36-00.png)
