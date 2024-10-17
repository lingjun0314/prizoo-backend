package controllers

import (
	"GinApiGateway/models"
	"GinApiGateway/proto/activity"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v5/logger"
)

type ActivityController struct {
	Data models.FileAccess
}

var activityService activity.ActivityService

func init() {
	activityService = activity.NewActivityService("activity", models.MicroClient)
}

func NewActivityController(file models.FileAccess) *ActivityController {
	return &ActivityController{
		Data: file,
	}
}
// @Summary 獲取所有活動的資料
// @Tags Activity
// @version 1.0
// @produce application/json
// @Success 200 {object} []map[string]interface{} "{"data": "活動內容"}"
// @Router /activity [get]
func (con *ActivityController) GetActivity(ctx *gin.Context) {
	res, err := activityService.GetActivity(context.Background(), &activity.GetActivityRequest{})
	if err != nil {
		logger.Error("micro error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "獲取失敗，請重試",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": res.Activity,
	})

	jsonData, err := json.Marshal(res.Activity)
	if err != nil {
		log.Println("error by get activity json: ", err)
		return
	}

	//	Set local random source
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	//	Set min and max hours
	min := 10
	max := 20
	//	Calculate expiration time
	expirationTime := (rng.Intn(max-min+1) + min) * int(time.Minute)

	_, err = models.RedisClient.Set(context.Background(), "activity", string(jsonData), time.Duration(expirationTime)).Result()
	if err != nil {
		log.Println("Error set activity: ", err)
	}
}
// @Summary 獲取特定 id 的活動
// @Tags Activity
// @version 1.0
// @produce application/json
// @param id path string true "活動 id"
// @Success 200 {object} map[string]interface{} "{"data": "活動內容"}"
// @Failure 500 {object} map[string]interface{} "{"message": "獲取失敗，請重試"}"
// @Failure 400 {object} map[string]interface{} "{"message": "錯誤訊息"}"
// @Router /activity/{id} [get]
func (con *ActivityController) GetActivityById(ctx *gin.Context) {
	//	Get param
	id := ctx.Param("id")

	//	Call function and get response
	res, err := activityService.GetActivityById(context.Background(), &activity.GetActivityByIdRequest{
		Id: id,
	})
	if err != nil {
		logger.Error("micro error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "獲取失敗，請重試",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": res.Data,
	})

	jsonData, err := json.Marshal(res.Data)
	if err != nil {
		log.Println("error by GetActivityById json: ", err.Error())
		return
	}

	//	Set local random source
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	//	Set min and max hours
	min := 10
	max := 24
	//	Calculate expiration time
	expirationTime := (rng.Intn(max-min+1) + min) * int(time.Hour)

	_, err = models.RedisClient.Set(context.Background(), "activity_"+id, string(jsonData), time.Duration(expirationTime)).Result()
	if err != nil {
		log.Println("Error set activity by id: ", err)
	}
}

// @Summary 新建一個活動
// @Tags Activity
// @version 1.0
// @produce application/json
// @param title formData string true "活動名稱"
// @param detial formData string true "活動細節介紹"
// @param startTime formData string true "活動開始時間(unix time stamp)"
// @param endTime formData string true "活動結束時間(unix time stamp)"
// @param partner formData string true "提供活動的廠商名稱"
// @param prizeName formData string true "獎品名稱"
// @param image formData file true "獎品的圖片"
// @param addPrizeName formData string true "加碼的獎品名稱"
// @param addPrizeImage formData file true "加碼獎品的圖片"
// @param addPrizeTicketAmount formData string true "多少票券加入之後開啟加碼"
// @Success 200 {object} activity.CreateActivityResponse
// @Failure 400 {object} activity.CreateActivityResponse "輸入錯誤"
// @Failure 500 {object} activity.CreateActivityResponse "系統錯誤"
// @Router /activity [post]
func (con *ActivityController) CreateActivity(ctx *gin.Context) {
	title := ctx.PostForm("title")
	detial := ctx.PostForm("detial")

	startTime, err := strconv.Atoi(ctx.PostForm("startTime"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error time format",
		})
		return
	}

	endTime, err := strconv.Atoi(ctx.PostForm("endTime"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error time format",
		})
		return
	}

	if timeValid := models.CheckTimeValid(startTime, endTime); timeValid != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": timeValid,
		})
		return
	}

	partner := ctx.PostForm("partner")
	prizeName := ctx.PostForm("prizeName")
	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error by receive file 1",
		})
		return
	}
	addPrizeName := ctx.PostForm("addPrizeName")
	addPrizeTicketAmount, err := strconv.Atoi(ctx.PostForm("addPrizeTicketAmount"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error ticket amount",
		})
		return
	}
	addPrizeImage, err := ctx.FormFile("addPrizeImage")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error by receive file 2",
		})
		return
	}

	//	Check file format
	if !models.VerifyFileExt(image) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type",
		})
		return
	}
	if !models.VerifyFileExt(addPrizeImage) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type",
		})
		return
	}

	//	Upload file
	imgPath, objName, err := con.Data.UploadFile(image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occur while upload file: " + err.Error(),
		})
		return
	}
	addImgPath, objName1, err := con.Data.UploadFile(addPrizeImage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occur while upload file: " + err.Error(),
		})
		return
	}

	//	Create activity service
	res, err := activityService.CreateActivity(context.Background(), &activity.CreateActivityRequest{
		Title:     title,
		Detial:    detial,
		StartTime: int64(startTime),
		EndTime:   int64(endTime),
		Partner:   partner,
		Prize: &activity.PrizeModule{
			Name:                        prizeName,
			ImagePath:                   imgPath,
			AdditionalPrizeName:         addPrizeName,
			AdditionalPrizeTicketAmount: int32(addPrizeTicketAmount),
			AdditionalPrizeImagePath:    addImgPath,
		},
	})

	if err != nil {
		logger.Error("micro error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "新增失敗，請重試",
			"success": false,
		})

		//	Delete file in file server
		con.Data.DeleteFile(objName)
		con.Data.DeleteFile(objName1)

		return
	}

	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
		"success": res.Success,
	})
}

// @Summary 刪除特定 id 的活動
// @Tags Activity
// @version 1.0
// @produce application/json
// @param id path string true "活動 id"
// @Success 200 {object} activity.DeleteActivityResponse
// @Failure 400 {object} activity.DeleteActivityResponse "錯誤的 id"
// @Failure 500 {object} activity.DeleteActivityResponse "系統錯誤"
// @Router /activity/{id} [delete]
func (con *ActivityController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := activityService.DeleteActivity(ctx, &activity.DeleteActivityRequest{Id: id})
	if err != nil {
		logger.Error("micro error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "刪除失敗，請重試",
			"success": false,
		})
		return
	}

	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
		"success": res.Success,
	})
}

// @Summary 更新指定 id 的活動內容
// @Tags Activity
// @version 1.0
// @produce application/json
// @param id path string true "要更新的活動 id"
// @param title formData string true "活動名稱"
// @param detial formData string true "活動細節介紹"
// @param startTime formData string true "活動開始時間(unix time stamp)"
// @param endTime formData string true "活動結束時間(unix time stamp)"
// @param partner formData string true "提供活動的廠商名稱"
// @param prizeName formData string true "獎品名稱"
// @param image formData file true "獎品的圖片"
// @param addPrizeName formData string true "加碼的獎品名稱"
// @param addPrizeTicketAmount formData string true "多少票券加入之後開啟加碼"
// @Success 200 {object} activity.UpdateActivityResponse
// @Failure 400 {object} activity.UpdateActivityResponse "錯誤的 id"
// @Failure 500 {object} activity.UpdateActivityResponse "系統錯誤"
// @Router /activity/{id} [patch]
func (con *ActivityController) UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	title := ctx.PostForm("title")
	detial := ctx.PostForm("detial")
	startTime, err := strconv.Atoi(ctx.PostForm("startTime"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error time format",
		})
		return
	}
	endTime, err := strconv.Atoi(ctx.PostForm("endTime"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error time format",
		})
		return
	}

	timeValid := models.CheckTimeValid(startTime, endTime)
	if timeValid != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": timeValid,
		})
		return
	}

	//	Get prize data
	var imgPath string
	var addImgPath string
	prizeName := ctx.PostForm("prizeName")
	addPrizeTicketAmount, err := strconv.Atoi(ctx.PostForm("addPrizeTicketAmount"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Prize ticket amount is not a number",
		})
		return
	}
	image, err := ctx.FormFile("image")
	if err != nil {
		imgPath = ""
	}
	addPrizeName := ctx.PostForm("addPrizeName")
	addImage, err := ctx.FormFile("addImage")
	if err != nil {
		addImgPath = ""
	}

	var objName, objName1 string

	//	Check file valid
	if image != nil {
		if !models.VerifyFileExt(image) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid file format",
			})
			return
		}

		imgPath, objName, err = con.Data.UploadFile(image)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occur while upload file: " + err.Error(),
			})
			return
		}
	}

	if addImage != nil {
		if !models.VerifyFileExt(addImage) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid file format",
			})
			return
		}

		addImgPath, objName1, err = con.Data.UploadFile(addImage)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occur while upload file: " + err.Error(),
			})
			return
		}
	}

	res, err := activityService.UpdateActivity(ctx, &activity.UpdateActivityRequest{
		Title:     title,
		Detial:    detial,
		StartTime: int64(startTime),
		EndTime:   int64(endTime),
		Prize: &activity.PrizeModule{
			Name:                        prizeName,
			ImagePath:                   imgPath,
			AdditionalPrizeName:         addPrizeName,
			AdditionalPrizeImagePath:    addImgPath,
			AdditionalPrizeTicketAmount: int32(addPrizeTicketAmount),
		},
		Id: id,
	})
	if err != nil {
		log.Println(err)
		//	Delete file
		con.Data.DeleteFile(objName)
		con.Data.DeleteFile(objName1)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新失敗，請重試",
		})
		return
	}

	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
		"success": res.Success,
	})

	//	Delete cache
	_, err = models.RedisClient.Del(context.Background(), "activity_"+id).Result()
	if err != nil {
		log.Println("Error by delete history in cache: ", err)
	}
	_, err = models.RedisClient.Del(context.Background(), "activity_history_"+id).Result()
	if err != nil {
		log.Println("Error by delete history in cache: ", err)
	}
}

// @Summary 查看特定活動 id 的修改歷史記錄
// @Tags Activity
// @version 1.0
// @produce application/json
// @param id path string true "要查看的活動 id"
// @Success 200 {object} map[string]interface{} "{"data": "歷史記錄陣列"}"
// @Router /activity/history/{id} [get]
func (con *ActivityController) GetActivityUpdateHistoryById(ctx *gin.Context) {
	id := ctx.Param("id")
	
	res, err := activityService.GetActivityUpdateHistoryById(ctx, &activity.GetActivityUpdateHistoryByIdRequest{
		Id: id,
	})
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "查詢失敗，請重試",
		})
		return
	}

	//	Return result
	ctx.JSON(int(res.StatusCode), gin.H{
		"message": res.Message,
		"data":    res.Data,
	})

	if res.StatusCode != http.StatusOK {
		return
	}
	//	Store result in cache
	jsonData, err := json.Marshal(res.Data)
	if err != nil {
		log.Println("error by GetActivityById json: ", err.Error())
		return
	}

	_, err = models.RedisClient.Set(context.Background(), "activity_history_"+id, string(jsonData), 0).Result()
	if err != nil {
		log.Println("Error set activity by id: ", err)
	}
}
