package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	// http.Server: HTTPリクエストを処理するサーバーを構築するための構造体
	s := &http.Server{
		Addr: ":18080",
		// http.HandlerFunc: http.Handlerを実装するための関数型
		// 関数を http.Handler インタフェースを実装するオブジェクトに変換する
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, World! Section55 %s", r.URL.Path)
		}),
	}

	// errgroup.WithContext: 複数のゴルーチンで発生するエラーをグループ化し、コンテキストを指定してキャンセル処理を行うための関数
	eg, ctx := errgroup.WithContext(ctx)

	//  eg.Go: errgroup.WithContextでグループ化された複数のゴルーチンを起動するための関数
	eg.Go(func() error {
		// http.Serverで定義されたサーバーを開始するための関数です。
		// http.ErrServerClosed: サーバーが正常にシャットダウンされたことを示すエラー
		if err := s.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to listen and serve: %v", err)
			return err
		}
		return nil
	})

	// <-ctx.Done(): コンテキストがキャンセルされたことを待つための構文
	<-ctx.Done()

	// http.Server.Shutdown: http.Serverで定義されたサーバーをシャットダウンするための関数
	// context.Background(): 空のコンテキストを生成するための関数
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}

	// eg.Wait: errgroup.WithContextでグループ化された複数のゴルーチンの終了を待つための関数
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to run: %v", err)
	}
}
