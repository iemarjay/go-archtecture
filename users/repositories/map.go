package repositories

import "archtecture/users/logic"

type Map struct {
	users []*logic.UserData
}

func (m *Map) Find(id string) (*logic.UserData, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, nil
}

func NewMap(users ...*logic.UserData) *Map {
	return &Map{users: users}
}

func (m *Map) StoreUser(data *logic.UserData) (*logic.UserData, error) {
	if err := data.EncryptPassword(); err != nil {
		return nil, err
	}
	m.users = append(m.users, data)
	return data, nil
}

func (m *Map) EmailOrPhoneExists(email string, phone string) (bool, error) {
	for _, user := range m.users {
		if user.Email == email || user.Phone == phone {
			return true, nil
		}
	}

	return false, nil
}

func (m *Map) FindByUniqueFields(field string) (*logic.UserData, error) {
	for _, user := range m.users {
		if user.Email == field || user.Phone == field {
			return user, nil
		}
	}

	return nil, nil
}
