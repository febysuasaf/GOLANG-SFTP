# GOLANG-SFTP
Membuat Project Golang Sftp Client Menggunakan Framework Gin + Sftp

Buat Project Baru Jalankan perintah Berikut:
```
go get -u github.com/gin-gonic/gin
```
Setelah Berhasil ditambahkan jalankan perinta Berikut :

```
go get -u github.com/pkg/sftp
```
Setelah framework gin dan middleware sftp sudah di tambahkan buatlah file :
> main.go
dan tuliskan barisan function seperti berikut :
```
package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	r := gin.Default()
	r.POST("/test_connection", TestingConnection) // route testing connection to server sftp
	r.POST("/send_sftp", SendFiles) // route send file to server sftp

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func SendFiles(c *gin.Context) {

	//Check Form Request
	HostName := c.PostForm("host") // body form-data hostname server (text)
	Username := c.PostForm("user") // body form-data username server (text)
	Password := c.PostForm("password") // body form-data password server (text)
	Port := c.PostForm("port") // body form-data port server (text integer)
	Path := c.PostForm("path") // body form-data path directory send file sftp (text)
	file, err := c.FormFile("files") // body form-data select data file (file)

	if err != nil {
		panic(err.Error())
	}

	intPort, err := strconv.Atoi(Port) // convert port string to int

	SendToSftp(c, HostName, Username, Password, intPort, Path, file) // Call function SendtoSftp

}
```
Buat File Baru dan tambahkan baris code function berikut :
> sftp.go
```
type SftpClient struct {
	Host     string `json:"host" form:"host"`
	User     string `json:"user" form:"user"`
	Password string `json:"-" form:"password"`
	Port     int    `json:"port" form:"port"`
	*sftp.Client
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
```
