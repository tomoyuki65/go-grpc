FROM golang:1.24-alpine3.21

WORKDIR /go/src

# # インストール可能なパッケージ一覧の更新
RUN apk update && \
    apk upgrade && \
    # パッケージのインストール（--no-cacheでキャッシュ削除）
    apk add --no-cache \
            # git \
            # curl
            protobuf-dev \
            protoc

COPY ./src .

# go.modがあれば依存関係をインストール
RUN if [ -f ./go.mod ]; then \
      go install; \
    fi

# Protocol BuffersのGoとgRPCのプラグインをインストール
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5

# バリデーションプラグインをインストール
RUN go install github.com/envoyproxy/protoc-gen-validate@v1.2.1

# ドキュメント生成用ライブラリをインストール
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1

# 開発用のライブラリをインストール
# RUN go install github.com/air-verse/air@latest
# RUN go install honnef.co/go/tools/cmd/staticcheck@latest
# RUN go install go.uber.org/mock/mockgen@latest
# RUN go install github.com/swaggo/swag/cmd/swag@latest

# マイグレーション用のatlasをインストール
# RUN curl -sSf https://atlasgo.sh | sh

EXPOSE 8080