package controllers

import (
	"MoZaki-Organization-Manager/pkg/database/mongodb/models"
	"MoZaki-Organization-Manager/pkg/database/mongodb/repository"
	"MoZaki-Organization-Manager/pkg/utils"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) (err error) {
	var user models.User

	if err = c.BindJSON(&user); err != nil {
		return
	}

	err = Validate.Struct(user)

	if err != nil {
		return
	}

	count, err := repository.CountUsersByEmail(user.Email)
	if err != nil {
		err = errors.New("error occured while checking the email")
		return
	}

	if count > 0 {
		err = errors.New("this email already exists")
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, err := utils.GenerateAllTokens(*user.Email, *user.Name, user.User_id)

	if err != nil {
		err = errors.New("error occured generating tokens")
		return
	}
	user.Token = &token
	user.Refresh_Token = &refreshToken
	err = repository.CreateUser(user)

	return

}

// Function that extracts user data from context and compares it to database if its valid, Tokens are regenrated for the user and are returned
func Login(c *gin.Context) (token string, refreshToken string, err error) {
	var user models.User
	if err = c.BindJSON(&user); err != nil {
		return
	}

	foundUser, err := repository.GetUserByEmail(user.Email)

	if err != nil {
		return
	}
	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

	if !passwordIsValid {
		err = errors.New(msg)
		return
	}
	token, refreshToken, err = utils.GenerateAllTokens(*foundUser.Email, *foundUser.Name, foundUser.User_id)

	if err != nil {
		return
	}
	_, _, err = utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	return

}

// Function that accepts context, extracts refreshToken from it and returns new authorization token and refresh token.
func Refresh_Token(c *gin.Context) (token string, refreshToken string, err error) {
	// temporary struct to capture the refresh token from context
	type Input struct {
		Token_ string `json:"refresh_token"`
	}

	var in Input
	if err = c.BindJSON(&in); err != nil {
		return
	}

	refreshToken = in.Token_
	claims, err := utils.ValidateToken(refreshToken)

	if err != nil {
		return
	}
	//Capture user details extracted from the refresh token provided
	name, email, uid := claims["Name"].(string), claims["Email"].(string), claims["Uid"].(string)
	token, refreshToken, err = utils.GenerateAllTokens(name, email, uid)

	if err != nil {
		return
	}
	token, refreshToken, err = utils.UpdateAllTokens(token, refreshToken, uid)

	return

}

// Function that hashes passwords for users
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// Function that compares both provided password and saved password to verify the user
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Email or password is incorrect"
		check = false
	}
	return check, msg
}
