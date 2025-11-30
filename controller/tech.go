package controller

import "github.com/gin-gonic/gin"

func TechGoHandler(c *gin.Context) {
	//TODO:
	c.HTML(200, "go.html", nil)
}
func TechJavaHandler(c *gin.Context) {
	c.HTML(200, "java.html", nil)
}
func TechFrontendHandler(c *gin.Context) {
	c.HTML(200, "frontend.html", nil)
}
