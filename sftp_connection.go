package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestingConnection(c *gin.Context) {
	sftp := new(SftpClient)
	if err := c.Bind(sftp); err != nil {
		panic(err.Error())
	}

	if sftp.Connect() != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Not Connected",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Connected",
		})
	}
}

func ConnectionSftp(c *gin.Context) bool {
	sftp := new(SftpClient)
	if err := c.Bind(sftp); err != nil {
		panic(err.Error())
		return false
	}

	if sftp.Connect() != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Not Connected",
		})
		return false
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Connected",
		})
		return true
	}
}
