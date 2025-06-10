package interceptor

import (
	"context"
	"fmt"
	"strings"

	utilCtx "go-grpc/internal/util/context"
	"go-grpc/internal/util/logger"
	pb "go-grpc/pb/sample"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Streamでコンテキストを共有させるためのラッパー構造体（対象メソッドのオーバーライド）
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

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

	// ステータスコードを取得
	st, ok := status.FromError(err)
	if !ok {
		return nil, fmt.Errorf("ステータスコードを取得できませんでした。")
	}

	// ctxにstatusとstatusCodeを設定
	ctx = context.WithValue(ctx, utilCtx.Status, st.Code().String())
	ctx = context.WithValue(ctx, utilCtx.StatusCode, fmt.Sprintf("%d", int(st.Code())))

	// トレーラーに情報追加
	if err != nil {
		tStatus := metadata.New(map[string]string{"status": fmt.Sprintf("ERROR ( %d %s )", int(st.Code()), st.Code().String())})
		if err := grpc.SetTrailer(ctx, tStatus); err != nil {
			return nil, err
		}
	} else {
		tStatus := metadata.New(map[string]string{"status": fmt.Sprintf("SUCCESS ( %d %s )", int(st.Code()), st.Code().String())})
		if err := grpc.SetTrailer(ctx, tStatus); err != nil {
			return nil, err
		}
	}

	// リクエスト終了のログ出力
	logger.Info(ctx, "finish request")

	return res, err
}

// リクエスト用のStreamインターセプター
func RequestStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// ctxにx-request-idを設定
	requestId := uuid.New().String()
	ctx := ss.Context()
	ctx = context.WithValue(ctx, utilCtx.XRequestId, requestId)

	// レスポンスのメタデータにx-request-idを追加
	headerMD := metadata.New(map[string]string{string(utilCtx.XRequestId): requestId})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return err
	}

	// メタデータを取得
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("メタデータを取得できません。")
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
		return err
	}

	// リクエスト開始のログ出力
	logger.Info(ctx, "start stream request")

	err := handler(srv, &wrappedServerStream{ss, ctx})

	// ステータスコードを取得
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("ステータスコードを取得できませんでした。")
	}

	// ctxにstatusとstatusCodeを設定
	ctx = context.WithValue(ctx, utilCtx.Status, st.Code().String())
	ctx = context.WithValue(ctx, utilCtx.StatusCode, fmt.Sprintf("%d", int(st.Code())))

	// トレーラーに情報追加
	if err != nil {
		tStatus := metadata.New(map[string]string{"status": fmt.Sprintf("ERROR ( %d %s )", int(st.Code()), st.Code().String())})
		if err := grpc.SetTrailer(ctx, tStatus); err != nil {
			return err
		}
	} else {
		tStatus := metadata.New(map[string]string{"status": fmt.Sprintf("SUCCESS ( %d %s )", int(st.Code()), st.Code().String())})
		if err := grpc.SetTrailer(ctx, tStatus); err != nil {
			return err
		}
	}

	// リクエスト終了のログ出力
	logger.Info(ctx, "finish stream request")

	return err
}

// 認証用のUnaryインターセプター
func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 対象外のメソッドを設定
	skipMethods := []string{
		pb.SampleService_Hello_FullMethodName,
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
		return nil, status.Errorf(codes.InvalidArgument, "%s", "認証用トークンが設定されていません。")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "認証用トークンが設定されていません。")
	}

	// TODO: 認証チェック処理を追加

	// 処理を実行
	return handler(ctx, req)
}

// 認証用のStreamインターセプター
func AuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 対象外のメソッドを設定
	skipMethods := []string{
		// pb.SampleService_HelloServerStream_FullMethodName,
	}

	// 対象外メソッドの場合はスキップ
	for _, method := range skipMethods {
		if info.FullMethod == method {
			return handler(srv, ss)
		}
	}

	// authorizationからトークンを取得
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("メタデータを取得できません。")
	}
	authHeader, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "%s", "認証用トークンが設定されていません。")
	}
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == "" {
		return status.Errorf(codes.InvalidArgument, "%s", "認証用トークンが設定されていません。")
	}

	// TODO: 認証チェック処理を追加

	return handler(srv, &wrappedServerStream{ss, ctx})
}
