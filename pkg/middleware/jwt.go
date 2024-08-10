package middleware

import (
	"context"
	"net/http"

	"github.com/OscarMoya/Glubber/pkg/authentication"
)

type ClaimsKeyType string

const ClaimsKey ClaimsKeyType = "claims"

func JWTMiddleware(next http.Handler, authenticator authentication.DriverAuthenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		claims, err := authenticator.ValidateDriverJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
