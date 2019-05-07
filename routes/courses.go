package routes

import (
	"coursify-api/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type courseLister interface {
	GetList(limit int, offset int) []models.Course
}

func listCourses(model courseLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		list := model.GetList(limit, offset)

		c.JSON(http.StatusOK, gin.H{
			"courses": list,
			"meta": gin.H{
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}

func SetUpCourses(group *gin.RouterGroup, db *sql.DB) {
	m := models.NewCourseModel(db)

	group.GET("/", listCourses(m))
	//group.POST("/", createCourse(model))
	// ...
}
