package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"go-grpc/internal/interceptor"
	mw "go-grpc/internal/middleware"
	sGrpcSample "go-grpc/internal/servers/grpc/sample"
	sGwSample "go-grpc/internal/servers/gw/sample"
	pbSample "go-grpc/pb/sample"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// gRPC-Gatewayのサーバー起動用の関数
func grpcGateway(grpcPort, gatewayPort string) error {
	// gRPC-Serverへのエンドポイント設定
	ctx := context.Background()
	mux := runtime.NewServeMux()
	grpcServer := fmt.Sprintf("localhost:%s", grpcPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pbSample.RegisterSampleServiceHandlerFromEndpoint(ctx, mux, grpcServer, opts); err != nil {
		return err
	}

	// gRPC-Gatewayのエンドポイント設定
	if err := pbSample.RegisterSampleServiceHandlerServer(ctx, mux, sGwSample.NewSampleApi()); err != nil {
		return err
	}

	// ミドルウェアの設定（muxをラップ）
	handler := mw.RequestMiddleware(mw.AuthMiddleware(mux))

	// HTTPサーバーの起動
	listener := fmt.Sprintf(":%s", gatewayPort)
	return http.ListenAndServe(listener, handler)
}

func main() {
	// .env ファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		slog.Error(fmt.Sprintf(".envファイルの読み込みに失敗しました。: %v", err))
	}

	// ENVの設定
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// gRPC用のポート番号の設定
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// gRPC-Gateway用のポート番号の設定
	gatewayPort  := os.Getenv("GATEWAY_PORT")
	if gatewayPort  == "" {
		gatewayPort  = "8080"
	}

	// Listenerの設定
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		slog.Error(fmt.Sprintf("Listenerの設定に失敗しました。: %v", err))
	}

	// gRPCサーバーの作成
	s := grpc.NewServer(
		// インターセプターの適用
		grpc.ChainUnaryInterceptor(
			interceptor.RequestUnaryInterceptor,
			interceptor.AuthUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			interceptor.RequestStreamInterceptor,
			interceptor.AuthStreamInterceptor,
		),
	)

	// サービス設定
	pbSample.RegisterSampleServiceServer(s, sGrpcSample.NewSample())

	// リフレクション設定
	reflection.Register(s)

	// gRPCサーバーの起動（非同期）
	slog.Info(fmt.Sprintf("[ENV=%s] start gRPC-Server port: %s", env, grpcPort))
	go func() {
		if err := s.Serve(listener); err != nil {
			slog.Error(fmt.Sprintf("gRPC-Server の起動に失敗しました。: %v", err))
		}
	}()

	// gRPC-Gatewayの起動
	slog.Info(fmt.Sprintf("[ENV=%s] start gRPC-Gateway port: %s", env, gatewayPort))
	if err := grpcGateway(grpcPort, gatewayPort); err != nil {
		slog.Error(fmt.Sprintf("gRPC-Gatewayの起動に失敗しました。: %v", err))
	}
}
