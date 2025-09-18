package tester

import (
	"context"
	"fmt"
	"go-api-arch-mvc-template/app/models"
	"go-api-arch-mvc-template/configs"
	"log"
	"net"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBMySQLSuite struct {
	suite.Suite
	mySQLContainer testcontainers.Container
	ctx            context.Context
}

func CheckPort(host string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if conn != nil {
		conn.Close()
		return false
	}
	if err != nil {
		return true
	}
	return false
}

// ポートが使用可能になるまで待機する
func WaitForPort(host string, port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if CheckPort(host, port) {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

// テスト用のMySQLを立ち上げる
func (suite *DBMySQLSuite) SetupTestContainers() (err error) {
	WaitForPort(configs.Config.DBHost, configs.Config.DBPort, 10*time.Second)
	suite.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image: "mysql:8",
		Env: map[string]string{
			"MYSQL_DATABASE":             configs.Config.DBName,
			"MYSQL_USER":                 configs.Config.DBUser,
			"MYSQL_PASSWORD":             configs.Config.DBPassword,
			"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{fmt.Sprintf("%d:3306/tcp", configs.Config.DBPort)},
		WaitingFor:   wait.ForLog("port: 3306 MySQL Community Server"),
	}

	suite.mySQLContainer, err = testcontainers.GenericContainer(
		suite.ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}

func (suite *DBMySQLSuite) SetupSuite() {
	err := suite.SetupTestContainers()
	suite.Assert().Nil(err)

	err = models.SetDatabase(models.InstanceMySQL)
	suite.Assert().Nil(err)

	for _, model := range models.GetModels() {
		err := models.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBMySQLSuite) TearDownSuite() {
	if suite.mySQLContainer == nil {
		return
	}
	err := suite.mySQLContainer.Terminate(suite.ctx)
	suite.Assert().Nil(err)
}
