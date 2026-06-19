package postgres

import (
	"context"

	"github.com/Vanak/Backend/internal/domain/entities"
)

type UserRepository interface {
	FindUserByPhone(phone string) (*entities.User, error)
	Create(user *entities.User) error
	FindUserById(userId uint) (*entities.User, error)
}

type UserProfileRepository interface {
	FindUserById(userId uint) (*entities.User, error)
	UpdateProfile(ctx context.Context, userId uint, updateData map[string]any) error
	CheckNationalCodeExists(ctx context.Context, nationalCodee string, excludeUserID uint) bool
	CheckShabaExists(ctx context.Context, shaba string, excludeUserID uint) bool
}
