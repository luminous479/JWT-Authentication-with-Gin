package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/luminous479/JWT-Authentication-with-Gin/database"
	"github.com/luminous479/JWT-Authentication-with-Gin/helpers"
	"github.com/luminous479/JWT-Authentication-with-Gin/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUserById() gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User

		err := helpers.MatchUserTypetoUid(c, id)
		if err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = userCollection.FindOne(ctx, bson.M{"user_id": id}).Decode(&user)
		if err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(200, user)
	}
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func VerifyPassword(userPass string, hashedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPass))
	if err != nil {
		return false, "Email or password is incorrect"
	}
	return true, ""
}

func SignUp() gin.HandlerFunc {
 
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		var user models.User	
		c.BindJSON(&user); err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for the email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		password := HashPassword(user.Password)
		user.Password = &password
		
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for the phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		
		password := HashPassword(user.Password)
		user.Password = password
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex( )
		
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken
		
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"result": resultInsertionNumber})
	}
}

func SignIn() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}		

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserId)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)
		err := userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)	

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while signing in"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refreshToken})
	}

}
func GetAllUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPage, err := strconv.Atoi(c.Query("recordPage"))
		if err != nil || recordPage <= 0 {
			recordPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		skip := (page - 1) * recordPage

		findOptions := options.Find().
			SetLimit(int64(recordPage)).
			SetSkip(int64(skip))

		cursor, err := userCollection.Find(
			ctx,
			bson.M{},
			findOptions,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while fetching users",
			})
			return
		}
		defer cursor.Close(ctx)

		var users []models.User

		for cursor.Next(ctx) {
			var user models.User

			if err := cursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error while decoding user data",
				})
				return
			}

			users = append(users, user)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}