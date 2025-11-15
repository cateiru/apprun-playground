# ビルドステージ
FROM golang:1.25-alpine AS builder

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
