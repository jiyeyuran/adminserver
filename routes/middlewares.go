package routes

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// timeout middleware wraps the request context with a timeout
func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				c.AbortWithError(http.StatusRequestTimeout, errors.New("请求超时"))
			}
			cancel()
		}()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// 格式化错误JSON输出，{"error": "msg"}
func errorMiddleware(c *gin.Context) {
	c.Next()

	if err := c.Errors; err != nil {
		if strings.HasPrefix(c.ContentType(), "application/json") {
			// log.Error().Err(err.Last()).Send()
		}
		status := c.Writer.Status()
		if status <= 0 {
			status = 500
		}
		c.JSON(status, err.Last())
	}
}
