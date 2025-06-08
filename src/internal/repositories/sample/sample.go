package sample

import (
	"os"
)

// インターフェース定義
type SampleRepository interface {
	Hello() string
}

// 構造体定義
type sampleRepository struct{}

// インスタンス生成用関数
func NewSampleRepository() SampleRepository {
	return &sampleRepository{}
}

// メソッド定義
func (r *sampleRepository) Hello() string {
	env := os.Getenv("ENV")

	res := "Hello World !!"
	if env == "testing" {
		res = "Testing Hello World !!"
	}

	return res
}
