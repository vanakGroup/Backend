package postgres

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/logging"
	"gorm.io/gorm"
)

type UserRepository struct {
	cfg    *bootstrap.Config
	logger logging.Logger
	db     *gorm.DB
}

func NewUserRepository(cfg *bootstrap.Config, db *gorm.DB) *UserRepository {
	logger := logging.NewLogger(cfg)

	return &UserRepository{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (userRepository *UserRepository) FindUserByPhone(phone string) (*entities.User, error) {
	var foundUser entities.User
	result := userRepository.db.Where("phone = ?", phone).First(&foundUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &foundUser, nil
}

func (userRepository *UserRepository) Create(user *entities.User) error {
	return userRepository.db.Create(user).Error
}
