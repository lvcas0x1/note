---
date: 2026-09-22
---
## REST

REST: Representational State Transfer(分散型システムにおける設計原則群)  
RESTful: RESTで求められる原則に従っているのこと

## URI設計

- 短く入力しやすい
- 人間が理解できる
- すべて小文字
- 単語はハイフンで繋げる
- 単語は複数形を利用する
- エンコードを必要とする文字を使わない
- サーバ側のアーキテクチャが反映されていない
- 改造しやすい
- ルールが統一されている

**クエリパラメータとパスパラメータ**

- 一意なリソースを表すかつ、条件を追加する場合: パスパラメータ
  - GET http://api.example.com/users?page=3
- 一意なリソースを表す: クエリパラメータ
  - GET http://api.example.com/users/123

## ステータスコード

| ステータスコード | 意味                               |
| ---------------- | ---------------------------------- |
| 100              | 情報                               |
| 200              | 成功                               |
| 300              | リダイレクト                       |
| 400              | クライアントサイドに起因するエラー |
| 500              | サーバサイドに起因するエラー       |

**200番台**
|種類|説明|
|-|-|
|200 OK|リクエストが成功した、本文にデータ含まれる|
|201 Created|リクエストが成功し、新しいリソースが作成された。ヘッダーのLocationに新しいリソースへのURLを含める|
|202 Accepted|非同期ジョブを受け付けた、実際の処理結果は別途受け取る|
|204 No Content|リクエスト成功した、レスポンスデータがない|

**300番台**  
API利用者はリダイレクトを実装していないことが多いのでREST APIは基本3xxを利用しない

**400番台**
|種類|説明|
|-|-|
|400 Bad Request|その他エラー、エラー起因はクライアントサイドのリクエスト|
|401 Unauthorized|認証されていない|
|403 Forbidden|リソースに対するアクセスが許可されていない|
|404 Not Found|リクエストされたリソースが存在しない|
|409 Conflict|リソースが競合して処理が完了できなかった|
|429 Too Many Requests|アクセス回数が制限回数を超えたため処理できなかった|

**500番台**
|種類|説明|
|-|-|
|500 Internal Server Error|サーバサイドのアプリケーションエラーが発生した|
|503 Service Unavailable|サービスが一時的に利用できないこと（メンテナンス期間、過負荷で応答できないなど）|

## HTTPメソッドとステータスコード

**GET**
|成功|失敗|
|-|-|
|200 OK|400 Bad Request|
|304 Not Modified|401 Unauthorized|
||403 Forbidden|
||404 Not Found|
||429 Too Many Requests|
||500 Internal Server Error|
||503 Service Unavailable|

**POST**
|成功|失敗|
|-|-|
|200 OK|400 Bad Request|
|201 Created|401 Unauthorized|
|202 Accepted|403 Forbidden|
||409 Conflict|
||429 Too Many Requests|
||500 Internal Server Error|
||503 Service Unavailable|

**PUT**
|成功|失敗|
|-|-|
|200 OK|400 Bad Request|
|201 Created|401 Unauthorized|
|204 No Content|403 Forbidden|
||404 Not Found|
||409 Conflict|
||429 Too Many Requests|
||500 Internal Server Error|
||503 Service Unavailable|

**DELETE**
|成功|失敗|
|-|-|
|200 OK|400 Bad Request|
|202 Accepted|401 Unauthorized|
|204 No Content|403 Forbidden|
||404 Not Found|
||429 Too Many Requests|
||500 Internal Server Error|
||503 Service Unavailable|

## データフォーマット

- XML: Content-Type: application/xml
- JSON: Content-Type: application/json
- JSONP(JSON with Padding): Content-Type: application/javascript
  - クロスドメインでデータ受け渡しできる

リクエストヘッダーでフォーマット指定: Accept: application/json

## データの内部構造

- エンベロープは使わない
- オブジェクトはできるだけフラットにする
- ページネーションをサポートする情報を返す
- プロパティの命名規則はAPI全体で統一する
  - snake_case
  - camelCase
  - PascalCase
- 日付はRFC3339(W3c-DTF)形式を使う
- 大きな数値(64bit整数)は文字列で返す

## エラー表現

- 足りないエラー情報はレスポンスボディに追加
  - エラーコード、エラーメッセージ
- サービス閉塞時は"503" + "Retry-After"を返す

## APIバージョン

パス: http://api.example.com/v1/users  
クエリ: http://api.example.com/users?version=1  
ヘッダー: GET http://api.example.com/users X-Api-Version: 1 非推奨

バージョンの付け方

- セマンティックバージョニング: メジャー.マイナー.パッチ
  - メジャー: 後方互換しない修正
  - マイナー: 後方互換する機能追加
  - パッチ: 後方互換するバグ修正

APIは基本後方互換しないタイミングでつける、メジャーのみ

## Authorizationヘッダー

Authorization: `<Type>` `<Credentials>`

`<Type>`

- Basic: ベーシック認証、IDとパスワードを平文で送信
- Bearer: Oauth2.0(JWT利用する場合)
- Digest: ダイジェスト認証(IDとパスワードをハッシュ化して送信)
- Oauth: Oauth1.0

`<Credentials>`

認証情報を記述する。JWTの場合はJWTの実態が記載されている

## 大量アクセス対策

レートリミット

- 誰に対して: APIキー、ユーザID
- 何に対して: 単一機能、機能群
- 制限回数: 10回、100回
- 単位時間: 10分、1時間

アルゴリズム

- Fixed Window: 特定の期間で倍のアクセスできてしまう。
- Sliding Log: 過去のリクエストログ時間をチェックする。ログ大量保管、破棄を考える必要がある
- Sliding Window: 比率で計算している。

レートリミットを緩和する措置を考える必要がある。

## キャッシュ制御

### ヘッダー

| 分類               | ヘッダー             |
| ------------------ | -------------------- |
| 有効期限による制御 | Expires              |
|                    | Cache-Control + Date |
| 検証による制御     | Last-Modified + Etag |

**Expires**

Expires: Sun, 03 May 2020 12:30:00 GMT

- キャッシュとしていつまで利用可能かの期限を指定
- 過去日を指定すると「リソースが有効期限切れ」であることを意味する
- Cache-Controlヘッダーが同時指定されている場合、Expiresは無視される

**Cache-Control + Date**

Cache-Control: public, max-age=604800  
Date: Sun, 03 May 2020 12:30:00 GMT

Cache-Controlでキャッシュの可否、期限を指定
||||
|-|-|-|
|キャッシュ可否|Public|通信経路上のどこでも保存できる|
||Private|クライアント端末のみ保存できる|
||no-cache|クライアント側に保存される、必ずサーバ側と有効性の確認を行う|
||no-store|クライアント側で保存不可、毎回サーバと通信して全量取得|
|キャッシュ期限|max-age=秒|Dateヘッダーを基準として、新しいとみなせる時間(秒)|

**Last-Modified + Etag**

Last-Modified: Sun, 03 May 2020 12:30:00 GMT
ETag: "33a859234fb89437fe908432da903284f5"

- Last-Modifiedにリソースの最終更新日時を指定
- Etagに特定バージョンを示す文字列を指定
  - コンテンツのハッシュ
  - バージョン番号
  - 最終更新日時のハッシュ

### キャッシュさせる単位

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 212
Cache-COntrol: public
Content-Language: ja
Date: Sun, 03 May 2020 13:30:00 GMT
Server: api.example.com
Vary: Content-Language
```

Vary: キャッシュ判断に利用するヘッダーを指定

## セキュリティ

- XSS: 悪意のある埋め込み脆弱性
  - X-XSS-Protection = 1でXSSフィルタリング有効化
  - X-Frame-Options = DENYでFrameタグ呼び出しを拒否
  - X-Content-Type-Options = nosniffでIE脆弱性対応
- CSRF: 不正のサイトからリクエストを正規サイトへ送信される
  - 許可されていないアクセス元からのリクエストを拒否
    - X-API-Keyでシステム単位で実行可否を判断
    - Authentication Headerで正規のユーザかどうかを判断
  - 攻撃者に推測されにくいトークンの発行/照合処理を実装
    - X-CSRF-TOKENトークンを使って実行可否判断
- HTTP
  - 常時HTTPSで通信
  - SSLは安全に通信を行うためのプロトコル(2015年に使用禁止)
  - TLSは安全に通信を行うためのプロトコル、SSLの後継
  - HTTPS = HTTP + SSL/TLS。 Webで安全に通信するプロトコル
- JSON Web Token(JWT)
  - ヘッダーの"alg": "ES256"など利用して暗号化
  - ペイロードの"aud": "https://api.example.com/" に想定する利用者を指定して、受信時に検証する

## OpenAPI&Swagger

### OpenAPI

|                       |                                                                                           |
| --------------------- | ----------------------------------------------------------------------------------------- |
| OpenAPI               | WSDLやXMLと比較されるような”フォーマット”を意味する。JSONやYAMLでREST API仕様が記述できる |
| OpenAPI Specification | OpenAPIを記述するための書式ルール                                                         |
| Swagger               | OpenAPIを作成、表示、利用するツール群                                                     |
|                       | Swagger Editor: OpenAPIを記述するためのエディター                                         |
|                       | Swagger Codegen: 作成したOpenAPIを読み取って、Stub、クライアントのソースコード生成する    |
|                       | Swagger UI: OpenAPIの仕様書を読み取って、可視化する                                       |

### OpenAPI Specification

https://github.com/OAI/OpenAPI-Specification?tab=readme-ov-file

### Hello World

```yaml
openapi: "3.1.0"

info:
  title: "sample api"
  version: "1.0.0"

paths:
  "/message":
    get:
      summary: "Sample API get operation"
      description: "Sample API get operation"
      responses:
        "200":
          description: "Success operation"
          content:
            application/json:
              schema:
                type: string
                example: "Hello World"
```

### データ型

- integer
  - int32
  - int64
- number
  - float
  - double
- string: JSON Schemaで定義されているformatが利用可能
  - string
  - byte
  - binary
  - date
  - date-time
  - email
  - hostname
  - ipv4
  - uri
  - uuid
- boolean
- object
- array

サンプルSchema:

```yaml
openapi: "3.1.0"

info:
  title: "sample api"
  version: "1.0.0"

paths: {}

components:
  schemas:
    SampleInt32:
      type: integer
      format: int32
    SampleDateString:
      type: string
      format: date
```

### ルートオブジェクト

- openapi: "3.1.0" バージョン指定。必須
- info: APIのメタデータを定義する。必須
- servers: APIを提供するサーバを定義する
- tags: APIを分類するタグを定義する
- paths: APIとして利用可能なパスおよび操作を定義する。必須
- security: API全体にかけるセキュリティ要件
- components: OpenAPI中で利用するオブジェクトをコンポーネント化して再利用可能にする

#### Info

infoオブジェクトでMetadata定義する

- title: APIのタイトル
- description: APIの詳細説明、Markdown使える
- termsOfService:サービス利用規約のUrl
- contact:連絡先情報
- license:ライセンス(https://spdx.org/licenses/)
- version: APIドキュメントのバージョン情報

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  description: |
    # Featrures
    - Get reviews
    - Post reviews
  termsOfService: "https://example.com/terms/"
  contact:
    name: "Customer Support"
    url: "https://example.com/support/"
    email: "support@example.com"
  license:
    name: "MIT License"
  version: "1.0.0"

paths: {}
```

#### Servers

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  description: |
    # Featrures
    - Get reviews
    - Post reviews
  termsOfService: "https://example.com/terms/"
  contact:
    name: "Customer Support"
    url: "https://example.com/support/"
    email: "support@example.com"
  license:
    name: "MIT License"
  version: "1.0.0"

servers:
  - url: "http://localhost:{port}"
    description: "Local development"
    variables:
      port:
        enum: ["3000", "8080"]
        default: "3000"
  - url: "https://integration.example.com"
    description: "Integration Env"
  - url: "https://staging.example.com"
    description: "Staging Env"
  - url: "https://api.example.com"
    description: "Production Env"

paths: {}
```

#### Paths

paths

- path: url
- method: get,post,put,delete
  - summary: 操作概要説明
  - description: 操作詳細説明
  - tags: ["users"] API整理のため
  - deprecated: 操作が廃止したかどうか
  - parameters: リクエストパラメータ
  - requestBody: リクエストボディー
  - responses: レスポンス
  - security: セキュリティ

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  description: |
    # Featrures
    - Get reviews
    - Post reviews
  termsOfService: "https://example.com/terms/"
  contact:
    name: "Customer Support"
    url: "https://example.com/support/"
    email: "support@example.com"
  license:
    name: "MIT License"
  version: "1.0.0"

paths:
  "/shops/{shopId}/reviews":
    get:
      summary: "Get specified shop reviews"
      description: "(description): Get specified shop reviews"
      tags: ["shops"]
      deprecated: false
      parameters:
        - name: shopId
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          description: "Success operation"
```

#### Request Parameters

parameters

- name: パラメータ名
- in: query, header, path, cookie
- description: パラメータの説明
- required: 必須かどうか
- schema: 型情報
- example: サンプルデータ

APIキーをヘッダーに要求する  
ワンタイムトークンをクッキーに要求する

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  version: "1.0.0"

paths:
  "/shops/{shopId}/reviews":
    post:
      summary: "Create specified shop reviews"
      description: "(description): Create specified shop reviews"
      tags: ["shops"]
      deprecated: false
      parameters:
        - name: shopId
          in: path
          required: true
          schema: { type: string }
          example: "nike"
        - name: X-Api-Key
          in: header
          description: "request user's identifier"
          schema: { type: string }
          example: "xxxx-xxxx-xxxx-xxxx"
        - name: token
          in: cookie
          description: "one-time token"
          schema: { type: string }
          example: "xxxx-xxxx-xxxx-xxxx"
      responses:
        "201":
          description: "Success operation"
```

#### Request Body

resquestBody

- description: リクエストボディの説明
- required: リクエストボディが必須かどうか
- content: リクエストボディの内容
  - media-type: application/jsonなど
  - schema: スキーマ定義
  - example: サンプルデータ

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  version: "1.0.0"

paths:
  "/shops/{shopId}/reviews":
    post:
      summary: "Create specified shop reviews"
      description: "(description): Create specified shop reviews"
      tags: ["shops"]
      deprecated: false
      parameters:
        - name: shopId
          in: path
          required: true
          schema: { type: string }
          example: "nike"
      requestBody:
        description: "Contents of review"
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                score: { type: integer, example: 3 }
                comment: { type: string, example: "good" }

      responses:
        "201":
          description: "Success operation"
```

#### Responses

responses

- statuscode
  - description: レスポンスの説明
  - headers: 特別ヘッダーあれば定義
  - content: レスポンスボディー
    - media-type: application/jsonなど
    - schema: データ型

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  version: "1.0.0"

paths:
  "/shops/{shopId}/reviews":
    post:
      summary: "Create specified shop reviews"
      description: "(description): Create specified shop reviews"
      tags: ["shops"]
      deprecated: false
      parameters:
        - name: shopId
          in: path
          required: true
          schema: { type: string }
          example: "nike"
      requestBody:
        description: "Contents of review"
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                score: { type: integer, example: 3 }
                comment: { type: string, example: "good" }

      responses:
        "201":
          description: "Success operation"
          headers:
            location:
              description: "New url of created review"
              schema: { type: string, format: url }
        "400":
          description: "Client error"
          content:
            application/json:
              schema:
                type: object
                properties:
                  code: { type: string, example: "1234" }
                  type: { type: string, example: "InvalidDataFormat" }
                  message: { type: string, example: "Invalid data" }
                  errors:
                    type: array
                    items:
                      type: object
                      properties:
                        field: { type: string, example: "comment" }
                        code: { type: string, example: "1234-001" }
```

#### schema

データ型に依存せずに定義できる項目

- description: データ型の説明
- default: デフォルト値
- nullable: null許容するかどうか
- example: サンプルデータ
- deprecated: Schemaが廃止したかどうか

> integer, number

- multipleOf: 指定された数値の倍数になっているかどうか
- maximum: 最大値
- exclusiveMaximum: 最大値を含まないかどうか
- minimum: 最小値
- exclusiveMinimum: 最小値を含まないかどうか

> string

- minLength: 最小文字数。以上
- maxLength: 最大文字数。以下

> boolean

- なし

> object

- properties: プロパティの項目を定義
- additionalProperties: スキーマ以外のプロパティを許すかどうか
- required: propertiesで定義した中の、必須プロパティ項目を指定
- minProperties: 最小プロパティ数
- maxProperties: 最大プロパティ数

> array

- items: 配列に入れられるスキーマを指定
- minItems: 最小アイテム数
- maxItems: 最大アイテム数
- uniqueItems: 配列内の値の重複を許すかどうか

> enum

- 指定した型で、選択可能な値を設定

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  version: "1.0.0"

paths: {}

components:
  schemas:
    SampleInt:
      type: integer
      minimum: 0
      exclusiveMinimum: false
      maximum: 100
      exclusiveMaximum: true
    SampleString:
      type: string
      nullable: true
      minLength: 10
      maxLength: 30
    SampleBoolean:
      type: boolean
    SampleObject:
      type: object
      properties:
        shop:
          type: object
          properties:
            name: { type: string }
            place: { type: string }
        reviews:
          type: array
          items:
            type: object
            properties:
              score:
                type: integer
                emun: [1, 2, 3, 4, 5]
              comment:
                type: string
```

#### Tags

ルートオブジェクトでTagsを定義し、Pathsの中で利用して分類する

```yaml
openapi: "3.1.0"

info:
  title: "Shop Review api"
  version: "1.0.0"

tags:
  - name: reviews
    description: "Shop review operation"

paths:
  "/shops/{shopId}/reviews":
    parameters:
      - name: shopId
        in: path
        required: true
        schema: { type: string }
    get:
      summary: "Get specified shop reviews"
      description: "(description): Get specified shop reviews"
      tags: ["reviews"]
      responses:
        "200":
          description: "Success operation"

    post:
      summary: "Create specified shop reviews"
      description: "(description): Create specified shop reviews"
      tags: ["reviews"]
      responses:
        "201":
          description: "Success operation"
```

#### Components

components

- schemas
- parameters
- requestBody
- responses
- headers
- securitySchemes

```yaml
openapi: "3.0.0"

info:
  title: "Shop Review API"
  description: |
    # Features
    - Get reviews.
    - Post review.
  termsOfService: "https://example.com/terms/"
  contact:
    name: "Customer Support"
    url: "https://example.com/support/"
    email: "support-desk@example.com"
  license:
    name: "MIT License"
    url: "https://opensource.org/licenses/MIT"
  version: "1.0.0"

servers:
  - url: "http://localhost:{port}"
    description: "Local Development"
    variables:
      port:
        enum: ["3000", "8080"]
        default: "3000"
  - url: "https://it1.example.com/v1"
    description: "Integration Test"
  - url: "https://stg.example.com/v1"
    description: "Staging"
  - url: "https://api.example.com/v1"
    description: "Production"

tags:
  - name: "reviews"
    description: "Shop review operation"

paths:
  "/shops/{shopId}/reviews":
    get:
      summary: "Get shop review"
      description: "Get specified shop's reviews"
      tags: ["reviews"]
      parameters:
        - $ref: "#/components/parameters/ShopId"
      responses:
        "200":
          description: "Success operation"
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Review"
    post:
      summary: "Create new review"
      description: "Create specified shop's new review"
      tags: ["reviews"]
      parameters:
        - $ref: "#/components/parameters/ShopId"
      requestBody:
        description: "Contents of review"
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Review"
      responses:
        "201":
          description: "Success operation"
          headers:
            location:
              description: "New URL of created review"
              schema: { type: string, format: url }
        "400":
          $ref: "#/components/responses/400-BadRequest"

components:
  schemas:
    Review:
      type: object
      properties:
        score: { type: integer }
        comment: { type: string }
    ClientError:
      type: object
      properties:
        code: { type: string }
        type: { type: string }
        message: { type: string }
        errors:
          type: array
          items:
            type: object
            properties:
              field: { type: string }
              code: { type: string }
  parameters:
    ShopId:
      name: shopId
      in: path
      description: "Shop identifier"
      required: true
      schema: { type: string }
  responses:
    400-BadRequest:
      description: "Client side error"
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ClientError"
```

#### Security

| 種別          | 詳細   | 説明                  |
| ------------- | ------ | --------------------- |
| http          | Basic  | Basic認証             |
|               | Bearer | JWTなどを利用した認可 |
| apiKey        | header | API Keyを利用した認可 |
|               | cookie | ログインセッション    |
| oauth2        | -      | Oauth2.0              |
| openIdConnect | -      | Open ID Connect       |

SecuritySchemesの共通要素

```yaml
components:
  securitySchemes:
    sample_auth:
      description: "Sample authentication"
      type: http
```

Basic

```yaml
components:
  securitySchemes:
    sample_auth:
      description: "Sample authentication"
      type: http
      scheme: basic
```

JWT

```yaml
components:
  securitySchemes:
    sample_jwt_auth:
      description: "JWT authentication"
      type: http
      scheme: bearer
      bearerFormat: JWT
```

API-Key(header)

```yaml
components:
  securitySchemes:
    sample_apiKey_auth:
      description: "API-Key authentication"
      type: apiKey
      in: header
      name: X-Api-Key
```

API-Key(cookie)

```yaml
components:
  securitySchemes:
    sample_cookie_auth:
      description: "Login Session authentication"
      type: apiKey
      in: cookie
      name: JSESSIONID
```

Oauth2.0

```yaml
components:
  securitySchemes:
    sample_oauth2_auth:
      description: "Oauth2 Authorization"
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: "https://oauth.sample.com/auth"
          tokenUrl: "https://oauth.sample.com/token"
          scopes:
            "create_review": "Post new revieww."
```

Open ID Connect

```yaml
components:
  securitySchemes:
    sample_oidc_auth:
      description: "Open ID Connect"
      type: openIdConnect
      openIdConnectUrl: "https://oidc.sample.com/signin"
```

個別APIに適用

```yaml
paths:
  "/samples":
    get:
      summary: "Get all sample data"
      responses:
        "200":
          description: "Success operation"
      security:
        - sample_jwt_auth: [] #Scope必要な場合記述する
components:
  securitySchemes:
    sample_jwt_auth:
      description: "JWT authentication"
      type: http
      scheme: bearer
      bearerFormat: JWT
```

全体に適用して個別に解除

```yaml
paths:
  "/samples":
    get:
      summary: "Get all sample data"
      responses:
        "200":
          description: "Success operation"
      security: []
security:
  - sample_jwt_auth: [] #Scope必要な場合記述する
components:
  securitySchemes:
    sample_jwt_auth:
      description: "JWT authentication"
      type: http
      scheme: bearer
      bearerFormat: JWT
```

個別APIに適用

```yaml
openapi: "3.0.0"

info:
  title: "Shop Review API"
  description: |
    # Features
    - Get reviews.
    - Post review.
  termsOfService: "https://example.com/terms/"
  contact:
    name: "Customer Support"
    url: "https://example.com/support/"
    email: "support-desk@example.com"
  license:
    name: "MIT License"
    url: "https://opensource.org/licenses/MIT"
  version: "1.0.0"

servers:
  - url: "http://localhost:{port}"
    description: "Local Development"
    variables:
      port:
        enum: ["3000", "8080"]
        default: "3000"
  - url: "https://it1.example.com/v1"
    description: "Integration Test"
  - url: "https://stg.example.com/v1"
    description: "Staging"
  - url: "https://api.example.com/v1"
    description: "Production"

tags:
  - name: "reviews"
    description: "Shop review operation"

paths:
  "/shops/{shopId}/reviews":
    get:
      summary: "Get shop review"
      description: "Get specified shop's reviews"
      tags: ["reviews"]
      parameters:
        - name: shopId
          in: path
          required: true
          schema: { type: string }
      responses:
        "200":
          description: "Success operation"
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Review"
    post:
      summary: "Create new review"
      description: "Create specified shop's new review"
      tags: ["reviews"]
      parameters:
        - name: shopId
          in: path
          description: "Shop identifier"
          required: true
          schema: { type: string }
      requestBody:
        description: "Contents of review"
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Review"
      responses:
        "201":
          description: "Success operation"
          headers:
            location:
              description: "New URL of created review"
              schema: { type: string, format: url }
        "400":
          $ref: "#/components/responses/400-BadRequest"
      security:
        - apikey_auth: []

components:
  schemas:
    Review:
      type: object
      properties:
        score: { type: integer }
        comment: { type: string }
    Error:
      type: object
      properties:
        code: { type: string }
        type: { type: string }
        message: { type: string }
        errors:
          type: array
          items:
            type: object
            properties:
              field: { type: string }
              code: { type: string }
  responses:
    400-BadRequest:
      description: "Client side error"
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/Error"
  securitySchemes:
    apikey_auth:
      description: "API Key authorization"
      type: apiKey
      name: "X-Api-Key"
      in: header
```

全体に適用して個別に解除

```yaml
openapi: "3.0.0"

info:
  title: "Shop Review API"
  description: |
    # Features
    - Get reviews.
    - Post review.
  termsOfService: "https://example.com/terms/"
  contact:
    name: "Customer Support"
    url: "https://example.com/support/"
    email: "support-desk@example.com"
  license:
    name: "MIT License"
    url: "https://opensource.org/licenses/MIT"
  version: "1.0.0"

servers:
  - url: "http://localhost:{port}"
    description: "Local Development"
    variables:
      port:
        enum: ["3000", "8080"]
        default: "3000"
  - url: "https://it1.example.com/v1"
    description: "Integration Test"
  - url: "https://stg.example.com/v1"
    description: "Staging"
  - url: "https://api.example.com/v1"
    description: "Production"

tags:
  - name: "reviews"
    description: "Shop review operation"

paths:
  "/shops/{shopId}/reviews":
    get:
      summary: "Get shop review"
      description: "Get specified shop's reviews"
      tags: ["reviews"]
      parameters:
        - name: shopId
          in: path
          required: true
          schema: { type: string }
      responses:
        "200":
          description: "Success operation"
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Review"
      security: []
    post:
      summary: "Create new review"
      description: "Create specified shop's new review"
      tags: ["reviews"]
      parameters:
        - name: shopId
          in: path
          description: "Shop identifier"
          required: true
          schema: { type: string }
      requestBody:
        description: "Contents of review"
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/Review"
      responses:
        "201":
          description: "Success operation"
          headers:
            location:
              description: "New URL of created review"
              schema: { type: string, format: url }
        "400":
          $ref: "#/components/responses/400-BadRequest"
security:
  - apikey_auth: []
components:
  schemas:
    Review:
      type: object
      properties:
        score: { type: integer }
        comment: { type: string }
    Error:
      type: object
      properties:
        code: { type: string }
        type: { type: string }
        message: { type: string }
        errors:
          type: array
          items:
            type: object
            properties:
              field: { type: string }
              code: { type: string }
  responses:
    400-BadRequest:
      description: "Client side error"
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/Error"
  securitySchemes:
    apikey_auth:
      description: "API Key authorization"
      type: apiKey
      name: "X-Api-Key"
      in: header
```
