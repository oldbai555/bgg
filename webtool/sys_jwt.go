package webtool

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

var defaultJwtAppKey = []byte("lbserver-jwt-app-key")

const (
	TokenClaimsKey = "token-claim"
)

type Claims struct {
	UserId uint64
	jwt.StandardClaims
}

func ParseToken(token string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return defaultJwtAppKey, nil
	})
	return parseToken, claims, err
}

func GenToken(ctx context.Context, uid uint64) (string, error) {
	boundIP, err := utils.GetOutBoundIP()
	if err != nil {
		log.Errorf("err is : %v", err)
		return "", err
	}
	expireTime := time.Now().Add(time.Hour)
	claims := &Claims{
		UserId: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    boundIP,             // 签名颁发者
			Subject:   "lbuser_auth_token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(defaultJwtAppKey)
}
