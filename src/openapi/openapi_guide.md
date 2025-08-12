# OpenAPI 仕様書（解説）

このドキュメントは `src/openapi/openapi.yaml` の内容を初心者向けにわかりやすく解説したものです。

---

## OpenAPI とは？

Web API の仕様（エンドポイント、リクエスト・レスポンスの型、認証など）を機械可読な YAML/JSON で記述するための標準フォーマットです。

---

### OpenAPI→Go コードの自動生成手順

```
# openapi.yaml から Go サーバースタブを生成
docker run --rm -v "$PWD":/work -w /work \
  hidori/oapi-codegen \
  -generate types,server -package api \
  -o backend/internal/api/openapi.gen.go \
  ./openapi/openapi.yaml

```

---

## 基本情報

- **title**: Gym Utilization & Maintenance API (MVP, No-Reservation)
- **version**: 0.1.0
- **description**: 小規模ジム向けの利用状況・設備管理 API。予約機能なしの MVP。

---

## サーバー情報

- `servers` で API のベース URL を定義します。
  - 本番: `https://api.example.com/v1`
  - ローカル: `http://localhost:8080/v1`

---

## 認証

- `securitySchemes` で `bearerAuth`（JWT トークン認証）を定義。
- すべての API で `Authorization: Bearer <token>` ヘッダが必要です。

---

## パラメータ

- `X-Gym-ID`（ヘッダ）: 利用中のジム ID。全 API で必要。
- ページング用: `limit`, `cursor`
- 日時フィルタ: `from`, `to`

---

## 主なデータ型（schemas）

- **User**: 会員ユーザー情報
- **Gym**: ジム情報
- **Zone**: ジム内のエリア
- **ResourceGroup**: 設備グループ（例: パワーラック）
- **Resource**: 個別設備（例: ラック A）
- **Member**: 会員
- **UtilizationSession**: 利用セッション（誰がいつどの設備を使ったか）
- **DashboardNow**: 現在の利用状況ダッシュボード
- **AuditLog**: 操作履歴
- **Error**: エラー応答

---

## 主な API エンドポイント

### 認証

- `POST /auth/login` : ログイン（JWT 発行）
- `POST /auth/refresh` : トークン再発行

### ジム・会員

- `GET /gyms/me` : 現在のジム情報取得
- `GET /members` : 会員一覧
- `POST /members` : 会員登録

### ゾーン・設備

- `GET /zones` : ゾーン一覧
- `POST /zones` : ゾーン作成
- `PATCH /zones/{zone_id}` : ゾーン更新
- `GET /resource-groups` : 設備グループ一覧
- `POST /resource-groups` : 設備グループ作成
- `PATCH /resource-groups/{group_id}` : 設備グループ更新
- `GET /resources` : 設備一覧
- `POST /resources` : 設備作成
- `PATCH /resources/{resource_id}` : 設備更新

### 利用セッション

- `GET /sessions` : 利用履歴一覧
- `POST /sessions` : 利用開始
- `PATCH /sessions/{session_id}/end` : 利用終了

### ダッシュボード・監査

- `GET /dashboard/now` : 現在の利用状況
- `GET /audit-logs` : 操作履歴

---

## 使い方の流れ（例）

1. ログインして JWT トークンを取得
2. `X-Gym-ID` ヘッダを付けて API を呼び出す
3. 会員・設備・利用セッションなどを管理

---

## 生成コードと実装コードの対応関係

```
backend/internal/api/openapi.gen.go   # 自動生成（編集禁止）
backend/cmd/server/handlers_impl.go   # 自分で実装するハンドラ
backend/cmd/server/handlers_stub.go   # 未実装部分のダミー
```

- openapi.gen.go にある ServerInterface をすべて実装する
- 実装は handlers_impl.go に書き、handlers_stub.go は最小限に

---

## よく使う用語

- **エンドポイント**: API の URL ごとの機能
- **スキーマ**: データ型の定義
- **パラメータ**: API 呼び出し時に指定する値
- **レスポンス**: API から返るデータ
- **JWT**: 認証用トークン

---

## 参考

- OpenAPI 公式: https://swagger.io/specification/
- Redoc: https://redocly.com/

---

このドキュメントは `openapi.yaml` の内容を初心者向けに要約・解説したものです。詳細は YAML 本体や自動生成ドキュメントも参照してください。
