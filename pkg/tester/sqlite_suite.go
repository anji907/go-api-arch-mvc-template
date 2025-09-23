package tester

import (
	"os"

	"github.com/stretchr/testify/suite"

	"go-api-arch-mvc-template/app/models"
	"go-api-arch-mvc-template/configs"
)

type DBSQLiteSuite struct {
	suite.Suite
}

// テスト前に自動で実行されるsuiteパッケージの機能
func (suite *DBSQLiteSuite) SetupSuite() {
	configs.Config.DBName = "unittest.sqlite"
	err := models.SetDatabase(models.InstanceSqlLite)
	suite.Assert().Nil(err)

	for _, model := range models.GetModels() {
		// モデルに応じたテーブルを作成
		err := models.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

// テスト後に実行されるsuiteパッケージの機能
func (suite *DBSQLiteSuite) TearDownSuite() {
	// データベースファイルを削除
	err := os.Remove(configs.Config.DBName)
	suite.Assert().Nil(err)
}
