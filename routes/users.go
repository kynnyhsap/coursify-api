package routes

import (
	"coursify-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
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

func ListUsers(model models.IUserLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		decodedSearchQuery, _ := url.QueryUnescape(c.DefaultQuery("search", ""))

		list := model.GetList(limit, offset, decodedSearchQuery)
		total := model.Count()

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"limit":  limit,
				"offset": offset,
				"total":  total,
			},
			"users": list,
		})
	}
}
