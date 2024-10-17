package handler

import (
	"context"
	"net/http"
	"prizeService/proto/prize"
	"time"

	"cloud.google.com/go/firestore"
)

type Prize struct {
	DataAccess[*firestore.DocumentSnapshot]
	CheckIdInterface
	NameInterface
	ImagePathInterface
}

func NewPrize(data DataAccess[*firestore.DocumentSnapshot], CheckId CheckIdInterface, name NameInterface, path ImagePathInterface) *Prize {
	return &Prize{
		DataAccess:         data,
		CheckIdInterface:   CheckId,
		NameInterface:      name,
		ImagePathInterface: path,
	}
}

func (e *Prize) GetPrizeById(ctx context.Context, req *prize.GetPrizeByIdRequest, res *prize.GetPrizeByIdResponse) error {
	id := req.Id

	if !e.DataAccess.IsIdInDatabase(id) {
		res.AdditionalPrizeImagePath = ""
		res.AdditionalPrizeName = ""
		res.AdditionalPrizeTicketAmount = 0
		res.ImagePath = ""
		res.Name = ""
		res.TicketList = nil
		return nil
	}

	doc, err := e.DataAccess.GetPrizeDataById(req.Id)
	if err != nil {
		return err
	}
	res.AdditionalPrizeImagePath = doc.Data()["AdditionalPrizeImagePath"].(string)
	res.AdditionalPrizeName = doc.Data()["AdditionalPrizeName"].(string)
	res.AdditionalPrizeTicketAmount = doc.Data()["additionalPrizeTicketAmount"].(int32)
	res.ImagePath = doc.Data()["ImagePath"].(string)
	res.Name = doc.Data()["Name"].(string)

	return nil
}

func (e *Prize) CreatePrize(ctx context.Context, req *prize.CreatePrizeRequest, res *prize.CreatePrizeResponse) error {
	//	Check name
	if e.NameInterface.IsPrizeNameEmpty(req.Name) {
		res.Message = "prize name is empty"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if e.NameInterface.IsPrizeNameEmpty((req.AdditionalPrizeName)) {
		res.Message = "additional prize name is empty"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	//	Check image
	if e.ImagePathInterface.IsPathEmpty(req.ImagePath) {
		res.Message = "No image access URL"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if e.ImagePathInterface.IsPathEmpty(req.AdditionalPrizeImagePath) {
		res.Message = "No image access URL"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	//	Check ticket amount
	if req.AdditionalPrizeTicketAmount <= 0 {
		res.Message = "additional prize ticket amount cannot smaller than 0"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	docRef, err := e.DataAccess.CreateData("prize", req)
	if err != nil {
		return err
	}
	res.Message = docRef.ID
	res.StatusCode = http.StatusOK

	return nil
}

func (e *Prize) UpdatePrize(ctx context.Context, req *prize.UpdatePrizeRequest, res *prize.UpdatePrizeResponse) error {
	//	Check name
	if e.NameInterface.IsPrizeNameEmpty(req.Name) {
		res.Message = "prize name is empty"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if e.NameInterface.IsPrizeNameEmpty((req.AdditionalPrizeName)) {
		res.Message = "additional prize name is empty"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	//	Check ticket amount
	if req.AdditionalPrizeTicketAmount <= 0 {
		res.Message = "additional prize ticket amount cannot smaller than 0"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Get document from path
	doc, err := e.DataAccess.GetPrizeDataById(req.DocumentPath)
	if err != nil {
		res.Message = "prize id not exists"
		res.StatusCode = http.StatusNotFound
		return nil
	}

	//	Get the current version
	version := doc.Data()["Version"].(int64)
	_, err = doc.Ref.Collection("versions").Doc(time.Now().Format(time.RFC3339)).Set(ctx, doc.Data())
	if err != nil {
		return err
	}

	//	Update data
	updates := []firestore.Update{
		{Path: "Name", Value: req.Name},
		{Path: "AdditionalPrizeName", Value: req.AdditionalPrizeName},
		{Path: "AdditionalPrizeTicketAmount", Value: req.AdditionalPrizeTicketAmount},
		{Path: "Version", Value: version + 1},
	}
	//	If user upload image, update image path
	if !e.ImagePathInterface.IsPathEmpty(req.ImagePath) {
		updates = append(updates, firestore.Update{Path: "imagePath", Value: req.ImagePath})
	}
	if !e.ImagePathInterface.IsPathEmpty(req.AdditionalPrizeImagePath) {
		updates = append(updates, firestore.Update{Path: "AdditionalPrizeImagePath", Value: req.AdditionalPrizeImagePath})
	}

	_, err = doc.Ref.Update(ctx, updates)
	if err != nil {
		return err
	}

	res.StatusCode = http.StatusOK
	return nil
}
