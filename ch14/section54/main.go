package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	/**
	http.ListenAndServe とは、HTTP サーバーを起動する関数です。
	第一引数には、ポート番号を指定します。
	第二引数には、リクエストが来たときに実行するハンドラを指定します。
	ここでは、http.HandlerFunc という型の値を渡しています。
	この型は、ハンドラを表す型です。
	*/
	err := http.ListenAndServe(
		":18080",
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			/**
			レスポンスに "Hello, World!" と書き込みます。
			ここでは、fmt.Fprintf を使っています。
			この関数は、フォーマットされた文字列を書き込む関数です。
			第一引数には、io.Writer を指定します。
			第二引数には、フォーマット文字列を指定します。
			第三引数以降には、フォーマット文字列に埋め込む値を指定します。
			ここでは、r.URL.Path を指定しています。
			これは、リクエストされた URL のパスを表す値です。
			例えば、http://localhost:18080/foo/bar という URL がリクエストされた場合、
			r.URL.Path は "/foo/bar" という値になります。
			*/
			fmt.Fprintf(w, "Hello, World! %s", r.URL.Path)
		}))

	if err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
