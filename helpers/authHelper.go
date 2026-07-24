package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("userType")

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	return nil
}

func MatchUserTypetoUid(c *gin.Context, role string) (err error) {
	userType := c.GetString("userType")
	uid := c.GetString("uid")

	err = nil

	if userType != role {

		err = errors.New("unauthorized to access this resource")
		return err

	}

	if userType == "user" && uid != c.Param("id") {
		err = errors.New("unauthorized to access this resource")
		return err
	}
    err = CheckUserType(c, userType)
	return err

}

