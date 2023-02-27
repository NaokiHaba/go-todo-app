# アプリケーションをビルドするための環境を構築
FROM golang:1.18.2-bullseye as deploy-builder

# 作業ディレクトリを /app に設定
WORKDIR /app

# go.mod と go.sum をコピー
COPY go.mod go.sum ./

#　依存パッケージをダウンロード
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド
# -trimpath: ビルド時にパスを削除
# -ldflags "-w -s": ビルド時にデバッグ情報を削除
# -o app: ビルドしたアプリケーションを app という名前で出力
RUN go build -trimpath -ldflags "-w -s" -o app

# ---------------------------------------------------

FROM debian:bullseye-slim as deploy

RUN apt-get update

# ビルド済みの実行可能ファイル(app)をコピー
COPY --from=deploy-builder /app/app .

# app を実行
CMD ["./app"]

# ---------------------------------------------------

FROM golang:1.18.2 as dev
WORKDIR /app

# 開発時にはコンテナ内でアプリケーションを起動し、ファイルの変更を検知して自動でビルド・リロードすることができます。
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]