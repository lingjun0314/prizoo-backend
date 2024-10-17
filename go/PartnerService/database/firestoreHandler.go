package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"go-micro.dev/v5/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Firestore struct {
	client *firestore.Client
}

func (e *Firestore) Connect() error {
	sa := option.WithCredentialsFile("./conf/lucky-draw-1bfec-firebase-adminsdk-o5ax8-0bf161d318.json")

	client, err := firestore.NewClient(context.Background(), "lucky-draw-1bfec", sa)
	if err != nil {
		logger.Errorf("Fail to create database client: %v", err)
		return err
	}
	e.client = client
	return nil
}

func (e *Firestore) GetData(collectionName string) *firestore.DocumentIterator {
	return e.client.Collection(collectionName).Where("DeleteStatus", "==", false).Documents(context.Background())
}

func (e *Firestore) GetDataById(id string) (*firestore.DocumentSnapshot, error) {
	doc, err := e.client.Collection("partner").Doc(id).Get(context.Background())
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *Firestore) GetDataByName(name string) (*firestore.DocumentSnapshot, error) {
	doc, err := e.client.Collection("partner").Where("BrandName", "==", name).Limit(1).Documents(context.Background()).Next()
	if err != nil {
		if err == iterator.Done {
			return nil, fmt.Errorf("404")
		}
		return nil, err
	}
	return doc, nil
}

func (e *Firestore) IsIdInDatabase(id string) bool {
	doc, err := e.client.Collection("partner").Doc(id).Get(context.Background())
	return err == nil && doc.Exists()
}

func (e *Firestore) CreateData(collectionName string, data any) (*firestore.DocumentRef, error) {
	docRef, _, err := e.client.Collection(collectionName).Add(context.Background(), data)
	if err != nil {
		return nil, err
	}

	return docRef, nil
}

// @dev: updates 的資料在調用時就要定好 []firestore.Update 格式再傳入
func (e *Firestore) UpdateData(collectionName, id string, updates any) error {
	//	Get documnet reference
	docRef := e.client.Collection(collectionName).Doc(id)
	if docRef == nil {
		return fmt.Errorf("error collection name or id")
	}

	//	Change type to []firestore.Update
	updatesData, ok := updates.([]firestore.Update)
	if !ok {
		return fmt.Errorf("error updates format")
	}

	//	Update data
	_, err := docRef.Update(context.Background(), updatesData)
	if err != nil {
		return err
	}

	return nil
}

func (e *Firestore) DeleteData(id string) error {
	doc, err := e.GetDataById(id)
	if err != nil {
		return err
	}

	_, err = doc.Ref.Update(context.Background(), []firestore.Update{
		{Path: "DeleteStatus", Value: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Firestore) Close() error {
	return e.client.Close()
}
