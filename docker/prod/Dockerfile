####################
# ビルドステージ
####################
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /go/src

COPY ./src .

# 依存関係をインストール
RUN go install

# ビルド
RUN go build -o main .

####################
# 実行ステージ
####################
FROM alpine:3.21 AS runner

WORKDIR /go/src

# コンテナ用ユーザー作成
RUN addgroup --system --gid 1001 appuser && \
    adduser --system --uid 1001 appuser

# ビルドステージで作成したバイナリをコピー
COPY --from=builder --chown=appuser:appuser ./go/src/main .
COPY --from=builder --chown=appuser:appuser ./go/src/.env.production ./.env

# ポートを設定
EXPOSE 8080

# コンテナ起動ユーザー設定
USER appuser

# APIサーバー起動コマンド
CMD ["./main"]