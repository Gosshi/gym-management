
# gym-management

ジムの利用状況・設備管理・会員管理を行うWebアプリケーションのサンプルプロジェクトです。

## 概要

- **バックエンド**: Go + Gin + OpenAPI (oapi-codegen)
- **フロントエンド**: Node.js/React（`web/` ディレクトリ、詳細は未実装または別途追加）
- **API設計**: OpenAPI 3.0（`src/openapi/openapi.yaml`）
- **認証**: Bearerトークン（JWT, 開発用はダミー）
- **DBマイグレーション**: `src/backend/migrations/`
- **CI/CD・Lint**: Spectral, Redocly, openapi-typescript など

## ディレクトリ構成

- `src/backend/`  
  - `cmd/server/` : サーバー起動・APIハンドラ（Gin）
  - `internal/api/` : oapi-codegenによるAPIインターフェース自動生成
  - `internal/handlers/` : ビジネスロジック（未実装/今後追加）
  - `internal/repo/` : DBアクセス層（未実装/今後追加）
  - `migrations/` : DBマイグレーションファイル
- `infra/` : インフラ構成（例: Docker, IaC等。現状空）
- `openapi/` : OpenAPI仕様書・ドキュメント
- `web/` : フロントエンド（React/Next.js等を想定）

## セットアップ

### バックエンド

```sh
cd src/backend
go mod tidy
go run ./cmd/server
```
- サーバーは `localhost:8080` で起動
- APIエンドポイントは `/v1/` 配下

### フロントエンド

```sh
cd web
npm install
npm run dev
```
- `localhost:3000` などで起動（フレームワークにより異なる）

### OpenAPIドキュメント

- `src/openapi/openapi.yaml` を編集
- ドキュメントHTML生成:  
  ```sh
  npm run api:docs
  ```
- 型定義生成（TypeScript用）:  
  ```sh
  npm run api:ts
  ```

## 主なAPIエンドポイント例

- `/v1/auth/login` : ログイン（JWT発行）
- `/v1/gyms/me` : 現在のジム情報取得
- `/v1/zones` : ゾーン一覧・作成・更新
- `/v1/resource-groups` : 設備グループ一覧・作成・更新
- `/v1/resources` : 設備（個別）一覧・作成・更新
- `/v1/members` : 会員一覧・作成
- `/v1/sessions` : 利用セッション一覧・開始・終了
- `/v1/dashboard/now` : 現在の利用状況ダッシュボード
- `/v1/audit-logs` : 監査ログ取得

## 開発Tips

- API仕様を変更したら `oapi-codegen` でGoコードを再生成してください
- GinのミドルウェアでOpenAPIバリデーション・認証を実装
- 認証は開発用ダミー（本番はJWT検証を実装）
- DBや永続化層は未実装。`internal/repo/` で拡張可能
- フロントエンドはAPIクライアント自動生成（openapi-typescript）を推奨

## 依存ツール

- Go 1.23+
- Node.js 18+
- oapi-codegen
- @redocly/cli, spectral, openapi-typescript など

## ライセンス

MIT
