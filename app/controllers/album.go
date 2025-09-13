package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api-arch-mvc-template/api"
	"go-api-arch-mvc-template/app/models"
	"go-api-arch-mvc-template/pkg/logger"
)

type AlbumHandler struct{}

// アルバムの作成
func (a *AlbumHandler) CreateAlbum(c *gin.Context) {
	var requestBody api.CreateAlbumJSONRequestBody
	// JSON形式のリクストボディをGoの構造体にマッピング
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// エラーが発生したら警告レベルでログ出力
		logger.Warn(err.Error())
		// HTTPステータスとメッセージをHTTPレスポンスに書き込む
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	// Album modelの関数を呼び出す
	createdAlbum, err := models.CreateAlbum(
		requestBody.Title,
		requestBody.ReleaseDate.Time,
		string(requestBody.Category.Name),
	)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAlbum)
}

// アルバムの取得
func (a *AlbumHandler) GetAlbumById(c *gin.Context, ID int) {
	album, err := models.GetAlbum(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

// アルバムの更新
func (a *AlbumHandler) UpdatedAlbumById(c *gin.Context, ID int) {
	var requestBody api.UpdatedAlbumByIdJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	album, err := models.GetAlbum(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	if requestBody.Category != nil {
		album.Category.Name = string(requestBody.Category.Name)
	}

	if requestBody.Title != nil {
		album.Title = *requestBody.Title
	}
	if err := album.Save(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

// アルバムの削除
func (a *AlbumHandler) DeleteAlbumById(c *gin.Context, ID int) {
	album := models.Album{ID: ID}
	if err := album.Delete(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
