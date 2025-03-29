package usecase

import (
	"go_vdot_api/model"
	"go_vdot_api/repository"
	"go_vdot_api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	LogIn(user model.User) (string, error)

	UpdateUser(user model.User) (model.UserResponse, error)
	DeleteUser(userId uint) error
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) LogIn(user model.User) (string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"name":    storedUser.Name,
		"email":   storedUser.Email,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, errerr := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if errerr != nil {
		return "", errerr
	}
	return tokenString, nil
}

func (uu *userUsecase) UpdateUser(user model.User) (model.UserResponse, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByID(&storedUser, user.ID); err != nil {
		return model.UserResponse{}, err
	}

	// 名前・メール・パスワードを個別に更新
	if user.Name != "" {
		storedUser.Name = user.Name
	}
	if user.Email != "" {
		storedUser.Email = user.Email
	}
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return model.UserResponse{}, err
		}
		storedUser.Password = string(hash)
	}

	// ユーザー情報更新
	if err := uu.ur.UpdateUser(&storedUser); err != nil {
		return model.UserResponse{}, err
	}

	// 更新後のユーザー情報を返す
	resUser := model.UserResponse{
		ID:    storedUser.ID,
		Name:  storedUser.Name,
		Email: storedUser.Email,
	}
	return resUser, nil
}


func (uu *userUsecase) DeleteUser(userId uint) error {
	if err := uu.ur.DeleteUser(userId); err != nil {
		return err
	}
	return nil
}
