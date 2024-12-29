package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	"github.com/rusystem/crm-api/pkg/auth"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/rusystem/crm-api/pkg/logger"
	"github.com/rusystem/crm-api/tools"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"time"
)

type Auth interface {
	SignIn(c *gin.Context, input domain.SignIn) (domain.TokenResponse, error)
	SignOut(c *gin.Context, userId, companyId int64) error
	SignUp(c *gin.Context, input domain.SignUp, info domain.JWTInfo) (int64, bool, error)
	RefreshTokens(c *gin.Context, refreshToken string) (domain.TokenResponse, error)
	ValidateAccessToken(c *gin.Context, token, userAgent, ip string) (domain.JWTInfo, bool, error)
}

type AuthServices struct {
	cfg          *config.Config
	repo         *repository.Repository
	tokenManager auth.TokenManager
}

func NewAuthServices(cfg *config.Config, repo *repository.Repository, tokenManager auth.TokenManager) *AuthServices {
	return &AuthServices{
		cfg:          cfg,
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (as *AuthServices) SignIn(c *gin.Context, input domain.SignIn) (domain.TokenResponse, error) {
	user, err := as.repo.User.GetByUsername(c.Request.Context(), input.Username)
	if err != nil {
		logger.Info(fmt.Sprintf("failed to get user by username, %+v", err))
		return domain.TokenResponse{}, domain.ErrLoginCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		logger.Info("failed to compare hash and password")
		return domain.TokenResponse{}, domain.ErrLoginCredentials
	}

	if !user.IsActive {
		return domain.TokenResponse{}, domain.ErrUserIsNotActive
	}

	if !user.IsApproved {
		return domain.TokenResponse{}, domain.ErrUserIsNotApproved
	}

	// check if user has full access to the company
	if input.CompanyId != 0 && tools.IsFullAccessSection(user.Sections) {
		isExist, err := as.repo.Company.IsExist(c.Request.Context(), input.CompanyId)
		if err != nil {
			return domain.TokenResponse{}, err
		}

		if !isExist {
			return domain.TokenResponse{}, domain.ErrCompanyNotFound
		}

		user.CompanyID = input.CompanyId
	}

	return as.createSession(c, user)
}

func (as *AuthServices) SignOut(c *gin.Context, userId, companyId int64) error {
	if err := as.repo.Auth.DeleteUserTokens(c.Request.Context(), userId, companyId); err != nil {
		return domain.ErrSignOut
	}

	return nil
}

func (as *AuthServices) SignUp(c *gin.Context, input domain.SignUp, info domain.JWTInfo) (int64, bool, error) {
	authUser, err := as.repo.User.GetById(c.Request.Context(), info.UserId)
	if err != nil {
		return 0, false, domain.ErrRoleNotAllowed
	}

	if !authUser.IsActive {
		return 0, false, domain.ErrUserIsNotActive
	}

	if !authUser.IsApproved {
		return 0, false, domain.ErrUserIsNotApproved
	}

	// проверяем чтобы не создавали админ пользователей простые смертные
	if input.Role == domain.AdminRole && !tools.IsFullAccessSection(authUser.Sections) {
		return 0, false, domain.ErrRoleNotAllowed
	}

	if input.Role == "" {
		input.Role = domain.UserRole
	}

	if input.CompanyId != 0 && !tools.IsFullAccessSection(authUser.Sections) {
		return 0, false, domain.ErrRoleNotAllowed
	}

	if input.CompanyId == 0 {
		input.CompanyId = info.CompanyId
	}

	if input.CompanyId != 0 {
		isExist, err := as.repo.Company.IsExist(c.Request.Context(), input.CompanyId)
		if err != nil {
			if errors.Is(err, domain.ErrCompanyNotFound) {
				return 0, false, domain.ErrCompanyNotFound
			}

			return 0, false, err
		}

		if !isExist {
			return 0, false, domain.ErrCompanyNotFound
		}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, false, domain.ErrCreateUser
	}

	if tools.IsFullAccessSection(input.Sections) {
		return 0, false, domain.ErrSectionsNotAllowed
	}

	id, err := as.repo.User.Create(c.Request.Context(), domain.User{
		CompanyID:                input.CompanyId,
		Username:                 input.Username,
		Name:                     input.Name,
		Email:                    input.Email,
		Phone:                    input.Phone,
		PasswordHash:             string(hashedPass),
		LastLogin:                sql.NullTime{},
		IsActive:                 input.IsActive,
		Role:                     input.Role,
		Language:                 input.Language,
		Country:                  input.Country,
		IsApproved:               input.IsApproved,
		IsSendSystemNotification: input.IsSendSystemNotification,
		Sections:                 input.Sections,
		Position:                 input.Position,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return 0, false, domain.ErrUserAlreadyExists
		}

		return 0, false, domain.ErrCreateUser
	}

	return id, false, nil
}

func (as *AuthServices) RefreshTokens(c *gin.Context, refreshToken string) (domain.TokenResponse, error) {
	ip, err := tools.GetIPAddress(c)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrGetIpAddress
	}

	valid, session, err := as.isValidRefreshToken(c, refreshToken, ip)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrInvalidRefreshToken
	}

	if err = as.repo.Auth.DeleteToken(c.Request.Context(), session.UserID, refreshToken); err != nil {
		return domain.TokenResponse{}, domain.ErrRefreshToken
	}

	if !valid {
		return domain.TokenResponse{}, domain.ErrInvalidRefreshToken
	}

	sections, err := as.repo.User.GetSections(c.Request.Context(), session.UserID)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrRefreshToken
	}

	resp, err := as.createSession(c, domain.User{
		ID:        session.UserID,
		Role:      session.Role,
		CompanyID: session.CompanyID,
		Sections:  sections,
	})
	if err != nil {
		return domain.TokenResponse{}, domain.ErrRefreshToken
	}

	return resp, nil
}

func (as *AuthServices) ValidateAccessToken(c *gin.Context, token, userAgent, ip string) (domain.JWTInfo, bool, error) {
	return as.validateAccessToken(c, token, userAgent, ip)
}

func (as *AuthServices) createSession(ctx *gin.Context, user domain.User) (domain.TokenResponse, error) {
	var (
		res domain.TokenResponse
		err error
	)

	ip, err := tools.GetIPAddress(ctx)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	userAgent := tools.GetUserAgent(ctx)

	fingerprint, err := tools.GetHashedFingerprint(ip, userAgent)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	res.AccessToken, err = as.tokenManager.NewJWT(
		domain.JWTInfo{
			UserId:      user.ID,
			Role:        user.Role,
			CompanyId:   user.CompanyID,
			Fingerprint: fingerprint,
			Sections:    user.Sections,
		}, as.cfg.Auth.AccessTokenTTL)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	res.RefreshToken, err = as.tokenManager.NewRefreshToken()
	if err != nil {
		return domain.TokenResponse{}, err
	}

	res.ExpiresIn = time.Now().UTC().Add(as.cfg.Auth.RefreshTokenTTL).Unix()

	if err = as.repo.Auth.CreateToken(ctx, domain.RefreshSession{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Role:      user.Role,
		Token:     res.RefreshToken,
		ExpiresAt: time.Now().UTC().Add(as.cfg.Auth.RefreshTokenTTL),
		Ip:        ip,
	}); err != nil {
		return domain.TokenResponse{}, err
	}

	if err = as.repo.User.UpdateLastLogin(ctx, user.ID); err != nil {
		return domain.TokenResponse{}, err
	}

	return res, nil
}

func (as *AuthServices) isValidRefreshToken(c *gin.Context, refreshToken, ip string) (bool, domain.RefreshSession, error) {
	var valid bool

	if refreshToken == "" {
		return false, domain.RefreshSession{}, domain.ErrRefreshTokenNotFound
	}

	session, err := as.repo.Auth.GetSessionToken(c.Request.Context(), refreshToken)
	if err != nil {
		return false, domain.RefreshSession{}, err
	}

	if reflect.DeepEqual(session, domain.RefreshSession{}) {
		return false, domain.RefreshSession{}, domain.ErrRefreshTokenNotFound
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		if err = as.repo.Auth.DeleteToken(c.Request.Context(), session.UserID, session.Token); err != nil {
			return false, domain.RefreshSession{}, err
		}

		return false, domain.RefreshSession{}, domain.ErrExpiredRefreshToken
	}

	whiteIp, err := as.repo.Auth.GetUserWhiteIp(c.Request.Context(), session.UserID)
	if err != nil {
		return false, domain.RefreshSession{}, err
	}

	//whiteIp = append(whiteIp, domain.OfficeIp, domain.VPNIp) // todo по необходимости добавить белые ip чтобы можно было получить рефреш токен

	valid = tools.StringExists(whiteIp, ip)

	/*	if !valid {
		as.SuspiciousActivityLog(c.Request.Context(), session.UserID, session.CompanyID, telegram.Message{
			Header:    fmt.Sprintf("При обновлении рефреш токена обнаружен невалидный ip - %s", ip),
			Datetime:  time.Now().UTC().String(),
			Payload:   fmt.Sprintf("user id - %s, company id - %s", session.UserID, session.CompanyID),
			UserAgent: tools.GetUserAgent(c),
			Ip:        ip,
		})
	}*/ // todo добавить впоследствии уведомления в телеграм

	return valid, session, nil
}

func (as *AuthServices) validateAccessToken(ctx context.Context, token, userAgent, ip string) (domain.JWTInfo, bool, error) {
	info, err := as.tokenManager.Parse(token)
	if err != nil {
		return domain.JWTInfo{}, false, err
	}

	/*
		fingerprint, err := tools.GetHashedFingerprint(ip, userAgent)
		if err != nil {
			return domain.JWTInfo{}, false, err
		}

		if info.Fingerprint != fingerprint {
			/*		as.SuspiciousActivityLog(ctx, info.UserId, info.CompanyId, telegram.Message{
					Header: fmt.Sprintf("У пользователя с id - %s, company id - %s не совпадает fingerprint",
						info.UserId, info.CompanyId),
					Datetime:  time.Now().UTC().String(),
					Payload:   fmt.Sprintf("token - %s", token),
					UserAgent: userAgent,
					Ip:        ip,
				}) // todo добавить по необходимости уведомления в телеграм

			return domain.JWTInfo{}, false, domain.ErrInvalidAccessToken
		}*/ //todo вернуть в прод после тестирования

	return info, true, nil
}
