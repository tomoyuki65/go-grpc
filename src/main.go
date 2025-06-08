package main

import (
	"fmt"
	"log/slog"
	"os"
	"net"

	"go-grpc/internal/interceptor"
	sGrpcSample "go-grpc/internal/servers/grpc/sample"
	pbSample "go-grpc/pb/sample"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// .env ファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		slog.Error(".envファイルの読み込みに失敗しました。")
	}

	// ポート番号の設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Listenerの設定
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーの作成
	s := grpc.NewServer(
		// インターセプターの適用
		grpc.ChainUnaryInterceptor(
			interceptor.RequestUnaryInterceptor,
			interceptor.AuthUnaryInterceptor,
		),
	)

	// サービス設定
	pbSample.RegisterSampleServiceServer(s, sGrpcSample.NewSample())

	// リフレクション設定
	reflection.Register(s)

	// サーバー起動
	fmt.Println("start gRPC server port:", port)
	s.Serve(listener)
}
