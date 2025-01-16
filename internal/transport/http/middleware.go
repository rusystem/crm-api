package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func corsMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")

	origin = strings.TrimSuffix(origin, "/")

	allowedOrigins := map[string]bool{
		"http://localhost":          true,
		"http://127.0.0.1:3000":     true,
		"http://91.243.71.100:3000": true,
		"http://91.243.71.100:5173": true,
		"http://91.243.71.100":      true,
	}

	isAllowed := allowedOrigins[origin]

	if isAllowed {
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		c.Header("Access-Control-Allow-Origin", "")
	}

	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
	c.Header("Access-Control-Expose-Headers", "Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == http.MethodOptions {
		if isAllowed {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
		return
	}

	c.Next()
}

func trailingSlashMiddleware(c *gin.Context) {
	path := c.Request.URL.Path
	if len(path) > 1 && path[len(path)-1] != '/' {
		c.Request.URL.Path += "/"
	}
	c.Next()
}
