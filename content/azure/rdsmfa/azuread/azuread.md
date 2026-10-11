---
date: 2026-09-22
---
## Integrate Remote Desktop Gateway using the NPS extension&Azure AD 構築手順

### 目次

1. 全体構成図 - 認証の流れ
2. RemoteGateway 構築 - サービスインストール
3. RemoteGateway 構築 - CAP サーバ設定
4. RemoteGateway 構築 - 自己証明書生成＆適用
5. RemoteGateway 構築 - RAP ポリシー作成
6. RemoteGateway 構築 - Network Policy Server 設定
7. NPS Server 構築 - サービスインストール
8. NPS Server 構築 - NPS Extension for Azure MFA インストール
9. NPS Server 構築 - Network Policy Server 設定
10. ログイン確認
11. ログ確認 - RemoteGateway
12. ログ確認 - NPS Server
13. まとめ（注意事項）

### 1. 全体構成図 - 認証の流れ

![image](./img/Snipaste_2022-11-17_20-37-13.png)

### 2. RemoteGateway 構築 - サービスインストール

RemoteGateway サーバにログインし、IE Enhanced Security をオフに変更する
![image](./img/Snipaste_2022-11-17_17-20-00.png)
LocalAD サーバのプライベート IP を DNS に設定する
![image](./img/Snipaste_2022-11-17_17-21-10.png)
LocalAD ドメインに参加する
![image](./img/Snipaste_2022-11-17_17-21-40.png)
再起動後、RemoteDesktopService をインストールする
![image](./img/Snipaste_2022-11-17_17-32-41.png)
RemoteDesktopGateway を選択する
![image](./img/Snipaste_2022-11-17_17-33-05.png)
デフォルト値のままで進み、インストールが正常に完了したことを確認する
![image](./img/Snipaste_2022-11-17_17-34-38.png)
今後のため、NetworkPolicyServer の起動オプションを自動に変更する
![image](./img/Snipaste_2022-11-17_17-35-43.png)
今後のため、RemoteDesktopGateway の起動オプションを自動に変更する
![image](./img/Snipaste_2022-11-17_17-36-00.png)

### 3. RemoteGateway 構築 - CAP サーバ設定

RDGatewayManager のプロパティを開き、 RDCAPStore を Central に変更する  
IP アドレスは NPS サーバのプライベート IP を指定する
![image](./img/Snipaste_2022-11-17_17-37-52.png)
ShardSecret を入力する
![image](./img/Snipaste_2022-11-17_17-38-12.png)

### 4. RemoteGateway 構築 - 自己証明書生成＆適用

SSL Certificate タブを開き、自己証明書を作成し、OK をクリックする（ドメイン名も可）
![image](./img/Snipaste_2022-11-17_17-38-40.png)
自己証明書が正しくインポートされている
![image](./img/Snipaste_2022-11-17_17-39-05.png)
CentralNetworkPolicyServer は先ほど設定した NPS サーバが表示されている
![image](./img/Snipaste_2022-11-17_17-39-22.png)

### 5. RemoteGateway 構築 - RAP ポリシー作成

RAP ポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-39-39.png)
任意
![image](./img/Snipaste_2022-11-17_17-39-48.png)
Domain Users
![image](./img/Snipaste_2022-11-17_17-40-01.png)
Allow users to connect to any ...
![image](./img/Snipaste_2022-11-17_17-40-06.png)
3389
![image](./img/Snipaste_2022-11-17_17-40-12.png)
Finish
![image](./img/Snipaste_2022-11-17_17-40-18.png)
以下のように RAP が作られる
![image](./img/Snipaste_2022-11-17_17-40-32.png)

### 6. RemoteGateway 構築 - Network Policy Server 設定

NPS（local）を開き、新規 RADIUS Client を作成する  
FriendlyName：任意  
IP アドレス：NPS サーバのプライベート IP  
SharedSecret：先ほどと同じもの
![image](./img/Snipaste_2022-11-17_17-41-39.png)
既存の RmoteRADIUSServerGroups を開き、Edit をクリックする
![image](./img/Snipaste_2022-11-17_17-42-01.png)
以下のように Timeout 時間を指定する
![image](./img/Snipaste_2022-11-17_17-42-13.png)
ConnectionRequestPolicies を開き、新規名前は MFA Request Forward の接続ポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-43-17.png)
以下のように接続ポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-44-56.png)
もう一つ名前は MFA Request No Forward の接続ポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-45-14.png)
以下のように接続ポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-46-03.png)
デフォルトの TS Gateway ...のポリシーを無効にする
![image](./img/Snipaste_2022-11-17_17-46-12.png)
NPS サーバをドメインに登録する
![image](./img/Snipaste_2022-11-17_17-46-25.png)

### 7. NPS Server 構築 - サービスインストール

NPS Server サーバにログインし、IE Enhanced Security をオフに変更する
![image](./img/Snipaste_2022-11-17_17-09-07.png)
LocalAD サーバのプライベート IP を DNS に設定する
![image](./img/Snipaste_2022-11-17_17-09-16.png)
LocalAD ドメインに参加する
![image](./img/Snipaste_2022-11-17_17-09-41.png)
再起動後、Network Policy and Access Services をインストールする
![image](./img/Snipaste_2022-11-17_17-22-30.png)
デフォルト値のままで進み、インストールが正常に完了したことを確認する
![image](./img/Snipaste_2022-11-17_17-23-32.png)
今後のため、NetworkPolicyServer の起動オプションを自動に変更する
![image](./img/Snipaste_2022-11-17_17-24-03.png)

### 8. NPS Server 構築 - NPS Extension for Azure MFA インストール

以下の Url から NPS Extension for Azure MFA をダウンロードする  
https://www.microsoft.com/en-us/download/details.aspx?id=54688
![image](./img/Snipaste_2022-11-17_17-25-05.png)
インストールする
![image](./img/Snipaste_2022-11-17_17-25-33.png)
インストール成功
![image](./img/Snipaste_2022-11-17_17-25-43.png)
以下図のようにディレクトリに進み、スクリプトを実行する
![image](./img/Snipaste_2022-11-17_17-26-14.png)
Y で同意
![image](./img/Snipaste_2022-11-17_17-26-41.png)
A で全て同意
![image](./img/Snipaste_2022-11-17_17-26-59.png)
アクセス許可求められる場合は許可する（数回聞かれることあります、スクリプト実行失敗の場合は再度実行すれば問題ない）
![image](./img/Snipaste_2022-11-17_17-27-30.png)
AzureAD のグローバル管理者アカウントでログイン
![image](./img/Snipaste_2022-11-17_17-28-56.png)
Azure Portal からテナント ID コピーして貼り付ければ、インストール完了する
![image](./img/Snipaste_2022-11-17_17-29-23.png)

### 9. NPS Server 構築 - Network Policy Server 設定

NPS（local）を開き、新規 RADIUS クライアントを作成する
![image](./img/Snipaste_2022-11-17_17-47-43.png)
名前は RDGW、Add をクリックする
![image](./img/Snipaste_2022-11-17_17-48-08.png)
ServerIP は RemoteGateway のプライベート IP を入力する
![image](./img/Snipaste_2022-11-17_17-48-23.png)
以下のように Timeout 時間を設定する
![image](./img/Snipaste_2022-11-17_17-48-32.png)
新規名前は FromRDGW の接続ポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-48-54.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-50-00.png)
新規名前は ToRDGW のポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-50-14.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-51-05.png)
NetworkPolicy を作成する  
既存のポリシーをコピーして利用する
![image](./img/Snipaste_2022-11-17_17-52-22.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-52-50.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-53-04.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-53-13.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-53-26.png)
以下のようにポリシーを作成する
![image](./img/Snipaste_2022-11-17_17-53-47.png)
デフォルト値のままで OK をクリックする
![image](./img/Snipaste_2022-11-17_17-53-56.png)
ポリシーが正しく作成される
![image](./img/Snipaste_2022-11-17_17-54-09.png)
NPS サーバをドメインに登録する
![image](./img/Snipaste_2022-11-17_17-54-26.png)

### 10. ログイン確認

自己証明書利用したため、証明書をログイン元の PC に登録する
![image](./img/Snipaste_2022-11-17_17-55-20.png)
以下図のように RemoteGateway サーバを指定する
![image](./img/Snipaste_2022-11-17_17-56-37.png)
ログイン先サーバのプライベート IP とユーザー名を指定する
![image](./img/Snipaste_2022-11-17_17-56-44.png)
RemoteGateway サーバの認証情報を入力する
![image](./img/Snipaste_2022-11-17_17-57-03.png)
携帯アプリにプッシュ通知が送信される  
Approve する  
![image](./img/Snipaste_2022-11-17_18-54-15.png ":size=200")  
ログイン先サーバの認証情報を入力する
![image](./img/Snipaste_2022-11-17_17-57-20.png)
Yes
![image](./img/Snipaste_2022-11-17_17-57-30.png)
正しくログインできることを確認する
![image](./img/Snipaste_2022-11-17_17-57-37.png)

### 11. ログ確認 - RemoteGateway

![image](./img/Snipaste_2022-11-17_17-59-45.png)
![image](./img/Snipaste_2022-11-17_17-59-56.png)
![image](./img/Snipaste_2022-11-17_18-00-07.png)
![image](./img/Snipaste_2022-11-17_18-02-21.png)

### 12. ログ確認 - NPS Server

![image](./img/Snipaste_2022-11-17_18-01-26.png)

### 13. まとめ（注意事項）

1. AzureAD Connector 構築済み
2. AzureAD GlobalAdminstrator ユーザー、Tenat ID 必要
3. NPS サーバは AD サーバと同居できる（公式ドキュメント）、今回の手順は AD、RDGW、NPS サーバ三つ分けて構築している
4. すべてのログイン認証を RDGW⇨NPS⇨AzureAD に転送するため、ユーザーによって除外ができない／除外したい場合はもう一台 RDGateway を用意する？
5. MFA は MicrosoftAuthenticator を利用する＆携帯 Push 通知有効にする必要ある

※Azure でデプロイする場合の注意事項

1. デフォルトファイアウォールがオンになっているため、オフにすること
2. DNS 設定は Windows OS 内ではなく、ポータルより変更すること（MS 推奨）

※NPS サーバ経由で MFA 疎通できない場合こちらの記事を参考にレジストリを修正する

https://learn.microsoft.com/en-us/azure/active-directory/authentication/how-to-mfa-number-match
