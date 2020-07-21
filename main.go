package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
	"jhmeeting.com/adminserver/routes"
)

func main() {
	r := gin.Default()
	app := app.NewApp()

	routes.Setup(r, app)

	r.Run(fmt.Sprintf(":%d", app.Config().Port))
}
