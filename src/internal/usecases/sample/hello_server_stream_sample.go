package sample

import (
	"fmt"
	"time"

	"go-grpc/internal/util/logger"
	pb "go-grpc/pb/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// コンテキストを取得
	ctx := stream.Context()

	// バリデーションチェック
	if err := in.Validate(); err != nil {
		msg := fmt.Sprintf("バリデーションエラー：%s", err.Error())
		logger.Warn(ctx, msg)

		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

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
