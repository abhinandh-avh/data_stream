package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("")

	r.POST("")

	return r

}
