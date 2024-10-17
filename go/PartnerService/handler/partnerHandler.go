package handler

import (
	"context"
	"net/http"
	"partnerService/dependent"
	"partnerService/models"
	"partnerService/proto/partner"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/structpb"
)

type PartnerHandler struct {
	Database  dependent.DataAccess
	DataCheck dependent.DataCheck
}

func NewPartnerHandler(database dependent.DataAccess, dataCheck dependent.DataCheck) *PartnerHandler {
	return &PartnerHandler{
		Database:  database,
		DataCheck: dataCheck,
	}
}

func (con *PartnerHandler) GetPartnerByName(ctx context.Context, req *partner.GetPartnerByNameRequest, res *partner.GetPartnerByNameResponse) error {
	//	Check data valid
	if con.DataCheck.IsEmpty(req.BrandName) {
		res.Message = "請提供廠商名稱"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Get data
	doc, err := con.Database.GetDataByName(req.BrandName)
	if err != nil {
		if err.Error() == "404" {
			res.Message = "沒有此廠商，請先建立廠商"
			res.StatusCode = http.StatusNotFound
			return nil
		}
		return err
	}

	//	Return data
	res.Message = doc.Ref.ID
	res.StatusCode = http.StatusOK

	return nil
}

func (con *PartnerHandler) GetPartnerById(ctx context.Context, req *partner.GetPartnerByIdRequest, res *partner.GetPartnerByIdResponse) error {
	//	Check Id
	if !con.Database.IsIdInDatabase(req.Id) {
		res.Message = "請提供有效的 id"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	doc, err := con.Database.GetDataById(req.Id)
	if err != nil {
		return err
	}
	res.Partner, err = models.FromJsonToPbStruct(ctx, 0, doc.Data())
	if err != nil {
		return err
	}
	res.Partner.Fields["Id"] = structpb.NewStringValue(doc.Ref.ID)
	res.StatusCode = http.StatusOK
	return nil
}

func (con *PartnerHandler) CreatePartner(ctx context.Context, req *partner.CreatePartnerRequest, res *partner.CreatePartnerResponse) error {
	//	Check name
	if con.DataCheck.IsEmpty(req.Partner.BrandName) {
		res.Message = "請提供品牌名稱"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if con.DataCheck.IsEmpty(req.Partner.PersonInCharge) {
		res.Message = "請提供負責人名稱"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if con.DataCheck.IsEmpty(req.Partner.ContactPerson) {
		res.Message = "請提供窗口人姓名"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check phone
	if con.DataCheck.IsEmpty(req.Partner.Contact.Phone) && con.DataCheck.IsEmpty(req.Partner.Contact.Tel) {
		res.Message = "請提供手機號碼或電話"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if !con.DataCheck.IsPhoneValidFormat(req.Partner.Contact.Phone) {
		res.Message = "請輸入正確的手機格式"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check email
	if con.DataCheck.IsEmpty(req.Partner.Contact.Email) {
		res.Message = "請提供 email"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if !con.DataCheck.IsEmailValidFormat(req.Partner.Contact.Email) {
		res.Message = "請輸入正確的 email 格式"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check age
	if con.DataCheck.IsEmpty(req.Partner.Customer.AgeRange) {
		res.Message = "請提供客戶年齡范圍"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check sex
	if !con.DataCheck.IsSexValidFormat(req.Partner.Customer.Sex) {
		res.Message = "錯誤的性別"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check social media
	if !con.DataCheck.IsSocialMediaOptionOne(
		req.Partner.SocialMedia.Instagram,
		req.Partner.SocialMedia.Facebook,
		req.Partner.SocialMedia.Threads,
		req.Partner.SocialMedia.Twitter,
	) {
		res.Message = "至少輸入一個社群帳號連結"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check name exist
	if doc, _ := con.Database.GetDataByName(req.Partner.BrandName); doc != nil {
		res.Message = "此廠商已經存在"
		res.StatusCode = http.StatusConflict
		return nil
	}

	//	Create data
	docRef, err := con.Database.CreateData("partner", req.Partner)
	if err != nil {
		return err
	}
	res.Message = docRef.ID
	res.StatusCode = http.StatusOK

	return nil
}

func (con *PartnerHandler) UpdatePartner(ctx context.Context, req *partner.UpdatePartnerRequest, res *partner.UpdatePartnerResponse) error {
	//	Check name
	if con.DataCheck.IsEmpty(req.Partner.BrandName) {
		res.Message = "請提供品牌名稱"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if con.DataCheck.IsEmpty(req.Partner.PersonInCharge) {
		res.Message = "請提供負責人名稱"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if con.DataCheck.IsEmpty(req.Partner.ContactPerson) {
		res.Message = "請提供窗口人姓名"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check phone
	if con.DataCheck.IsEmpty(req.Partner.Contact.Phone) && con.DataCheck.IsEmpty(req.Partner.Contact.Tel) {
		res.Message = "請提供手機號碼或電話"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if !con.DataCheck.IsPhoneValidFormat(req.Partner.Contact.Phone) {
		res.Message = "請輸入正確的手機格式"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check email
	if con.DataCheck.IsEmpty(req.Partner.Contact.Email) {
		res.Message = "請提供 email"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	if !con.DataCheck.IsEmailValidFormat(req.Partner.Contact.Email) {
		res.Message = "請輸入正確的 email 格式"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check age
	if con.DataCheck.IsEmpty(req.Partner.Customer.AgeRange) {
		res.Message = "請提供客戶年齡范圍"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check sex
	if !con.DataCheck.IsSexValidFormat(req.Partner.Customer.Sex) {
		res.Message = "錯誤的性別"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check social media
	if !con.DataCheck.IsSocialMediaOptionOne(
		req.Partner.SocialMedia.Instagram,
		req.Partner.SocialMedia.Facebook,
		req.Partner.SocialMedia.Threads,
		req.Partner.SocialMedia.Twitter,
	) {
		res.Message = "至少輸入一個社群帳號連結"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	//	Check Id
	if !con.Database.IsIdInDatabase(req.Id) {
		res.Message = "請提供有效的 id"
		res.StatusCode = http.StatusBadRequest
		return nil
	}
	doc, _ := con.Database.GetDataByName(req.Partner.BrandName)
	docNow, _ := con.Database.GetDataById(req.Id)
	if doc != nil {
		//	若有查到 partner，且此 partner 並不是當前 partner
		//	此時代表有其他 partner 的名字是現在要更改的名字，又可以保證自己名字不改的情況下邏輯可以過
		if doc.Ref.ID != docNow.Ref.ID {
			res.Message = "此廠商已經存在"
			res.StatusCode = http.StatusConflict
			return nil
		}
	}

	updates := []firestore.Update{
		{Path: "BrandName", Value: req.Partner.BrandName},
		{Path: "PersonInCharge", Value: req.Partner.PersonInCharge},
		{Path: "CompanyName", Value: req.Partner.CompanyName},
		{Path: "ContactPerson", Value: req.Partner.ContactPerson},
		{Path: "Contact", Value: req.Partner.Contact},
		{Path: "Customer", Value: req.Partner.Customer},
		{Path: "SocialMedia", Value: req.Partner.SocialMedia},
	}

	err := con.Database.UpdateData("partner", req.Id, updates)
	if err != nil {
		return err
	}

	res.Message = "更新成功"
	res.StatusCode = http.StatusOK

	return nil
}

func (con *PartnerHandler) DeletePartner(ctx context.Context, req *partner.DeletePartnerRequest, res *partner.DeletePartnerResponse) error {
	if !con.Database.IsIdInDatabase(req.Id) {
		res.Message = "請提供有效的 id"
		res.StatusCode = http.StatusBadRequest
		return nil
	}

	err := con.Database.DeleteData(req.Id)
	if err != nil {
		return err
	}
	res.StatusCode = http.StatusOK
	res.Message = "刪除成功"

	return nil
}

func (con *PartnerHandler) GetPartners(ctx context.Context, req *partner.GetPartnersRequest, res *partner.GetPartnersResponse) error {
	iter := con.Database.GetData("partner")
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
		res.Partner = append(res.Partner, data)
	}
	return nil
}
