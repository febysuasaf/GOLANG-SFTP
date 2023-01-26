package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net/http"
	"time"
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

func (sc *SftpClient) Connect() (err error) {
	config := &ssh.ClientConfig{
		User:            sc.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(sc.Password)},
		Timeout:         30 * time.Second,
	}

	// connet to ssh
	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	// create sftp client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	sc.Client = client
	return nil
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
