package interfaces

import (
	"context"
	"project/src/entities"
)

type UserRepository interface {
	GetById(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
}
