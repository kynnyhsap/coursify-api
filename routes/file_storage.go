package routes

import (
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

func PostImageFile(dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageData, err := c.GetRawData()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		name := guuid.New().String() + ".jpeg"
		path := dir + "/images/" + name

		err = ioutil.WriteFile(path, imageData, 0777)
		if err != nil {
			c.String(http.StatusInternalServerError, "")
		}

		imageURL := "http://localhost:8080/fs/images/" + name
		c.String(http.StatusOK, imageURL)
	}
}
