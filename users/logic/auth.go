package logic

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepository interface {
	FindByUniqueFields(field string) (*UserData, error)
}

var UserNotFound = errors.New("user with credentials not found")

type Auth struct {
	repository authRepository
}

func NewAuth(repository authRepository) *Auth {
	return &Auth{repository: repository}
}

func (a *Auth) AttemptLogin(input map[string]string) (*UserData, error) {
	user, err := a.repository.FindByUniqueFields(input["identifier"])
	if err == mongo.ErrNoDocuments || user == nil {
		return nil, UserNotFound
	}
	if err != nil {
		return nil, err
	}

	if !user.PasswordMatch(input["password"]) {
		return nil, UserNotFound
	}

	return user, nil
}
