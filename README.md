# Cloud Tasks Go - MAX_CONCURRENT_DISPATCHES 検証

Cloud TasksのMAX_CONCURRENT_DISPATCHESが期待通りに動作するかを検証するためのGoサーバー。

## 構成

- **Goサーバー**: 意図的にsleepを入れたHTTPハンドラー（`/task`エンドポイント）
- **Docker Compose**: ローカル開発環境のセットアップ
- **Cloud Tasks クライアント**: タスク作成用のGoクライアント

## セットアップ

1. 依存関係のダウンロード:
```bash
go mod tidy
```

2. Docker Composeで環境起動:
```bash
docker-compose up --build
```

3. GCPプロジェクト設定:
   - `client.go`でプロジェクトID、ロケーション、キュー名を設定
   - 実際のサーバーURLを設定

## 使用方法

### タスク作成（クライアント）

```bash
# 10個のタスクを作成
go run cmd/client/main.go 10

# 5個のタスクを作成
go run cmd/client/main.go 5
```

### サーバーログの確認

```bash
docker-compose logs -f app
```

## 検証のポイント

1. **同時実行数制限**: MAX_CONCURRENT_DISPATCHESで設定した数を超えて同時実行されないか
2. **順次処理**: 設定した最大数まで並列実行し、完了次第新しいタスクが開始されるか
3. **タスク完了**: 全てのタスクが確実に処理されるか

## エンドポイント

- `GET /health`: ヘルスチェック
- `POST /task?sleep=N`: タスク処理（N秒間sleep）
