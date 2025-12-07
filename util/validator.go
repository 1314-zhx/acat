package util

import (
	"regexp"
)

// 校验是不是邮件
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}

// 校验是不是电话
func IsPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	return regexp.MustCompile(pattern).MatchString(phone)
}
