package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	ヘルスチェック処理
*/

// HTTPリクエストを受け取ってレスポンスを返却
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
