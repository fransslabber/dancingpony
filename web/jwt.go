package rest_api

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"
	sqldb "wo-infield-service/db"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateECDSAKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func GenerateJWT(userid uint32, device_uuid string, secretkey *ecdsa.PrivateKey) (string, error) {

	interval, err := sqldb.GetJWTSessionInterval()
	if err == nil {
		token := jwt.New(jwt.SigningMethodES256)
		claims := token.Claims.(jwt.MapClaims)

		claims["authorized"] = true
		claims["userid"] = strconv.FormatInt(int64(userid), 10)
		claims["device_uuid"] = device_uuid
		claims["exp"] = time.Now().Add(time.Minute * time.Duration(interval)).Unix()

		tokenString, err := token.SignedString(secretkey)

		if err != nil {
			return "", fmt.Errorf("Something Went Wrong: %s", err.Error())
		}
		return tokenString, nil
	} else {
		return "", err
	}
}

func AuthenticateJWT(tokenstr, device_uuid string, secretkey *ecdsa.PrivateKey) (uint32, error) {

	token, err := jwt.Parse(tokenstr,

		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Authorization failed. Wrong signing.")
			}
			return secretkey.Public(), nil
		},
		jwt.WithLeeway(5*time.Second))

	if err != nil {
		return 0, fmt.Errorf("Token cryptography failed: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["device_uuid"] != device_uuid {
			return 0, fmt.Errorf("JWT Auth failed")
		}
		userid, err := strconv.ParseUint(claims["userid"].(string), 10, 32)
		return uint32(userid), err
		// 	return fmt.Errorf("JWT Auth failed")
		// }
	}
	return 0, fmt.Errorf("JWT token invalid or could retreive claims map")
}
