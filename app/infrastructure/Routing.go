package infrastructure

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"net/http"
	"os"
	"time"
)

type Routing struct {
	Gin          *gin.Engine
	AbsolutePath string
}

func NewRouting() *Routing {
	c, _ := NewConfig()
	r := &Routing{
		Gin:          gin.Default(),
		AbsolutePath: c.AbsolutePath,
	}
	r.loadTemplates()
	r.setHeader()
	r.setRouting()
	return r
}

func (r *Routing) loadTemplates() {
	r.Gin.Use(favicon.New("./dist/assets/favicon.ico"))
	r.Gin.Static("/assets", r.AbsolutePath+"/dist/assets")
	r.Gin.LoadHTMLGlob(r.AbsolutePath + "/app/interfaces/presenters/*")
}

func (r *Routing) setHeader() {
	r.Gin.Use(cors.New(cors.Config{
		AllowOrigins: []string {
			"*",
		},
		AllowCredentials: false,
		AllowHeaders: []string {
			"Content-Type",
		},
		AllowMethods: []string {
			"GET",
			"HEAD",
			"OPTIONS",
		},
		MaxAge: time.Duration(86400),
	}))
}

func (r *Routing) setRouting() {
	r.Gin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Gin.GET("/api", func(c *gin.Context) {
		from := c.Query("from")
		to := c.Query("to")

		ride := newRide(from, to)
		busInformation := ride.fetch()

		c.JSON(200, gin.H{
			"body": len(body),
		})
	})
}

func (r *Routing) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return r.Gin.Run(":" + port)
}
