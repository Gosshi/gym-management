## いま実装できているもの（Go 初心者向けまとめ）

### 全体像

- **目的**：ジムの混雑状況を“予約なし”で把握する API（MVP）。
- **やったこと**：OpenAPI 定義から**型とルーター**を自動生成し、**Gin**でサーバを起動。まずは**メモリ上の簡易 DB**で `resource-groups`（設備グループ）と `sessions`（利用開始/終了）を動かし、**ダッシュボード**で現在の稼働数が見えるようにした。

---

### 技術スタック（ザックリ）

- **Go 1.23**
- **Gin**（Web フレームワーク）
- **oapi-codegen（types, gin）**：OpenAPI から Go の型・ルーターを生成
- **kin-openapi バリデータ**：OpenAPI に沿って**リクエストを自動検証**
- **開発用 Auth**：`Authorization: Bearer ...` があれば通す簡易チェック
- **Base URL**：すべて **`/v1`** 配下

---

### リクエストが処理される流れ

1. クライアントが **`/v1/...`** にアクセス（ヘッダ `X-Gym-ID` 必須、開発中は `Authorization: Bearer DUMMY` 必須）。
2. **OpenAPI バリデータ**が、パラメータ形式や必須ヘッダをチェック。
3. **生成コード**が、該当エンドポイントの**Go ハンドラ**（あなたが書いた関数）にルーティング。
4. ハンドラは **メモリストア（構造体）** を読んだり書いたりして、JSON を返す。

---

### 主要ファイル（だいたいここだけ触れば OK）

- `openapi/openapi.yaml`：API 仕様の元。これを元にコード生成。
- `backend/internal/api/openapi.gen.go`：**生成コード**（型・ルーター定義など）。
- `backend/cmd/server/main.go`：サーバ起動、バリデータ設定、`/v1` グループ設定。
- `backend/cmd/server/store_mem.go`：**メモリ上の簡易 DB**（map とロック）。
- `backend/cmd/server/handlers_impl.go`：**実装済みハンドラ**（200 を返すやつ）。
- `backend/cmd/server/handlers_stub.go`：**未実装のスタブ**（501 を返すやつ）。

---

### メモリストア（簡易 DB）の中身

- `groups[gymID][groupID] = ResourceGroup`  
  設備グループ（例：パワーラック、キャパ 2 台）
- `sessions[gymID][sessionID] = UtilizationSession`  
  利用セッション（開始時間、終了時間、ステータス）

> **集計**：`activeCount(gymID, groupID)` でそのグループの**稼働中セッション数**をカウント。  
> **容量チェック**：`capacityOf(...)` を見て、超えていたら `409 capacity_exceeded`。

---

### 実装済みエンドポイント（MVP）

- `POST /v1/resource-groups`：グループ作成（名前・キャパ・任意ルール）
- `GET  /v1/resource-groups`：グループ一覧
- `PATCH /v1/resource-groups/{group_id}`：グループ更新（名前・キャパ等）
- `POST /v1/sessions`：利用開始（capacity 超過なら**409**）
- `PATCH /v1/sessions/{session_id}/end`：利用終了（終了時刻と経過秒を計算）
- `GET  /v1/sessions`：履歴一覧（from/to, group_id, member_id で簡易フィルタ）
- `GET  /v1/dashboard/now`：各グループの **現在の稼働数** と（満員なら）**簡易待ち時間**を返す

> それ以外（zones/members/resources/audit 等）は**スタブ（501 Not Implemented）**で待機中。

---

### 使い方（最低限の起動とテスト）

1. **起動**
   ```bash
   cd backend
   go run ./cmd/server
   ```
2. **ダッシュボード（最初は空）**
   ```
   curl -H "X-Gym-ID: gym_demo" -H "Authorization: Bearer DUMMY" http://localhost:8080/v1/dashboard/now
   ```
3. **グループ作成 → 確認**

   ```
   curl -X POST -H "X-Gym-ID: gym_demo" -H "Authorization: Bearer DUMMY" \
    -H "Content-Type: application/json" \
    -d '{"name":"Power Rack","capacity":2}' \
    http://localhost:8080/v1/resource-groups

   curl -H "X-Gym-ID: gym_demo" -H "Authorization: Bearer DUMMY" http://localhost:8080/v1/resource-groups
   ```

4. **セッション開始 → ダッシュボードで稼働数が増える**

   ```
   # <GROUP_ID> は上の応答で得た id
   curl -X POST -H "X-Gym-ID: gym_demo" -H "Authorization: Bearer DUMMY" \
   -H "Content-Type: application/json" \
   -d "{\"group_id\":\"<GROUP_ID>\"}" \
   http://localhost:8080/v1/sessions

   ```

5. **終了**
   ```
   # <SESSION_ID> は開始APIの応答で得た id
   curl -X PATCH -H "X-Gym-ID: gym_demo" -H "Authorization: Bearer DUMMY" \
   -H "Content-Type: application/json" -d '{}' \
   http://localhost:8080/v1/sessions/<SESSION_ID>/end
   ```

### ここが“Go らしい”ポイント

- **インターフェース満たす=実装完了**  
  生成された `ServerInterface` にある関数を**同じシグネチャ**（関数名・引数・戻り値すべて一致）で実装すれば、自動的にルーターに組み込まれる。
- **並行安全**  
  メモリストアは `sync.RWMutex` で**読み書きロック**しており、複数の HTTP リクエストが同時に来てもデータが壊れにくい。
- **型が守ってくれる**  
  OpenAPI から生成された**構造体・列挙型**をそのまま使うため、JSON の構造やフィールド名のミスがほぼ発生しない。IDE 補完も効く。

---

### 次にやると良いこと（順序の提案）

1. **待ち時間ロジック改善**
   - 直近の平均利用時間から推定する
   - 満員時に「あと ◯ 分で空きそう」と返す
2. **ユニットテスト追加**
   - handlers は最小限にして、ストアの関数を中心にテスト
   - capacity 超過や終了処理の挙動確認
3. **PostgreSQL に置き換え**
   - Docker Compose で DB を起動
   - GORM や sqlc で型安全なクエリ実装
4. **フロント雛形作成**
   - Vite + React + TypeScript
   - `openapi-typescript`で型を自動生成してフロントから叩く

---

### トラブル時のコツ

- **エラー全文**と**叩いたコマンド**をセットで記録  
  → 原因特定が早くなる
- `/v1` が付かずに 404 → **ルート登録**が `/v1` グループになっているか確認
- 認証エラー → `Authorization: Bearer DUMMY` を付け忘れていないか確認
- バリデーションエラー → `X-Gym-ID` と `Content-Type: application/json` を確認
