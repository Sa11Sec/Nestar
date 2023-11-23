package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func TlSHandler(port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + port,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			return
		}
		c.Next()
	}
}
