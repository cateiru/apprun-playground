# apprun-playground

さくらインターネットの AppRun にデプロイするための検証用リポジトリです。

## 概要

- Go と Echo フレームワークを使用した API サーバー
- `GET /` にアクセスすると `{"status": "ok"}` を返します
- GitHub Actions で自動的に ghcr.io に Docker イメージを push します

## 開発環境

- Go 1.23+
- Docker

## ローカル実行

```bash
go run main.go
```

## Docker ビルド

```bash
docker build -t apprun-playground .
docker run -p 8080:8080 apprun-playground
```
