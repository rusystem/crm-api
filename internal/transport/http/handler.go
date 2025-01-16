package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/rusystem/crm-api/docs/v1"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/service"
	v1 "github.com/rusystem/crm-api/internal/transport/http/v1"
	"github.com/rusystem/crm-api/pkg/auth"
	"github.com/rusystem/crm-api/pkg/limiter"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/unrolled/secure"
	"net/http"
)

type Handler struct {
	service      *service.Service
	tokenManager auth.TokenManager
	cfg          *config.Config
}

func NewHandler(service *service.Service, tokenManager auth.TokenManager, cfg *config.Config) *Handler {
	return &Handler{
		service:      service,
		tokenManager: tokenManager,
		cfg:          cfg,
	}
}

func (h *Handler) Init() *gin.Engine {
	if h.cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin handler
	router := gin.Default()

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		//ContentSecurityPolicy: "default-src 'self'",
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			if err != nil {
				c.Abort()
				return
			}

			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	//router.RedirectTrailingSlash = false
	//router.RemoveExtraSlash = false

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
		secureFunc,
		limiter.Limit(h.cfg.Limiter.RPS, h.cfg.Limiter.Burst, h.cfg.Limiter.TTL),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service, h.tokenManager, h.cfg)
	api := router.Group("/api/web-api-gateway")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		handlerV1.Init(api)
	}
}
