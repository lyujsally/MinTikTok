package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var mySecret = []byte("字节和心脏都要一直跳动")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func GenToken(userID int64, userName string) (string, error) {
	c := MyClaims{
		userID,
		userName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expires")) * time.Hour).Unix(),
			Issuer: viper.GetString("name"),
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, err
	}
	return nil, errors.New("invalid token")
}
