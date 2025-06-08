package sample

import (
	"context"
	"testing"

	pb "go-grpc/pb/sample"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// カバレッジの対象から除外
func TestExcludeFromCoverage(t *testing.T) {
	s := NewSample()
	_, _ = s.Hello(context.Background(), &pb.Empty{})
	_, _ = s.HelloAddText(context.Background(), &pb.HelloAddTextRequestBody{Text: "Add World !"})
}

func TestSampleHello(t *testing.T) {
	// gRPCクライアントの設定
	conn, err := grpc.NewClient("dns:///localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewSampleServiceClient(conn)

	// テストの実行
	res, err := client.Hello(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Failed to call Hello: %v", err)
	}

	// 検証
	assert.Equal(t, "Testing Hello World !!", res.Message)
}

func TestSampleHelloAddText(t *testing.T) {
	// gRPCクライアントの設定
	conn, err := grpc.NewClient("dns:///localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewSampleServiceClient(conn)

	// メタデータにauthorizationを追加
	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": "Bearer token"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// テストの実行
	res, err := client.HelloAddText(ctx, &pb.HelloAddTextRequestBody{Text: "Add World !"})
	if err != nil {
		t.Fatalf("Failed to call HelloAddText: %v", err)
	}

	// 検証
	assert.Equal(t, "Hello Add World !", res.Message)
}
