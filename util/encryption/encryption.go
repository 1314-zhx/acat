// utils/password.go
package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 加密明文密码
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10（推荐值，平衡安全与性能）
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证明文密码是否匹配哈希值
func CheckPassword(password, hashedPassword string) bool {
	if hashedPassword == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
