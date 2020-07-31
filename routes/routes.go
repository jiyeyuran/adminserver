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

		passport := admin.Group("/passport", authMiddleware(app))
		{
			server := server.NewPassportServer(app)
			passport.POST("/signup", server.Signup)
			passport.POST("/login", server.Login)
			passport.POST("/logout", server.Logout)
		}

		roomGroup := admin.Group("/room", authMiddleware(app))
		{
			roomServer := server.NewRoomServer(app)
			roomGroup.POST("/room/info", roomServer.Info)
			roomGroup.POST("/room/list", roomServer.List)
			roomGroup.POST("/room/create", roomServer.Create)
			roomGroup.POST("/room/modify", roomServer.Modify)
			roomGroup.POST("/room/delete", roomServer.Delete)
			roomGroup.POST("/room/token", roomServer.Token)
		}

		conferenceGroup := admin.Group("/conference", authMiddleware(app))
		{
			conferenceServer := server.NewConferenceServer(app)
			conferenceGroup.POST("/conference/info", conferenceServer.Info)
			conferenceGroup.POST("/conference/runing", conferenceServer.Runing)
			conferenceGroup.POST("/conference/dispose", conferenceServer.Dispose)
			conferenceGroup.POST("/conference/lock", conferenceServer.Lock)
			conferenceGroup.POST("/conference/unlock", conferenceServer.Unlock)

			conferenceGroup.POST("/conference/history", conferenceServer.History)
		}
	}
}

func handleCaptchaId(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": captcha.NewLen(6),
	})
}
