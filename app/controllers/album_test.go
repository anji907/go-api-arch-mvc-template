package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-api-arch-mvc-template/api"
	"go-api-arch-mvc-template/app/models"
	"go-api-arch-mvc-template/pkg/tester"
)

type AlbumControllersSuite struct {
	tester.DBSQLiteSuite
	albumHandler AlbumHandler
	originalDB   *gorm.DB
}

func TestAlbumControllersTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumControllersSuite))
}

func (suite *AlbumControllersSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.albumHandler = AlbumHandler{}
	suite.originalDB = models.DB
}

func (suite *AlbumControllersSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	models.DB = mockGormDB
	return mock
}

func (suite *AlbumControllersSuite) AfterTest(suiteName, testName string) {
	models.DB = suite.originalDB
}

func (suite *AlbumControllersSuite) TestCreate() {
	request, _ := api.NewCreateAlbumRequest("/api/v1", api.CreateAlbumJSONRequestBody{
		Title:       "test",
		Category:    api.Category{Name: "sports"},
		ReleaseDate: api.ReleaseDate{Time: time.Now()},
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.CreateAlbum(ginContext)

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse api.AlbumResponse
	err := json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, w.Code)
	suite.Assert().Equal("test", albumGetResponse.Title)
	suite.Assert().Equal("sports", string(albumGetResponse.Category.Name))
	suite.Assert().NotNil(albumGetResponse.ReleaseDate)
}

func (suite *AlbumControllersSuite) TestCreateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("POST", "/api/v1/album", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.albumHandler.CreateAlbum(ginContext)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *AlbumControllersSuite) TestCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`name` = ? ORDER BY `categories`.`id` LIMIT ?")).WithArgs("sports", 1).WillReturnError(errors.New("create error"))

	request, _ := api.NewCreateAlbumRequest("/api/v1", api.CreateAlbumJSONRequestBody{
		Title:       "test",
		Category:    api.Category{Name: "sports"},
		ReleaseDate: api.ReleaseDate{Time: time.Now()},
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.CreateAlbum(ginContext)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().True(strings.Contains(w.Body.String(), "create error"))
}

func (suite *AlbumControllersSuite) TestGet() {
	createdAlbum, _ := models.CreateAlbum("test", time.Now(), "sports")

	request, _ := api.NewGetAlbumByIdRequest("/api/v1", createdAlbum.ID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.GetAlbumById(ginContext, createdAlbum.ID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse api.AlbumResponse
	err := json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("test", albumGetResponse.Title)
	suite.Assert().Equal("sports", string(albumGetResponse.Category.Name))
	suite.Assert().NotNil(albumGetResponse.ReleaseDate)
}

func (suite *AlbumControllersSuite) TestGetNoAlbumFailure() {
	doesNotExistAlbumID := 1111
	deletedAlbum, err := models.GetAlbum(doesNotExistAlbumID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedAlbum)

	request, _ := api.NewGetAlbumByIdRequest("/api/v1", doesNotExistAlbumID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.GetAlbumById(ginContext, doesNotExistAlbumID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse api.AlbumResponse
	err = json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
}

func (suite *AlbumControllersSuite) TestUpdate() {
	createdAlbum, _ := models.CreateAlbum("test", time.Now(), "sports")

	title := "updated"
	category := api.Category{
		Name: "food",
	}
	request, _ := api.NewUpdatedAlbumByIdRequest("/api/v1", createdAlbum.ID,
		api.UpdatedAlbumByIdJSONRequestBody{
			Title:    &title,
			Category: &category,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.UpdatedAlbumById(ginContext, createdAlbum.ID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse api.AlbumResponse
	err := json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("updated", albumGetResponse.Title)
	suite.Assert().Equal("food", string(albumGetResponse.Category.Name))
	suite.Assert().NotNil(albumGetResponse.ReleaseDate)
}

func (suite *AlbumControllersSuite) TestUpdateRequestBodyFailure() {
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PATCH", "/api/v1/album", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.albumHandler.CreateAlbum(ginContext)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *AlbumControllersSuite) TestUpdateNoAlbumFailure() {
	doesNotExistAlbumID := 1111
	deletedAlbum, err := models.GetAlbum(doesNotExistAlbumID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedAlbum)

	title := "updated"
	category := api.Category{
		Name: "food",
	}
	request, _ := api.NewUpdatedAlbumByIdRequest("/api/v1", doesNotExistAlbumID,
		api.UpdatedAlbumByIdJSONRequestBody{
			Title:    &title,
			Category: &category,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.UpdatedAlbumById(ginContext, doesNotExistAlbumID)
	bodyBytes, _ := io.ReadAll(w.Body)
	var albumGetResponse api.AlbumResponse
	err = json.Unmarshal(bodyBytes, &albumGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
}

func (suite *AlbumControllersSuite) TestUpdateFailure() {
	mockDB := suite.MockDB()

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `albums` WHERE `albums`.`id` = ? ORDER BY `albums`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("update error"))

	title := "updated"
	category := api.Category{
		Name: "food",
	}
	request, _ := api.NewUpdatedAlbumByIdRequest("/api/v1", 1,
		api.UpdatedAlbumByIdJSONRequestBody{
			Title:    &title,
			Category: &category,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.albumHandler.UpdatedAlbumById(ginContext, 1)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().True(strings.Contains(w.Body.String(), "update error"))
}

func (suite *AlbumControllersSuite) TestDelete() {
	createdAlbum, _ := models.CreateAlbum("test", time.Now(), "sports")

	request, _ := api.NewDeleteAlbumByIdRequest("/api/v1", createdAlbum.ID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.DeleteAlbumById(ginContext, createdAlbum.ID)
	suite.Assert().Equal(http.StatusNoContent, w.Code)

	deletedAlbum, err := models.GetAlbum(createdAlbum.ID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedAlbum)
}

func (suite *AlbumControllersSuite) TestDeleteNoAlbumFailure() {
	doesNotExistAlbumID := 1111
	deletedAlbum, err := models.GetAlbum(doesNotExistAlbumID)
	suite.Assert().NotNil(err)
	suite.Assert().Nil(deletedAlbum)

	request, _ := api.NewDeleteAlbumByIdRequest("/api/v1", doesNotExistAlbumID)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.DeleteAlbumById(ginContext, doesNotExistAlbumID)
	suite.Assert().Equal(http.StatusNoContent, w.Code)
}

func (suite *AlbumControllersSuite) TestDeleteAlbumFailure() {
	mockDB := suite.MockDB()

	mockDB.ExpectBegin()
	mockDB.ExpectExec("DELETE FROM `albums`").WillReturnError(errors.New("delete error"))
	mockDB.ExpectCommit()

	request, _ := api.NewDeleteAlbumByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request
	suite.albumHandler.DeleteAlbumById(ginContext, 1)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().True(strings.Contains(w.Body.String(), "delete error"))
}
