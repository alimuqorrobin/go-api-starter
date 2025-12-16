package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// IsValidPassword validates password strength
// Min 6 characters, at least 1 uppercase, 1 number
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper || hasNumber
}

// TrimString removes leading/trailing spaces
func TrimString(s string) string {
	return strings.TrimSpace(s)
}

// StringInSlice checks if string exists in slice
func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// GenerateRandomString generates random string of given length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

// ParseInt safely parses string to int
func ParseInt(s string, defaultVal int) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return num
}

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// NewPaginationParams creates pagination params from page and limit
func NewPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset calculates offset from page and limit
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
