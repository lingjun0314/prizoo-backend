package routers

import (
	"GinApiGateway/controllers"
	"GinApiGateway/middlewares"
	"GinApiGateway/models"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

var (
	activityController *controllers.ActivityController
	activityCache      *middlewares.ActivityCache
	loginController    *controllers.Login
	partnerController  *controllers.PartnerController
	partnerCache       *middlewares.PartnerCache
)

func init() {
	activityController = controllers.NewActivityController(&models.FireStorage{})
}

func InitApiRouters(r *gin.Engine) {
	//	swagger api
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// login router
	r.POST("/login", loginController.SetUserRole)

	//	Activity routers
	r.GET("/activity", activityCache.AllActivityMiddleware, activityController.GetActivity)
	r.POST("/activity", activityController.CreateActivity)
	r.GET("/activity/:id", activityCache.ActivityDetialMiddleware, activityController.GetActivityById)
	r.PATCH("/activity/:id", activityController.UpdateActivity)
	r.DELETE("/activity/:id", activityController.DeleteActivity)
	r.GET("/activity/history/:id", activityCache.ActivityHistoryMiddleware, activityController.GetActivityUpdateHistoryById)

	//	Partner routers
	r.GET("/partner", partnerCache.AllPartnerMiddleware, partnerController.GetPartner)
	r.POST("/partner", partnerController.CreatePartner)
	r.GET("/partner/:id", partnerCache.PartnerDetialMiddleware, partnerController.GetPartnerById)
	r.PATCH("/partner/:id", partnerController.UpdatePartner)
	r.DELETE("/partner/:id", partnerController.DeletePartner)

}
