package helpers

import (
	"log"
	"os"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserType  string
	Uid       string
	jwt.StandardClaims
}

var secretKey = os.Getenv("SECRET_KEY")

func GenerateAllTokens(
	email string,
	firstName string,
	lastName string,
	userType string,
	uid string,
) (signedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(168 * time.Hour).Unix(), // 7 days
		},
	}

	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error generating access token:", err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	).SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", err
	}

	return token, refreshToken, nil
}