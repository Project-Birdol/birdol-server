# birdol API Specification (暫定) （2021/9/21）

## ルーティング

- `/api/v2`
    - `/user` : AccessTokenを使用しない（発行以前の）ユーザ－関連処理
        - ここでは署名検証が入らない
    - `/auth` : AccessTokenを使用したユーザー関連処理
        - RSA暗号による署名検証が有効（詳細は以下）
    - `/gamedata` : ゲームデータの操作（ここ以下は未実装）
        - RSA暗号による署名検証が有効
        - セッションによる複数デバイスでの連携管理が有効

## リクエスト仕様

### アカウント作成

- **Endpoint**: `/api/v2/user`
- **Method**: `HTTP PUT`
- **Content-Type**: `application/json`
- **Header**: 追加パラメータなし
- **Body**:
    - `name` : ユーザー名
    - `public_key` : RSA公開鍵（base64 encoded XML）
    - `device_id` : クライアントで生成したUUID
    ```json
    # SAMPLE
    {
        "name": "<-- ユーザー名 -->",
        "public_key": "<-- RSA公開鍵(base64 encoded XML) -->",
        "device_id": "<-- Client's UUID -->"
    }
    ```
- **Response**: 
    - `result` : `success` or `failed`
    - `error` : エラーの詳細
    - `user_id` : ユーザーID(DB上の)
    - `account_id` : ユーザーID(自動生成)
    - `access_token` : アクセストークン
    - `refresh_token` : リフレッシュトークン
    ```json
    # SAMPLE
    {
        "result": "success",
        "user_id": 102,
        "account_id": "jkz#ckb2nsozxna&@21sf"
        "access_token" : "z9sfho*^$dck$jc@v"
        "refresh_token" : "3js9&bd%aszlx#hxo$"
    }
    ```
        
### アカウント連携
- **Endpoint**: `/api/v2/user`
- **Method**: `HTTP POST`
- **Content-Type**: `application/json`
- **Header**: 追加パラメータなし
- **Body**:
    - `account_id` : ユーザーID（自動生成のほう）
    - `password` : パスワード
    - `device_id` : クライアントで生成したUUID
    - `public_key` : RSA公開鍵（base64 encoded XML）
    ```json
    # SAMPLE
    {
        "account_id": "<-- ユーザーID -->",
        "password": "<-- パスワード -->",
        "device_id": "<-- Client's UUID -->",
        "public_key": "<-- RSA公開鍵(base64 encoded XML) -->"
    }
    ```
- **Response**:
    - `result` : `success` or `failed`
    - `error` : エラーの詳細
    - `user_id` : ユーザーID(DB上の)
    - `access_token` : アクセストークン
    - `refresh_token` : リフレッシュトークン
    ```json
    # SAMPLE
    {
        "result": "ok",
        "user_id": 102,
        "access_token" : "z9sfho*^$dck$jc@v"
        "refresh_token" : "3js9&bd%aszlx#hxo$"
    }
    ```
        
### ログイン処理
- **Endpoint**: `/api/v2/auth`
- **Method**: `HTTP GET`
- **Content-Type**: `None`
- **Header**: 独自パラメータあり
    - `Authorization` : `Bearer <AccessToken>` 
    - `X-Birdol-Signature` : 署名（詳細は後述）
    - `X-Birdol-TimeStamp` : タイムスタンプ 署名で使用
    - `device_id` : 登録したUUID
- **Body**: NULL
- **Response**: 
    - `result` : `success` or `failed`
    - `error` : エラーの詳細
    - `session_id` : セッションID

### ログアウト（リンク解除）処理
- **Endpoint**: `/api/v2/auth`
- **Method**: `HTTP DELETE`
- **Content-Type**: `None`
- **Header**: 独自パラメータあり
    - `Authorization` : `Bearer <AccessToken>`
    - `X-Birdol-Signature` : 署名（詳細は後述）
    - `X-Birdol-TimeStamp` : タイムスタンプ 署名で使用
    - `device_id` : 登録したUUID
- **Body**: NULL
- **Response**: 
    - `result` : `success` or `failed`
    - `error` : エラーの詳細

### データリンク設定（パスワード登録）
- **Endpoint**: `/api/v2/auth`
- **Method**: `HTTP PUT`
- **Content-Type**: `application/json`
- **Header**: 独自パラメータあり
    - `Authorization` : `Bearer <AccessToken>`
    - `X-Birdol-Signature` : 署名（詳細は後述）
    - `X-Birdol-TimeStamp` : タイムスタンプ 署名で使用
    - `device_id` : 登録したUUID
- **Body**:
    - `password` : パスワード
- **Response**:
    - `result` : `success` or `failed`
    - `error` : エラーの詳細
    - `expire_date` : パスワードの有効期限

### トークンリフレッシュ
- **Endpoint**: `/api/v2/auth/refresh?refresh_token=xxxxxxxx`
- **Method**: `HTTP GET`
- **Content-Type**: `None`
- **Header**: 独自パラメータあり
    - `Authorization` : `Bearer <AccessToken>`
    - `X-Birdol-Signature` : 署名（詳細は後述）
    - `X-Birdol-TimeStamp` : タイムスタンプ 署名で使用
    - `device_id` : 登録したUUID
- **Body**: NULL
- **Response**:
    - `result` : `success` or `failed`
    - `error` : エラーの詳細
    - `token` : 新しいアクセストークン
    - `refresh_token` : 新しいリフレッシュトークン
    - `session_id` : 新たなセッションID

### Template
- **Endpoint**: 
- **Method**: 
- **Content-Type**: 
- **Header**: 
- **Body**:
- **Response**:

## 署名について
AccessTokenを使用するすべてのリクエストは，端末登録時に生成したRSA鍵により署名を付与する．署名に関係するパラメータはすべてリクエストヘッダに入るので，[リクエスト仕様](#リクエスト仕様)を参照

### クライアント側の処理

1. デバイス登録時（アカウント生成 or アカウント連携時）にクライアント側でRSAキーペアを生成し，公開鍵部分をリクエストにくっつけて送る
    - C#の`RSACryptoServiceProvider`の使用を想定（というかこれを使ってください）
        - 鍵はXMLでエクスポートできるので，XML化したら気休め程度ですが`base64`でエンコードして扱う
        - 秘密鍵はbase64 encoded XMLのままPlayerPrefsなりに入れるか，使えるなら`RSACryptoServiceProvider`のキーコンテナを使いたい
            - 誰かUnityで使えるか検証してもらえると･･･

        - アカウント生成 or アカウント連携時に生成して，ログアウト（連携解除　要はユーザーデータの紐付け解除　データリセット）しない限りは同一の鍵を使用
        - 公開鍵はbase64 encoded XMLでリクエストにくっつけて送りつける
3. リクエストを送る際，ベース文字列 `｛v1|v2}:<タイムスタンプ>:<リクエストボディ>` を形成(v1かv2かはAPIバージョン)
    - タイムスタンプは`YYYY-MM-DD-hh-mm-ss`を想定
    - リクエストボディがない時（GETなど）は空文字列扱いとする
        - [x] 【優先事項】Server側の，空のBody読み出しの際の挙動が未検証

2. ベース文字列を`SHA-256`でハッシュ化
3. ハッシュ化して得られたバイト列を秘密鍵で暗号化　→　これが署名

### サーバー側の処理（実装済み・挙動未確認）
1. デバイス登録時（アカウント生成 or アカウント連携時）のリクエストで送りつけられた公開鍵（base64 encoded XML）をDBに保管しておく（正直保存方法がガバガバな気がするけど，公開鍵暗号ならまだ許されるだろう･･･との判断）
2. リクエストが来たら，主処理の前に検証が入る
3. まずベース文字列 `v2:<タイムスタンプ>:<リクエストボディ>` を形成
    - タイムスタンプはリクエストヘッダにあり，リクエストボディは読み込んで文字列化
4. ベース文字列を`SHA-256`でハッシュ化
5. RSA PKCS#1 v1.5 の署名を検証
    - [x] C#側での署名との整合性がとれているかが未検証

## セッション管理について
おそらくタイトル画面からホーム画面に遷移する段階でだと思うが，アクセストークンを使用して簡易的な[ログイン処理](#ログイン処理)のようなものをしてもらう．（クライアント側でGETリクエストを投げるだけ）

このとき，レスポンスでセッションIDが返ってくるので，ログイン後のAPIアクセスの際はbodyかheaderに必ずこれをつける
- ボディのないリクエストはクエリパラメータに `?session_id=xxxxxx` を入れる
- ボディがあるリクエストはボディに `session_id` をいれる