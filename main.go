package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

type Model struct {
	db *sql.DB
}

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Avatar      []byte `json:"avatar"`
	OwnerID     int    `json:"owner_id"`
}

type CourseModel struct {
	Model
}

func (m *CourseModel) GetList(limit int, offset int) []Course {
	courses := make([]Course, 0)

	rows, err := m.db.Query(`SELECT * FROM courses LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		//
		return courses
	}

	for rows.Next() {
		var course Course

		err = rows.Scan(&course.ID, &course.Title, &course.Description, &course.Avatar, &course.OwnerID)
		if err != nil {
			//
			return courses
		}

		courses = append(courses, course)
	}

	return courses
}

const (
	DB_USER        = "root"
	DB_PASS        = "qjuehn123"
	DB_NAME        = "coursify"
	dataSourceName = DB_USER + ":" + DB_PASS + "@/" + DB_NAME + "?charset=utf8&parseTime=true"
)

func main() {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()

	courseModel := CourseModel{Model{db}}
	coursesRoute := r.Group("/courses")
	coursesRoute.GET("/", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		list := courseModel.GetList(limit, offset)

		c.JSON(http.StatusOK, gin.H{
			"courses": list,
			"meta": gin.H{
				"limit":  limit,
				"offset": offset,
			},
		})
	})

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
