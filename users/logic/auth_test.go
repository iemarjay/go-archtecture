package logic_test

import (
	"archtecture/users/logic"
	"archtecture/users/repositories"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestAuth_AttemptLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	phone := "+234808894"
	email := "auth@example.test"
	user := &logic.UserData{
		ID:        "",
		Lastname:  "",
		Firstname: "",
		Phone:     phone,
		Email:     email,
		Password:  "secret",
	}
	_ = user.EncryptPassword()

	identifiers := []string{email, phone}

	for _, identifier := range identifiers {
		repository := logic.NewMockauthRepository(ctrl)
		repository.EXPECT().FindByUniqueFields(gomock.Eq(identifier)).Return(user, nil).Times(1)

		auth := logic.NewAuth(repository)

		userData, _ := auth.AttemptLogin(map[string]string{
			"identifier": identifier,
			"password":   "secret",
		})

		assert.NotEqual(t, userData, nil)
		assert.Equal(t, userData.Email, email)
		assert.Equal(t, userData.Phone, phone)
	}
}

func TestAuth_AttemptLoginWithWrongCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	inputs := [][]string{
		{"auth@example.test", "secret"},
		{"+23480889424323", "secret"},
		{"+234808894", "password"},
		{"auth@example.test", "password"},
	}
	phone := "+234808894"
	email := "auth@example.test"
	user := &logic.UserData{
		ID:        "",
		Lastname:  "",
		Firstname: "",
		Phone:     phone,
		Email:     email,
		Password:  "secret",
	}
	_ = user.EncryptPassword()

	for _, data := range inputs {
		repository := logic.NewMockauthRepository(ctrl)
		repository.EXPECT().FindByUniqueFields(gomock.Eq(data[0])).Return(nil, nil).Times(1)

		auth := logic.NewAuth(repository)
		userData, err := auth.AttemptLogin(map[string]string{
			"identifier": data[0],
			"password":   data[1],
		})

		assert.NotEqual(t, err, nil)
		assert.Equal(t, err, logic.UserNotFound)
		assert.Equal(t, userData, nil)
	}
}

func makeMapRepository() *repositories.Map {
	user := &logic.UserData{
		ID:        "",
		Lastname:  "",
		Firstname: "",
		Phone:     "+234808894",
		Email:     "auth@example.test",
		Password:  "secret",
	}
	_ = user.EncryptPassword()

	return repositories.NewMap(user)
}
