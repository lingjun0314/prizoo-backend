package middlewares

import (
	"GinApiGateway/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ActivityCache struct{}

func getOneFromRedis(ctx *gin.Context, cacheKey string) {
	data, err := models.RedisClient.Get(context.Background(), cacheKey).Result()

	if err == redis.Nil {
		return
	} else if err != nil {
		log.Println("error1: ", err)
		return
	}

	var result map[string]interface{}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println("error2: ", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
	ctx.Abort()
}

func getManyFromRedis(ctx *gin.Context, cacheKey string) {
	data, err := models.RedisClient.Get(context.Background(), cacheKey).Result()

	if err == redis.Nil {
		return
	} else if err != nil {
		log.Println("error1:", err)
		return
	}

	var result []map[string]interface{}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println("error2: ", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
	ctx.Abort()
}

func (con *ActivityCache) ActivityDetialMiddleware(ctx *gin.Context) {
	id := ctx.Param("id")
	cacheKey := "activity_" + id
	getOneFromRedis(ctx, cacheKey)
}

func (con *ActivityCache) AllActivityMiddleware(ctx *gin.Context) {
	cacheKey := "activity"
	getManyFromRedis(ctx, cacheKey)
}

func (con *ActivityCache) ActivityHistoryMiddleware(ctx *gin.Context) {
	id := ctx.Param("id")
	cacheKey := "activity_history_" + id
	getManyFromRedis(ctx, cacheKey)
}
