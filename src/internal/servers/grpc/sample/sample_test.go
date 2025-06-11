package sample

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
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

// HelloClientStreamのstream用のモック構造体
type mockHelloClientStream struct {
	grpc.ClientStreamingServer[pb.HelloClientStreamRequestBody, pb.HelloClientStreamResponseBody]
	ctx        context.Context
	recvData   []*pb.HelloClientStreamRequestBody
	recvIndex  int
	recvError  error
	sendResult *pb.HelloClientStreamResponseBody
	sendError  error
}

func (m *mockHelloClientStream) Recv() (*pb.HelloClientStreamRequestBody, error) {
	if m.recvError != nil {
		return nil, m.recvError
	}
	if m.recvIndex >= len(m.recvData) {
		return nil, io.EOF
	}
	data := m.recvData[m.recvIndex]
	m.recvIndex++
	return data, nil
}

func (m *mockHelloClientStream) SendAndClose(resp *pb.HelloClientStreamResponseBody) error {
	m.sendResult = resp
	return m.sendError
}

func (m *mockHelloClientStream) Context() context.Context {
	return m.ctx
}

// HelloBidirectionalStreamのstream用のモック構造体
type mockHelloBidirectionalStream struct {
	grpc.BidiStreamingServer[pb.HelloBidirectionalStreamRequestBody, pb.HelloBidirectionalStreamResponseBody]
	ctx       context.Context
	sent      []*pb.HelloBidirectionalStreamResponseBody
	recv      []*pb.HelloBidirectionalStreamRequestBody
	recvIndex int
	sendError error
	recvError error
}

func (m *mockHelloBidirectionalStream) Send(resp *pb.HelloBidirectionalStreamResponseBody) error {
	if m.sendError != nil {
		return m.sendError
	}
	m.sent = append(m.sent, resp)
	return nil
}

func (m *mockHelloBidirectionalStream) Recv() (*pb.HelloBidirectionalStreamRequestBody, error) {
	if m.recvError != nil {
		return nil, m.recvError
	}
	if m.recvIndex >= len(m.recv) {
		return nil, io.EOF
	}
	req := m.recv[m.recvIndex]
	m.recvIndex++
	return req, nil
}

func (m *mockHelloBidirectionalStream) Context() context.Context {
	return m.ctx
}

// カバレッジの対象から除外
func TestExcludeFromCoverage(t *testing.T) {
	s := NewSample()
	_, _ = s.Hello(context.Background(), &pb.Empty{})
	_, _ = s.HelloAddText(context.Background(), &pb.HelloAddTextRequestBody{Text: "Add World !"})

	// サーバーストリーミング
	in := &pb.HelloServerStreamRequestBody{Text: "World !"}
	mockServerStream := &mockHelloServerStream{}
	_ = s.HelloServerStream(in, mockServerStream)

	// クライアントストリーミング
	mockClientStream := &mockHelloClientStream{}
	_ = s.HelloClientStream(mockClientStream)

	// 双方向ストリーミング
	mockBidirectionalStream := &mockHelloBidirectionalStream{}
	_ = s.HelloBidirectionalStream(mockBidirectionalStream)
}

func TestSampleHello(t *testing.T) {
	// gRPCクライアントの設定
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	target := fmt.Sprintf("dns:///localhost:%s", grpcPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	target := fmt.Sprintf("dns:///localhost:%s", grpcPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	target := fmt.Sprintf("dns:///localhost:%s", grpcPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		if errors.Is(err, io.EOF) {
			break
		}
		msgs = append(msgs, res.Message)
	}

	// 検証
	assert.Equal(t, []string{"[0]Hello, World !", "[1]Hello, World !", "[2]Hello, World !"}, msgs)
}

func TestSampleHelloClientStream(t *testing.T) {
	// gRPCクライアントの設定
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	target := fmt.Sprintf("dns:///localhost:%s", grpcPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	stream, err := client.HelloClientStream(ctx)
	if err != nil {
		t.Fatalf("Failed to call client.HelloClientStream: %v", err)
	}

	sendCount := 3
	for i := 0; i < sendCount; i++ {
		if err := stream.Send(&pb.HelloClientStreamRequestBody{Text: strconv.Itoa(i)}); err != nil {
			t.Fatalf("Failed to stream.Send: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to close and stream.CloseAndRecv: %v", err)
	}

	assert.Equal(t, "Hello, [0 1 2]!", res.Message)
}

func TestSampleHelloBidirectionalStream(t *testing.T) {
	// gRPCクライアントの設定
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	target := fmt.Sprintf("dns:///localhost:%s", grpcPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	stream, err := client.HelloBidirectionalStream(ctx)
	if err != nil {
		t.Fatalf("Failed to call client.HelloBidirectionalStream: %v", err)
	}

	var sendEnd, recvEnd bool
	for !(sendEnd && recvEnd) {
		// 送信処理
		if !sendEnd {
			if err := stream.Send(&pb.HelloBidirectionalStreamRequestBody{Text: "Tanaka"}); err != nil {
				t.Fatalf("Failed to stream.Send: %v", err)
			}
			if err := stream.CloseSend(); err != nil {
				t.Fatalf("Failed to stream.CloseSend: %v", err)
			}
			sendEnd = true
		}

		// 受信処理
		if !recvEnd {
			res, err := stream.Recv()
			if err != nil {
				t.Fatalf("Failed to stream.Recv: %v", err)
			}

			// 検証
			assert.Equal(t, "Hello, Tanaka !!", res.Message)
			recvEnd = true
		}
	}
}
