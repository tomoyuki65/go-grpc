package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	utilCtx "go-grpc/internal/util/context"
	"go-grpc/internal/util/logger"

	"github.com/google/uuid"
)

// エラーレスポンス用の構造体
type errorResponse struct {
	Message string `json:"message"`
}

// リクエスト用のミドルウェア（gRPC-Gateway用）
func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctxにx-request-idを設定
		requestId := uuid.New().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, utilCtx.XRequestId, requestId)

		// レスポンスヘッダーにx-request-idを設定
		w.Header().Set(string(utilCtx.XRequestId), requestId)

		// リクエストヘッダーからx-request-sourceを取得
		requestSource := r.Header.Get(string(utilCtx.XRequestSource))
		if requestSource == "" {
			requestSource = "-"
		}

		// ctxにx-request-sourceを設定
		ctx = context.WithValue(ctx, utilCtx.XRequestSource, requestSource)

		// リクエスト開始のログ出力
		logger.Info(ctx, "start gRPC-Gateway request")

		// コンテキストを更新
		r = r.WithContext(ctx)

		// 処理を実行
		next.ServeHTTP(w, r)

		// リクエスト終了のログ出力
		logger.Info(ctx, "finish gRPC-Gateway request")
	})
}

// 認証用ミドルウェア（gRPC-Gateway用）
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 対象外のメソッドとエンドポイントを設定
		skipMethodAndEndpoints := []map[string]string{
			{
				"method":   http.MethodGet,
				"endpoint": "/api/v1/hello",
			},
		}

		// 特定のメソッドかつエンドポイントの場合にスキップ
		for _, target := range skipMethodAndEndpoints {
			if r.Method == target["method"] && r.URL.Path == target["endpoint"] {
				next.ServeHTTP(w, r)
				return
			}
		}

		ctx := r.Context()

		// リクエストヘッダーからAuthorizationを取得
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			// レスポンスをJSON形式で構築
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			// エラーメッセージをJSONでエンコード
			errRes := errorResponse{Message: "認証用トークンが設定されていません。"}
			if err := json.NewEncoder(w).Encode(errRes); err != nil {
				// ctxにstatusとstatusCodeを設定
				ctx = context.WithValue(ctx, utilCtx.Status, "Internal Server Error")
				ctx = context.WithValue(ctx, utilCtx.StatusCode, "500")

				logger.Error(ctx, fmt.Sprintf("Internal Server Error: %v", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

			// ctxにstatusとstatusCodeを設定
			ctx = context.WithValue(ctx, utilCtx.Status, "Unauthorized")
			ctx = context.WithValue(ctx, utilCtx.StatusCode, "401")

			logger.Warn(ctx, "認証用トークンが設定されていません。")

			return
		}

		// トークンを取得
		token := strings.TrimPrefix(authorization, "Bearer ")
		if token == "" {
			// レスポンスをJSON形式で構築
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			// エラーメッセージをJSONでエンコード
			errRes := errorResponse{Message: "認証用トークンが設定されていません。"}
			if err := json.NewEncoder(w).Encode(errRes); err != nil {
				// ctxにstatusとstatusCodeを設定
				ctx = context.WithValue(ctx, utilCtx.Status, "Internal Server Error")
				ctx = context.WithValue(ctx, utilCtx.StatusCode, "500")

				logger.Error(ctx, fmt.Sprintf("Internal Server Error: %v", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

			// ctxにstatusとstatusCodeを設定
			ctx = context.WithValue(ctx, utilCtx.Status, "Unauthorized")
			ctx = context.WithValue(ctx, utilCtx.StatusCode, "401")

			logger.Warn(ctx, "認証用トークンが設定されていません。")

			return
		}

		// TODO: 認証チェック処理を追加

		// 処理を実行
		next.ServeHTTP(w, r)
	})
}
