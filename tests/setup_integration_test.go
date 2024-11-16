package integration

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"test-task/internal/config"
	"test-task/internal/handlers"
	"test-task/internal/repositories"
	"test-task/internal/router"
	"test-task/internal/services"
	"testing"
	"time"
)

var ginEngine *gin.Engine
var dbContainer testcontainers.Container

func setupRoutesForTests() *gin.Engine {

	gin.SetMode(gin.TestMode)
	engine := gin.New()

	cfg := config.Get()
	dbContext, err := repositories.NewDbContext(cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("setupRoutesForTests: error create dbContext: %v", err)
	}

	walletRepository := repositories.NewWalletsRepository(dbContext.DB)
	walletService := services.NewWalletsService(walletRepository)
	walletHandler := handlers.NewWalletHandler(walletService)

	router.Setup(engine, walletHandler)
	return engine
}

func upEnvironment() {

	ctx := context.Background()

	db := "test_db"
	user := "postgres"
	password := "postgres"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       db,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
	}

	var err error
	dbContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("could not start PostgreSQL container: %s", err)
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("could not get port for PostgreSQL container: %s", err)
	}

	conn := fmt.Sprintf("host=localhost dbname=%s port=%d user=%s password=%s sslmode=disable", db, port.Int(), user, password)
	err = os.Setenv("DB_CONNECTION_STRING", conn)

	if err != nil {
		log.Fatalf("could not set environment variable DB_CONNECTION_STRING: %s", err)
	}

	cfg := config.Get()
	dbContext, err := repositories.NewDbContext(cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("upEnvironment: error create dbContext: %v", err)
	}
	defer dbContext.Close()

	err = dbContext.Migrate()
	if err != nil {
		log.Fatalf("error run database migration: %v", err)
	}
	log.Infof("database migration complete")
}

func downEnvironment() {
	ctx := context.Background()
	if err := dbContainer.Terminate(ctx); err != nil {
		fmt.Printf("Could not terminate PostgreSQL container: %s", err)
	}
}

func TestMain(m *testing.M) {

	err := os.Chdir("../") //project root to resolve correctly relative paths in code
	if err != nil {
		log.Fatal(err)
	}

	upEnvironment()
	ginEngine = setupRoutesForTests()

	code := m.Run()

	downEnvironment()

	os.Exit(code)
}
