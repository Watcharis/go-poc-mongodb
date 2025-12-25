package mongodb

import (
	"context"
	"time"
	"watcharis/go-poc-mongodb/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return &userRepository{
		client: client,
	}
}

func (r *userRepository) UsersCollection() *mongo.Collection {
	return r.client.Database(MONGODB_DATABASE_NAME).Collection(USERS_COLLECTION)
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]bson.M, error) {
	cursor, err := r.UsersCollection().Find(ctx, bson.D{},
		options.Find().SetProjection(bson.D{{Key: "user_id", Value: 1}, {Key: "email", Value: 1}, {Key: "phone_on", Value: 1}}),
		options.Find().SetSkip(0),
		options.Find().SetLimit(2))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserById ... use struct model
func (r *userRepository) GetUserById(ctx context.Context, id string) (models.Users, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Users{}, err
	}

	var user models.Users
	if err := r.UsersCollection().FindOne(ctx, bson.D{{Key: "_id", Value: objectId}}).Decode(&user); err != nil {
		return models.Users{}, err
	}

	return user, nil
}

// GetUserByIdUseBson ... use bson.M
func (r *userRepository) GetUserByIdUseBson(ctx context.Context, id string) (bson.M, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user bson.M
	if err := r.UsersCollection().FindOne(ctx, bson.D{{Key: "_id", Value: objectId}}).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUserPhoneOnById(ctx context.Context, id string, phoneOn string) (*mongo.UpdateResult, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := r.UsersCollection().UpdateOne(ctx,
		bson.D{
			{
				Key:   "_id",
				Value: objectId,
			},
		},
		bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{
						Key:   "phone_on",
						Value: phoneOn,
					},
				},
			},
		},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
	return result, nil
}

func (r *userRepository) AggregateUsers(ctx context.Context) ([]bson.M, error) {

	pipeline := mongo.Pipeline{
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key:   "user_id",
						Value: "4980f3a6fae54e5aa14780617bb2f045",
					},
				},
			},
		},
		bson.D{
			{
				Key: "$addFields",
				Value: bson.D{
					{
						Key: "first_name_mail",
						Value: bson.D{
							{
								Key:   "$concat",
								Value: bson.A{"$first_name", " ", "$email"},
							},
						},
					},
					{
						Key: "born_date",
						Value: bson.D{
							{
								Key: "$dateFromString",
								Value: bson.D{
									{
										Key:   "dateString",
										Value: time.Now().Format(time.DateTime),
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{
				Key: "$replaceRoot",
				Value: bson.D{
					{
						Key: "newRoot",
						Value: bson.D{
							{
								Key:   "$mergeObjects",
								Value: bson.A{bson.D{{Key: "cid", Value: "1199022901564"}}, "$$ROOT"},
							},
						},
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{
						Key:   "_id",
						Value: 1,
					},
					{
						Key:   "user_id",
						Value: 1,
					},
					{
						Key:   "first_name_mail",
						Value: 1,
					},
					{
						Key:   "cid",
						Value: 1,
					},
					{
						Key:   "born_date",
						Value: 1,
					},
				},
			},
		},
	}

	cursor, err := r.UsersCollection().Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
