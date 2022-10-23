package logic

import (
	"archtecture/app/events"
)

const UserRegisteredEvent events.Name = "UserRegistered"

type UserRegistered struct {
	user *UserData
}

func NewUserRegistered(user *UserData) *UserRegistered {
	return &UserRegistered{user: user}
}

func (u *UserRegistered) User() *UserData {
	return u.user
}
