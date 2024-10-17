package handler

import "cloud.google.com/go/firestore"

type Data interface {
	~*firestore.DocumentSnapshot | ~string
}

type DataAccess[T Data] interface {
	Connect() error
	GetPrizeDataById(string) (T, error)
	IsIdInDatabase(id string) bool
	CreateData(collectionName string, data any) (*firestore.DocumentRef, error)
	Close() error
}

type CheckIdInterface interface {
	IsIdEmpty(id string) bool
	IsIdExist(id string) bool
}

type NameInterface interface {
	IsPrizeNameEmpty(name string) bool
}

type ImagePathInterface interface {
	IsPathEmpty(path string) bool
}
