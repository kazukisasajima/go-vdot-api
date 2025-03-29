package repository

import (
	"go_vdot_api/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	GetUserByID(user *model.User, userId uint) error
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userId uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetUserByID(user *model.User, userId uint) error {
	if err := ur.db.First(user, userId).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUser(user *model.User) error {
	result := ur.db.Model(user).Where("id=?", user.ID).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (ur *userRepository) DeleteUser(userId uint) error {
	if err := ur.db.Delete(&model.User{}, userId).Error; err != nil {
		return err
	}
	return nil
}
