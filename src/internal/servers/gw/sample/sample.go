package sample

import (
	"context"

	repositorySample "go-grpc/internal/repositories/sample"
	serviceSample "go-grpc/internal/services/sample"
	usecaseSample "go-grpc/internal/usecases/sample"
	pb "go-grpc/pb/sample"
)

type sampleApi struct {
	pb.UnimplementedSampleServiceServer
}

func NewSampleApi() *sampleApi {
	return &sampleApi{}
}

// gRPC-Gateway（GETメソッド）
func (s *sampleApi) HelloApi(ctx context.Context, in *pb.Empty) (*pb.HelloResponseBody, error) {
	// インスタンス生成
	sampleRepository := repositorySample.NewSampleRepository()
	sampleService := serviceSample.NewSampleService(sampleRepository)
	sampleUsecase := usecaseSample.NewSampleHelloUsecase(sampleService)

	// ユースケースを実行
	return sampleUsecase.Exec(ctx, in)
}

// gRPC-Gateway（POSTメソッド）
func (s *sampleApi) HelloAddTextApi(ctx context.Context, in *pb.HelloAddTextRequestBody) (*pb.HelloAddTextResponseBody, error) {
	// インスタンス生成
	sampleUsecase := usecaseSample.NewSampleHelloAddTextUsecase()

	// ユースケースの実行
	return sampleUsecase.Exec(ctx, in)
}
