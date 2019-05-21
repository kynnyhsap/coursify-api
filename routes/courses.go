package routes

import (
	"coursify-api/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func listCourses(model models.ICourseLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		//userID, _ := strconv.Atoi(c.DefaultQuery("userId", "0"))

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

func createCourse(model models.ICourseCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		inputData := models.CourseCreateInput{}
		err := c.ShouldBindJSON(&inputData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		id := model.Create(inputData)
		course := model.Get(id)

		c.JSON(http.StatusCreated, course)
	}
}

func getCourse(model models.ICourseGetter) gin.HandlerFunc {
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

func updateCourse(model models.ICourseUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		course := model.Get(id)
		// TODO: check if exists

		err = c.ShouldBindJSON(&course)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		model.Update(course)

		c.JSON(http.StatusOK, model.Get(id))
	}
}

func deleteCourse(model models.ICourseDeleter) gin.HandlerFunc {
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

func SetUpCourses(group *gin.RouterGroup, db *sql.DB) {
	m := models.NewCourseModel(db)

	group.GET("/", listCourses(m))
	group.POST("/", createCourse(m))
	group.GET("/:id", getCourse(m))
	group.DELETE("/:id", deleteCourse(m))
	group.PUT("/:id", updateCourse(m))
}
