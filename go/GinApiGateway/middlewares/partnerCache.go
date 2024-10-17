package middlewares

import "github.com/gin-gonic/gin"

type PartnerCache struct{}

func (con *PartnerCache) AllPartnerMiddleware(ctx *gin.Context) {
	cacheKey := "partner"
	getManyFromRedis(ctx, cacheKey)
}

func (con *PartnerCache) PartnerDetialMiddleware(ctx *gin.Context) {
	id := ctx.Param("id")
	cacheKey := "partner_" + id
	getOneFromRedis(ctx, cacheKey)
}
