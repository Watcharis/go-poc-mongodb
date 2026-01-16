package mongodb

import (
	"context"
	"watcharis/go-poc-mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MONGODB_DATABASE_NAME = "go-poc-mongodb"
	USERS_COLLECTION      = "users"
)

type UserRepository interface {
	GetUserByIdUseBson(ctx context.Context, id string) (bson.M, error) // use bson.M
	GetUserById(ctx context.Context, id string) (models.Users, error)  // use struct model
	GetAllUsers(ctx context.Context) ([]bson.M, error)
	UpdateUserPhoneOnById(ctx context.Context, id string, phoneOn string) (*mongo.UpdateResult, error)
	AggregateUsers(ctx context.Context) ([]bson.M, error)
	InsertUser(ctx context.Context, user models.Users) (*mongo.InsertOneResult, error)
	RemoveUserById(ctx context.Context, id string) (*mongo.DeleteResult, error)
}
