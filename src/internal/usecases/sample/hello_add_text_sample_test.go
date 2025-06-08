package sample

import (
	"context"
	"testing"

	pb "go-grpc/pb/sample"

	"github.com/stretchr/testify/assert"
)

func TestSampleHelloAddTextOK(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloAddTextUsecase()

	// パラメータの設定
	ctx := context.Background()
	in := &pb.HelloAddTextRequestBody{Text: "Add World !"}

	// テストの実行
	res, err := sampleUsecase.Exec(ctx, in)

	// 検証
	assert.Equal(t, "Hello Add World !", res.Message)
	assert.Equal(t, nil, err)
}

func TestSampleHelloAddTextValidateErr(t *testing.T) {
	// ユースケースのインスタンス化
	sampleUsecase := NewSampleHelloAddTextUsecase()

	// パラメータの設定
	ctx := context.Background()
	in := &pb.HelloAddTextRequestBody{}

	// テストの実行
	res, err := sampleUsecase.Exec(ctx, in)

	// 検証
	assert.Equal(t, "", res.GetMessage())
	assert.NotEqual(t, nil, err)
}
