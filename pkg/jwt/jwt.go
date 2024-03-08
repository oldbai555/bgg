package jwt

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

var defaultJwtAppKey = []byte(TokenClaimsKey)

const (
	TokenClaimsKey = "token-claim"
	CtxWithClaim   = "claim"
)

var ErrGetLoginFail = lberr.NewInvalidArg("获取登录信息失败")

type Claims struct {
	Sid    string
	UserId uint64
	jwt.StandardClaims
}

func (c *Claims) GetSid() string {
	return c.Sid
}

func (c *Claims) GetUserId() uint64 {
	return c.UserId
}

func GetClaimsWithCtx(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(CtxWithClaim).(*Claims)
	if !ok {
		log.Errorf("err is : %v", ErrGetLoginFail)
		return nil, ErrGetLoginFail
	}
	return claims, nil
}

func SetClaimsWithCtx(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, CtxWithClaim, c)
}

func ParseToken(token string) (*jwt.Token, *Claims, error) {
	var claims Claims
	parseToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return defaultJwtAppKey, nil
	})
	return parseToken, &claims, err
}

func GenToken(ctx context.Context, uid uint64, sid string) (string, error) {
	boundIP, err := utils.UdpLocalIP()
	if err != nil {
		log.Errorf("err is : %v", err)
		return "", err
	}
	expireTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Sid:    sid,
		UserId: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    boundIP,        // 签名颁发者
			Subject:   TokenClaimsKey, //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(defaultJwtAppKey)
}
