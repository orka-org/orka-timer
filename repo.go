package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db interface {
	CreateTimer(timer Timer) (string, error)
	GetTimer(id string) (*Timer, error)
	UpdateTimer(id string, timer Timer) error
	DeleteTimer(id string) error
	ListTimers() ([]Timer, error)
}

type Timer struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	TimeInterval TimeInterval       `bson:"time_interval"`
	Pauses       []TimeInterval     `bson:"pauses"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongo(client *mongo.Client, database, collection string) *Mongo {
	return &Mongo{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (m *Mongo) Collection() *mongo.Collection {
	return m.client.Database(m.database).Collection(m.collection)
}

func (m *Mongo) CreateTimer(timer Timer) (string, error) {
	log.Println("CREATE TIMER ")
	timer.CreatedAt = time.Now()
	timer.UpdatedAt = time.Now()

	log.Println("CREAT TIMER ...")
	result, err := m.Collection().InsertOne(context.Background(), timer)
	log.Println("CREATE TIMER ", "result", result, "err", err)
	if err != nil {
		return "", err
	}
	log.Println("CREATE TIMER ", "result", result, "err", err)

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *Mongo) GetTimer(id string) (*Timer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var timer Timer
	err = m.Collection().FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&timer)
	if err != nil {
		return nil, err
	}

	return &timer, nil
}

func (m *Mongo) UpdateTimer(id string, timer Timer) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	timer.UpdatedAt = time.Now()
	update := bson.M{
		"$set": timer,
	}

	_, err = m.Collection().UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)
	return err
}

func (m *Mongo) DeleteTimer(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Collection().DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (m *Mongo) ListTimers() ([]Timer, error) {
	fmt.Println(m.Collection().Name())
	cursor, err := m.Collection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var timers []Timer
	for cursor.Next(context.Background()) {
		val := Timer{}
		if err := cursor.Decode(&val); err != nil {
			return []Timer{}, err
		}
		timers = append(timers, val)
	}
	if err = cursor.All(context.Background(), &timers); err != nil {
		fmt.Println("ERRORR", err)
		return nil, err
	}
	fmt.Println("ALL GOOD", timers)

	return timers, nil
}
