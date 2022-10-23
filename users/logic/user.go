package logic

import (
	appEvent "archtecture/app/events"
	"archtecture/app/hash"
	"archtecture/app/validation"
)

type userRepository interface {
	StoreUser(*UserData) (*UserData, error)
	EmailOrPhoneExists(string, string) (bool, error)
}

type event interface {
	Emit(appEvent.Name, interface{})
}

type UserData struct {
	ID        string `json:"id,omitempty"`
	Lastname  string `json:"lastname,omitempty" validate:"required"`
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Phone     string `json:"phone" validate:"required,e164"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" validate:"required"`
}

func (u *UserData) PasswordMatch(password string) bool {
	return hash.Compare(password, u.Password)
}

func (u *UserData) EncryptPassword() error {
	var err error
	u.Password, err = hash.Generate(u.Password)
	return err
}

func (u *UserData) RouteNotificationForMail() string {
	return u.Email
}

func (u *UserData) RouteNotificationForSms() string {
	return u.Phone
}

func (u *UserData) Key() string {
	return u.ID
}

func (u *UserData) GetFirstname() string {
	return u.Firstname
}

type User struct {
	repository userRepository
	validator  validator
	event      event
}

func NewUser(
	repository userRepository,
	validator validator,
	event event) *User {
	return &User{
		repository: repository,
		validator:  validator,
		event:      event,
	}
}

func (u *User) Register(input *UserData) (*UserData, error) {
	if fails, err := u.validator.Validate(input); fails {
		return nil, err
	}
	if err := input.EncryptPassword(); err != nil {
		return nil, err
	}
	err := u.userWithEmailOrPhoneExists(input)
	if err != nil {
		return nil, err
	}

	user, err := u.repository.StoreUser(input)
	if err != nil {
		return nil, err
	}

	u.event.Emit(UserRegisteredEvent, NewUserRegistered(user))

	return user, nil
}

func (u *User) userWithEmailOrPhoneExists(input *UserData) error {
	exists, err := u.repository.EmailOrPhoneExists(input.Email, input.Phone)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	errBag := validation.NewErrorBag()
	errBag.Add(&validation.Error{
		Namespace: "UserData.Email",
		Field:     "Email",
		Tag:       "unique",
		Value:     input.Email,
	})
	errBag.Add(&validation.Error{
		Namespace: "UserData.Phone",
		Field:     "Phone",
		Tag:       "unique",
		Value:     input.Phone,
	})
	return errBag
}
