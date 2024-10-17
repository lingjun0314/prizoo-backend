package handler

import (
	"activity/database"
	"activity/models"
	pb "activity/proto/activity"
	"activity/proto/partner"
	"activity/proto/prize"
	"context"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"go-micro.dev/v5/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/structpb"
)

type Activity struct{}

var (
	prizeService   prize.PrizeService
	partnerService partner.PartnerService
)

func InitService() {
	prizeService = prize.NewPrizeService("prize", models.MicroClient)
	partnerService = partner.NewPartnerService("partner", models.MicroClient)
}

func (e *Activity) GetActivity(ctx context.Context, req *pb.GetActivityRequest, res *pb.GetActivityResponse) error {
	var result []*structpb.Struct

	iter := database.FirestoreClient.Collection("activity").Where("DeleteStatus", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		//	End for condition
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		data, err := models.FromJsonToPbStruct(ctx, 0, doc.Data())
		if err != nil {
			return err
		}
		data.Fields["Id"] = structpb.NewStringValue(doc.Ref.ID)
		result = append(result, data)
	}

	res.Activity = result

	return nil
}

func (e *Activity) GetActivityById(ctx context.Context, req *pb.GetActivityByIdRequest, rsp *pb.GetActivityByIdResponse) error {
	//	Get document id
	id := req.GetId()

	if id == "" {
		return errors.New("no id in request")
	}

	//	Find data in firestore
	doc, err := database.FirestoreClient.Collection("activity").Doc(id).Get(ctx)
	if err != nil {
		return err
	}

	//	Change data to google.protobuf.Struct
	data, err := models.FromJsonToPbStruct(ctx, 0, doc.Data())
	if err != nil {
		return err
	}

	data.Fields["Id"] = structpb.NewStringValue(id)

	rsp.Data = data

	return nil
}

func (e *Activity) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest, res *pb.CreateActivityResponse) error {
	// Define channel
	prizePathChan := make(chan *prize.CreatePrizeResponse, 1)
	errChan := make(chan error, 1)

	//	Use goroutine to call prize service(save about 60% response time)
	go func() {
		prizePath, err := prizeService.CreatePrize(ctx, &prize.CreatePrizeRequest{
			Name:                        req.Prize.Name,
			ImagePath:                   req.Prize.ImagePath,
			AdditionalPrizeName:         req.Prize.AdditionalPrizeName,
			AdditionalPrizeTicketAmount: req.Prize.AdditionalPrizeTicketAmount,
			AdditionalPrizeImagePath:    req.Prize.AdditionalPrizeImagePath,
			Version:                     1,
		})
		if err != nil {
			errChan <- err
			return
		}
		prizePathChan <- prizePath
	}()

	// Get partner id
	partnerRes, err := partnerService.GetPartnerByName(ctx, &partner.GetPartnerByNameRequest{
		BrandName: req.Partner,
	})
	if err != nil {
		return err
	}
	//	No this partner
	if partnerRes.StatusCode != http.StatusOK {
		res.Message = partnerRes.Message
		res.StatusCode = partnerRes.StatusCode
		return nil
	}

	partnerRef := database.FirestoreClient.Collection("partner").Doc(partnerRes.Message)

	// Handle goroutine response result
	var prizePath *prize.CreatePrizeResponse
	select {
	case prizePath = <-prizePathChan:
		if prizePath.StatusCode != http.StatusOK {
			res.Message = prizePath.Message
			res.StatusCode = prizePath.StatusCode
			return nil
		}
	case err = <-errChan:
		if err != nil {
			return err
		}
	}
	// 獲取獎品文件引用
	prizeRef := database.FirestoreClient.Collection("prize").Doc(prizePath.Message)

	// 創建活動
	docRef, _, err := database.FirestoreClient.Collection("activity").Add(ctx, models.Activity{
		Title:        req.Title,
		Detial:       req.Detial,
		StartTime:    time.Unix(req.StartTime, 0),
		EndTime:      time.Unix(req.EndTime, 0),
		Partner:      partnerRef,
		Prize:        prizeRef,
		DeleteStatus: false,
		Version:      1,
	})
	if err != nil {
		return err
	}

	// 返回結果
	res.Success = true
	res.StatusCode = http.StatusOK
	res.Message = docRef.ID
	return nil
}

func (e *Activity) DeleteActivity(ctx context.Context, req *pb.DeleteActivityRequest, res *pb.DeleteActivityResponse) error {
	//	Query id exist or not
	docRef := database.FirestoreClient.Collection("activity").Doc(req.Id)
	_, err := docRef.Get(ctx)
	if err != nil {
		res.Success = false
		res.Message = "找不到 id"
		res.StatusCode = http.StatusNotFound
		return nil
	}

	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "DeleteStatus", Value: true},
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	res.Success = true
	res.StatusCode = http.StatusOK
	res.Message = "刪除成功"

	return nil
}

func (e *Activity) UpdateActivity(ctx context.Context, req *pb.UpdateActivityRequest, res *pb.UpdateActivityResponse) error {
	docRef := database.FirestoreClient.Collection("activity").Doc(req.Id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		res.Message = "找不到 id"
		res.Success = false
		res.StatusCode = http.StatusNotFound
		return nil
	}
	version := doc.Data()["Version"].(int64)

	//	Backup update version
	_, err = docRef.Collection("versions").Doc(time.Now().Format(time.RFC3339)).Set(ctx, doc.Data())
	if err != nil {
		return err
	}

	prizeRef := doc.Data()["Prize"].(*firestore.DocumentRef)

	prizeRes, err := prizeService.UpdatePrize(ctx, &prize.UpdatePrizeRequest{
		Name:                        req.Prize.Name,
		ImagePath:                   req.Prize.ImagePath,
		AdditionalPrizeName:         req.Prize.AdditionalPrizeName,
		AdditionalPrizeTicketAmount: req.Prize.AdditionalPrizeTicketAmount,
		AdditionalPrizeImagePath:    req.Prize.AdditionalPrizeImagePath,
		DocumentPath:                prizeRef.ID,
	})
	if err != nil {
		return err
	}
	if prizeRes.StatusCode != http.StatusOK {
		res.Message = prizeRes.Message
		res.StatusCode = prizeRes.StatusCode
		return nil
	}

	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "Title", Value: req.Title},
		{Path: "Detial", Value: req.Detial},
		{Path: "StartTime", Value: time.Unix(req.StartTime, 0)},
		{Path: "EndTime", Value: time.Unix(req.EndTime, 0)},
		{Path: "Version", Value: version + 1},
	})
	if err != nil {
		return err
	}

	res.Message = "更新成功"
	res.StatusCode = http.StatusOK
	res.Success = true

	return nil
}

func (e *Activity) GetActivityUpdateHistoryById(ctx context.Context, req *pb.GetActivityUpdateHistoryByIdRequest, res *pb.GetActivityUpdateHistoryByIdResponse) error {
	docRef := database.FirestoreClient.Collection("activity").Doc(req.Id)
	_, err := docRef.Get(ctx)
	if err != nil {
		res.Message = "找不到 id"
		res.StatusCode = http.StatusNotFound
		res.Data = nil
		return nil
	}

	var result []*structpb.Struct

	iter := docRef.Collection("versions").Documents(ctx)
	for {
		doc, err := iter.Next()
		//	End for condition
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		//	Get this document version
		version := doc.Data()["Version"].(int64)
		data, err := models.FromJsonToPbStruct(ctx, version, doc.Data())
		if err != nil {
			return err
		}
		data.Fields["Id"] = structpb.NewStringValue(doc.Ref.ID)
		result = append(result, data)
	}

	res.Data = result
	res.StatusCode = http.StatusOK

	return nil
}
