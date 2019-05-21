package main

import (
	"database/sql"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Credentials
}

func searchCredential(authValue string, db *sql.DB) (string, bool) {
	if authValue == "" {
		return "", false
	}

	auth := strings.SplitN(authValue, " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return "", false
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		return "", false
	}

	login := pair[0]
	password := pair[1]

	result := db.QueryRow(`SELECT password_hash FROM users WHERE email = ?`, login)

	storedPasswordHash := ""
	err := result.Scan(&storedPasswordHash)
	if err != nil {
		return "", false
	}

	if password != storedPasswordHash {
		return "", false
	}

	return login, true
}

// createBasicAuthMiddleware returns a Basic HTTP Authorization middleware.
func createBasicAuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Header.Get("Authorization"))

		user, found := searchCredential(c.Request.Header.Get("Authorization"), db)
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(gin.AuthUserKey, user)
	}
}

func registretionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		newUser := RegisterInput{}

		err := c.ShouldBindJSON(&newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		// check for unique `email` field

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		stmt, err := db.Prepare(`INSERT INTO users(first_name, last_name, email, password_hash) VALUE (?, ?, ?, ?)`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		_, err = stmt.Exec(newUser.FirstName, newUser.LastName, newUser.Email, hashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, "successfully registered")
	}
}
