# GoのgRPCによるAPI開発のサンプル
Go言語（Golang）のgRPCによるバックエンドAPI開発用サンプルです。  
  
<br />
  
## 要件
・Goのバージョンは<span style="color:green">1.24.x</span>です。  
  
<br />
  
## ローカル開発環境構築
### 1. 環境変数ファイルをリネーム
```
cp ./src/.env.example ./src/.env
```  
  
### 2. コンテナのビルドと起動
```
docker compose build --no-cache
docker compose up -d
```  
> ※テストコードを実行させる際はテスト用環境変数ファイルを使うため、「docker compose --env-file ./src/.env.testing up -d」で起動すること。  
  
### 3. コンテナの停止・削除
```
docker compose down
```  
  
<br />
  
## コード修正後に使うコマンド
ローカルサーバー起動中に以下のコマンドを実行可能です。  
  
### 1. go.modの修正
```
docker compose exec grpc go mod tidy
```  
  
### 2. フォーマット修正
```
docker compose exec grpc go fmt ./internal/...
```  
  
### 3. コード解析チェック
```
docker compose exec grpc staticcheck ./internal/...
```  
  
### 4. モック用ファイル作成（例）
```
docker compose exec grpc mockgen -source=./internal/repositories/XXX/XXX.go -destination=./internal/repositories/XXX/mock_XXX/mock_XXX.go
```  
  
### 5. テストコードの実行
事前にテスト用環境変数を設定したローカルサーバーを起動（docker compose --env-file ./src/.env.testing up -d）してから以下のコマンドを実行して下さい。  
```
docker compose exec grpc go test -v -cover ./internal/servers/...
docker compose exec grpc go test -v -cover ./internal/usecases/...
```  
> ※オプション「-cover」でカバレッジを確認できます。カバレッジは80%以上推薦です。必要に応じてサービス層のテストコードも作成して下さい。  
  
### 6. テストコードのカバレッジ対象確認用のファイル出力
必要に応じて以下のコマンドを実行し、出力されるファイルからカバレッジ対象のコードを確認して下さい。  
  
・サーバー層のカバレッジ確認ファイルを出力  
```
docker compose exec grpc go test -v -coverprofile=internal/servers/coverage.out ./internal/servers/...

docker compose exec grpc go tool cover -html=internal/servers/coverage.out -o=internal/servers/coverage.html
```  
  
・ユースケース層のカバレッジ確認ファイルを出力
```
docker compose exec grpc go test -v -coverprofile=internal/usecases/coverage.out ./internal/usecases/...

docker compose exec grpc go tool cover -html=internal/usecases/coverage.out -o=internal/usecases/coverage.html
```  
  
<br />
  
## .protoファイルからコード生成
ローカルサーバー起動中に以下のコマンドを実行可能です。  
  
### 1. pb（Protocol Buffers）ファイルを生成
```
docker compose exec grpc protoc -I=.:../pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.2.1:../pkg/mod/github.com/googleapis/googleapis@v0.0.0-20250610203048-111b73837522:../pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3 --go_out=. --go-grpc_out=. --validate_out=lang=go:. ./proto/sample/sample.proto
```  
> ※ファイルの出力先は.protoファイル内の「option go_package="pb/sample";」の部分で指定 
  
> ※オプション「-I」でライブラリのパス（コンテナ内のパス）を指定  
  
### 2. ドキュメントファイルを生成
```
docker compose run --rm grpc protoc -I=.:../pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.2.1:../pkg/mod/github.com/googleapis/googleapis@v0.0.0-20250610203048-111b73837522:../pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.26.3 --doc_out=./doc --doc_opt=markdown,docs.md --openapiv2_out=allow_merge=true,merge_file_name=./openapi:./doc ./proto/sample/sample.proto
```  
> ※対象の.protoファイルが複数ある場合、末尾にパスを追加して下さい。  
  
> ※gRPC-Gateway用のOpenAPI仕様書として「openapi.swagger.json」も出力されます。  
  
<br />
  
## 本番環境用のコンテナについて
本番環境用コンテナをローカルでビルドして確認したい場合は、以下の手順で行って下さい。  
  
### 1. .env.productionの修正
本番環境用の機密情報を含まない環境変数の設定には「.env.production」を使います。
ローカルで確認する場合は必要に応じて内容を修正して下さい。  
  
### 2. コンテナのビルド
以下のコマンドを実行し、コンテナをビルドします。  
```
docker build --no-cache -f ./docker/prod/Dockerfile -t go-grpc:latest .
```  
  
### 3. コンテナの起動
以下のコマンドを実行し、コンテナを起動します。  
```
docker run -d -p 80:8080 -p 50051:50051 go-grpc:latest
```  
  