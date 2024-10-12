package repositories

import (
	"context"
	"errors"
	"project/src/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(connectionString string) (*PostgresUserRepository, error) {

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = tryMigrate(db); err != nil {
		return nil, err
	}

	return &PostgresUserRepository{db: db}, nil
}

func (repo *PostgresUserRepository) GetById(ctx context.Context, id string) (*entities.User, error) {

	var user entities.User
	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	return repo.db.WithContext(ctx).Save(user).Error
}

func tryMigrate(db *gorm.DB) error {

	if db.Migrator().HasTable(&entities.User{}) {
		return nil
	}

	if err := db.AutoMigrate(&entities.User{}); err != nil {
		return err
	}

	user1 := entities.User{
		ID:    "1",
		Email: "example@mail.com",
	}

	user2 := entities.User{
		ID:    "2",
		Email: "example2@mail.com",
	}

	_ = db.Create(&user1)
	_ = db.Create(&user2)
	return nil
}
