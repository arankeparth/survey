package handlers

import (
	"context"
	"survey-service/db"
)

type Handler struct {
	DataLayer *db.DataLayer
}

func NewHandler(ctx context.Context) (*Handler, error) {
	dl, err := db.NewDataLayer(ctx)
	if err != nil {
		return nil, err
	}

	newHandler := &Handler{
		DataLayer: dl,
	}

	return newHandler, nil
}
