---
date: 2026-09-22
---
[OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)

OP(OIDC Provider)  
RP(Relying Party)

```mermaid
sequenceDiagram
  autonumber
  App-->>RP(Relying Party): request protected resource
  RP(Relying Party)-->>App: send authorization code request
  App-->>OP(OIDC Provider): redirect authorization code request
  OP(OIDC Provider)-->>App: generate and send authorization code response
  App-->>RP(Relying Party): send authorization code
  RP(Relying Party)-->>OP(OIDC Provider): send access token request with authorization code and clientid & secret
  OP(OIDC Provider)-->>OP(OIDC Provider): validate authorization code, clientid, secret
  OP(OIDC Provider)-->>RP(Relying Party): generate and send access token response
  Note over OP(OIDC Provider): Token Endpoint
  RP(Relying Party)-->>OP(OIDC Provider): request userinfo
  OP(OIDC Provider)-->>OP(OIDC Provider): validate access token
  OP(OIDC Provider)-->>RP(Relying Party): response userinfo
  Note over OP(OIDC Provider): UserInfo Endpoint
```

OAuth2.0 より response_type の仕様を拡張した

- id_token 追加
- none 追加
- code, token, id_token の任意の組み合わせを許可
- scope パラメータに openid 含められるように

### response_type=code (scope に openid を含まない)

OAuth2.0 認可コードフローと同じ

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate authorization code
    Note over Authorization Server: Authorization Endpoint
    App-->>Authorization Server: show authorization code
    Note over Authorization Server: Token Endpoint
    Authorization Server-->>App: generate access token
```

### response_type=code (scope に openid を含む)

ID Token と Access Token が一緒に発行される

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate authorization code
    Note over Authorization Server: Authorization Endpoint
    App-->>Authorization Server: show authorization code
    Note over Authorization Server: Token Endpoint
    Authorization Server-->>App: generate access token and id token
```

### response_type=token (scope に openid を含まない)

OAuth2.0 インプリシットフローと同じ

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate access token
    Note over Authorization Server: Authorization Endpoint
```

### response_type=id_token (scope に openid を含む)

ID Token のみ発行される

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate id token
    Note over Authorization Server: Authorization Endpoint
```

### response_type=id_token token (scope に openid を含む)

ID Token と Acess Token が一緒に発行される

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate id token and access token
    Note over Authorization Server: Authorization Endpoint
```

### response_type=code id_token (scope に openid を含む)

一時的に authorization_code が発行される  
ID Token と Access Token が発行される
ID Token が 2 回発行されるが、同じでもいい、異なってもいい

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate authorization code and id token
    Note over Authorization Server: Authorization Endpoint
    App-->>Authorization Server: show authorization code
    Note over Authorization Server: Token Endpoint
    Authorization Server-->>App: generate access token
```

### response_type=code token (scope に openid を含まない)

一時的に authorization_code が発行される  
Access Token が 2 回発行されるが、同じでもいい、異なってもいい

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate authorization code and access token
    Note over Authorization Server: Authorization Endpoint
    App-->>Authorization Server: show authorization code
    Note over Authorization Server: Token Endpoint
    Authorization Server-->>App: generate access token
```

### response_type=code token (scope に openid を含む)

一時的に authorization_code が発行される  
ID Token と Access Token が発行される
ID Token が 2 回発行されるが、同じでもいい、異なってもいい

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate authorization code and access token
    Note over Authorization Server: Authorization Endpoint
    App-->>Authorization Server: show authorization code
    Note over Authorization Server: Token Endpoint
    Authorization Server-->>App: generate id token and access token
```

### response_type=code id_token token (scope に openid を含む)

一時的に authorization_code が発行される  
ID Token と Access Token が 2 回発行されるが、同じでもいい、異なってもいい

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate authorization code and access token and id token
    Note over Authorization Server: Authorization Endpoint
    App-->>Authorization Server: show authorization code
    Note over Authorization Server: Token Endpoint
    Authorization Server-->>App: generate id token and access token
```

response_type=none (scope に openid を含まない)

何も生成されない、ログインのみチェックするなどに利用

```mermaid
sequenceDiagram
    autonumber
    App-->>Authorization Server: send authorization request
    Authorization Server-->>App: generate none
    Note over Authorization Server: Authorization Endpoint
```
