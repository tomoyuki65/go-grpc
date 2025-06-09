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

func TestSampleHelloBidirectionaStreamOK(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloBidirectionalStreamUsecase()

	// パラメータの設定
	mockBidirectionalStream := &mockHelloBidirectionalStream{
		ctx:       context.Background(),
		recv:      []*pb.HelloBidirectionalStreamRequestBody{{Text: "Tanaka"}},
		recvError: nil,
		sendError: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(mockBidirectionalStream)

	// 検証
	assert.Equal(t, nil, err)
}

func TestSampleHelloBidirectionaStreamRecvErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloBidirectionalStreamUsecase()

	// パラメータの設定
	mockBidirectionalStream := &mockHelloBidirectionalStream{
		ctx:       context.Background(),
		recv:      []*pb.HelloBidirectionalStreamRequestBody{{Text: "Tanaka"}},
		recvError: fmt.Errorf("error"),
		sendError: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(mockBidirectionalStream)

	// 検証
	assert.NotEqual(t, nil, err)
}

func TestSampleHelloBidirectionaStreamValidateErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloBidirectionalStreamUsecase()

	// パラメータの設定
	mockBidirectionalStream := &mockHelloBidirectionalStream{
		ctx:       context.Background(),
		recv:      []*pb.HelloBidirectionalStreamRequestBody{{Text: ""}},
		recvError: nil,
		sendError: nil,
	}

	// テストの実行
	err := sampleUsecase.Exec(mockBidirectionalStream)

	// 検証
	assert.NotEqual(t, nil, err)
}

func TestSampleHelloBidirectionaStreamSendErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloBidirectionalStreamUsecase()

	// パラメータの設定
	mockBidirectionalStream := &mockHelloBidirectionalStream{
		ctx:       context.Background(),
		recv:      []*pb.HelloBidirectionalStreamRequestBody{{Text: "Tanaka"}},
		recvError: nil,
		sendError: fmt.Errorf("error"),
	}

	// テストの実行
	err := sampleUsecase.Exec(mockBidirectionalStream)

	// 検証
	assert.NotEqual(t, nil, err)
}
