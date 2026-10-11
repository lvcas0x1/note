---
date: 2026-09-22
---
## AzureAD Connector 構築手順

### 目次

1. IESec オフにする
2. AD Role インストールする
3. AzureAD カスタムドメイン追加
4. AzureAD 同期用ローカルユーザー作成
5. AzureAD Connector インストール&同期設定
6. 動作確認 - ローカルユーザー MFA 有効化
7. 動作確認 - Azure Portal ログイン
8. まとめ（注意事項）

### 1. IESec オフにする

サーバマネージャーを起動する
![image](./img/AzureADConnector_2022-11-16_22-27-44.png)
LocalServer を開く
![image](./img/AzureADConnector_2022-11-16_22-27-50.png)
IE Enhanced Security Configuration をオフに変更する
![image](./img/AzureADConnector_2022-11-16_22-27-57.png)
IE Enhanced Security Configuration の設定を確認する
![image](./img/AzureADConnector_2022-11-16_22-28-07.png)

### 2. AD Role インストールする

Dashboard に戻り、Addroles and features を開き、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-28-17.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-28-22.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-28-27.png)
ActiveDirectoryDomainService をチェックする
![image](./img/AzureADConnector_2022-11-16_22-28-39.png)
Add Feature をクリックする
![image](./img/AzureADConnector_2022-11-16_22-28-45.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-28-54.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-28-59.png)
Install をクリックする
![image](./img/AzureADConnector_2022-11-16_22-29-06.png)
Promote this server to a domain controller をクリックする
![image](./img/AzureADConnector_2022-11-16_22-30-01.png)
Add a new forest を選択し、ドメイン名を入力する
![image](./img/AzureADConnector_2022-11-16_22-30-20.png)
DSRM パスワードを入力し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-31-00.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-31-05.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-31-21.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-31-31.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-31-36.png)
Install をクリックする
![image](./img/AzureADConnector_2022-11-16_22-31-49.png)
Install 完了後、サーバ自動ログオフされる
![image](./img/AzureADConnector_2022-11-16_22-38-10.png)

### 3. Azure カスタムドメイン追加

Azure AD にカスタムドメインを追加する
![image](./img/AzureADConnector_2022-11-16_22-40-45.png)
DNS サービスにてドメイン認証を行う
![image](./img/AzureADConnector_2022-11-16_22-40-59.png)

### 4. 同期用ローカルユーザー作成

Active Directory Users and Computers を開く
![image](./img/AzureADConnector_2022-11-16_22-42-36.png)
後ほど Azure AD 同期用に、新規 OU を作成する
![image](./img/AzureADConnector_2022-11-16_22-43-03.png)
名前を AD Connect にし、OK をクリックする
![image](./img/AzureADConnector_2022-11-16_22-43-14.png)
作成した OU 内に新規ユーザーを作成する
![image](./img/AzureADConnector_2022-11-16_22-43-26.png)
名前を Localuser にし、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-43-41.png)
パスワードを入力し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-43-55.png)
Finish をクリックする
![image](./img/AzureADConnector_2022-11-16_22-44-00.png)

### 5. AzureAD Connector インストール&同期設定

AzureADConnector をダウンロードする  
https://www.microsoft.com/en-us/download/details.aspx?id=47594
![image](./img/AzureADConnector_2022-11-16_22-45-48.png)
ダウンロードしたファイルを実行する
![image](./img/AzureADConnector_2022-11-16_22-46-20.png)
カスタムをクリックする
![image](./img/AzureADConnector_2022-11-16_22-46-27.png)
必要に応じて設定変更し、Install をクリックする
![image](./img/AzureADConnector_2022-11-16_22-46-41.png)
必要に応じて設定変更し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-48-11.png)
Azure AD の GlobalAdministrator 権限を持つユーザー情報を入力する
![image](./img/AzureADConnector_2022-11-16_22-48-55.png)
Add Directory をクリックする
![image](./img/AzureADConnector_2022-11-16_22-49-14.png)
ローカルユーザーの Enterprise Admins 権限を持つユーザー情報を入力する
![image](./img/AzureADConnector_2022-11-16_22-49-36.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-49-47.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-49-56.png)
同期対象グループを先ほどの OU を選択し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-50-18.png)
必要に応じて設定変更し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-50-26.png)
必要に応じて設定変更し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-50-33.png)
必要に応じて設定変更し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-50-38.png)
必要に応じて設定変更し、Install をクリックする
![image](./img/AzureADConnector_2022-11-16_22-50-47.png)
Exit をクリックする
![image](./img/AzureADConnector_2022-11-16_22-52-53.png)

### 6. 動作確認 - ローカルユーザー MFA 有効化

Azure AD を開き、選択されたユーザーが同期されていることを確認し、Per-user MFA を開く
![image](./img/AzureADConnector_2022-11-16_22-54-11.png)
MFA を Enable にする
![image](./img/AzureADConnector_2022-11-16_22-54-41.png)
enable multi-factor auth をクリックする
![image](./img/AzureADConnector_2022-11-16_22-55-01.png)

### 7. 動作確認 - Azure Portal ログイン

portal.azure.com にアクセスし、ユーザー名を入力する
![image](./img/AzureADConnector_2022-11-16_22-55-51.png)
パスワードを入力する
![image](./img/AzureADConnector_2022-11-16_22-56-03.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-56-15.png)
MobileApp を選択し、SetUp をクリックする
![image](./img/AzureADConnector_2022-11-16_22-56-37.png)
QR コードをスキャンし、MFA を登録する
![image](./img/AzureADConnector_2022-11-16_22-56-54.png)
MFA 認証を行い、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-58-01.png)
携帯番号を登録し、Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-58-29.png)
Done をクリックする
![image](./img/AzureADConnector_2022-11-16_22-58-40.png)
Next をクリックする
![image](./img/AzureADConnector_2022-11-16_22-58-46.png)
携帯と EMail の認証を行う
![image](./img/AzureADConnector_2022-11-16_22-59-06.png)
Finish をクリックする
![image](./img/AzureADConnector_2022-11-16_23-00-03.png)
Yes をクリックする
![image](./img/AzureADConnector_2022-11-16_23-00-10.png)
Azure Portal にログインできることを確認する
![image](./img/AzureADConnector_2022-11-16_23-00-19.png)

### 8. まとめ（注意事項）

1. Azure AD のカスタムドメインにローカル AD ドメインを追加  
   ■ ドメイン認証必要のため、実在するドメイン取得し、ローカル AD を構築する
2. AzureADConnector 設定の際  
   ■ AzureAD ユーザー GlobalAdmins 権限のユーザーを指定  
   ■ ローカルユーザー EnterpriseAdmins グループに所属するユーザーを指定
3. Azure MFA 利用のための P1/P2 ライセンス、または MFA 機能を含む M365 ライセンス

※Azure でデプロイする場合の注意事項

1. デフォルトファイアウォールがオンになっているため、オフにすること
2. DNS 設定は Windows OS 内ではなく、ポータルより変更すること（MS 推奨）
