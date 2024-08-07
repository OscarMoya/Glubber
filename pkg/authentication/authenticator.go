// Package: authentication defines authentication lifecycle of the service
// for now it only provides basic JWT lifecycle for drivers without checking
// anything in the driver db.

// Having this, showcases how the system can still process messages with a certain
// level of trust, without having to implement all the boilerplate for the authentication
// service. This can be done later on, when the system is more mature and the requirements
package authentication

import "github.com/dgrijalva/jwt-go"

// DriverClaims is the struct that will be used to store the claims of the driver JWT
// This struct will be used to generate and validate the JWT
// Even if there are some repeated fieds in the standard claims, it is better to have them
// in the struct to avoid any confusion or if the claims need to be extended in the future
type DriverClaims struct {
	DriverID  string `json:"driver_id"`
	Timestamp int64  `json:"timestamp"`
	jwt.StandardClaims
}

// DriverAuthenticator is the interface that will be used to generate and validate the driver JWT
// This interface will be used to generate and validate the JWT
// Later On, there may be a need to add more methods to this interface to handle more complex
// authentication methods
type DriverAuthenticator interface {
	GenerateDriverJWT(driverID string) (string, error)
	ValidateDriverJWT(tokenString string) (*DriverClaims, error)
}
