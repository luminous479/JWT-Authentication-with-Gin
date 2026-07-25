package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("userType")
	if userType != role {
		return errors.New("unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypetoUid(c *gin.Context, role string) error {
	userType := c.GetString("userType")
	uid := c.GetString("uid")

	if userType != role {
		return errors.New("unauthorized to access this resource")
	}

	if userType == "USER" && uid != c.Param("id") {
		return errors.New("unauthorized to access this resource")
	}

	return nil
}
