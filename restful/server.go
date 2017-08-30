package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getting(c *gin.Context)  {}
func posting(c *gin.Context)  {}
func putting(c *gin.Context)  {}
func deleting(c *gin.Context) {}
func patching(c *gin.Context) {}
func head(c *gin.Context)     {}
func options(c *gin.Context)  {}

//StartAPIServer 开启RESTful API服务
func StartAPIServer(c chan error) {
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.POST("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	router.GET("/Get", getting)
	router.POST("/Post", posting)
	router.PUT("/Put", putting)
	router.DELETE("/Delete", deleting)
	router.PATCH("/Patch", patching)
	router.HEAD("/Head", head)
	router.OPTIONS("/Options", options)
	//router.Run(":8080")
	c <- http.ListenAndServe(":8080", router)
}
