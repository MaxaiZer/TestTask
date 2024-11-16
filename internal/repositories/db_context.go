package repositories

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type DbContext struct {
	DB *sqlx.DB
}

func NewDbContext(connectionString string) (*DbContext, error) {

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	return &DbContext{DB: db}, nil
}

func (c *DbContext) Migrate() error {
	return runMigrations(c.DB.DB)
}

func (c *DbContext) Close() error {
	return c.DB.Close()
}
