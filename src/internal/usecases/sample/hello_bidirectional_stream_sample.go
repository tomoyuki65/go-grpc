package sample

import (
	"errors"
	"fmt"
	"io"

	pb "go-grpc/pb/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// インターフェースの定義
type SampleHelloBidirectionalStreamUsecase interface {
	Exec(stream grpc.BidiStreamingServer[pb.HelloBidirectionalStreamRequestBody, pb.HelloBidirectionalStreamResponseBody]) error
}

// 構造体の定義
type sampleHelloBidirectionalStreamUsecase struct{}

// インスタンス生成用関数の定義
func NewSampleHelloBidirectionalStreamUsecase() SampleHelloBidirectionalStreamUsecase {
	return &sampleHelloBidirectionalStreamUsecase{}
}

func (u *sampleHelloBidirectionalStreamUsecase) Exec(stream grpc.BidiStreamingServer[pb.HelloBidirectionalStreamRequestBody, pb.HelloBidirectionalStreamResponseBody]) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		// バリデーションチェック
		if err := req.Validate(); err != nil {
			return status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		msg := fmt.Sprintf("Hello, %s !!", req.GetText())
		if err := stream.Send(&pb.HelloBidirectionalStreamResponseBody{Message: msg}); err != nil {
			return err
		}
	}
}
