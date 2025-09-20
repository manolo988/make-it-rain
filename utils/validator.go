package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(strings.ToLower(email))
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func ValidatePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func ValidateSortOrder(order string) string {
	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		return "desc"
	}
	return order
}
