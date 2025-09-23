# コンテナイメージを指定
FROM golang:1.25.1-alpine3.22
# curlをインストール
RUN apk --no-cache add curl
# ワークディレクトリを指定
WORKDIR /go/src/web
# ファイルやディレクトリをコピー
COPY .. .
# go.modに記載されているモジュールをダウンロード
RUN go mod download
# 環境変数の設定
ENV GO111MODULE=on
# main.goをコンパイルして実行可能ファイルを生成
RUN go build main.go
# コンテナ起動時に実行するコマンド
ENTRYPOINT ./main
