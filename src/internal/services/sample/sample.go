package sample

import (
	"fmt"

	"go-grpc/internal/repositories/sample"
)

// インターフェース定義
type SampleService interface {
	Sample() (string, error)
}

// 構造体定義
type sampleService struct {
	sampleRepository sample.SampleRepository
}

// インスタンス生成用関数
func NewSampleService(
	sampleRepository sample.SampleRepository,
) SampleService {
	return &sampleService{
		sampleRepository: sampleRepository,
	}
}

func (s *sampleService) Sample() (string, error) {
	text := s.sampleRepository.Hello()
	if text == "" {
		return "", fmt.Errorf("textが空です。")
	}

	return text, nil
}
