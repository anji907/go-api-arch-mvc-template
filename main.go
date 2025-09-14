package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"github.com/gin-contrib/timeout"
	ginzap "github.com/gin-contrib/zap"

	"go-api-arch-mvc-template/api"
	"go-api-arch-mvc-template/app/controllers"
	"go-api-arch-mvc-template/app/models"
	"go-api-arch-mvc-template/configs"
	"go-api-arch-mvc-template/pkg/logger"
)

func main() {
	// データベースをセット
	if err := models.SetDatabase(models.InstanceMySQL); err != nil {
		logger.Fatal(err.Error())
	}

	// HTTPリクエストを振り分けるためのルーターをセット
	router := gin.Default()

	// Swaggerの仕様を取得
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	// 開発環境でのSwaggerのパスを設定
	if configs.Config.IsDevelopment() {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// APIのルーティング
	apiGroup := router.Group("/api")
	{
		// 2秒以内に処理が完了しなかったらタイムアウト
		apiGroup.Use(timeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			// リクエスト形式のチェック
			v1.Use(middleware.OapiRequestValidator(swagger))
			albumHandler := &controllers.AlbumHandler{}
			// albumに関するルーターの登録
			api.RegisterHandlers(v1, albumHandler)
		}
	}

	// ヘルスチェックようのルーティング
	router.GET("/health", controllers.Health)

	/*
		サーバーの起動処理
	*/
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		// ListenAndServe()でサーバーを起動してリクエストを待ち受ける
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	/*
		サーバーの終了処理
	*/

	// OSのシグナルを受け取るためのチャネルを作成
	quit := make(chan os.Signal, 1)
	// SIGINT, SIGTERMのシグナルがあればquitチャネルに送信
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	defer logger.Sync()

	// 2秒のタイムアウトを持つコンテキストとキャンセルする関数を作成
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}
	<-ctx.Done()
	logger.Info("Shutdown. bye bye...")

	// 許可するオリジンのリストを渡す
	router.Use(corsMiddleware(configs.Config.APICorsAllowOrigins))
	// リクエストやレスポンスの情報をログに出力
	router.Use(ginzap.Ginzap(logger.ZapLogger, time.RFC3339, true))
	// パニックが発生した場合にエラーとスタックトレースをログに出力
	router.Use(ginzap.RecoveryWithZap(logger.ZapLogger, true))
}

// CORSを設定するミドルウェア
func corsMiddleware(allowOrigins []string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = allowOrigins
	return cors.New(config)
}

// タイムアウト処理
func timeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(
				http.StatusRequestTimeout,
				api.ErrorResponse{Message: "timeout"},
			)
			c.Abort()
		}),
	)
}
