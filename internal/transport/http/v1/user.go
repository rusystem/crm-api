package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	user := api.Group("/user")
	{
		user.GET("/info", h.userIdentity, h.getUserInfo)
		user.PUT("/profile", h.userIdentity, h.updateProfile)

		// only admin can create, update, delete user
		user.GET("/:id", h.adminIdentity, h.getUser)
		user.PUT("/:id", h.adminIdentity, h.updateUser)
		user.DELETE("/:id", h.adminIdentity, h.deleteUser)
		user.GET("/company", h.adminIdentity, h.getUsers)
	}
}

// @Summary Get user info
// @Security ApiKeyAuth
// @Tags user
// @Description Получение информации о пользователе.
// @ID get-user-info
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /user/info [GET]
func (h *Handler) getUserInfo(c *gin.Context) {
	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := h.services.User.GetById(c, info.UserId, info)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data: domain.UserResponse{
			ID:                       userInfo.ID,
			CompanyID:                userInfo.CompanyID,
			Username:                 userInfo.Username,
			Name:                     userInfo.Name,
			Email:                    userInfo.Email,
			Phone:                    userInfo.Phone,
			CreatedAt:                userInfo.CreatedAt,
			UpdatedAt:                userInfo.UpdatedAt,
			LastLogin:                userInfo.LastLogin.Time,
			IsActive:                 userInfo.IsActive,
			Role:                     userInfo.Role,
			Language:                 userInfo.Language,
			Country:                  userInfo.Country,
			IsApproved:               userInfo.IsApproved,
			IsSendSystemNotification: userInfo.IsSendSystemNotification,
			Position:                 userInfo.Position,
			Sections:                 userInfo.Sections,
		},
		TotalCount: 1,
	})
}

// @Summary Update user profile
// @Security ApiKeyAuth
// @Tags user
// @Description Обновление информации о пользователе.
// @Description Необходимо передавать только измененные данные.
// @ID update-user-profile
// @Accept  json
// @Produce  json
// @Param id path int true "user id" example(1)
// @Param request body domain.UserProfileUpdate true "request body"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /user/profile [PUT]
func (h *Handler) updateProfile(c *gin.Context) {
	var req domain.UserProfileUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	req.ID, err = parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = h.services.User.UpdateProfile(c, req, info); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Get user
// @Security ApiKeyAuth
// @Tags user
// @Description Получение информации о пользователе по id.
// @Description Только super admin может получать информацию по id пользователя.
// @ID get-user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID" example(1)
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /user/{id} [GET]
func (h *Handler) getUser(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := h.services.User.GetById(c, id, info)
	if err != nil {
		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data: domain.UserResponse{
			ID:                       userInfo.ID,
			CompanyID:                userInfo.CompanyID,
			Username:                 userInfo.Username,
			Name:                     userInfo.Name,
			Email:                    userInfo.Email,
			Phone:                    userInfo.Phone,
			CreatedAt:                userInfo.CreatedAt,
			UpdatedAt:                userInfo.UpdatedAt,
			LastLogin:                userInfo.LastLogin.Time,
			IsActive:                 userInfo.IsActive,
			Role:                     userInfo.Role,
			Language:                 userInfo.Language,
			Country:                  userInfo.Country,
			IsApproved:               userInfo.IsApproved,
			IsSendSystemNotification: userInfo.IsSendSystemNotification,
			Position:                 userInfo.Position,
			Sections:                 userInfo.Sections,
		},
		TotalCount: 1,
	})
}

// @Summary Update user
// @Security ApiKeyAuth
// @Tags user
// @Description Обновление информации о пользователе по id.
// @Description Необходимо передавать только измененные данные.
// @Description Только super admin может обновлять информацию по любому id пользователя.
// @Description Только admin может обновлять информацию по id пользователя в рамках своей компании.
// @Description Только super admin может менять role для пользователя
// @ID update-user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID" example(1)
// @Param request body domain.UserUpdate true "request body"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /user/{id} [PUT]
func (h *Handler) updateUser(c *gin.Context) {
	var req domain.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		newBindingErrorResponse(c, err)
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	*req.ID, err = parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err = h.services.User.Update(c, req, info); err != nil {
		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Удаление пользователя по id.
// @Description Только super admin может удалить пользователя по id.
// @Description Только admin может удалить пользователя в рамках своей компании.
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /user/{id} [DELETE]
func (h *Handler) deleteUser(c *gin.Context) {
	id, err := parseIdIntPathParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err = h.services.User.Delete(c, id, info); err != nil {
		if errors.Is(err, domain.ErrNotAllowed) {
			newErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessOkResponse(c)
}

// @Summary Get users
// @Security ApiKeyAuth
// @Tags user
// @Description Получение списка пользователей компании.
// @ID get-users
// @Accept  json
// @Produce  json
// @Param sort query string true "Sort order" Enums(asc, desc)
// @Param sort_field query string true "Field to sort by" Enums(id, username, name, email, phone, created_at, updated_at, last_login, is_active, role, country, is_approved, position) default(name)
// @Param limit query int true "limit query param"
// @Param offset query int true "offset query param"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /user/company [GET]
func (h *Handler) getUsers(c *gin.Context) {
	info, err := getUserInfo(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	sort, field, err := parseSortParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	limit, err := parseLimitQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	offset, err := parseOffsetQueryParam(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	users, count, err := h.services.User.GetListByCompanyId(c, info.CompanyId, domain.Param{
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
		SortField: field,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, http.StatusOK, domain.SuccessResponse{
		Data:       users,
		TotalCount: count,
	})
}
