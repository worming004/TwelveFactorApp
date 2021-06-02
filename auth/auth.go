package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(email string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

type CreateTokenRequest struct {
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (j *JwtWrapper) CreateToken(w http.ResponseWriter, r *http.Request) {
	var request CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Password != os.Getenv("TWELVE_AUTH_PASSWORD") {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	tok, err := j.GenerateToken(request.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := CreateTokenResponse{
		Token: tok,
	}

	json.NewEncoder(w).Encode(res)
}

func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}

func GetDefaultJwtWrapper() JwtWrapper {
	return JwtWrapper{
		SecretKey:       os.Getenv("TWELVE_AUTH_SECRET"),
		Issuer:          "Mail",
		ExpirationHours: 24,
	}
}
