package sample

import (
	"context"

	pb "go-grpc/pb/sample"

	serviceSample "go-grpc/internal/services/sample"
)

// インターフェースの定義
type SampleHelloUsecase interface {
	Exec(ctx context.Context, in *pb.Empty) (*pb.HelloResponseBody, error)
}

// 構造体の定義
type sampleHelloUsecase struct {
	sampleService serviceSample.SampleService
}

// インスタンス生成用関数の定義
func NewSampleHelloUsecase(
	sampleService serviceSample.SampleService,
) SampleHelloUsecase {
	return &sampleHelloUsecase{
		sampleService: sampleService,
	}
}

func (u *sampleHelloUsecase) Exec(ctx context.Context, in *pb.Empty) (*pb.HelloResponseBody, error) {
	text, err := u.sampleService.Sample()
	if err != nil {
		return nil, err
	}

	return &pb.HelloResponseBody{Message: text}, nil
}
