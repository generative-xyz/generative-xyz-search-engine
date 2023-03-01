package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id primitive.ObjectID `json:"_id" bson:"_id" gorm:"-"`
	// Id_ uint               `json:"-" bson:"id" gorm:"column:id;primaryKey"`

	DateCreated    time.Time  `json:"date_created,omitempty" bson:"date_created" gorm:"date_created;autoCreateTime"`
	DateModified   time.Time  `json:"date_modified,omitempty" bson:"date_modified" gorm:"date_modified;autoCreateTime"`
	ModifiedUserId string     `json:"modified_user_id,omitempty" bson:"modified_user_id" gorm:"modified_user_id"`
	CreatedUserId  string     `json:"created_user_id,omitempty" bson:"created_user_id" gorm:"created_user_id"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at" gorm:"deleted_at"`
}
