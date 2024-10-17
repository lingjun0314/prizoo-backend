package dependent

import "cloud.google.com/go/firestore"

type DataAccess interface {
	Connect() error
	GetData(collectionName string) *firestore.DocumentIterator
	GetDataById(string) (*firestore.DocumentSnapshot, error)
	GetDataByName(string) (*firestore.DocumentSnapshot, error)
	IsIdInDatabase(id string) bool
	CreateData(collectionName string, data any) (*firestore.DocumentRef, error)
	UpdateData(collectionName, id string, data any) error
	DeleteData(id string) error
	Close() error
}

type DataCheck interface {
	IsEmpty(string) bool
	IsPhoneValidFormat(string) bool
	IsEmailValidFormat(string) bool
	IsSexValidFormat(string) bool
	IsSocialMediaOptionOne(instagram, facebook, threads, twitter string) bool
}
