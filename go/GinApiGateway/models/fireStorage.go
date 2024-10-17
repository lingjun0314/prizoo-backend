package models

import (
	"context"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

type FireStorage struct{}

func (*FireStorage) UploadFile(file *multipart.FileHeader) (string, string, error) {
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	objectName := "uploads/" + file.Filename
	wc := StorageBucket.Object(objectName).NewWriter(context.Background())
	if _, err = io.Copy(wc, f); err != nil {
		return "", "", err
	}
	if err = wc.Close(); err != nil {
		return "", "", err
	}

	//	Set this image can be read by all users
	if err := StorageBucket.Object(objectName).ACL().Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		return "", "", err
	}

	//	Get image link
	attrs, err := StorageBucket.Object(objectName).Attrs(context.Background())
	if err != nil {
		return "", "", err
	}
	downloadURL := attrs.MediaLink

	return downloadURL, objectName, nil
}

func (*FireStorage) DeleteFile(filePath string) error {
	//	Create a reference to the file to delete
	obj := StorageBucket.Object(filePath)

	//	Delete the file
	if err := obj.Delete(context.Background()); err != nil {
		return err
	}

	return nil
}
