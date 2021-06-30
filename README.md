# proxy-server

new proxy server

## Content

- [proxy-server](#proxy-server)
  - [Content](#content)
  - [Setting](#setting)
    - [Config(config.yml)](#configconfigyml)
    - [Authorization(auth.yml)](#authorizationauthyml)
    - [LoadBalancer(loadbalancer.yml)](#loadbalancerloadbalanceryml)
    - [Menu(menu.yml)](#menumenuyml)
    - [OpenID Connect(oidc.yml)](#openid-connectoidcyml)

## Setting

### Config(config.yml)

| Key             |   Type   | Overview                            |
| :-------------- | :------: | :---------------------------------- |
| redis_endpoints | string[] | セッション格納用redisエンドポイント |
| redis_username  |  string  | ユーザー名                          |
| redis_password  |  string  | パスワード                          |
| cert_file       |  string  | 証明書                              |
| key_file        |  string  | 証明書の秘密鍵                      |
| port            |   int    | ポート番号                          |
| log_level       |  string  | ログ出力レベル                      |

### Authorization(auth.yml)

| Key      |  Type  | Overview                         |
| :------- | :----: | :------------------------------- |
| path     | string | グループエンドポイント           |
| login    | string | ログインエンドポイント           |
| callback | string | コールバックエンドポイント       |
| logout   | string | ログアウトエンドポイント         |
| userinfo | string | ユーザー情報取得用エンドポイント |

### LoadBalancer(loadbalancer.yml)

| Key                  |    Type    | Overview                                                   |
| :------------------- | :--------: | :--------------------------------------------------------- |
| domain               |   string   | 待ち受けドメイン名                                         |
| locations            | Location[] | ロケーション情報                                           |
| - path               |   string   | URLパス                                                    |
| - token_type         |   string   | バックエンドヘ流すトークンタイプ(id_token or access_token) |
| - upstream           |  string[]  | バックエンドサービス名(Upstream[].nameと一致すること)      |
| - rewrite            |    map     | 書き換えるURLのマップ(Key: string, Value: string)          |
| - Plugins            |  Plugins   | このLocationで使用するプラグイン                           |
| -- request_transport |  string[]  | リクエスト時に実行されるプラグイン                         |
| -- response_modify   |  string[]  | レスポンス時に実行されるプラグイン                         |
| Upstream             | Upstream[] |                                                            |
| - name               |   string   | バックエンドサービス名                                     |
| - url                |   string   | バックエンドサービスのエンドポイント                       |
| - weight             |   float    | アクセス比重                                               |

### Menu(menu.yml)

| Key         |  Type  | Overview       |
| :---------- | :----: | :------------- |
| name        | string | メニュー名     |
| description | string | 説明           |
| path        | string | エンドポイント |
| thumbnail   | string | サムネイルURL  |

### OpenID Connect(oidc.yml)

| Key           |   Type   | Overview                     |
| :------------ | :------: | :--------------------------- |
| scopes        | string[] | スコープ                     |
| provider      |  string  | IDP情報                      |
| client_id     |  string  | クライアントID               |
| client_secret |  string  | クライアントシークレット情報 |
| callback_url  |  string  | コールバックURL              |
| logout        |  string  | IDPログアウトURL             |
| audiences     | string[] | Authorization Audience       |

<!-- ### Certificate Authority(ca.yml)

| Key                 |   Type   | Overview             |
| :------------------ | :------: | :------------------- |
| country             | string[] | 国                   |
| province            | string[] | 都道府県             |
| locality            | string[] | 地域                 |
| organizational_unit | string[] | 部門名               |
| organization        | string[] | 組織・団体名         |
| common_name         |  string  | ドメインコモンネーム |
| years               |   int    | 証明書の発行期限     |
| Serial              |   int    | シリアル番号         |
| refresh             |   bool   | 再作成               | --> |
