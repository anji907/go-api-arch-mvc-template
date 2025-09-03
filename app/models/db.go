package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-api-arch-mvc-template/configs"
)

const (
	InstanceSqlLite int = iota
	InstanceMySQL
)

var (
	DB                            *gorm.DB
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

// modelの一覧をDBから返す
func GetModels() []interface{} {
	return []interface{}{&Album{}, &Category{}}
}

/*
引数で指定されたデータベースのインスタンスに応じて
データベースに接続するための *gorm.DB型の値を作成
*/
func NewDatabaseSQLFactory(instance int) (db *gorm.DB, err error) {
	switch instance {
	// MySQLへの接続処理
	case InstanceMySQL:
		// data source nameを作成
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			configs.Config.DBUser,
			configs.Config.DBPassword,
			configs.Config.DBHost,
			configs.Config.DBPort,
			configs.Config.DBName)
		// データベースに接続
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// SQLiteへの接続処理
	case InstanceSqlLite:
		db, err = gorm.Open(sqlite.Open(configs.Config.DBName),
			&gorm.Config{})
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
	return db, err
}

/*
引数で受け取ったインスタンスをNewDatabaseSQLFactory関数に
渡してグローバル変数DBに接続しているデータベースの情報を格納
*/
func SetDatabase(instance int) (err error) {
	db, err := NewDatabaseSQLFactory(instance)
	if err != nil {
		return err
	}
	DB = db
	return err
}
