# Go TODO アプリ実装計画

## プロジェクト概要

Go言語を使用してモダンUIを持つWeb TODOアプリケーションを開発する。
シンプルで使いやすく、レスポンシブデザインのTODO管理システムを構築する。

## 技術スタック

### バックエンド
- **Go 1.21+** - メインのプログラミング言語
- **Gin** - 軽量で高性能なWebフレームワーク
- **GORM** - Go用ORM（Object-Relational Mapping）
- **SQLite** - 開発・本番環境用データベース
- **JWT** - 認証・認可システム
- **Air** - ホットリロード開発ツール

### フロントエンド
- **HTML5** - マークアップ
- **CSS3** - スタイリング（CSS Grid、Flexbox使用）
- **Vanilla JavaScript** - インタラクティブ機能
- **Tailwind CSS** - ユーティリティファーストCSSフレームワーク
- **Alpine.js** - 軽量リアクティブフレームワーク

### 開発ツール
- **Docker** - コンテナ化
- **Docker Compose** - 開発環境管理
- **Git** - バージョン管理
- **Make** - ビルド自動化

## プロジェクト構造

```
gotodo/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── todo.go
│   │   └── user.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── models/
│   │   ├── todo.go
│   │   └── user.go
│   ├── database/
│   │   └── database.go
│   └── utils/
│       ├── jwt.go
│       └── validation.go
├── web/
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/
│       ├── index.html
│       ├── login.html
│       └── register.html
├── docs/
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## 機能要件

### 基本機能
1. **ユーザー認証**
   - ユーザー登録
   - ログイン・ログアウト
   - JWT トークンベース認証

2. **TODO管理**
   - TODO項目の作成
   - TODO項目の表示（リスト形式）
   - TODO項目の編集
   - TODO項目の削除
   - 完了・未完了の切り替え

3. **UI/UX機能**
   - レスポンシブデザイン
   - リアルタイム更新
   - ドラッグ&ドロップでの並び替え
   - フィルタリング（全て・完了・未完了）
   - 検索機能

### 拡張機能（後期実装）
1. **カテゴリ機能**
   - TODOのカテゴリ分け
   - カテゴリ別表示

2. **優先度機能**
   - 優先度設定（高・中・低）
   - 優先度別ソート

3. **期限機能**
   - 期限設定
   - 期限順ソート
   - 期限通知

## データベース設計

### Usersテーブル
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Todosテーブル
```sql
CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 1,
    due_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## API設計

### 認証エンドポイント
- `POST /api/auth/register` - ユーザー登録
- `POST /api/auth/login` - ログイン
- `POST /api/auth/logout` - ログアウト
- `GET /api/auth/me` - ユーザー情報取得

### TODOエンドポイント
- `GET /api/todos` - TODO一覧取得
- `POST /api/todos` - TODO作成
- `GET /api/todos/:id` - TODO詳細取得
- `PUT /api/todos/:id` - TODO更新
- `DELETE /api/todos/:id` - TODO削除
- `PATCH /api/todos/:id/toggle` - TODO完了状態切り替え

## 実装フェーズ

### フェーズ 1: 基盤構築（1-2週間）
1. プロジェクト初期化
2. 基本的なプロジェクト構造作成
3. データベース設定
4. 基本的なWebサーバー構築

### フェーズ 2: 認証システム（1週間）
1. ユーザーモデル実装
2. JWT認証システム実装
3. 登録・ログイン画面作成
4. 認証ミドルウェア実装

### フェーズ 3: TODO基本機能（1-2週間）
1. TODOモデル実装
2. CRUD API実装
3. 基本的なUI作成
4. フロントエンド・バックエンド連携

### フェーズ 4: UI/UX改善（1週間）
1. レスポンシブデザイン実装
2. リアルタイム更新機能
3. フィルタリング・検索機能
4. ドラッグ&ドロップ機能

### フェーズ 5: 拡張機能（1-2週間）
1. カテゴリ機能実装
2. 優先度機能実装
3. 期限機能実装
4. 通知機能実装

### フェーズ 6: テスト・デプロイ（1週間）
1. ユニットテスト実装
2. 統合テスト実装
3. Docker化
4. デプロイ準備

## 開発環境セットアップ

### 前提条件
- Go 1.21以上
- Node.js 18以上（フロントエンドツール用）
- Git
- Docker（オプション）

### セットアップ手順
1. リポジトリクローン
2. 依存関係インストール（`go mod tidy`）
3. データベース初期化
4. 開発サーバー起動（`make dev`）

## セキュリティ考慮事項

1. **パスワード**
   - bcryptによるハッシュ化
   - 最小文字数制限

2. **JWT トークン**
   - 短い有効期限設定
   - リフレッシュトークン実装

3. **入力検証**
   - SQLインジェクション対策
   - XSS対策
   - CSRF対策

4. **API レート制限**
   - IP基準のレート制限
   - ユーザー基準のレート制限

## パフォーマンス最適化

1. **データベース**
   - インデックス適用
   - クエリ最適化

2. **静的ファイル**
   - CSS/JS圧縮
   - 画像最適化
   - CDN使用（本番環境）

3. **キャッシュ**
   - HTTP キャッシュヘッダー
   - インメモリキャッシュ

## 監視・ログ

1. **ログ管理**
   - 構造化ログ（JSON形式）
   - ログレベル分け
   - ログローテーション

2. **メトリクス**
   - アプリケーションメトリクス
   - パフォーマンスメトリクス
   - エラー率監視

## デプロイメント

1. **開発環境**
   - Docker Compose使用
   - ホットリロード対応

2. **本番環境**
   - Docker コンテナ化
   - リバースプロキシ（Nginx）
   - HTTPS対応
   - 自動デプロイ（CI/CD）