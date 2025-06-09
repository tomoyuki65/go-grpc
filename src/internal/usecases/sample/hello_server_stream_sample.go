package sample

import (
	"fmt"
	"time"

	pb "go-grpc/pb/sample"

	"google.golang.org/grpc"
)

// インターフェースの定義
type SampleHelloServerStreamUsecase interface {
	Exec(in *pb.HelloServerStreamRequestBody, stream grpc.ServerStreamingServer[pb.HelloServerStreamResponseBody]) error
}

// 構造体の定義
type sampleHelloServerStreamUsecase struct{}

// インスタンス生成用関数の定義
func NewSampleHelloServerStreamUsecase() SampleHelloServerStreamUsecase {
	return &sampleHelloServerStreamUsecase{}
}

func (u *sampleHelloServerStreamUsecase) Exec(in *pb.HelloServerStreamRequestBody, stream grpc.ServerStreamingServer[pb.HelloServerStreamResponseBody]) error {
	resCount := 3
	for i := 0; i < resCount; i++ {
		if err := stream.Send(
			&pb.HelloServerStreamResponseBody{
				Message: fmt.Sprintf("[%d]Hello, %s", i, in.GetText()),
			},
		); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}
