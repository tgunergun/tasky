package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeffthorne/tasky/auth"
	"github.com/jeffthorne/tasky/database"
	"github.com/jeffthorne/tasky/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTodo(c *gin.Context, db database.DBClient) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	id := c.Param("id")
	todo, err := db.GetTodo(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
	}

	defer cancel()
	c.JSON(http.StatusOK, todo)
}

func ClearAll(c *gin.Context, db database.DBClient) {
	session := auth.ValidateSession(c)
	if !session {
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := db.ClearTodos(ctx, c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "All todos deleted."})

}

func GetTodos(c *gin.Context, db database.DBClient) {
	session := auth.ValidateSession(c)
	if !session {
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	todos, err := db.GetTodos(ctx, c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"FindError": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func DeleteTodo(c *gin.Context, db database.DBClient) {
	session := auth.ValidateSession(c)
	if !session {
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	msg, err := db.DeleteTodo(ctx, c.Param("id"), c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": msg})

}

func UpdateTodo(c *gin.Context, db database.DBClient) {
	session := auth.ValidateSession(c)
	if !session {
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var newTodo models.Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.UpdateTodo(ctx, &newTodo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, newTodo)
}

func AddTodo(c *gin.Context, db database.DBClient) {
	session := auth.ValidateSession(c)
	if !session {
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var todo models.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.ID = primitive.NewObjectID()
	todo.UserID = c.Param("userid")

	err := db.AddTodo(ctx, &todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"insertedId": todo.ID})
}
