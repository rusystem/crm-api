package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func corsMiddleware(c *gin.Context) {
	// Получаем Origin из заголовков запроса
	origin := c.Request.Header.Get("Origin")

	// Список разрешённых источников
	allowedOrigins := map[string]bool{
		"http://localhost:3000":     true, // Локальный фронтенд
		"http://localhost":          true, // Локальный фронтенд
		"http://127.0.0.1:3000":     true, // Альтернативный локальный адрес
		"http://91.243.71.100:3000": true, // IP фронтенда
		"http://91.243.71.100:5173": true, // Vite DevServer
		"http://91.243.71.100":      true, // Vite DevServer
	}

	// Проверяем, разрешён ли Origin
	isAllowed := allowedOrigins[origin]

	// Устанавливаем CORS-заголовки
	if isAllowed {
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		c.Header("Access-Control-Allow-Origin", "") // Заблокировать неразрешённые Origin
	}

	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
	c.Header("Access-Control-Expose-Headers", "Authorization") // Позволяет клиенту читать этот заголовок
	c.Header("Access-Control-Allow-Credentials", "true")       // Поддержка авторизационных данных (куки/токены)

	// Обработка preflight-запросов (OPTIONS)
	if c.Request.Method == http.MethodOptions {
		if isAllowed {
			c.AbortWithStatus(http.StatusOK) // Успешный ответ для разрешённого Origin
		} else {
			c.AbortWithStatus(http.StatusForbidden) // 403 для неразрешённого Origin
		}
		return
	}

	// Продолжаем выполнение для других методов
	c.Next()
}

func trailingSlashMiddleware(c *gin.Context) {
	path := c.Request.URL.Path

	if len(path) > 1 && path[len(path)-1] != '/' {
		c.Request.URL.Path = path + "/"
	}

	c.Next()
}
