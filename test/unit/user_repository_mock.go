package unit

import (
	"context"
	"fmt"
	"project/src/entities"
)

type MockUserRepository struct {
	Users []entities.User
}

func NewMockUserRepository(users []entities.User) *MockUserRepository {
	return &MockUserRepository{Users: users}
}

func (r MockUserRepository) GetById(_ context.Context, id string) (*entities.User, error) {

	for _, user := range r.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, nil
}

func (r MockUserRepository) Update(_ context.Context, user *entities.User) error {

	for i := 0; i < len(r.Users); i++ {
		if r.Users[i].ID == user.ID {
			r.Users[i] = *user
			return nil
		}
	}

	return fmt.Errorf("user not found")
}
