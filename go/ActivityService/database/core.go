package database

import (
	"context"

	"cloud.google.com/go/firestore"
	"go-micro.dev/v5/logger"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func ConnectFirestore() error {
	sa := option.WithCredentialsFile("./conf/lucky-draw-1bfec-firebase-adminsdk-o5ax8-0bf161d318.json")

	client, err := firestore.NewClient(context.Background(), "lucky-draw-1bfec", sa)
	if err != nil {
		logger.Errorf("Fail to create database client: %v", err)
		return err
	}

	FirestoreClient = client

	return nil
}

func CloseFirestore() {
	FirestoreClient.Close()
}