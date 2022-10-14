package logic

import (
	"archtecture/app/hash"
	"archtecture/app/validation"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUser_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	result := &UserData{Email: "test@test", Password: "password"}
	_ = result.EncryptPassword()
	input := &UserData{
		Lastname:  "Joseph",
		Firstname: "Emma",
		Phone:     "+234988998898",
		Email:     "test@test.test",
		Password:  "password",
	}

	repository := NewMockuserRepository(ctrl)
	repository.EXPECT().EmailOrPhoneExists(gomock.Any(), gomock.Any()).Return(false, nil).MinTimes(1)
	repository.EXPECT().StoreUser(gomock.Any()).Return(result, nil).MinTimes(1)

	event := NewMockevent(ctrl)
	event.EXPECT().Emit(gomock.Eq(UserRegisteredName), gomock.Any()).MinTimes(1)

	validator := validation.NewValidator()

	logic := NewUser(repository, validator, event)
	user, err := logic.Register(input)

	assert.Equal(t, user, result)
	assert.Equal(t, err, nil)
	assert.Equal(t, hash.Compare("password", result.Password), true)
}

func TestUser_CannotRegisterUserWhenValidationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	result := &UserData{Email: "test@test", Password: "password"}
	_ = result.EncryptPassword()
	input := &UserData{
		Lastname: "Joseph",
		Email:    "test@test.test",
		Password: "password",
	}

	repository := NewMockuserRepository(ctrl)
	repository.EXPECT().EmailOrPhoneExists(gomock.Any(), gomock.Any()).Return(false, nil).Times(0)
	repository.EXPECT().StoreUser(gomock.Any()).Return(result, nil).Times(0)

	event := NewMockevent(ctrl)
	event.EXPECT().Emit(gomock.Eq(UserRegisteredName), gomock.Any()).Times(0)

	validator := validation.NewValidator()

	logic := NewUser(repository, validator, event)
	user, err := logic.Register(input)

	assert.Equal(t, user, nil)
	assert.NotEqual(t, err, nil)
}

func TestUser_CannotRegisterUserWhenEmailOrPhoneExits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	result := &UserData{Email: "test@test", Password: "password"}
	_ = result.EncryptPassword()
	input := &UserData{
		Lastname:  "Joseph",
		Firstname: "Emma",
		Phone:     "+234988998898",
		Email:     "test@test.test",
		Password:  "password",
	}

	repository := NewMockuserRepository(ctrl)
	repository.EXPECT().EmailOrPhoneExists(gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
	repository.EXPECT().StoreUser(gomock.Any()).Return(result, nil).Times(0)

	event := NewMockevent(ctrl)
	event.EXPECT().Emit(gomock.Eq(UserRegisteredName), gomock.Any()).Times(0)

	validator := validation.NewValidator()

	logic := NewUser(repository, validator, event)
	user, err := logic.Register(input)

	assert.Equal(t, user, nil)
	assert.NotEqual(t, err, nil)
}

func TestUser_CannotRegisterUserWhenRepositoryReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	result := &UserData{Email: "test@test", Password: "password"}
	_ = result.EncryptPassword()
	input := &UserData{
		Lastname:  "Joseph",
		Firstname: "Emma",
		Phone:     "+234988998898",
		Email:     "test@test.test",
		Password:  "password",
	}

	err := errors.New("error")
	repository := NewMockuserRepository(ctrl)
	repository.EXPECT().EmailOrPhoneExists(gomock.Any(), gomock.Any()).Return(false, nil).MinTimes(1)
	repository.EXPECT().StoreUser(gomock.Any()).Return(result, err).MinTimes(1)

	event := NewMockevent(ctrl)
	event.EXPECT().Emit(gomock.Eq(UserRegisteredName), gomock.Any()).MinTimes(0)

	validator := validation.NewValidator()

	logic := NewUser(repository, validator, event)
	user, returnErr := logic.Register(input)

	assert.Equal(t, user, nil)
	assert.Equal(t, returnErr, err)
}
