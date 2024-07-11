package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty"`
	Status    bool               `json:"status,omitempty"`
	AddedOn   time.Time          `json:"addedOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty"`
}
