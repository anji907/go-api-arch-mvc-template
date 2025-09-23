package tester

import (
	"go-api-arch-mvc-template/pkg/logger"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mockClock.Now()で固定された時刻を返す
type mockClock struct {
	t time.Time
}

func NewMockClock(t time.Time) mockClock {
	return mockClock{t}
}

func (m mockClock) Now() time.Time {
	return m.t
}

// データベースの操作をモック化する
func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	// モックデータベースの作成
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		logger.Fatal(err.Error())
	}

	// モックデータベースに接続
	mockGormDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "mock_db",
		DriverName:                "mysql",
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mock, mockGormDB
}
