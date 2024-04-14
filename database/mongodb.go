package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jeffthorne/tasky/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client         *mongo.Client
	userCollection *mongo.Collection
	todoCollection *mongo.Collection
}

func (m *MongoDBClient) FindExistingUsers(ctx context.Context, user models.User) (int64, error) {
	return m.userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
}

func (m *MongoDBClient) AddUser(ctx context.Context, user *models.User) ([]byte, error) {
	user.ID = primitive.NewObjectID()
	resultInsertionNumber, insertErr := m.userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		return []byte{}, fmt.Errorf("user item was not created")
	}
	result, err := json.Marshal(&resultInsertionNumber)
	if err != nil {
		return []byte{}, err
	}
	user.ID = primitive.NewObjectID()
	return result, nil
}

func (m *MongoDBClient) GetUser(ctx context.Context, email string) (models.User, error) {
	var found models.User
	if err := m.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&found); err != nil {
		return found, err
	}
	return found, nil
}

func (m *MongoDBClient) GetTodo(ctx context.Context, id string) (models.Todo, error) {
	var todo models.Todo
	objId, _ := primitive.ObjectIDFromHex(id)
	if err := m.todoCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&todo); err != nil {
		return todo, err
	}
	return todo, nil
}

func (m *MongoDBClient) GetTodos(ctx context.Context, userid string) ([]models.Todo, error) {
	var todos []models.Todo
	findResult, err := m.todoCollection.Find(ctx, bson.M{"userid": userid})
	if err != nil {
		return todos, err
	}
	for findResult.Next(ctx) {
		var todo models.Todo
		err := findResult.Decode(&todo)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (m *MongoDBClient) ClearTodos(ctx context.Context, userid string) error {
	_, err := m.todoCollection.DeleteMany(ctx, bson.M{"userid": userid})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDBClient) DeleteTodo(ctx context.Context, id string, userid string) (string, error) {
	objId, _ := primitive.ObjectIDFromHex(id)
	deleteResult, err := m.todoCollection.DeleteOne(ctx, bson.M{"_id": objId, "userid": userid})
	if err != nil {
		return "", err
	}
	if deleteResult.DeletedCount == 0 {
		return "", fmt.Errorf("No todo with id : %v was found, no deletion occurred.", id)
	}
	return fmt.Sprintf("todo with id: %v was deleted successfully.", id), nil
}

func (m *MongoDBClient) UpdateTodo(ctx context.Context, newTodo *models.Todo) error {
	_, err := m.todoCollection.UpdateOne(ctx, bson.M{"_id": newTodo.ID, "userid": newTodo.UserID}, bson.M{"$set": newTodo})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDBClient) AddTodo(ctx context.Context, userid string) (todo models.Todo, err error) {
	todo.ID = primitive.NewObjectID()
	_, err = m.todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return todo, err
	}
	return todo, nil
}

func NewMongoDBClient() *MongoDBClient {
	m := MongoDBClient{}
	m.client = CreateMongoClient()
	m.todoCollection = OpenCollection(m.client, "todos")
	m.userCollection = OpenCollection(m.client, "users")
	return &m
}

var Client *mongo.Client = CreateMongoClient()

func CreateMongoClient() *mongo.Client {
	godotenv.Overload()
	MongoDbURI := os.Getenv("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDbURI))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer cancel()
	fmt.Println("Connected to MONGO -> ", MongoDbURI)
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	MongoDbDatabaseName := os.Getenv("MONGODB_DB")
	if len(MongoDbDatabaseName) == 0 {
		MongoDbDatabaseName = "go-mongodb"
	}
	return client.Database(MongoDbDatabaseName).Collection(collectionName)
}
