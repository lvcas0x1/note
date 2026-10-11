---
date: 2026-09-22
---
##### App Service Docker Timeout

```shell
WEBSITES_PORT = 3000 (default 80)
WEBSITES_CONTAINER_START_TIME_LIMIT = 1800 (default 230 sec)
```

##### App Service with Easy Auth

Get access token and refresh token

```
Microsoft: In https://resources.azure.com, do the following steps:
1.At the top of the page, select Read/Write.
2.In the left browser, navigate to subscriptions > <subscription_name> > resourceGroups > <resource_group_name> >
  providers > Microsoft.Web > sites > <app_name> > config > authsettingsV2.
3.Click Edit.
4.Modify the following property.
```

```json
"identityProviders": {
  "azureActiveDirectory": {
    "login": {
      "loginParameters": ["scope=openid profile email offline_access user.read"]
    }
  }
}
```

```
Once your provider is configured, you can find the refresh token and the expiration time for the access token in the token store.
To refresh your access token at any time, just call /.auth/refresh in any language.
The following snippet uses jQuery to refresh your access tokens from a JavaScript client.
```

```javascript
function refreshTokens() {
  let refreshUrl = "/.auth/refresh";
  $.ajax(refreshUrl)
    .done(function () {
      console.log("Token refresh completed successfully.");
    })
    .fail(function () {
      console.log("Token refresh failed. See application logs for details.");
    });
}
```

```
If a user revokes the permissions granted to your app, your call to /.auth/me may fail with a 403 Forbidden response.
To diagnose errors, check your application logs for details.
```

##### App Service with Next-Auth

App Registrations

| URI                                                                  | Description |
| -------------------------------------------------------------------- | ----------- |
| http://localhost:3000/api/auth/callback/azure-ad                     | local       |
| https://appservice-name.azurewebsites.net/api/auth/callback/azure-ad | Production  |

App Service

| Variable               | Value                                     |
| ---------------------- | ----------------------------------------- |
| AZURE_AD_CLIENT_ID     | -                                         |
| AZURE_AD_CLIENT_SECRET | -                                         |
| AZURE_AD_TENANT_ID     | -                                         |
| NEXTAUTH_URL           | https://appservice-name.azurewebsites.net |

##### API 監視手順

1. 「App Service」 > 「Application Insights」を開き、有効になっていることを確認（無効となっている場合は新規 ApplicationInsights を有効にする）
2. 「Link to an Application Insights resource」配下のリソース名をクリック
3. 「Availability」を開く。
4. 「Add Standard test」をクリック
5. 以下項目を入力する

| Item                                          | Value | Description                                                                      |
| --------------------------------------------- | ----- | -------------------------------------------------------------------------------- |
| Test name                                     |       | テストの名前                                                                     |
| Url                                           |       | https://< AppServiceName >.azurewebsites.net                                     |
| Parse dependent requests                      |       | img,css など全ての request が正しくロードされることを確認する。default: disabled |
| Enable retries for availability test failures |       | リトライ default: enabled                                                        |
| Enable SSL certificate validity               |       | SSL 証明書チェック default: enabled                                              |
| Proactive lifetime check                      |       | SSL 証明書有効期間チェックう default: 7 days                                     |
| Test frequency                                |       | テスト間隔 5(default), 10, 15 mins                                               |
| Test locations                                |       | テストロケーション                                                               |

Standard test info(option):

| Item              | Value | Description                                                                  |
| ----------------- | ----- | ---------------------------------------------------------------------------- |
| HTTP Request verb |       | HTTP Request メソッド: GET(default), HEAD, POST, PUT, PATCH, DELETE, OPTIONS |
| Add custom header |       | カスタムヘッダ追加                                                           |

Success criteria(option):

| Item          | Value | Description                               |
| ------------- | ----- | ----------------------------------------- | ------------------------------------------------ |
| Test Timeout  |       | タイムアウト時間: 30, 60, 90, 120 seconds |
| HTTP response |       | レスポンスコード定義 default: 200         |
| Content match |       |                                           | レスポンスコンテンツマッチング default: disabled |

Alerts(option):

| Item         | Value | Description               |
| ------------ | ----- | ------------------------- |
| Alert status |       | アラート default: enabled |

6. 可用性テスト作成後、対象のテストをクリックし、「Open Rules(Alerts) page」をクリックする
7. 先ほど作成したアラートルールを開く
8. 「Edit」をクリック
9. 「ActionsGroup」を追加する

| Item              | Value | Description                      |
| ----------------- | ----- | -------------------------------- |
| Resource group    |       | 保存先のリソースグループ         |
| Region            |       | 作成レベル                       |
| Action group name |       | ActionGroup の名前               |
| Display name      |       | 表示名                           |
| Notifications     |       | 通知方法及び通知先を選択         |
| Actions           |       | アラート発生時のアクションを定義 |
| Tags              |       | タグ                             |

10. 設定後、通知先に設定完了通知が届く
