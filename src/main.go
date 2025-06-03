package main

import (
	"fmt"
	"net"

	pbSample "go-grpc/internal/pb/sample"
	serverSample "go-grpc/internal/servers/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// ポート番号の設定
	port := 8080
	
	// Listenerの設定
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーの作成
	s := grpc.NewServer()

	// サービス設定
	pbSample.RegisterSampleServiceServer(s, serverSample.NewSample())

	// リフレクション設定
	reflection.Register(s)

	// サーバー起動
	fmt.Println("start gRPC server port:", port)
	s.Serve(listener)
}