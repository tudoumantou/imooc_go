package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"shop/datamodels"
	"shop/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewService(repositories repositories.IUserRepository) IUserService {
	return &UserService{UserRepository: repositories}
}

func (u *UserService)IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService)AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword=string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string)([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil{
		return false, errors.New("密码对比错误")
	}
	return true, nil
}