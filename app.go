package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var state int

func GetState() string {
	return strconv.Itoa(state)
}

func IncrementState() string {
	state++
	return GetState()
}

//localhost:port/ (возвращается текущее значение счётчика, инкремента не происходит);
func HandleStateRender(c *gin.Context) {
	c.String(http.StatusOK, GetState())
}

//localhost:port/stat (возвращается текущее значение счётчика, и происходит инкремент);
func HandleStatIncrement(c *gin.Context) {
	c.String(http.StatusOK, IncrementState())
	// or if literally:
	//
	// HandleStateRender(c)
	// state++
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
