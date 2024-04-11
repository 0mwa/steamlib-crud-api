package service

import (
	"TestProject/internal/repository"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
)

var ErrUserExists = errors.New("User already exists")

type Auth struct {
	UsersRepo *repository.Users
	SessRepo  *repository.Sessions
	Logger    *zap.SugaredLogger
}

func NewAuthService(u *repository.Users, s *repository.Sessions, l *zap.SugaredLogger) *Auth {
	return &Auth{u, s, l}
}

func (a Auth) Auth(login string, passwd string) (string, error) {
	passwd = fmt.Sprintf("%x", sha256.Sum256([]byte(passwd)))
	userId, err := a.UsersRepo.GetUserId(login, passwd)
	if err != nil {
		return "", err
	}
	if userId == 0 {
		return "", err
	}
	var token string
	token, err = a.SessRepo.CheckSessionExp(userId)
	if err != nil {
		errstr := err.Error()
		if errstr == "no token" {
			token = strconv.Itoa(userId+rand.Int()) + randomdata.FirstName(randomdata.RandomGender)
			err = a.SessRepo.AddSession(userId, token)
			if err != nil {
				return "", err
			}
			return token, nil
		}
		return "", err
	}
	return token, nil
}

func (a Auth) CreateUser(login string, passwd string) error {
	ok, err := a.UsersRepo.CheckUserByLogin(login)
	if err != nil {
		return err
	}
	if ok {
		return ErrUserExists
	}
	passwd = fmt.Sprintf("%x", sha256.Sum256([]byte(passwd)))
	err = a.UsersRepo.CreateUser(login, passwd)
	if err != nil {
		return err
	}
	a.Logger.Infof("New user created with login : %s", login)
	return nil
}
