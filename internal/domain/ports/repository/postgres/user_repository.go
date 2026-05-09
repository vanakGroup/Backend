package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type UserRepository interface {
	FindUserByPhone(phone string) (*entities.User, error)
	Create(user *entities.User) error
}
