package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	cfg *bootstrap.Config
	db  *gorm.DB
}

func NewUserRepository(cfg *bootstrap.Config, db *gorm.DB) *UserRepository {

	return &UserRepository{
		cfg: cfg,
		db:  db,
	}
}

func (ur *UserRepository) FindUserByPhone(phone string) (*entities.User, error) {
	var foundUser entities.User
	result := ur.db.Where("phone = ?", phone).First(&foundUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &foundUser, nil
}

func (ur *UserRepository) Create(user *entities.User) error {
	return ur.db.Create(user).Error
}
