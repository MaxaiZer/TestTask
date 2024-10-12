package integration

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"project/src/server"
	"testing"
)

var dbContainer testcontainers.Container

func setupRoutesForTests() *gin.Engine {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	server.InitializeRoutes(router)

	return router
}

func upEnvironment() {

	ctx := context.Background()

	db := "test_db"
	user := "postgres"
	password := "postgres"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       db,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	var err error
	dbContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Sprintf("Could not start PostgreSQL container: %s", err))
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(fmt.Sprintf("Could not get port for PostgreSQL container: %s", err))
	}

	conn := fmt.Sprintf("host=localhost port=%d dbname=%s user=%s password=%s", port.Int(), db, user, password)

	err = os.Setenv("DB_CONNECTION_STRING", conn)
	if err != nil {
		panic(fmt.Sprintf("Could not set environment variable DB_CONNECTION_STRING: %s", err))
	}
}

func downEnvironment() {
	ctx := context.Background()

	if err := dbContainer.Terminate(ctx); err != nil {
		fmt.Printf("Could not terminate PostgreSQL container: %s", err)
	}
}

func TestMain(m *testing.M) {

	err := os.Chdir("../../") //because config file in the root folder
	if err != nil {
		log.Fatal(err)
	}

	upEnvironment()

	code := m.Run()

	downEnvironment()

	os.Exit(code)
}
