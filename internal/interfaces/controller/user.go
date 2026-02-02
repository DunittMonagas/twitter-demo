package controller

import (
	"net/http"
	"strconv"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAllUsers(ctx gin.Context) ([]domain.User, error)
	GetUserByID(ctx gin.Context, id int) (domain.User, error)
	GetUserByEmail(ctx gin.Context, email string) (domain.User, error)
	GetUserByUsername(ctx gin.Context, username string) (domain.User, error)
	CreateUser(ctx gin.Context, user domain.User) (domain.User, error)
	UpdateUser(ctx gin.Context, id int, user domain.User) (domain.User, error)
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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)

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

	ctx.JSON(http.StatusOK, user)
}
