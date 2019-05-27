package main

import (
	"database/sql"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func searchCredential(authValue string, db *sql.DB) (int64, bool) {
	if authValue == "" {
		return 0, false
	}

	auth := strings.SplitN(authValue, " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return 0, false
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		return 0, false
	}

	name := pair[0]
	password := pair[1]

	result := db.QueryRow(`
		SELECT
		    password_hash,
		    id
		FROM users WHERE user_name = ?
	`, name)

	var id int64
	var storedPasswordHash string

	err := result.Scan(&storedPasswordHash, &id)
	if err != nil {
		return 0, false
	}

	if password != storedPasswordHash {
		return 0, false
	}

	return id, true
}

// createBasicAuthMiddleware returns a Basic HTTP Authorization middleware.
func createBasicAuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Header.Get("Authorization"))

		userID, found := searchCredential(c.Request.Header.Get("Authorization"), db)
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(gin.AuthUserKey, userID)
	}
}
