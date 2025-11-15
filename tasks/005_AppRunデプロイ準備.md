# 005: AppRun デプロイ準備

## 概要

さくらインターネットの AppRun にアプリケーションをデプロイするための準備と設定を行います。

## 前提条件

- [ ] さくらインターネットのアカウントを持っていること
- [ ] AppRun のサービスにアクセスできること
- [ ] ghcr.io に Docker イメージが push されていること

## AppRun とは

さくらインターネットが提供するコンテナベースのアプリケーション実行環境です。Docker イメージを指定するだけで、簡単にアプリケーションをデプロイできます。

## デプロイ手順

### 1. AppRun コンソールへのアクセス

1. さくらインターネットのコントロールパネルにログイン
2. AppRun のサービスページに移動

### 2. アプリケーションの作成

1. **新しいアプリケーションを作成** をクリック
2. 以下の情報を入力：
   - **アプリケーション名**: `apprun-playground`
   - **リージョン**: 任意（東京など）

### 3. コンテナイメージの設定

1. **コンテナイメージ** セクションで以下を設定：
   - **イメージ URL**: `ghcr.io/cateiru/apprun-playground:latest`
   - **ポート**: `8080`

### 4. 環境変数の設定（必要に応じて）

AppRun の環境変数設定で以下を追加（任意）：

| 変数名 | 値 | 説明 |
|--------|-----|------|
| `PORT` | `8080` | アプリケーションのポート（デフォルト値） |

### 5. リソースの設定

- **CPU**: 最小限（0.25 vCPU など）
- **メモリ**: 256MB〜512MB
- **インスタンス数**: 1（最初は1つで十分）

### 6. デプロイ

1. 設定を確認
2. **デプロイ** ボタンをクリック
3. デプロイが完了するまで待機

### 7. 動作確認

デプロイ完了後、AppRun が提供する URL にアクセス：

```bash
curl https://your-app-url.apprun.sakura.ne.jp/
```

期待されるレスポンス：

```json
{"status":"ok"}
```

## GitHub Container Registry のイメージを使用する場合

### Public イメージの場合

- 特別な設定は不要
- イメージ URL に `ghcr.io/cateiru/apprun-playground:latest` を指定するだけ

### Private イメージの場合

1. GitHub Personal Access Token (PAT) を作成
   - Scope: `read:packages`
2. AppRun のコンテナイメージ設定で認証情報を追加：
   - **ユーザー名**: GitHub ユーザー名
   - **パスワード**: Personal Access Token

## 継続的デプロイ

### 自動デプロイの設定

AppRun では、Webhook を使用して自動デプロイを設定できます：

1. AppRun のアプリケーション設定で **Webhook URL** を取得
2. GitHub リポジトリの **Settings** > **Webhooks** で追加
3. イベントを選択（Package publish など）

または、GitHub Actions から直接 AppRun にデプロイする方法もあります（AppRun の API を使用）。

## チェックリスト

- [ ] AppRun コンソールにアクセス
- [ ] アプリケーションを作成
- [ ] コンテナイメージを設定（`ghcr.io/cateiru/apprun-playground:latest`）
- [ ] ポート番号を設定（8080）
- [ ] 必要に応じて環境変数を設定
- [ ] リソース（CPU、メモリ）を設定
- [ ] デプロイを実行
- [ ] デプロイが成功したことを確認
- [ ] 提供された URL にアクセスして動作確認
- [ ] レスポンスが `{"status":"ok"}` であることを確認

## トラブルシューティング

### デプロイに失敗する場合

- イメージ URL が正しいか確認
- イメージが public になっているか、または認証情報が正しいか確認
- ポート番号が正しいか確認（8080）
- AppRun のログを確認

### アプリケーションが起動しない場合

- アプリケーションログを確認
- 環境変数 `PORT` が正しく設定されているか確認
- リソース（CPU、メモリ）が不足していないか確認

### ヘルスチェックが失敗する場合

- `/` エンドポイントが正しく実装されているか確認
- ポート番号が一致しているか確認
- アプリケーションが正常に起動しているかログで確認

## 補足

### スケーリング

トラフィックが増えた場合：

- インスタンス数を増やす
- CPU/メモリを増やす
- オートスケーリングを設定（AppRun の機能による）

### モニタリング

- AppRun のダッシュボードでメトリクスを確認
- ログを定期的に確認
- アラート設定（任意）

### コスト

- 使用したリソースに応じて課金
- 無料枠がある場合は活用
- 不要な時はインスタンスを停止

## 参考リンク

- [さくらのクラウド AppRun ドキュメント](https://manual.sakura.ad.jp/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [GitHub Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)

## 次のステップ

- カスタムドメインの設定
- HTTPS の設定（通常は自動）
- 追加のエンドポイントの実装
- データベースの接続
- ロギングとモニタリングの強化
