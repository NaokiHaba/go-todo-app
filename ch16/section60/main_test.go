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
	t.Skip("リファクタリング中")

	// TCPリスナーを作成する
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}

	// キャンセル可能なコンテキストを作成する
	ctx, cancel := context.WithCancel(context.Background())

	// ゴルーチンをグループ化する
	eg, ctx := errgroup.WithContext(ctx)

	// ゴルーチンを起動する
	eg.Go(func() error {
		return run(ctx)
	})

	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)

	// どんなポート番号でリッスンしているのか確認
	t.Logf("try request to %q", url)

	rsp, err := http.Get(url)

	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}

	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTPサーバーの戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	// run関数に終了通知を送信する。
	cancel()
	// run関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

}
