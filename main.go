package main // import "github.com/nzlov/taskplan"

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yinheli/static"

	"github.com/nzlov/taskplan/app"
)

func main() {
	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(static.Serve("/", static.LocalFile("./public/dist", false)))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-AppUser", "X-AppSign")
	config.AddAllowMethods("DELETE")
	g.Use(cors.New(config))

	app.Init(g)

	g.Run(":9005")
}
