package models

import (
	"time"

	"cloud.google.com/go/firestore"
	"go-micro.dev/v5/client"
)

var MicroClient client.Client

type Activity struct {
	Title        string                 `json:"title"`
	Detial       string                 `json:"detial"`
	StartTime    time.Time              `json:"startTime"`
	EndTime      time.Time              `json:"endTime"`
	Partner      *firestore.DocumentRef `json:"partner"`
	Prize        *firestore.DocumentRef `json:"prize"`
	DeleteStatus bool                   `json:"deleteStatus"`
	Version      int                    `json:"version"`
}
