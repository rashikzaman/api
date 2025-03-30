package integration_test

import (
	"log"
	"rashikzaman/api/application"
	"rashikzaman/api/config"
	"rashikzaman/api/db"
	"rashikzaman/api/models"
	"testing"
	_ "testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	application application.Application
}

type DummyData struct {
	User *models.User
}

func TestIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	suite.Run(t, new(TestSuite))
}

// run once, before test suite methods.
func (s *TestSuite) SetupSuite() {
	app := application.Application{}

	config, err := config.InitConfig("../../.env")
	if err != nil {
		log.Fatal("Failed reading config ", err)
	}

	db, err := db.InitDB(config.GetTestDBConfig())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	app.Config = config
	app.DB = db

	s.application = app
}

// run once, after test suite methods.
func (s *TestSuite) TearDownSuite() {
	s.T().Log("shutting off")
}
