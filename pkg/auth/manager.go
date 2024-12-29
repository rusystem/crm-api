package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rusystem/crm-api/pkg/domain"
	"math/rand"
	"strconv"
	"time"
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
type TokenManager interface {
	NewJWT(info domain.JWTInfo, ttl time.Duration) (string, error)
	Parse(accessToken string) (domain.JWTInfo, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(info domain.JWTInfo, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"exp":      time.Now().UTC().Add(ttl).Unix(),
		"jti":      strconv.Itoa(int(info.UserId)),
		"aud":      info.Role,
		"sub":      strconv.Itoa(int(info.CompanyId)),
		"iss":      info.Fingerprint,
		"sections": info.Sections,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (domain.JWTInfo, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return domain.JWTInfo{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.JWTInfo{}, fmt.Errorf("error get user claims from token")
	}

	// Инициализируем секции как пустой срез
	var sections []string
	if sectionsClaim, ok := claims["sections"]; ok {
		// Проверяем, что это срез интерфейсов (как это хранится в JWT) и преобразуем в срез строк
		if sectionsSlice, ok := sectionsClaim.([]interface{}); ok {
			for _, section := range sectionsSlice {
				if str, ok := section.(string); ok {
					sections = append(sections, str)
				}
			}
		} else {
			return domain.JWTInfo{}, fmt.Errorf("error parsing sections from token")
		}
	}

	userId := claims["jti"].(string)
	companyId := claims["sub"].(string)

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return domain.JWTInfo{}, fmt.Errorf("error parsing userId from token")
	}

	companyIdInt, err := strconv.Atoi(companyId)
	if err != nil {
		return domain.JWTInfo{}, fmt.Errorf("error parsing companyId from token")
	}

	return domain.JWTInfo{
		UserId:      int64(userIdInt),
		Role:        claims["aud"].(string),
		CompanyId:   int64(companyIdInt),
		Fingerprint: claims["iss"].(string),
		Sections:    sections,
	}, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
