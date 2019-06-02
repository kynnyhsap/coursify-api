package main

import (
	"coursify-api/models"
	"coursify-api/routes"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func getDataSourceName(user string, pass string, dbname string) string {
	return user + ":" + pass + "@/" + dbname + "?charset=utf8&parseTime=true"
}

func main() {
	db, err := sql.Open("mysql", getDataSourceName(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")))
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

	authMiddleware := createBasicAuthMiddleware(db)

	userModel := models.NewUserModel(db)
	courseModel := models.NewCourseModel(db)
	lessonModel := models.NewLessonModel(db)

	coursesGroup := r.Group("/courses", authMiddleware)
	lessonsGroup := r.Group("/lessons", authMiddleware)
	usersGroup := r.Group("/users", authMiddleware)

	lessonsGroup.GET("/", routes.ListLessons(lessonModel))
	lessonsGroup.POST("/", routes.CreateLesson(lessonModel))
	lessonsGroup.GET("/:id", routes.GetLesson(lessonModel))
	lessonsGroup.PUT("/:id", routes.UpdateLesson(lessonModel))
	lessonsGroup.DELETE("/:id", routes.DeleteLesson(lessonModel))

	coursesGroup.GET("/", routes.ListCourses(courseModel))
	coursesGroup.POST("/", routes.CreateCourse(courseModel))
	coursesGroup.GET("/:id", routes.GetCourse(courseModel))
	coursesGroup.DELETE("/:id", routes.DeleteCourse(courseModel))
	coursesGroup.PUT("/:id", routes.UpdateCourse(courseModel))

	usersGroup.GET("/self/", routes.GetSelf(userModel))

	r.POST("/register/", routes.RegisterUser(userModel))
	r.GET("/login/", authMiddleware, routes.LogInUser(userModel))

	r.POST("/fs/images/", routes.PostImageFile("file_storage"))
	r.StaticFS("/fs/images/", http.Dir("file_storage/images"))

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
