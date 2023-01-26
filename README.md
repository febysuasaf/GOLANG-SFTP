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
> Main.go
dan tuliskan barisan function seperti berikut :
```
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
```
