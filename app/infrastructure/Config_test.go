package infrastructure

import (
	"os"
	"testing"
)

func TestNewConfigSuccess(t *testing.T) {
	result, err := NewConfig()
	if err != nil {
		t.Fatal("絶対パスの取得に失敗しました")
	}
	absolutePath, _ := os.Getwd()
	if result.AbsolutePath != absolutePath {
		t.Fatal("絶対パスに誤りがあります")
	}
}
