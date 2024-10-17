package models

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go-micro.dev/v5/logger"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// structpb.Struct is a variable structure used to ensure that
// NoSQL data can still be read despite different data structures.
// Version = 0
func FromJsonToPbStruct(ctx context.Context, version int64, doc map[string]interface{}) (*structpb.Struct, error) {
	data := &structpb.Struct{
		Fields: make(map[string]*structpb.Value),
	}

	for k, v := range doc {
		var value *structpb.Value
		switch v := v.(type) {
		case time.Time:
			ts := timestamppb.New(v)
			value = structpb.NewNumberValue(float64(ts.Seconds))
		case *firestore.DocumentRef:
			//	Search reference document
			if version == 0 {
				docSnap, err := v.Get(ctx)
				if err != nil {
					return nil, err
				}
				subData, err := FromJsonToPbStruct(ctx, 0, docSnap.Data())
				if err != nil {
					return nil, err
				}
				value = structpb.NewStructValue(subData)
			} else {
				if v.Parent.ID == "partner" {
					continue
				}
				query := v.Collection("versions").Where("Version", "==", version).Limit(1)
				iter := query.Documents(ctx)
				doc, err := iter.Next()
				if err != nil {
					logger.Error(err)
					return nil, fmt.Errorf("no this version in prize")
				}
				subData, err := FromJsonToPbStruct(ctx, version, doc.Data())
				if err != nil {
					return nil, err
				}
				value = structpb.NewStructValue(subData)
			}

		default:
			var err error
			value, err = structpb.NewValue(v)
			if err != nil {
				return nil, fmt.Errorf("type error: %s", err.Error())
			}
		}
		data.Fields[k] = value
	}

	return data, nil
}
