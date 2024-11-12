// 用户访问令牌access_token和刷新令牌refresh_token生成工具
// 注意: 目前登录令牌和刷新令牌包含的内容相同, 实际还需要优化
package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义生成Token时包含的信息
type CustomClaims struct {
	TkUID       string `json:"tk_uid"`
	TkUsername  string `json:"tk_username"`
	TkRoleLevel string `json:"tk_role_level"`
	jwt.RegisteredClaims
}

// 创建令牌用于存入客户端Cookie进行权限控制
// 客户端还可以便捷使用Token中的用户标识来查询个人数据
func CreateToken(jwtKey string, accessExp, refreshExp time.Duration, uid string, username string, roleLevel string) (*string, *string, error) {
	// 创建访问令牌
	accessTokenClaims := &CustomClaims{
		TkUID:       uid,
		TkUsername:  username,
		TkRoleLevel: roleLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExp)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, nil, err
	}
	// 创建刷新令牌
	refreshTokenClaims := &CustomClaims{
		TkUID:       uid,
		TkUsername:  username,
		TkRoleLevel: roleLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExp)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, nil, err
	}
	return &accessTokenString, &refreshTokenString, nil
}

// 校验JWT令牌
func VerifyToken(jwtKey string, tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("令牌错误: %w", err)
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("无效的令牌")
	}
	return claims, nil
}
