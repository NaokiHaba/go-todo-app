package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	// os.Args: コマンドライン引数を取得するための変数
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <port>", os.Args[0])
		os.Exit(1)
	}

	p := os.Args[1]

	// net.Listen("tcp", ":"+p): 指定されたTCPポートでの接続待ちを行うための net.Listener オブジェクトを作成
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// run関数の外部から動的に選択したポートの接続待ちを行うための net.Listener オブジェクトを渡す このようにして、ポート番号を動的に選択できる
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to run: %v", err)
		os.Exit(1)
	}
}

// net.Listener を使用することで、ネットワーク上でクライアントからのリクエストを待ち受ける
func run(ctx context.Context, l net.Listener) error {
	// create a new http.Server
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, World! Section56 %s", r.URL.Path)
		}),
	}

	// create a new errgroup
	eg, ctx := errgroup.WithContext(ctx)

	// start a goroutine to serve the http server
	eg.Go(func() error {
		// start the http server
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to serve: %v", err)
			return err
		}
		return nil
	})

	// チャンネルからの終了を待つ
	<-ctx.Done()

	// HTTPサーバーをシャットダウンする
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}

	// 別ゴルーチンの終了を待つ
	return eg.Wait()
}
