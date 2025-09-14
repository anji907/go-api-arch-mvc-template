package models_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-api-arch-mvc-template/app/models"
	"go-api-arch-mvc-template/pkg/tester"
)

type AlbumTestSuite struct {
	tester.DBSQLiteSuite
	originalDB *gorm.DB
}

func TestAlbumTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumTestSuite))
}

func (suite *AlbumTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.originalDB = models.DB
}

// 各テストケースの後に実行
func (suite *AlbumTestSuite) AfterTest(suiteName, testName string) {
	// テスト前のデータベースの状態に戻す
	models.DB = suite.originalDB
}

func Str2time(t string) time.Time {
	parsedTime, _ := time.Parse("2006-01-02", t)
	return parsedTime
}

// AlbumモデルのCRUDをテスト
func (suite *AlbumTestSuite) TestAlbum() {
	// CreateAlbum関数をテスト
	createdAlbum, err := models.CreateAlbum("Test", time.Now(), "sports")
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", createdAlbum.Title)
	suite.Assert().NotNil(createdAlbum.ReleaseDate)
	suite.Assert().NotNil(createdAlbum.Category.ID)
	suite.Assert().Equal(createdAlbum.Category.Name, "sports")

	// GetAlbum関数のテスト
	getAlbum, err := models.GetAlbum(createdAlbum.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test", getAlbum.Title)
	suite.Assert().NotNil(getAlbum.ReleaseDate)
	suite.Assert().NotNil(getAlbum.Category.ID)
	suite.Assert().Equal(getAlbum.Category.Name, "sports")

	// Saveメソッド(Upate)のテスト
	getAlbum.Title = "updated"
	err = getAlbum.Save()
	suite.Assert().Nil(err)
	updatedAlbum, err := models.GetAlbum(createdAlbum.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedAlbum.Title)
	suite.Assert().NotNil(updatedAlbum.ReleaseDate)
	suite.Assert().NotNil(updatedAlbum.Category.ID)
	suite.Assert().Equal(updatedAlbum.Category.Name, "sports")

	// Deleteメソッドのテスト
	err = updatedAlbum.Delete()
	suite.Assert().Nil(err)
	deletedAlbum, err := models.GetAlbum(updatedAlbum.ID)
	suite.Assert().Nil(deletedAlbum)
	suite.Assert().True(strings.Contains("record not founc", err.Error()))
}

// JSONへの変換処理をテスト
func (suite *AlbumTestSuite) TestAlbumMarshal() {
	album := models.Album{
		Title:       "Test",
		ReleaseDate: Str2time("2023-01-01"),
		Category:    &models.Category{},
	}
	anniversary := time.Now().Year() - 2023
	albumJSON, err := album.MarshalJSON()
	suite.Assert().Nil(err)
	suite.Assert().JSONEq(fmt.Sprintf(`{
		"anniversary":%d,
		"category":{
			"id":0,"name":"sports"
		},
		"id":0,
		"releaseDate":"2023-01-01",
		"title":"Test"
	}`, anniversary), string(albumJSON))
}
