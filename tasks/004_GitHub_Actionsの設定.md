# 004: GitHub Actions の設定

## 概要

GitHub Actions を使用して、main ブランチへの push 時に自動的に Docker イメージをビルドし、GitHub Container Registry (ghcr.io) に push するワークフローを作成します。

## 実装内容

### 1. ワークフローディレクトリの作成

```bash
mkdir -p .github/workflows
```

### 2. ワークフローファイルの作成

`.github/workflows/docker-publish.yml` を作成します：

```yaml
name: Docker Image CI

on:
  push:
    branches:
      - main

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-
            type=raw,value=latest,enable={{is_default_branch}}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

## ディレクトリ構造

```
apprun-playground/
├── .github/
│   └── workflows/
│       └── docker-publish.yml
├── Dockerfile
├── main.go
├── go.mod
└── go.sum
```

## チェックリスト

- [ ] `.github/workflows/` ディレクトリを作成
- [ ] `docker-publish.yml` ワークフローファイルを作成
- [ ] GitHub リポジトリにコミット & プッシュ
- [ ] GitHub Actions のワークフローが実行されることを確認
- [ ] ghcr.io にイメージが push されることを確認
- [ ] GitHub リポジトリの Packages ページでイメージを確認

## ワークフローの説明

### トリガー

- `main` ブランチへの push 時に自動実行

### 環境変数

- `REGISTRY`: ghcr.io（GitHub Container Registry）
- `IMAGE_NAME`: リポジトリ名（自動取得）

### ステップ

1. **Checkout repository**: コードをチェックアウト
2. **Set up Docker Buildx**: Docker Buildx のセットアップ
3. **Log in to Container Registry**: ghcr.io にログイン（`GITHUB_TOKEN` を使用）
4. **Extract metadata**: イメージのタグとラベルを生成
5. **Build and push Docker image**: イメージをビルドして push

### イメージタグ

- `main`: main ブランチの最新
- `main-<commit-sha>`: 特定のコミット
- `latest`: デフォルトブランチの最新

### キャッシュ

- GitHub Actions のキャッシュ機能を使用してビルドを高速化

## 補足

### GitHub Container Registry の公開設定

イメージを public にする場合：

1. GitHub リポジトリの **Packages** ページに移動
2. 該当イメージを選択
3. **Package settings** で **Change visibility** をクリック
4. **Public** に変更

### イメージの Pull

```bash
docker pull ghcr.io/cateiru/apprun-playground:latest
```

### Secrets の設定

このワークフローでは `GITHUB_TOKEN` を使用しますが、これは GitHub Actions で自動的に提供されるため、手動での設定は不要です。

## トラブルシューティング

### ワークフローが失敗する場合

- GitHub リポジトリの **Actions** タブでログを確認
- `permissions` が正しく設定されているか確認
- Dockerfile が正しく配置されているか確認

### イメージが push されない場合

- GitHub リポジトリの **Settings** > **Actions** > **General** で、ワークフローの権限を確認
- `GITHUB_TOKEN` の権限が `packages: write` を含んでいるか確認

## 参考リンク

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Publishing Docker images](https://docs.github.com/en/actions/publishing-packages/publishing-docker-images)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [docker/build-push-action](https://github.com/docker/build-push-action)
