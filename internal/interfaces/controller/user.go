package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"twitter-demo/internal/interfaces/dto"
	"twitter-demo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAllUsers(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	// GetUserByEmail(ctx gin.Context, email string) (domain.User, error)
	// GetUserByUsername(ctx gin.Context, username string) (domain.User, error)
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
}

type User struct {
	userUsecase usecase.UserUsecase
}

func NewUser(userUsecase usecase.UserUsecase) User {
	return User{
		userUsecase: userUsecase,
	}
}

func (u User) GetAllUsers(ctx *gin.Context) {

	users, err := u.userUsecase.GetAllUsers(ctx)

	fmt.Println(users)
	fmt.Println(err)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.ToUserResponse(user)
	}

	ctx.JSON(http.StatusOK, userResponses)

}

func (u User) GetUserByID(ctx *gin.Context) {

	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := u.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (u User) CreateUser(ctx *gin.Context) {

	createUserRequest := dto.CreateUserRequest{}
	if err := ctx.ShouldBindJSON(&createUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := dto.ToUserDomain(createUserRequest)
	newUser, err := u.userUsecase.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToUserResponse(newUser))
}

func (u User) UpdateUser(ctx *gin.Context) {

	updateUserRequest := dto.UpdateUserRequest{}
	if err := ctx.ShouldBindJSON(&updateUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idString := ctx.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user := dto.ToUpdateUserDomain(updateUserRequest)
	updatedUser, err := u.userUsecase.UpdateUser(ctx, id, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToUserResponse(updatedUser))
}
