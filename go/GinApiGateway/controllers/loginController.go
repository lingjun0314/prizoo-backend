package controllers

import (
	"GinApiGateway/models"
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type Login struct{}

func setCustomUserClaims(userID string, role string) error {
	claims := map[string]interface{}{
		"role": role,
	}

	return models.AuthenticationClient.SetCustomUserClaims(context.Background(), userID, claims)
}

func (con *Login) SetUserRole(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "請先登入再執行操作",
		})
		return
	}

	role := "user"

	firebaseToken, ok := token.(*auth.Token)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token type",
		})
		return
	}

	if firebaseToken.UID == "TZr4TfGqTBM3mlkWCwAc4tkQzbF3" {
		role = "admin"
	}

	err := setCustomUserClaims(firebaseToken.UID, role)
	if err != nil {
		log.Println("Error by set custom user claims: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "設置身份失敗，請重試",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User role set to: " + role,
	})

}
