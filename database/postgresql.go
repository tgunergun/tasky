package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeffthorne/tasky/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostgresDBClient struct {
	conn *pgxpool.Pool
	ctx  context.Context
}

func (p *PostgresDBClient) FindExistingUsers(ctx context.Context, user models.User) (int64, error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM \"users\" WHERE email = $1", user.Email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var rowCount int64
	for rows.Next() {
		rowCount++
	}
	return rowCount, nil
}

func (p *PostgresDBClient) AddUser(ctx context.Context, user *models.User) ([]byte, error) {
	if _, err := p.conn.Query(ctx, "SELECT id FROM \"users\""); err != nil {
		return nil, err
	}

	user.ID = primitive.NewObjectID()
	_, err := p.conn.Exec(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.ID, user.Email, user.Password)
	if err != nil {
		return []byte{}, err
	}
	return []byte{}, nil
}

func (p *PostgresDBClient) GetUser(ctx context.Context, email string) (models.User, error) {
	var found models.User
	if err := p.conn.QueryRow(ctx, "SELECT * FROM \"users\" WHERE email = $1", email).Scan(&found); err != nil {
		return found, err
	}
	return found, nil
}

func (p *PostgresDBClient) GetTodo(ctx context.Context, id string) (models.Todo, error) {
	var todo models.Todo
	if err := p.conn.QueryRow(ctx, "SELECT * FROM todos WHERE id = $1", id).Scan(&todo); err != nil {
		return todo, err
	}
	return todo, nil
}

func (p *PostgresDBClient) GetTodos(ctx context.Context, userid string) ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := p.conn.Query(ctx, "SELECT * FROM todos WHERE userid = $1", userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (p *PostgresDBClient) ClearTodos(ctx context.Context, userid string) error {
	if _, err := p.conn.Exec(ctx, "DELETE FROM todos WHERE userid = $1", userid); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDBClient) DeleteTodo(ctx context.Context, id string, userid string) (string, error) {
	if _, err := p.conn.Exec(ctx, "DELETE FROM todos WHERE userid = $1 AND id = $2", userid, id); err != nil {
		return "", err
	}
	return "", nil
}

func (p *PostgresDBClient) UpdateTodo(ctx context.Context, newTodo *models.Todo) error {
	if _, err := p.conn.Exec(ctx, "UPDATE todos SET name = $1, status = $2 WHERE userid = $3 AND id = $4", newTodo.Name, newTodo.Status, newTodo.UserID, newTodo.ID); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDBClient) AddTodo(ctx context.Context, todo *models.Todo) (err error) {
	todo.ID = primitive.NewObjectID()
	if err != nil {
		return err
	}
	_, err = p.conn.Exec(ctx, "INSERT INTO todos (id, name, status, userid) VALUES ($1, $2, $3, $4)", todo.ID, todo.Name, todo.Status, todo.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDBClient) Close() error {
	p.conn.Close()
	return nil
}

func NewPostgresDBClient() *PostgresDBClient {
	p := PostgresDBClient{ctx: context.Background()}
	var err error
	p.ctx = context.Background()
	p.conn, err = pgxpool.New(p.ctx, os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}
	if err = createTableIfNotExist(&p, "users", "id text, email text, password text"); err != nil {
		panic(err)
	}
	if err := createTableIfNotExist(&p, "todos", "id text, name text, status text, userid text"); err != nil {
		panic(err)
	}
	return &p
}

func createTableIfNotExist(p *PostgresDBClient, tableName string, cols string) error {
	var exists bool
	if err := p.conn.QueryRow(p.ctx, "SELECT EXISTS (SELECT FROM pg_tables WHERE tablename = $1)", tableName).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}
	query := fmt.Sprintf("CREATE TABLE %s (%s)", pgx.Identifier.Sanitize([]string{tableName}), cols)
	if _, err := p.conn.Exec(p.ctx, query); err != nil {
		return err
	}
	return nil
}
