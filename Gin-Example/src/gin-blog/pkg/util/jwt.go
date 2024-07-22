package util

import (
	"Gin-Example/src/gin-blog/pkg/setting"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(username, password string) (string, error) {

	// 获取当前时间
	nowTime := time.Now()
	// 设置过期时间为3小时
	expireTime := nowTime.Add(3 * time.Hour)

	// 创建 Claims 对象
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	// 使用 Claims 对象创建 tokenClaims
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用 jwtSecret 对 tokenClaims 进行签名
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err

}

// ParseToken 解析JWT
func ParseToken(token string) (*Claims, error) {

	// 解析JWT
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回jwt密钥
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 如果tokenClaims存在，并且有效
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			// 返回Claims
			return claims, nil
		}
	}

	return nil, err
}
