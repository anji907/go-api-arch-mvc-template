package logger

import (
	"os"

	"go.uber.org/zap"
)

var (
	ZapLogger        *zap.Logger
	zapSugaredLogger *zap.SugaredLogger
)

func init() {
	// production環境用の設定を構築して変数cfgに代入
	cfg := zap.NewProductionConfig()
	// 環境変数 APP_LOG_FILE の値を取得
	logFile := os.Getenv("APP_LOG_FILE")
	if logFile != "" {
		// Config構造体の OutputPaths フィールドにログの出力先を設定
		// 標準出力エラーと環境変数で指定したファイルにエラーを出力
		cfg.OutputPaths = []string{"stderr", logFile}
	}

	// ロガーを作成
	// zap.Mustの引数に渡すとエラー発生時にpanicを起こしてプログラムが停止する
	ZapLogger = zap.Must(cfg.Build())
	if os.Getenv("APP_ENV") == "development" {
		ZapLogger = zap.Must(zap.NewDevelopment())
	}
	zapSugaredLogger = ZapLogger.Sugar()
}

func Sync() {
	// メモリ内のバッファに残っているログをファイルに出力
	err := zapSugaredLogger.Sync()
	if err != nil {
		zap.Error(err)
	}
}

func Info(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Infow(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Fatalw(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	zapSugaredLogger.Panicw(msg, keysAndValues...)
}
