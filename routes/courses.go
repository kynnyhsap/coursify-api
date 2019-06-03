package routes

import (
	"coursify-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

func ListCourses(model models.ICourseLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		listType := c.DefaultQuery("type", models.TypeAllCourses)

		var list []models.CourseForList
		var total int

		if listType == models.TypeAllCourses {
			searchQuery := c.DefaultQuery("search", "")
			decodedSearchQuery, _ := url.QueryUnescape(searchQuery)

			list = model.GetList(limit, offset, decodedSearchQuery)
			total = model.Count()
		} else {
			any, _ := c.Get(gin.AuthUserKey)
			selfID, _ := any.(int64)

			list = model.GetListForUser(limit, offset, selfID, listType == models.TypeAdminCourses)
			total = model.CountForUser(selfID, listType == models.TypeAdminCourses)
		}

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"limit":  limit,
				"offset": offset,
				"total":  total,
			},
			"courses": list,
		})
	}
}

func CreateCourse(model models.ICourseCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		inputData := models.CourseCreateInput{}
		err := c.ShouldBindJSON(&inputData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		any, _ := c.Get(gin.AuthUserKey)
		selfID, _ := any.(int64)

		id := model.Create(inputData, selfID)
		course := model.Get(id)

		c.JSON(http.StatusCreated, course)
	}
}

func GetCourse(model models.ICourseGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		course := model.Get(id)
		// TODO: check if exists
		if course.ID == 0 {
			c.String(http.StatusNotFound, "No course with id %d", id)
			return
		}

		c.JSON(http.StatusOK, course)
	}
}

func UpdateCourse(model models.ICourseUpdater) gin.HandlerFunc {
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

func DeleteCourse(model models.ICourseDeleter) gin.HandlerFunc {
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
