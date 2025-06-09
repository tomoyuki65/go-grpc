package sample

import (
	"errors"
	"fmt"
	"io"

	pb "go-grpc/pb/sample"

	"google.golang.org/grpc"
)

// インターフェースの定義
type SampleHelloClientStreamUsecase interface {
	Exec(stream grpc.ClientStreamingServer[pb.HelloClientStreamRequestBody, pb.HelloClientStreamResponseBody]) error
}

// 構造体の定義
type sampleHelloClientStreamUsecase struct{}

// インスタンス生成用関数の定義
func NewSampleHelloClientStreamUsecase() SampleHelloClientStreamUsecase {
	return &sampleHelloClientStreamUsecase{}
}

func (u *sampleHelloClientStreamUsecase) Exec(stream grpc.ClientStreamingServer[pb.HelloClientStreamRequestBody, pb.HelloClientStreamResponseBody]) error {
	textList := make([]string, 0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			msg := fmt.Sprintf("Hello, %v!", textList)
			return stream.SendAndClose(&pb.HelloClientStreamResponseBody{Message: msg})
		}
		if err != nil {
			return err
		}
		textList = append(textList, req.GetText())
	}
}
