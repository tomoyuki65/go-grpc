package sample

import (
	"context"
	"testing"

	repoSample "go-grpc/internal/repositories/sample"
	mockRepoSample "go-grpc/internal/repositories/sample/mock_sample"
	serviceSample "go-grpc/internal/services/sample"
	pb "go-grpc/pb/sample"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSampleHelloOK(t *testing.T) {
	// ユースケースのインスタンス化
	sampleRepository := repoSample.NewSampleRepository()
	sampleService := serviceSample.NewSampleService(sampleRepository)
	sampleUsecase := NewSampleHelloUsecase(sampleService)

	// パラメータ設定
	ctx := context.Background()
	in := &pb.Empty{}

	// テストの実行
	res, err := sampleUsecase.Exec(ctx, in)

	// 検証
	assert.Equal(t, "Testing Hello World !!", res.Message)
	assert.Equal(t, nil, err)
}

func TestSampleHelloErr(t *testing.T) {
	// リポジトリのモック化
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSampleRepository := mockRepoSample.NewMockSampleRepository(ctrl)
	mockSampleRepository.EXPECT().Hello().Return("")

	// ユースケースのインスタンス化
	sampleService := serviceSample.NewSampleService(mockSampleRepository)
	sampleUsecase := NewSampleHelloUsecase(sampleService)

	// パラメータ設定
	ctx := context.Background()
	in := &pb.Empty{}

	// テストの実行
	res, err := sampleUsecase.Exec(ctx, in)

	// 検証
	assert.Equal(t, "", res.GetMessage())
	assert.NotEqual(t, nil, err)
}
