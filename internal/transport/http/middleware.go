package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func corsMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	allowedOrigins := []string{
		"http://localhost",
		"http://127.0.0.1",
		"http://91.243.71.100:5173",
		"http://91.243.71.100:3000",
	}

	isAllowed := false
	for _, o := range allowedOrigins {
		if origin == o {
			isAllowed = true
			break
		}
	}

	if isAllowed {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	}

	if c.Request.Method == "OPTIONS" {
		if isAllowed {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}

func trailingSlashMiddleware(c *gin.Context) {
	path := c.Request.URL.Path

	if len(path) > 1 && path[len(path)-1] != '/' {
		c.Request.URL.Path = path + "/"
	}

	c.Next()
}
