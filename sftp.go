package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SftpClient struct {
	Host     string `json:"host" form:"host"`
	User     string `json:"user" form:"user"`
	Password string `json:"-" form:"password"`
	Port     int    `json:"port" form:"port"`
	*sftp.Client
}

func NewConn(host string, user string, password string, port int) (client *SftpClient, err error) {
	switch {
	case `` == strings.TrimSpace(host),
		`` == strings.TrimSpace(user),
		`` == strings.TrimSpace(password),
		0 >= port || port > 65535:
		return nil, errors.New("Invalid parameters")
	}

	client = &SftpClient{
		Host:     host,
		User:     user,
		Password: password,
		Port:     port,
	}

	if err = client.Connect(); nil != err {
		return nil, err
	}
	return client, nil
}

func SendToSftp(c *gin.Context, hostname string, username string, password string, port int, path string, file *multipart.FileHeader) {
	sftp := new(SftpClient)
	if err := c.Bind(sftp); err != nil {
		panic(err.Error())
	}
	// Connect to Server
	ftpClient, err := NewConn(hostname, username, password, port)
	if err != nil {
		log.Fatal(err)
	}

	fDestination, err := os.Create("./" + file.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fDestination.Close()

	fSource, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}

	// copy source file to destination file
	_, err = io.Copy(fDestination, fSource)
	if err != nil {
		log.Fatal(err)
	}
	ftpClient.Put(c, file.Filename, path)

}

// Upload file to sftp server
func (sc *SftpClient) Put(c *gin.Context, remoteFile string, pathDir string) (err error) {

	srcFile, err := os.Open("./" + remoteFile)
	if err != nil {
		panic(err.Error())
	}
	defer srcFile.Close()

	// Note: SFTP To Go doesn't support O_RDWR mode
	dstFile, err := sc.OpenFile(pathDir+filepath.Base(remoteFile), (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
		return
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Failed Send SFTP",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Success Send SFTP",
		})
		fmt.Fprintf(os.Stdout, "%d bytes copied\n", bytes)
	}

	return
}
