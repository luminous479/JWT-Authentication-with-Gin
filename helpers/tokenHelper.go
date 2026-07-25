package helpers

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/luminous479/JWT-Authentication-with-Gin/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserType  string
	Uid       string
	jwt.StandardClaims
}

var (
	secretKey = func() string {
		if key := os.Getenv("SECRET_KEY"); key != "" {
			return key
		}
		return "supersecretkey"
	}()
	userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
)

func GetSecretKey() []byte {
	return []byte(secretKey)
}

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
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(GetSecretKey())
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(GetSecretKey())
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(
	signedToken string,
	signedRefreshToken string,
	userID string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"token":         signedToken,
			"refresh_token": signedRefreshToken,
		},
	}

	filter := bson.M{
		"user_id": userID,
	}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating tokens:", err)
		return err
	}

	return nil
}
