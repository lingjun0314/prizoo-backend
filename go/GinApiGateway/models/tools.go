package models

import (
	"mime/multipart"
	"path"
	"strconv"
	"time"
)

func VerifyFileExt(file *multipart.FileHeader) bool {
	extName := path.Ext(file.Filename)
	//	Set allow file extention
	allowExtMap := map[string]bool{
		".jpg":  true,
		".JPG":  true,
		".png":  true,
		".jpeg": true,
		".JPEG": true,
		".PNG":  true,
	}

	if _, ok := allowExtMap[extName]; !ok {
		return false
	}

	return true
}

func CheckTimeValid(startTime, endTime int) string {
	if time.Unix(int64(startTime), 0).Before(time.Now()) {
		return "請選擇晚於現在的開始時間"
	}
	if time.Unix(int64(endTime), 0).Before(time.Unix(int64(startTime), 0)) {
		return "請選擇晚於開始時間的結束時間"
	}
	return ""
}

func CheckAgeValid(startAge, endAge string) string {
	ageStart, err := strconv.Atoi(startAge)
	if err != nil {
		return "錯誤的起始年齡類型"
	}
	ageEnd, err := strconv.Atoi(endAge)
	if err != nil {
		return "錯誤的結束年齡類型"
	}
	if ageEnd < ageStart {
		return "請輸入大於起始年齡的結束年齡"
	}
	return ""
}
