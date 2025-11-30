package auth

import (
	"acat/setting"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

// JwtClaims 自定义声明
type JwtClaims struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Identity int    `json:"identity"` // 1: 管理员, 0 或其他: 普通用户
	jwt.RegisteredClaims
}

func GenerateToken(id uint, userName string, identity int) (string, error) {
	var expiresTime time.Time

	if identity == 1 {
		expiresTime = time.Now().Add(30 * time.Minute) // 管理员：30 分钟
	} else {
		expiresTime = time.Now().Add(7 * 24 * time.Hour) // 普通用户：7 天
	}
	claims := JwtClaims{
		UserID:   id,
		Name:     userName,
		Identity: identity,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "acat",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString([]byte(setting.Conf.JwtSecret))
	if err != nil {
		zap.L().Error("生成 Token 失败", zap.Error(err))
		return "", err
	}
	return token, nil
}

// auth/token.go （或 auth/auth.go）
func ParseToken(tokenStr string) (*JwtClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token 不能为空")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(setting.Conf.JwtSecret), nil
	})
	if err != nil {
		zap.L().Warn("解析 Token 失败", zap.Error(err), zap.String("token", tokenStr))
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的 Token 声明或 Token 已失效")
}
