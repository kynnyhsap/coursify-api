package routes

import (
	"coursify-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUser(model models.IUserCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		inputData := models.UserCreateInput{}
		err := c.ShouldBindJSON(&inputData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		id := model.Create(inputData)
		user := model.Get(id)

		c.JSON(http.StatusCreated, user)
	}
}

func LogInUser(model models.IUserCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "")
	}
}

func GetSelf(model models.IUserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		any, _ := c.Get(gin.AuthUserKey)
		selfID, _ := any.(int64)

		user := model.Get(selfID)
		c.JSON(http.StatusOK, user)
	}
}
