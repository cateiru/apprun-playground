# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

**apprun-playground** は、さくらインターネットの AppRun プラットフォームへのデプロイを検証するための実験的なリポジトリです。Go + Echo フレームワークを使用したシンプルなマイクロサービスを題材として、Docker コンテナ化、GitHub Actions による CI/CD、ghcr.io へのイメージプッシュ、AppRun へのデプロイという一連のワークフローを構築することを目的としています。

## 技術スタック

- **言語**: Go 1.25.1
- **Web フレームワーク**: Echo v4.13.4
- **テストフレームワーク**: testify/assert v1.10.0
- **コンテナ**: Docker マルチステージビルド (alpine ベース)
- **CI/CD**: GitHub Actions
- **コンテナレジストリ**: GitHub Container Registry (ghcr.io)

## アーキテクチャ

### アプリケーション構造

- **main.go**: Echo サーバーの実装
  - `GET /` エンドポイント: ヘルスチェック用 (レスポンス: `{"status":"ok"}`)
  - ミドルウェア: Logger, Recover, CORS を使用
  - ポート番号: 環境変数 `PORT` で設定可能 (デフォルト: 8080)

- **main_test.go**: テストスイート
  - ユニットテスト: `handleHealth` 関数の動作確認
  - 統合テスト: エンドポイント全体の動作確認
  - ミドルウェアテスト: ミドルウェアチェーンの動作確認

### コンテナ戦略

マルチステージビルドにより、最小限のランタイムイメージを生成します：

- **ビルドステージ**: Go 1.23-alpine でバイナリをコンパイル (CGO 無効化)
- **ランタイムステージ**: alpine:latest で実行
- **セキュリティ**: 非 root ユーザー (UID 1000) での実行
- **最適化**: ビルドツールをランタイムに含めず、イメージサイズを最小化

### CI/CD パイプライン

GitHub Actions (`docker-publish.yml`) により、main ブランチへの push 時に自動的に：

1. Docker イメージをビルド
2. ghcr.io にプッシュ (タグ: `main`, `main-<commit-sha>`, `latest`)
3. GitHub Actions キャッシュでビルド時間を最適化

**セキュリティ**: すべての GitHub Actions は SHA ハッシュでピン留めされています（セキュリティベストプラクティス）

## 開発コマンド

### ローカル開発

```bash
# サーバーの実行
go run main.go

# 依存関係の整理
go mod tidy

# ヘルスチェック
curl http://localhost:8080/
```

### テスト

```bash
# 全テストを実行
go test -v ./...

# 単一のテストを実行
go test -v -run TestHandleHealth ./...

# カバレッジを確認
go test -cover ./...

# カバレッジレポートを生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Docker

```bash
# イメージのビルド
docker build -t apprun-playground:latest .

# コンテナの実行
docker run -p 8080:8080 apprun-playground:latest

# 環境変数でポート番号を変更
docker run -p 9000:9000 -e PORT=9000 apprun-playground:latest
```

## 重要な設計上の注意点

### ポート番号の柔軟性

`main.go` は環境変数 `PORT` を読み込み、指定がない場合は 8080 をデフォルトとします。これにより、AppRun のような外部プラットフォームでのポート割り当てに柔軟に対応できます。

### テスト戦略

`main_test.go` では単なるユニットテストだけでなく、ミドルウェアチェーン全体の動作確認も行います。新しい機能を追加する際は、同様に多層的なテストを実装してください。

### コンテナセキュリティ

Docker コンテナは非 root ユーザー (appuser, UID 1000) で実行されます。セキュリティベストプラクティスに従い、この設計を維持してください。

## プロジェクト進捗状況

### 完了済み

- ✅ Go モジュール初期化
- ✅ Echo サーバー実装 (main.go)
- ✅ テストスイート実装 (main_test.go)
- ✅ .gitignore 設定
- ✅ Dockerfile と .dockerignore の作成（マルチステージビルド実装）
- ✅ GitHub Actions ワークフロー設定（SHA pinning 実装済み）

### 計画段階 (tasks/ ディレクトリ参照)

- ⏳ AppRun デプロイメント手順 (タスク 005)

新しい機能を追加する際は、tasks/ ディレクトリ内のタスクファイルを参照してください。
