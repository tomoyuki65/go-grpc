package interceptor

import (
	"context"
	"fmt"
	"strings"

	utilCtx "go-grpc/internal/util/context"
	"go-grpc/internal/util/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// リクエスト用のUnaryインターセプター
func RequestUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// ctxにx-request-idを設定
	requestId := uuid.New().String()
	ctx = context.WithValue(ctx, utilCtx.XRequestId, requestId)

	// レスポンスのメタデータにx-request-idを追加
	headerMD := metadata.New(map[string]string{string(utilCtx.XRequestId): requestId})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}

	// メタデータを取得
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("メタデータを取得できません。")
	}

	// リクエストのメタデータからx-request-sourceを取得
	requestSource, ok := md[string(utilCtx.XRequestSource)]
	if !ok {
		requestSource = []string{"-"}
	}

	// ctxにx-request-sourceを設定
	ctx = context.WithValue(ctx, utilCtx.XRequestSource, requestSource[0])

	// レスポンスのメタデータにx-request-sourceを追加
	headerMD2 := metadata.New(map[string]string{string(utilCtx.XRequestSource): requestSource[0]})
	if err := grpc.SetHeader(ctx, headerMD2); err != nil {
		return nil, err
	}

	// リクエスト開始のログ出力
	logger.Info(ctx, "start request")

	// 処理を実行
	res, err := handler(ctx, req)

	// トレーラーに情報を追加
	if err != nil {
		trailerMD := metadata.New(map[string]string{"status": "ERROR"})
		if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
			return nil, err
		}
	} else {
		trailerMD := metadata.New(map[string]string{"status": "OK"})
		if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
			return nil, err
		}
	}

	// リクエスト終了のログ出力
	logger.Info(ctx, "finish request")

	return res, err
}

func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 対象外のメソッドを設定
	skipMethods := []string{
		"/sample.SampleService/Hello",
	}

	// 対象外メソッドの場合はスキップ
	for _, method := range skipMethods {
		if info.FullMethod == method {
			return handler(ctx, req)
		}
	}

	// authorizationからトークンを取得
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("メタデータを取得できません。")
	}
	authHeader, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("認証用トークンが設定されていません。")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == "" {
		return nil, fmt.Errorf("認証用トークンが設定されていません。")
	}

	// TODO: 認証チェック処理を追加

	// 処理を実行
	return handler(ctx, req)
}
