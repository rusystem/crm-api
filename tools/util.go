package tools

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/crm-api/pkg/domain"
	"github.com/segmentio/ksuid"
	"regexp"
	"strings"
	"time"
)

const fingerprintSalt = "pz_n3wdc9fan"

func GetEmailDomain(email string) (string, error) {
	components := strings.Split(email, "@")
	if len(components) != 2 {
		return "", errors.New("incorrect email")
	}

	if components[0] == "" || components[1] == "" {
		return "", errors.New("incorrect email")
	}

	return components[1], nil
}

func GenerateUUID() (string, error) {
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func GetIPAddress(c *gin.Context) (string, error) {
	ip := c.GetHeader("X-Real-Ip")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.ClientIP()
	}

	arr := strings.Split(ip, ":")
	if arr[0] == "" {
		return "", errors.New("failed to get ip address")
	}

	return arr[0], nil
}

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func GetHashedFingerprint(ip, userAgent string) (string, error) {
	if ip == "" {
		return "", errors.New("ip can`t be zero")
	}

	if userAgent == "" {
		return "", errors.New("user-agent can`t be zero")
	}

	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s%s", ip, userAgent)))

	return fmt.Sprintf("%x", hash.Sum([]byte(fingerprintSalt))), nil
}

func DateTimeFormat() string {
	return "2006-01-02 15:04:05"
}

func ParseTime(t string) (time.Time, error) {
	if t == "" {
		return time.Time{}, errors.New("parse time: nullable time string")
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", t)

	if parsedTime.IsZero() {
		return time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	}

	return parsedTime, err
}

func StringExists(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func DecodeBase64(input string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	decodedString := string(decodedBytes)

	return decodedString, nil
}

func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// IsFullAccessSection returns true if section is full access
func IsFullAccessSection(sections []string) bool {
	for _, section := range sections {
		if section == domain.SectionFullAllAccess {
			return true
		}
	}

	return false
}

// IsAllowedRole returns true if role is allowed
func IsAllowedRole(role string) bool {
	if role == "" {
		return true
	}

	for _, r := range domain.AllowedRoles {
		if r == role {
			return true
		}
	}

	return false
}

func RemoveFullAccessSection(slice []domain.Section, strToRemove string) []domain.Section {
	var result []domain.Section
	for _, str := range slice {
		if str.Name != strToRemove {
			result = append(result, str)
		}
	}
	return result
}

func IsValidCountryCode(code string) bool {
	re := regexp.MustCompile(`^[A-Z]{2}$`)
	return re.MatchString(code)
}

// IsValidAdminCode выполняет базовую валидацию adminCode
func IsValidAdminCode(adminCode string) bool {
	// Удаляем лишние пробелы
	adminCode = strings.TrimSpace(adminCode)

	// Проверка на пустое значение
	if adminCode == "" {
		return false
	}

	return true
}
