package database

import (
	"archtecture/app/env"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Cursor interface {
	Next(ctx context.Context) bool
	Decode(v interface{}) error
	Err() error
}

type Decoder interface {
	Decode(v interface{}) error
}

type MongoDatabaseConfig struct {
	DatabaseUrl  string
	DatabaseName string
}

func NewMongoDatabaseConfigFromEnv(e *env.Env) *MongoDatabaseConfig {
	return &MongoDatabaseConfig{
		DatabaseUrl:  e.DatabaseUrl,
		DatabaseName: e.DatabaseName,
	}
}

type MongoDatabase struct {
	connection *mongo.Database
	collection *mongo.Collection
}

func NewMongoDatabaseFromConfig(config *MongoDatabaseConfig) *MongoDatabase {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client()
	clientOptions = clientOptions.ApplyURI(config.DatabaseUrl)
	clientOptions = clientOptions.SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(config.DatabaseName)

	return NewMongoDatabase(database)
}

func NewMongoDatabase(connection *mongo.Database) *MongoDatabase {
	return &MongoDatabase{connection: connection}
}

func (db *MongoDatabase) Collection() *mongo.Collection {
	return db.collection
}

func (db *MongoDatabase) Table(name string) *MongoDatabase {
	db.collection = db.connection.Collection(name)

	return db
}

func (db *MongoDatabase) Insert(data interface{}) (Decoder, error) {
	ior, err := db.collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	insertedID := ior.InsertedID.(primitive.ObjectID)

	return db.findOneByID(insertedID), err
}

func (db *MongoDatabase) FindOneByID(id string) Decoder {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return db.findOneByID(objectId)
}

func (db *MongoDatabase) findOneByID(id primitive.ObjectID) Decoder {
	result := db.collection.FindOne(context.Background(), bson.M{"_id": id})
	return result
}

func (db *MongoDatabase) UpdateOneById(id string, data interface{}) (Decoder, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return db.UpdateOne(bson.M{"_id": objectId}, data)
}

func (db *MongoDatabase) UpdateOne(filter bson.M, data interface{}) (Decoder, error) {
	result, err := db.collection.UpdateOne(context.Background(), filter, data)

	_ = result.UpsertedID
	if err != nil {
		return nil, err
	}

	return db.FindOne(filter), nil
}

func (db *MongoDatabase) FindOne(filter interface{}) Decoder {
	result := db.collection.FindOne(context.Background(), filter)
	return result
}

func (db *MongoDatabase) Find(filter interface{}) (Cursor, error) {
	return db.collection.Find(context.Background(), filter)
}

func (db *MongoDatabase) Count(filter interface{}) (int64, error) {
	return db.collection.CountDocuments(context.Background(), filter)
}
