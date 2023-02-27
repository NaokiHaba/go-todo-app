package main

import (
	"context"
	"fmt"
	"github.com/NaokiHaba/go-todo-app/config"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminated server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// 環境変数をパースする
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// 環境変数で定義したポートで待ち受けるTCPリスナー(TCP接続を受け入れるためのプログラムや機能)を作成する
	// TCPは、データを複数のパケットに分割して送信し、パケットの到着を確認することで、データの信頼性を高めます
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	// l.Addr().String()は、リスナーが待ち受けているアドレスを返す
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	// http.Server構造体のポインタを作成
	s := &http.Server{
		// HTTPリクエストがあった際に呼び出されるハンドラー関数を定義
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	// 複数のゴルーチンで発生するエラーをグループ化して管理する
	eg, ctx := errgroup.WithContext(ctx)

	// errgroup.WithContextでグループ化された複数のゴルーチンを起動する
	eg.Go(func() error {
		// HTTPリクエストを受け付けるためにリスナーを監視
		// リクエストがあった場合には、事前に指定されたハンドラー関数を呼び出してレスポンスを生成し、クライアントに送信
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// ctxが終了するまで待機
	<-ctx.Done()

	// サーバーをシャットダウン
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %+v", err)
		return err
	}

	// 複数のゴルーチンが終了するまで、処理を待機
	return eg.Wait()
}
