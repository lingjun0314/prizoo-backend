package controllers

import (
	"GinApiGateway/models"
	"GinApiGateway/proto/partner"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v5/logger"
)

type PartnerController struct{}

var partnerService partner.PartnerService

func init() {
	partnerService = partner.NewPartnerService("partner", models.MicroClient)
}

func (con *PartnerController) GetPartner(ctx *gin.Context) {
	res, err := partnerService.GetPartners(ctx, &partner.GetPartnersRequest{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "系統錯誤，請重試",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": res.Partner,
	})
	jsonData, err := json.Marshal(res.Partner)
	if err != nil {
		log.Println("error by get partner json: ", err)
		return
	}

	_, err = models.RedisClient.Set(ctx, "partner", string(jsonData), 0).Result()
	if err != nil {
		log.Println("Error set partner: ", err)
	}
}

func (con *PartnerController) GetPartnerById(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := partnerService.GetPartnerById(ctx, &partner.GetPartnerByIdRequest{
		Id: id,
	})
	if err != nil {
		logger.Error("error by get partner by id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "獲取失敗，請重試",
		})
		return
	}

	if res.StatusCode != http.StatusOK {
		ctx.JSON(int(res.StatusCode), gin.H{
			"message": res.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": res.Partner,
	})

	//	Set data to redis
	jsonData, err := json.Marshal(res.Partner)
	if err != nil {
		log.Println("error by get partner json: ", err)
		return
	}

	_, err = models.RedisClient.Set(ctx, "partner_"+id, string(jsonData), 0).Result()
	if err != nil {
		log.Println("Error set partner: ", err)
	}
}

func (con *PartnerController) CreatePartner(ctx *gin.Context) {
	//	Basic information
	brandName := ctx.PostForm("brandName")
	personInCharge := ctx.PostForm("personInCharge")
	companyName := ctx.PostForm("companyName")
	contactPerson := ctx.PostForm("contactPerson")

	//	Contact information
	tel := ctx.PostForm("tel")
	phone := ctx.PostForm("phone")
	email := ctx.PostForm("email")
	address := ctx.PostForm("address")

	//	Customer information
	sex := ctx.PostForm("sex")
	region := ctx.PostForm("region")
	ageStart := ctx.PostForm("ageStart")
	ageEnd := ctx.PostForm("ageEnd")

	if ageValid := models.CheckAgeValid(ageStart, ageEnd); ageValid != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": ageValid,
		})
		return
	}
	//	拼接年齡範圍字串
	ageRange := ageStart + "~" + ageEnd

	//	Social media
	instagram := ctx.PostForm("instagram")
	facebook := ctx.PostForm("facebook")
	threads := ctx.PostForm("threads")
	twitter := ctx.PostForm("twitter")

	res, err := partnerService.CreatePartner(ctx, &partner.CreatePartnerRequest{
		Partner: &partner.PartnerModule{
			BrandName:      brandName,
			PersonInCharge: personInCharge,
			CompanyName:    companyName,
			ContactPerson:  contactPerson,
			DeleteStatus: false,
			Contact: &partner.ContactModule{
				Tel:     tel,
				Phone:   phone,
				Email:   email,
				Address: address,
			},
			Customer: &partner.CustomerModule{
				AgeRange: ageRange,
				Sex:      sex,
				Region:   region,
			},
			SocialMedia: &partner.SocialMediaModule{
				Instagram: instagram,
				Facebook:  facebook,
				Threads:   threads,
				Twitter:   twitter,
			},
		},
	})
	if err != nil {
		logger.Error("error by create partner: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "新增失敗，請重試",
		})
		return
	}

	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
	})
}

func (con *PartnerController) UpdatePartner(ctx *gin.Context) {
	id := ctx.Param("id")

	//	Basic information
	brandName := ctx.PostForm("brandName")
	personInCharge := ctx.PostForm("personInCharge")
	companyName := ctx.PostForm("companyName")
	contactPerson := ctx.PostForm("contactPerson")

	//	Contact information
	tel := ctx.PostForm("tel")
	phone := ctx.PostForm("phone")
	email := ctx.PostForm("email")
	address := ctx.PostForm("address")

	//	Customer information
	sex := ctx.PostForm("sex")
	region := ctx.PostForm("region")
	ageStart := ctx.PostForm("ageStart")
	ageEnd := ctx.PostForm("ageEnd")

	if ageValid := models.CheckAgeValid(ageStart, ageEnd); ageValid != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": ageValid,
		})
		return
	}
	//	拼接年齡範圍字串
	ageRange := ageStart + "~" + ageEnd

	//	Social media
	instagram := ctx.PostForm("instagram")
	facebook := ctx.PostForm("facebook")
	threads := ctx.PostForm("threads")
	twitter := ctx.PostForm("twitter")

	res, err := partnerService.UpdatePartner(ctx, &partner.UpdatePartnerRequest{
		Id: id,
		Partner: &partner.PartnerModule{
			BrandName:      brandName,
			PersonInCharge: personInCharge,
			CompanyName:    companyName,
			ContactPerson:  contactPerson,
			Contact: &partner.ContactModule{
				Tel:     tel,
				Phone:   phone,
				Email:   email,
				Address: address,
			},
			Customer: &partner.CustomerModule{
				AgeRange: ageRange,
				Sex:      sex,
				Region:   region,
			},
			SocialMedia: &partner.SocialMediaModule{
				Instagram: instagram,
				Facebook:  facebook,
				Threads:   threads,
				Twitter:   twitter,
			},
		},
	})

	if err != nil {
		logger.Error("error by update partner: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "新增失敗，請重試",
		})
		return
	}

	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
	})

	//	Delete cache
	_, err = models.RedisClient.Del(ctx, "partner_"+id).Result()
	if err != nil {
		log.Println("Error by delete partner id in cache: ", err)
	}
	_, err = models.RedisClient.Del(ctx, "partner").Result()
	if err != nil {
		log.Println("Error by delete partner id in cache: ", err)
	}
}

func (con *PartnerController) DeletePartner(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := partnerService.DeletePartner(ctx, &partner.DeletePartnerRequest{
		Id: id,
	})

	if err != nil {
		logger.Error("error by delete partner: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "刪除失敗，請重試",
		})
		return
	}

	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
	})
	//	Delete cache
	_, err = models.RedisClient.Del(ctx, "partner").Result()
	if err != nil {
		log.Println("Error by delete partner id in cache: ", err)
	}
}
