package sample

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	pb "go-grpc/pb/sample"

	"github.com/stretchr/testify/assert"
)

// カバレッジの対象から除外
func TestExcludeFromCoverage(t *testing.T) {
	s := NewSampleApi()
	_, _ = s.HelloApi(context.Background(), &pb.Empty{})
	_, _ = s.HelloAddTextApi(context.Background(), &pb.HelloAddTextRequestBody{Text: "Add World !"})
}

func TestSampleHelloApi(t *testing.T) {
	// リクエストURLの設定
	gatewayPort := os.Getenv("GATEWAY_PORT")
	if gatewayPort == "" {
		gatewayPort = "8080"
	}
	path := "/api/v1/hello"
	url := fmt.Sprintf("http://localhost:%s%s", gatewayPort, path)

	// リクエストの設定
	req, _ := http.NewRequest("GET", url, nil)

	// httpクライアントの初期化
	client := new(http.Client)

	// テストの実行
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call client.Do: %v", err)
	}
	defer res.Body.Close()

	// レスポンスボディをデコード
	var resBody pb.HelloResponseBody
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// 検証
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "Testing Hello World !!", resBody.Message)
}

func TestSampleHelloAddTextApi(t *testing.T) {
	// リクエストURLの設定
	gatewayPort := os.Getenv("GATEWAY_PORT")
	if gatewayPort == "" {
		gatewayPort = "8080"
	}
	path := "/api/v1/hello"
	url := fmt.Sprintf("http://localhost:%s%s", gatewayPort, path)

	// リクエストの設定
	reqBodyJsonByte := []byte(`{"text": "Add World !"}`)
	reqBody := bytes.NewBuffer(reqBodyJsonByte)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")

	// httpクライアントの初期化
	client := new(http.Client)

	// テストの実行
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call client.Do: %v", err)
	}
	defer res.Body.Close()

	// レスポンスボディをデコード
	var resBody pb.HelloAddTextResponseBody
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// 検証
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "Hello Add World !", resBody.Message)
}
