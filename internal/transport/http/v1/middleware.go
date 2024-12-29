package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	tools "github.com/rusystem/crm-api/tools"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userInfoCtx         = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	info, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userInfoCtx, info)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	info, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if info.Role != domain.AdminRole {
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	c.Set(userInfoCtx, info)
}

func (h *Handler) superAdminIdentity(c *gin.Context) {
	info, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if info.Role != domain.AdminRole {
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	if !tools.IsFullAccessSection(info.Sections) {
		newErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	c.Set(userInfoCtx, info)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (domain.JWTInfo, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return domain.JWTInfo{}, errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return domain.JWTInfo{}, errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return domain.JWTInfo{}, errors.New("token is empty")
	}

	ip, err := tools.GetIPAddress(c)
	if err != nil {
		return domain.JWTInfo{}, err
	}

	userAgent := tools.GetUserAgent(c)

	info, _, err := h.services.Auth.ValidateAccessToken(c, headerParts[1], userAgent, ip)
	if err != nil {
		return domain.JWTInfo{}, err
	}

	return info, nil
}

func getUserInfo(c *gin.Context) (domain.JWTInfo, error) {
	return getInfoByContext(c, userInfoCtx)
}

func getInfoByContext(c *gin.Context, context string) (domain.JWTInfo, error) {
	value, ex := c.Get(context)
	if !ex {
		return domain.JWTInfo{}, errors.New("userCtx not found")
	}

	info, ok := value.(domain.JWTInfo)
	if !ok {
		return domain.JWTInfo{}, errors.New("userCtx is of invalid type")
	}

	return info, nil
}
