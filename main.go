package main

import (
	"coursify-api/routes"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

	coursesRoute := r.Group("/courses")
	//lessonsRoute := r.Group("/lessons")
	//usersRoute := r.Group("/users")

	routes.SetUpCourses(coursesRoute, db)

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
