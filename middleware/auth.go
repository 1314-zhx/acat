// middleware/auth.go
package middleware

import (
	"acat/util/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func AuthUserHTML() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("AuthUserHTML")

		// 解析 token 并获取 claims
		claims, ok := validTokenForHTML(c)
		if !ok || claims == nil {
			c.HTML(http.StatusOK, "login_prompt.html", gin.H{
				"redirect": c.Request.URL.Path,
			})
			c.Abort()
			return
		}
		fmt.Println("验证成功，用户ID:", claims.UserID)
		c.Set("claims", claims)
		c.Next()
	}
}

// 修改为返回 *auth.JwtClaims 和 bool
func validTokenForHTML(c *gin.Context) (*auth.JwtClaims, bool) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		return nil, false
	}

	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		zap.L().Warn("HTML 中间件：Token 验证失败", zap.Error(err))
		return nil, false
	}
	return claims, true
}
