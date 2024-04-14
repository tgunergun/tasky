package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	controller "github.com/jeffthorne/tasky/controllers"
	"github.com/jeffthorne/tasky/database"
	"github.com/joho/godotenv"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

type ControllerMap struct {
	db database.DBClient
}

func (m *ControllerMap) GetTodos(c *gin.Context) {
	controller.GetTodos(c, m.db)
}

func (m *ControllerMap) GetTodo(c *gin.Context) {
	controller.GetTodo(c, m.db)
}

func (m *ControllerMap) AddTodo(c *gin.Context) {
	controller.AddTodo(c, m.db)
}

func (m *ControllerMap) DeleteTodo(c *gin.Context) {
	controller.DeleteTodo(c, m.db)
}

func (m *ControllerMap) SignUp(c *gin.Context) {
	controller.SignUp(c, m.db)
}

func (m *ControllerMap) Login(c *gin.Context) {
	controller.Login(c, m.db)
}

func (m *ControllerMap) ClearAll(c *gin.Context) {
	controller.ClearAll(c, m.db)
}

func (m *ControllerMap) UpdateTodo(c *gin.Context) {
	controller.UpdateTodo(c, m.db)
}

func (m *ControllerMap) Todo(c *gin.Context) {
	controller.Todo(c, m.db)
}

func main() {
	godotenv.Overload()

	db := database.CreateDBClientFromEnv()
	defer db.Close()
	m := ControllerMap{db: db}

	router := gin.Default()
	router.LoadHTMLGlob("assets/*.html")
	router.Static("/assets", "./assets")

	router.GET("/", index)
	router.GET("/todos/:userid", m.GetTodos)
	router.GET("/todo/:id", m.GetTodo)
	router.POST("/todo/:userid", m.AddTodo)
	router.DELETE("/todo/:userid/:id", m.DeleteTodo)
	router.DELETE("/todos/:userid", m.ClearAll)
	router.PUT("/todo", m.UpdateTodo)

	router.POST("/signup", m.SignUp)
	router.POST("/login", m.Login)
	router.GET("/todo", m.Todo)

	router.Run(":8080")

}
