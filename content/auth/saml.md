---
date: 2026-09-22
---
[OASIS-OPEN](https://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0.html),
[RFC7522](https://datatracker.ietf.org/doc/html/rfc7522)

```mermaid
sequenceDiagram
    autonumber
    Client-->>SP(Service Provider): access service
    SP(Service Provider)-->>Client: redirect sso url to idp
    Client-->>IdP(Identity Provider): client relays saml request to idp
    IdP(Identity Provider)-->>Client: generates saml assertion and send to client
    Client-->>SP(Service Provider): relays saml assertion to sp
    SP(Service Provider)-->>SP(Service Provider): verify saml assertion
    SP(Service Provider)-->>Client: if authenticated, sends securitycontext to client
    Client-->>SP(Service Provider): access protected resource
    SP(Service Provider)-->>Client: response with requested resource
```

SP(Service Provider): サービス提供エンティティー、通常はアプリケーション  
IdP(Identity Provider): 認証サービスを提供するエンティティー  
SP-Initiated(Service Provider Initiated): サービスプロバイダーから開始する SAML サインインフロー  
IdP-Initiated(Identity Provider Initiated): IdP から開始する SAML サインインフロー
RelayState: HTTP パラメータ、SAML リクエストとレスポンスに含まれる