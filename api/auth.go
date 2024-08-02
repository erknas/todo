package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeze322/todo/lib"
)

type signRequest struct {
	Password string `json:"password"`
}

type signResponse struct {
	Token string `json:"token"`
}

func (s *Server) handleSign(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req signRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if req.Password != s.password {
		return lib.WriteJSON(w, http.StatusUnauthorized, lib.ApiErr{Error: "invalid password"})
	}

	token, err := createJWT(req.Password)
	if err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, signResponse{Token: token})
}

func withJWTAuth(next http.HandlerFunc, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(password) > 0 {

			var tokenString string

			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}

			token, err := validateJWT(tokenString)
			if err != nil {
				permissionDenied(w)
				return
			}

			if !token.Valid {
				permissionDenied(w)
				return
			}

			passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				permissionDenied(w)
				return
			}

			if passwordHash != claims["passwordHash"] {
				permissionDenied(w)
				return
			}

		}
		next(w, r)
	}
}

func createJWT(password string) (string, error) {
	passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	claims := &jwt.MapClaims{
		"passwordHash": passwordHash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("TODO_SECRET")

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("TODO_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	lib.WriteJSON(w, http.StatusUnauthorized, lib.ApiErr{Error: "authentification required"})
}
