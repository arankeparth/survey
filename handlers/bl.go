package handlers

import (
	"context"
	"survey-service/db"

	"go.mongodb.org/mongo-driver/bson"
)

type BL struct {
	DL           *db.DL
	QuestionKeys []interface{}
}

func NewBL(ctx context.Context) (*BL, error) {
	dl, err := db.NewDL(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{}

	questionKeys, err := dl.QuestionCollection.Distinct(ctx, "key", filter)
	if err != nil {
		return nil, err
	}

	return &BL{
		DL:           dl,
		QuestionKeys: questionKeys,
	}, nil
}
