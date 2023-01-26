package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	r := gin.Default()
	r.POST("/test_connection", TestingConnection)
	r.POST("/send_sftp", SendFiles)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func SendFiles(c *gin.Context) {

	//Check Form Request
	HostName := c.PostForm("host")
	Username := c.PostForm("user")
	Password := c.PostForm("password")
	Port := c.PostForm("port")
	Path := c.PostForm("path")
	file, err := c.FormFile("files")

	if err != nil {
		panic(err.Error())
	}

	intPort, err := strconv.Atoi(Port) // convert string to int

	SendToSftp(c, HostName, Username, Password, intPort, Path, file)

}
