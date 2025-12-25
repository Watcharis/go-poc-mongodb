package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Users struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    string        `bson:"user_id,unique,omitempty"`
	LastName  string        `bson:"last_name,omitempty"`
	Email     string        `bson:"email,omitempty"`
	PhoneOn   string        `bson:"phone_on,omitempty"`
	FullName  string        `bson:"full_name,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}
