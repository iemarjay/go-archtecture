package repositories

import (
	"archtecture/app/database"
	"archtecture/users/logic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TableName = "users"

type db interface {
	FindOne(filter interface{}) database.Decoder
	FindOneByID(id string) database.Decoder
	Insert(data interface{}) (database.Decoder, error)
	Count(filter interface{}) (int64, error)
}

type mongoData struct {
	ID        string `bson:"_id,omitempty"`
	Firstname string `bson:"firstname"`
	Lastname  string `bson:"lastname"`
	Email     string `bson:"email"`
	Phone     string `bson:"phone"`
	Password  string `bson:"password"`
}

func newMongoData(d *logic.UserData) *mongoData {
	return &mongoData{
		ID:        d.ID,
		Firstname: d.Firstname,
		Lastname:  d.Lastname,
		Email:     d.Email,
		Phone:     d.Phone,
		Password:  d.Password,
	}
}

func (d *mongoData) toUserData() *logic.UserData {
	return &logic.UserData{
		ID:        d.ID,
		Firstname: d.Firstname,
		Lastname:  d.Lastname,
		Email:     d.Email,
		Phone:     d.Phone,
		Password:  d.Password,
	}
}

type Mongo struct {
	db db
}

func NewMongo(db db) *Mongo {
	return &Mongo{db: db}
}

func (m *Mongo) StoreUser(data *logic.UserData) (*logic.UserData, error) {
	d := newMongoData(data)
	decoder, err := m.db.Insert(d)
	if err != nil {
		return nil, err
	}

	user := &mongoData{}
	if err = decoder.Decode(user); err != nil {
		return nil, err
	}

	return user.toUserData(), nil
}

func (m *Mongo) Find(id string) (*logic.UserData, error) {
	user := &mongoData{}
	err := m.db.FindOneByID(id).Decode(user)
	if err != nil {
		return nil, err
	}

	return user.toUserData(), nil
}

func (m *Mongo) FindByUniqueFields(field string) (*logic.UserData, error) {
	user := &mongoData{}
	objectID, _ := primitive.ObjectIDFromHex(field)
	filter := bson.M{"$or": []bson.M{
		{"id": objectID},
		{"email": field},
		{"phone": field},
	}}
	err := m.db.FindOne(filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user.toUserData(), nil
}

func (m *Mongo) EmailOrPhoneExists(email string, phone string) (bool, error) {
	filter := bson.M{"$or": []bson.M{
		{"email": email},
		{"phone": phone},
	}}

	count, err := m.db.Count(filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
