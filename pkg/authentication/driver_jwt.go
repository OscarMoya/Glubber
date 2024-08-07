package authentication

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// driverSigningKey is the secret key used to sign the driver JWT this will become the location of the key
// via encrypted volume or the secret name via a secret manager. This can also be part of the struct later on
const driverSigningKey = "my_secret123" // TODO: Change this to a more secure secret key

type JWTDriverAuthenticationService struct {
}

func (d *JWTDriverAuthenticationService) GenerateDriverJWT(driverID string) (string, error) {
	claims := DriverClaims{
		DriverID:  driverID,
		Timestamp: time.Now().Unix(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(driverSigningKey)
}

func (d *JWTDriverAuthenticationService) ValidateDriverJWT(tokenString string) (*DriverClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return driverSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*DriverClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
