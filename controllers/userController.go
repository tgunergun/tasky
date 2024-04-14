package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeffthorne/tasky/auth"
	"github.com/jeffthorne/tasky/database"
	"github.com/jeffthorne/tasky/models"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func SignUp(c *gin.Context, db database.DBClient) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	emailCount, err := db.FindExistingUsers(ctx, user)
	defer cancel()

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	if emailCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists!"})
		return
	}
	msg, err := db.AddUser(ctx, &user)
	userId := user.ID.Hex()
	username := *user.Name

	token, err, expirationTime := auth.GenerateJWT(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "userID",
		Value:   userId,
		Expires: expirationTime,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "username",
		Value:   username,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, msg)

}
func Login(c *gin.Context, db database.DBClient) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bind error"})
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	foundUser, err := db.GetUser(ctx, *user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " email or password is incorrect"})
		return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()

	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found!"})
		return
	}
	userId := foundUser.ID.Hex()
	username := *foundUser.Name

	shouldRefresh, err, expirationTime := auth.RefreshToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token error"})
		return
	}

	if shouldRefresh {
		token, err, expirationTime := auth.GenerateJWT(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while generating token"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})

		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "userID",
			Value:   userId,
			Expires: expirationTime,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "username",
			Value:   username,
			Expires: expirationTime,
		})

	} else {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "userID",
			Value:   userId,
			Expires: expirationTime,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "username",
			Value:   username,
			Expires: expirationTime,
		})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "login successful"})
}

func Todo(c *gin.Context, db database.DBClient) {
	session := auth.ValidateSession(c)
	if session {
		c.HTML(http.StatusOK, "todo.html", nil)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}

	return check, msg
}
