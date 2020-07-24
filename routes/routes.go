package routes

import (
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"jhmeeting.com/adminserver/app"
	"jhmeeting.com/adminserver/server"
)

func Setup(r *gin.Engine, app *app.App) {
	admin := r.Group("/admin")
	admin.Use(errorMiddleware, timeoutMiddleware(5*time.Second))
	{

		admin.POST("/captcha-id", handleCaptchaId)
		admin.GET("/captcha-id", handleCaptchaId)
		admin.GET("/captcha/:id", gin.WrapH(captcha.Server(captcha.StdWidth, captcha.StdHeight)))

		passport := admin.Group("passport")
		{
			server := server.NewPassportServer(app)
			passport.POST("/signup", server.Signup)
			passport.POST("/login", server.Login)
			passport.POST("/logout", server.Logout)
		}
	}
}

func handleCaptchaId(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": captcha.NewLen(6),
	})
}
