package sample

import (
	"context"
	"fmt"

	pb "go-grpc/pb/sample"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// インターフェースの定義
type SampleHelloAddTextUsecase interface {
	Exec(ctx context.Context, in *pb.HelloAddTextRequestBody) (*pb.HelloAddTextResponseBody, error)
}

// 構造体の定義
type sampleHelloAddTextUsecase struct{}

// インスタンス生成用関数の定義
func NewSampleHelloAddTextUsecase() SampleHelloAddTextUsecase {
	return &sampleHelloAddTextUsecase{}
}

func (u *sampleHelloAddTextUsecase) Exec(ctx context.Context, in *pb.HelloAddTextRequestBody) (*pb.HelloAddTextResponseBody, error) {
	// バリデーションチェック
	if err := in.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	msg := fmt.Sprintf("Hello %s", in.Text)

	return &pb.HelloAddTextResponseBody{Message: msg}, nil
}
