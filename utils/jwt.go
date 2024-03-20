package utils

import (
	"fmt"
	"net/http"
	"time"
	"whats-app-clone-service/models"

	"github.com/dgrijalva/jwt-go"
)

const jwtKey = "secret-sandal-jepit-key"

// Define a struct to represent the claims (payload) in the JWT token
type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Image    string `json:"image"`
	jwt.StandardClaims
}

func GenerateJWT(user models.User) (string, error) {

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Image:    user.Image,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJwt parses the JWT token and returns the claims (pyaload) if valid
func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

	if err != nil {
		fmt.Errorf("err ParseWithClaims: ", err)
		return nil, err
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse the JWT token
		claims, err := ParseJWT(tokenString)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set the claims as headers in the HTTP request
		r.Header.Set("X-User-ID", claims.ID)
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-Phone", claims.Phone)
		r.Header.Set("X-Image", claims.Image)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
