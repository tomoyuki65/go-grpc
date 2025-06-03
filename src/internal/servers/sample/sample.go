package sample

import (
	"context"
	"fmt"

	pb "go-grpc/internal/pb/sample"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sample struct {
	pb.UnimplementedSampleServiceServer
}

func NewSample() *sample {
	return &sample{}
}

func (s *sample) Hello(ctx context.Context, in *pb.Empty) (*pb.HelloResponseBody, error) {
	return &pb.HelloResponseBody{Message: "Hello World !!"}, nil
}

func (s *sample) HelloAddText(ctx context.Context, in *pb.HelloAddTextRequestBody) (*pb.HelloAddTextResponseBody, error) {
	// バリデーションチェック
	if err := in.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	msg := fmt.Sprintf("Hello %s", in.Text)

	return &pb.HelloAddTextResponseBody{Message: msg}, nil
}