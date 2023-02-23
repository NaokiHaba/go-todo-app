package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net"
	"net/http"
	"testing"
)

func Test_run(t *testing.T) {
	// :18080で待ち受ける TCP リスナーを生成する
	l, err := net.Listen("tcp", ":18080")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}

	// キャンセル可能なコンテキストを生成する
	ctx, cancel := context.WithCancel(context.Background())

	// 複数のゴルーチンで発生するエラーをグループ化し、コンテキストを指定してキャンセル処理を行う
	eg, ctx := errgroup.WithContext(ctx)

	// グループ化された複数のゴルーチンを起動する
	eg.Go(func() error {
		return run(ctx, l)
	})

	in := "test"
	url := fmt.Sprintf("http://localhost:18080/%s", in)

	t.Logf("GET %s", url)

	// エンドポイントにGETリクエストを送信する
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("http.Get() error = %v", err)
	}

	// HTTPリクエストのレスポンスボディを閉じる
	defer rsp.Body.Close()

	// HTTPリクエストのレスポンスボディを完全に読み込む
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	want := fmt.Sprintf("Hello, World! Section56 /%s", in)
	if string(got) != want {
		t.Errorf("got = %v, want = %v", string(got), want)
	}

	// キャンセル可能なコンテキストをキャンセルする
	cancel()

	// グループ化された複数のゴルーチンの終了を待つ
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
