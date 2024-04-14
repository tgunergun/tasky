package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jeffthorne/tasky/models"
	"github.com/joho/godotenv"
)

type DBClient interface {
	FindExistingUsers(ctx context.Context, user models.User) (int64, error)
	AddUser(ctx context.Context, user *models.User) ([]byte, error)
	GetUser(ctx context.Context, email string) (models.User, error)
	GetTodo(ctx context.Context, id string) (models.Todo, error)
	GetTodos(ctx context.Context, id string) ([]models.Todo, error)
	ClearTodos(ctx context.Context, userid string) error
	DeleteTodo(ctx context.Context, id string, userid string) (string, error)
	UpdateTodo(ctx context.Context, newTodo *models.Todo) error
	AddTodo(ctx context.Context, todo *models.Todo) error
	Close() error
}

func CreateDBClientFromEnv() DBClient {
	godotenv.Overload()
	dbType := strings.ToLower(os.Getenv("DB_TYPE"))
	if len(dbType) == 0 {
		dbType = "mongodb"
	}
	switch dbType {
	case "mongodb":
		return NewMongoDBClient()
	case "postgresql":
		return NewPostgresDBClient()
	default:
		panic(fmt.Sprintf("This database type is unsupported: %s", dbType))
	}
}
