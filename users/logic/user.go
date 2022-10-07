package logic

import "archtecture/app/hash"

type UserData struct {
	Password  string `json:"password,omitempty"`
	ID        string `json:"id,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func (u *UserData) PasswordMatch(password string) bool {
	return hash.Compare(password, u.Password)
}

func (u *UserData) Key() string {
	return u.ID
}
