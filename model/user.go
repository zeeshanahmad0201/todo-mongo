package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         *string            `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=2,max=100"`
	Email        *string            `json:"email,omitempty" bson:"email,omitempty" validate:"email,required"`
	Password     *string            `validate:"required,min=6"`
	AddedOn      time.Time          `json:"added_on" bson:"added_on"`
	UpdatedOn    time.Time          `json:"updated_on" bson:"updated_on"`
	Token        *string            `json:"token" bson:"token"`
	RefreshToken *string            `json:"refresh_token" bson:"refresh_token"`
	UserID       string             `json:"user_id"`
}
