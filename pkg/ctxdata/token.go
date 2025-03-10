package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "hnz.com.cn"

func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	// 过期时间
	claims["exp"] = iat + seconds
	// 签发时间
	claims["iat"] = iat
	claims[Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
