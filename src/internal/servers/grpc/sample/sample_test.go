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

// HelloServerStreamのstream用のモック構造体
type mockHelloServerStream struct {
	grpc.ServerStreamingServer[pb.HelloServerStreamResponseBody]
	ctx     context.Context
	sent    []*pb.HelloServerStreamResponseBody
	sendErr error
}

func (m *mockHelloServerStream) Send(resp *pb.HelloServerStreamResponseBody) error {
	m.sent = append(m.sent, resp)
	return m.sendErr
}

func (m *mockHelloServerStream) Context() context.Context {
	return m.ctx
}

// カバレッジの対象から除外
func TestExcludeFromCoverage(t *testing.T) {
	s := NewSample()
	_, _ = s.Hello(context.Background(), &pb.Empty{})
	_, _ = s.HelloAddText(context.Background(), &pb.HelloAddTextRequestBody{Text: "Add World !"})

	// サーバーストリーミング
	in := &pb.HelloServerStreamRequestBody{Text: "World !"}
	mockStream := &mockHelloServerStream{
		ctx:     context.Background(),
		sent:    []*pb.HelloServerStreamResponseBody{},
		sendErr: nil,
	}
	_ = s.HelloServerStream(in, mockStream)
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

func TestSampleHelloServerStream(t *testing.T) {
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
	in := &pb.HelloServerStreamRequestBody{Text: "World !"}
	stream, err := client.HelloServerStream(ctx, in)
	if err != nil {
		t.Fatalf("Failed to call client.HelloServerStream: %v", err)
	}

	// ストリーミング処理のメッセージを取得
	var msgs []string
	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		msgs = append(msgs, res.Message)
	}

	// 検証
	assert.Equal(t, []string{"[0]Hello, World !", "[1]Hello, World !", "[2]Hello, World !"}, msgs)
}
