package routes

import (
	"coursify-api/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func listLessons(model models.ILessonLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		courseID, err := strconv.ParseInt(c.Query("courseId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		list := model.GetList(courseID, limit, offset)

		c.JSON(http.StatusOK, gin.H{
			"lessons": list,
			"meta": gin.H{
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}

func createLesson(model models.ILessonCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		inputData := models.LessonCreateInput{}
		err := c.ShouldBindJSON(&inputData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// TODO: check if course exits

		fmt.Println(inputData)

		id := model.Create(inputData)
		lesson := model.Get(id)

		c.JSON(http.StatusCreated, lesson)
	}
}

func getLesson(model models.ILessonGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		course := model.Get(id)
		// TODO: check if exists

		c.JSON(http.StatusOK, course)
	}
}

func updateLesson(model models.ILessonUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		lesson := model.Get(id)
		// TODO: check if exists

		err = c.ShouldBindJSON(&lesson)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		model.Update(lesson)

		c.JSON(http.StatusOK, model.Get(id))
	}
}

func deleteLesson(model models.ILessonDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		// TODO: check if exists
		model.Delete(id)

		c.JSON(http.StatusOK, gin.H{})
	}
}

func SetUpLessons(group *gin.RouterGroup, db *sql.DB) {
	m := models.NewLessonModel(db)

	group.GET("/", listLessons(m))
	group.POST("/", createLesson(m))
	group.GET("/:id", getLesson(m))
	group.PUT("/:id", updateLesson(m))
	group.DELETE("/:id", deleteLesson(m))
}
