package models

import "mime/multipart"

type FileAccess interface {
	UploadFile(*multipart.FileHeader) (string, string, error)
	DeleteFile(string) error
}
