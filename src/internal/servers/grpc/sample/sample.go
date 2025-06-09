package sample

import (
	"context"

	repositorySample "go-grpc/internal/repositories/sample"
	serviceSample "go-grpc/internal/services/sample"
	usecaseSample "go-grpc/internal/usecases/sample"
	pb "go-grpc/pb/sample"

	"google.golang.org/grpc"
)

type sample struct {
	pb.UnimplementedSampleServiceServer
}

func NewSample() *sample {
	return &sample{}
}

func (s *sample) Hello(ctx context.Context, in *pb.Empty) (*pb.HelloResponseBody, error) {
	// インスタンス生成
	sampleRepository := repositorySample.NewSampleRepository()
	sampleService := serviceSample.NewSampleService(sampleRepository)
	sampleUsecase := usecaseSample.NewSampleHelloUsecase(sampleService)

	// ユースケースを実行
	return sampleUsecase.Exec(ctx, in)
}

func (s *sample) HelloAddText(ctx context.Context, in *pb.HelloAddTextRequestBody) (*pb.HelloAddTextResponseBody, error) {
	// インスタンス生成
	sampleUsecase := usecaseSample.NewSampleHelloAddTextUsecase()

	// ユースケースの実行
	return sampleUsecase.Exec(ctx, in)
}

// サーバーストリーミングのメソッド
func (s *sample) HelloServerStream(in *pb.HelloServerStreamRequestBody, stream grpc.ServerStreamingServer[pb.HelloServerStreamResponseBody]) error {
	// インスタンス生成
	sampleUsecase := usecaseSample.NewSampleHelloServerStreamUsecase()

	// ユースケースの実行
	return sampleUsecase.Exec(in, stream)
}
