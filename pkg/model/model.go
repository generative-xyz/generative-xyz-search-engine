package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	ModifiedUserId string             `json:"modifiedUserId,omitempty" bson:"modified_user_id"`
	CreatedUserId  string             `json:"createdUserId,omitempty" bson:"created_user_id"`
	DeletedAt      *time.Time         `json:"deletedAt,omitempty" bson:"deleted_at"`
	CreatedAt      *time.Time         `bson:"created_at,omitempty" json:"createdAt"`
	UpdatedAt      *time.Time         `bson:"updated_at,omitempty" json:"updatedAt"`
}
