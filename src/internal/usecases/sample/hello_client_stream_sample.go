package sample

import (
	"errors"
	"fmt"
	"io"

	"go-grpc/internal/util/logger"
	pb "go-grpc/pb/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// コンテキストを取得
	ctx := stream.Context()

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

		// バリデーションチェック
		if err := req.Validate(); err != nil {
			msg := fmt.Sprintf("バリデーションエラー：%s", err.Error())
			logger.Warn(ctx, msg)

			return status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}

		textList = append(textList, req.GetText())
	}
}
