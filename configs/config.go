package configs

import (
	"go-api-arch-mvc-template/pkg/logger"
	"os"
	"strconv"

	"go.uber.org/zap"
)

// 環境変数の設定を保持
type ConfigList struct {
	Env                 string
	DBHost              string
	DBPort              int
	DBDriver            string
	DBName              string
	DBUser              string
	DBPassword          string
	APICorsAllowOrigins []string
}

// 開発環境かどうかをチェック
func (c *ConfigList) IsDevelopment() bool {
	return c.Env == "development"
}

func init() {
	if err := LoadEnv(); err != nil {
		logger.Error("Failed to load env: ", zap.Error(err))
		panic(err)
	}
}

// 環境変数の設定を格納する
var Config ConfigList

// 環境変数の設定を読み込んでConfig変数に代入
func LoadEnv() error {
	DBPort, err := strconv.Atoi(GetEnvDefault("MYSQL_PORT", "3306"))
	if err != nil {
		return err
	}

	Config = ConfigList{
		Env:                 GetEnvDefault("APP_ENV", "development"),
		DBDriver:            GetEnvDefault("DB_DRIVER", "mysql"),
		DBHost:              GetEnvDefault("DB_HOST", "0.0.0.0"),
		DBPort:              DBPort,
		DBUser:              GetEnvDefault("DB_USER", "app"),
		DBPassword:          GetEnvDefault("DB_PASSWORD", "password"),
		DBName:              GetEnvDefault("DB_NAME", "api_database"),
		APICorsAllowOrigins: []string{"http://0.0.0.0:8001"},
	}
	return nil
}

// 環境変数の値を取得
func GetEnvDefault(key, defVal string) string {
	// 環境変数の値検索と存在確認
	val, err := os.LookupEnv(key)
	if !err {
		return defVal
	}
	return val
}
