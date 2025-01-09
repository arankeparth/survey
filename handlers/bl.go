package handlers

import (
	"context"
	"survey-service/db"
)

type BL struct {
	DL *db.DL
}

func NewBL(ctx context.Context) (*BL, error) {
	dl, err := db.NewDL(ctx)
	if err != nil {
		return nil, err
	}
	return &BL{
		DL: dl,
	}, nil
}