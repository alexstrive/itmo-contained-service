package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

var isDbConnected = false
var db sql.DB

func ConnectDb() {
	connStr := "user=postgres dbname=postgres password=example sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		println(err.Error())
	}

	isDbConnected = true
}

func GetState() int {
	connStr := os.Getenv("DB_STR")
	if connStr == "" {
		connStr = "user=postgres dbname=postgres password=example sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		println(err.Error())
	}

	defer db.Close()

	var count int
	row := db.QueryRow("SELECT count(*) FROM table_Counter")
	err = row.Scan(&count)
	if err != nil {
		println(err.Error())
	}

	return count
}

func IncrementState(clientInfo string) {
	connStr := os.Getenv("DB_STR")
	if connStr == "" {
		connStr = "user=postgres dbname=postgres password=example sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		println(err.Error())
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO table_counter(datetime, client_info) VALUES (NOW(), $1)", clientInfo)

	if err != nil {
		println(err.Error())
	}
}

//localhost:port/ (возвращается текущее значение счётчика, инкремента не происходит);
func HandleStateRender(c *gin.Context) {
	c.String(http.StatusOK, strconv.Itoa(GetState()))
}

//localhost:port/stat (возвращается текущее значение счётчика, и происходит инкремент);
func HandleStatIncrement(c *gin.Context) {
	var userInfo = c.Request.UserAgent()
	IncrementState(userInfo)
	c.String(http.StatusOK, strconv.Itoa(GetState()))
}

//localhost:port/about (возвращается html-страничка ниже) как показано в app.py, функция hello в gist только с вашим именем <h3> Hello , _your_name_</h3>.
func HandleAboutPage(c *gin.Context) {

	c.HTML(http.StatusOK, "about.tmpl", gin.H{
		"hostname": c.Request.Host,
		"name":     "Alexey Novopashin",
	})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", HandleStateRender)
	r.GET("/stat", HandleStatIncrement)
	r.GET("/about", HandleAboutPage)
	r.Run(":3000")
}
