package controllers

import (
	"context"
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
