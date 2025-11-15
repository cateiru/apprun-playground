# 003: Dockerfile の作成

## 概要

マルチステージビルドを使用した効率的な Dockerfile を作成します。ビルドステージでアプリケーションをコンパイルし、ランタイムステージで最小限のイメージを作成します。

## 実装内容

### 1. Dockerfile の作成

プロジェクトルートに `Dockerfile` を作成します：

```dockerfile
# ビルドステージ
FROM golang:1.23-alpine AS builder

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ランタイムステージ
FROM alpine:latest

# セキュリティのための非 root ユーザー作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 作業ディレクトリの設定
WORKDIR /app

# CA 証明書のインストール（HTTPS 通信に必要）
RUN apk --no-cache add ca-certificates

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 所有者を変更
RUN chown -R appuser:appuser /app

# 非 root ユーザーに切り替え
USER appuser

# ポートの公開
EXPOSE 8080

# アプリケーションの実行
CMD ["./main"]
```

### 2. .dockerignore の作成

不要なファイルを Docker イメージに含めないよう、`.dockerignore` ファイルを作成します：

```dockerignore
# Git
.git
.gitignore

# IDE
.vscode
.idea
*.swp
*.swo

# Documentation
README.md
LICENSE

# Tasks
tasks/

# Build artifacts
bin/
dist/

# Environment files
.env
.env.local

# Tests
*_test.go

# CI/CD
.github/
```

### 3. ローカルでの Docker ビルドと実行

#### ビルド

```bash
docker build -t apprun-playground:latest .
```

#### 実行

```bash
docker run -p 8080:8080 apprun-playground:latest
```

#### 動作確認

```bash
curl http://localhost:8080/
```

期待されるレスポンス：

```json
{"status":"ok"}
```

## ディレクトリ構造

```
apprun-playground/
├── Dockerfile
├── .dockerignore
├── main.go
├── go.mod
├── go.sum
└── README.md
```

## チェックリスト

- [ ] `Dockerfile` を作成
- [ ] `.dockerignore` を作成
- [ ] `docker build` でイメージをビルド
- [ ] `docker run` でコンテナを起動
- [ ] `curl` でエンドポイントをテスト
- [ ] レスポンスが `{"status":"ok"}` であることを確認
- [ ] イメージサイズを確認（`docker images` コマンド）

## 補足

### マルチステージビルドの利点

- **イメージサイズの削減**: ビルドツールを含まない最小限のランタイムイメージ
- **セキュリティ向上**: 不要なツールやソースコードを含まない
- **ビルドの高速化**: レイヤーキャッシュの最適化

### セキュリティのベストプラクティス

- 非 root ユーザーでアプリケーションを実行
- Alpine Linux ベースで最小限のイメージ
- CA 証明書を含めて HTTPS 通信に対応

### ポート番号

- デフォルトで 8080 ポートを公開
- `PORT` 環境変数で変更可能（main.go で実装済み）

## 参考リンク

- [Docker Multi-stage builds](https://docs.docker.com/build/building/multi-stage/)
- [Best practices for writing Dockerfiles](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Docker security best practices](https://docs.docker.com/develop/security-best-practices/)
