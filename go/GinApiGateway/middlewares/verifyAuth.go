package middlewares

import (
	"GinApiGateway/models"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func VerifyIDToken(ctx *gin.Context) {
	//	Get Authorization data
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "請先登入再執行此操作",
		})
		ctx.Abort()
		return
	}

	//	Verify token
	token, err := models.AuthenticationClient.VerifyIDToken(ctx, authHeader)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "ID token 無效",
		})
		ctx.Abort()
		return
	}

	//	Set token in context
	ctx.Set("token", token)
	ctx.Next()
}

func AdminAPIVerify(ctx *gin.Context) {
	token, ok := ctx.Get("token")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "請先登入",
		})
		ctx.Abort()
		return
	}

	firebaseToken, ok := token.(*auth.Token)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token format",
		})
		ctx.Abort()
		return
	}

	role, ok := firebaseToken.Claims["role"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role type",
		})
		ctx.Abort()
		return
	}

	if role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "需要 admin 權限",
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}
