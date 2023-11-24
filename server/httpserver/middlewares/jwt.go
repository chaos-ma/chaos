package middlewares

/**
* created by mengqi on 2023/11/21
 */

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type CustomClaims struct {
	ID          uint `json:"userid"`
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT(signKey string) *JWT {
	return &JWT{
		[]byte(signKey), //可以设置过期时间
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		var expErr *jwt.TokenExpiredError
		var nbfErr *jwt.TokenNotValidYetError
		var fmtErr *jwt.MalformedTokenError
		if errors.As(err, &fmtErr) {
			return nil, TokenMalformed
		} else if errors.As(err, &expErr) {
			return nil, TokenExpired
		} else if errors.As(err, &nbfErr) {
			return nil, TokenNotValidYet
		} else {
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = jwt.At(time.Now().Add(1 * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
