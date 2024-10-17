package models

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var AuthenticationClient *auth.Client
var StorageBucket *storage.BucketHandle

const bucketName string = "lucky-draw-1bfec.appspot.com"

func InitFirebaseAuth() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./conf/lucky-draw-1bfec-firebase-adminsdk-o5ax8-0bf161d318.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatal("Error initializing firebase: ", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatal("Error initializing authentication client: ", err)
	}

	AuthenticationClient = client
}

func InitFirebaseStorage() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./conf/lucky-draw-1bfec-firebase-adminsdk-o5ax8-0bf161d318.json")
	config := &firebase.Config{
		StorageBucket: bucketName,
	}
	app, err := firebase.NewApp(ctx, config, sa)
	if err != nil {
		log.Fatal("Error initializing firebase: ", err)
	}
	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatal("Error initializing storage client: ", err)
	}
	StorageBucket, err = client.DefaultBucket()
	if err != nil {
		log.Fatal("Error getting default bucket: ", err)
	}
}
