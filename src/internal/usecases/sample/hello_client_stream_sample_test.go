package sample

import (
	"context"
	"fmt"
	"io"
	"testing"

	pb "go-grpc/pb/sample"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

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

func TestSampleHelloClientStreamOK(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloClientStreamUsecase()

	// パラメータの設定
	mockClientStream := &mockHelloClientStream{
		ctx:       context.Background(),
		recvData:  []*pb.HelloClientStreamRequestBody{{Text: "A"}, {Text: "B"}, {Text: "C"}},
		recvError: nil,
		sendError: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(mockClientStream)

	// 検証
	assert.Equal(t, nil, err)
}

func TestSampleHelloClientStreamRecvErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloClientStreamUsecase()

	// パラメータの設定
	mockClientStream := &mockHelloClientStream{
		ctx:       context.Background(),
		recvData:  []*pb.HelloClientStreamRequestBody{{Text: "A"}, {Text: "B"}, {Text: "C"}},
		recvError: fmt.Errorf("error"),
		sendError: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(mockClientStream)

	// 検証
	assert.NotEqual(t, nil, err)
}

func TestSampleHelloClientStreamValidateErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloClientStreamUsecase()

	// パラメータの設定
	mockClientStream := &mockHelloClientStream{
		ctx:       context.Background(),
		recvData:  []*pb.HelloClientStreamRequestBody{{Text: ""}, {Text: "B"}, {Text: "C"}},
		recvError: nil,
		sendError: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(mockClientStream)

	// 検証
	assert.NotEqual(t, nil, err)
}
