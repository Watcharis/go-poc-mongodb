package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id,unique,omitempty"`
	FirstName string             `bson:"first_name,omitempty"`
	LastName  string             `bson:"last_name,omitempty"`
	FullName  string             `bson:"full_name,omitempty"`
	Email     string             `bson:"email,omitempty"`
	PhoneOn   string             `bson:"phone_on,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `bson:"updated_at"`
}
