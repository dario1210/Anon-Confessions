package comments

import "github.com/gin-gonic/gin"

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		router: gin.Default(),
	}
}

func (r *Router) RegisterRoutes() {
	r.router.GET("/someGet", test)
}

func (r *Router) Run(addr string) error {
	return r.router.Run(addr)
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from /someGet",
	})
}
