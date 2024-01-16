package routes

import (
	"bookingHouses-server/models"
	"bookingHouses-server/storage"
	"bookingHouses-server/utils"
	"log"
	"strings"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx iris.Context) {
	var userInput RegisterUserInput
	if err := ctx.ReadJSON(&userInput); err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var newUser models.User

	if userExists, err := getAndHandleUserExists(&newUser, userInput.Email); err != nil {
		utils.CreateInternalServerError(ctx)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	} else if userExists {
		utils.CreateError(
			iris.StatusConflict,
			"Conflict",
			"Email already registered.",
			ctx,
		)

		return
	}

	hashedPassword, err := hashAndSaltPassword(userInput.Password)
	if err != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	newUser = models.User{
		FirstName:   userInput.FirstName,
		LastName:    userInput.LastName,
		Email:       userInput.Email,
		Password:    hashedPassword,
		SocialLogin: false,
	}
	if err := storage.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.Map{
		"ID":        newUser.ID,
		"firstName": newUser.FirstName,
		"lastName":  newUser.LastName,
		"email":     newUser.Email,
	})
}

func getAndHandleUserExists(user *models.User, email string) (exists bool, err error) {
	userExistsQuery := storage.DB.Where("email = ?", strings.ToLower(email)).Limit(1).Find(&user)

	if userExistsQuery.Error != nil {
		return false, userExistsQuery.Error
	}

	return userExistsQuery.RowsAffected > 0, nil
}

func hashAndSaltPassword(password string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

type RegisterUserInput struct {
	FirstName string `json:"firstName" validate:"required,max=255"`
	LastName  string `json:"lastName" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=255"`
}
