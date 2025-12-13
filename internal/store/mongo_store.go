package store

import (
	"context"
	"log"
	"time"

	"github.com/zeltbrennt/go-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
}

type mongoTask struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title"`
}

func taskToMongo(t models.Task) mongoTask {
	mt := mongoTask{
		Title: t.Title,
	}
	if t.ID != 0 {
		mt.ID = primitive.NewObjectIDFromTimestamp(time.Unix(int64(t.ID), 0))
	}
	return mt
}

func mongoToTask(m mongoTask) models.Task {
	return models.Task{
		ID:    int(m.ID.Timestamp().Unix()),
		Title: m.Title,
	}
}

func NewMongoStore(uri, dbname, collectionName string) (*mongoStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbname).Collection(collectionName)
	return &mongoStore{
		client:     client,
		collection: collection,
		ctx:        context.Background(),
	}, nil
}

func (m *mongoStore) GetAllTasks() ([]models.Task, error) {
	cursor, err := m.collection.Find(m.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		err = cursor.Close(m.ctx)
		if err != nil {
			log.Print("curser failed to close:", err)
		}
	}()

	var tasks []models.Task
	for cursor.Next(m.ctx) {
		var mt mongoTask
		if err := cursor.Decode(&mt); err != nil {
			return nil, err
		}
		tasks = append(tasks, mongoToTask(mt))
	}
	return tasks, nil
}

func (m *mongoStore) CreateTask(t models.Task) (models.Task, error) {
	mt := taskToMongo(t)
	res, err := m.collection.InsertOne(m.ctx, mt)
	if err != nil {
		return models.Task{}, err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	t.ID = int(oid.Timestamp().Unix())
	return t, nil
}
