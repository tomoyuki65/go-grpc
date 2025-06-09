package sample

import (
	"context"
	"fmt"
	"testing"

	pb "go-grpc/pb/sample"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

func TestSampleHelloServerStreamOK(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloServerStreamUsecase()

	// パラメータの設定
	in := &pb.HelloServerStreamRequestBody{Text: "World !"}
	mockStream := &mockHelloServerStream{
		ctx:     context.Background(),
		sent:    []*pb.HelloServerStreamResponseBody{},
		sendErr: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(in, mockStream)

	// 検証
	assert.Equal(t, nil, err)
}

func TestSampleHelloServerStreamErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloServerStreamUsecase()

	// パラメータの設定
	in := &pb.HelloServerStreamRequestBody{Text: "World !"}
	mockStream := &mockHelloServerStream{
		ctx:     context.Background(),
		sent:    []*pb.HelloServerStreamResponseBody{},
		sendErr: fmt.Errorf("error"),
	}

	// テストの実行
	err := sampleUsecase.Exec(in, mockStream)

	// 検証
	assert.NotEqual(t, nil, err)
}
