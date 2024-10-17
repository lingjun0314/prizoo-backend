package handler

import "cloud.google.com/go/firestore"

type CheckId struct {
	DataAccess[*firestore.DocumentSnapshot]
}

func NewCheckId(data DataAccess[*firestore.DocumentSnapshot]) *CheckId {
	return &CheckId{
		DataAccess: data,
	}
}

func (con *CheckId) IsIdEmpty(id string) bool {
	return id == ""
}

func (con *CheckId) IsIdExist(id string) bool {
	_, err := con.DataAccess.GetPrizeDataById(id)

	return err == nil
}
