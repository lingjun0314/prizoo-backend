package dependent

import (
	"context"

	"cloud.google.com/go/firestore"
	"go-micro.dev/v5/logger"
	"google.golang.org/api/option"
)

type FirestoreDataAccess struct {
	client *firestore.Client
}

func (f *FirestoreDataAccess) Connect() error {
	sa := option.WithCredentialsFile("./conf/lucky-draw-1bfec-firebase-adminsdk-o5ax8-0bf161d318.json")

	client, err := firestore.NewClient(context.Background(), "lucky-draw-1bfec", sa)
	if err != nil {
		logger.Errorf("Fail to create database client: %v", err)
		return err
	}
	f.client = client
	return nil
}

func (f *FirestoreDataAccess) GetPrizeDataById(id string) (*firestore.DocumentSnapshot, error) {
	doc, err := f.client.Collection("prize").Doc(id).Get(context.Background())
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (f *FirestoreDataAccess) IsIdInDatabase(id string) bool {
	doc, err := f.client.Collection("prize").Doc(id).Get(context.Background())
	return err == nil && doc.Exists()
}

func (f *FirestoreDataAccess) CreateData(collectionName string, data any) (*firestore.DocumentRef, error) {
	docRef, _, err := f.client.Collection(collectionName).Add(context.Background(), data)
	if err != nil {
		return nil, err
	}

	return docRef, nil
}

func (f *FirestoreDataAccess) Close() error {
	return f.client.Close()
}
