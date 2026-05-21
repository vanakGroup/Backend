package postgres

import (
	"context"

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

func (ur *UserRepository) UpdateProfile(ctx context.Context, userId uint, updateData map[string]any) error {

	result := ur.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userId).
		Updates(updateData)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ur *UserRepository) CheckNationalCodeExists(ctx context.Context, nationalCodee string, excludeUserId uint) bool {
	var Count int64
	ur.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("national_code = ?", nationalCodee).
		Where("id != ?", excludeUserId).
		Count(&Count)
	return Count > 0
}

func (ur *UserRepository) CheckShabaExists(ctx context.Context, shaba string, excludeUserId uint) bool {
	var Count int64
	ur.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("shaba_number = ?", shaba).
		Where("id != ?", excludeUserId).
		Count(&Count)
	return Count > 0
}

func (ur *UserRepository) FindUserById(userId uint) (*entities.User, error) {
	var user entities.User
	err := ur.db.First(&user, userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
