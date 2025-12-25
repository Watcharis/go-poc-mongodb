package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"watcharis/go-poc-mongodb/repository/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	MONGODB_URI           = "mongodb://root:example@0.0.0.0:27017/admin?retryWrites=true&w=majority&authSource=admin&maxPoolSize=20&connectTimeoutMS=10000"
	MONGODB_DATABASE_NAME = "go-poc-mongodb"
	DEV_MODE              = true
)

func main() {

	ctx := context.Background()

	var monitor *event.CommandMonitor
	if DEV_MODE {
		monitor = &event.CommandMonitor{
			Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
				log.Printf("Command started: %s, %v\n", evt.CommandName, evt.Command)
			},
			Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
				log.Printf("Command succeeded: %s\n", evt.CommandName)
			},
			Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
				log.Printf("Command failed: %s %v\n", evt.CommandName, evt.Failure)
			},
		}
	} else {
		monitor = nil
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(MONGODB_URI).
		SetMaxPoolSize(20).
		SetMaxConnecting(20).
		SetMaxConnIdleTime(10).
		SetMaxConnIdleTime(time.Duration(10000 * time.Millisecond)).
		SetConnectTimeout(time.Duration(10000 * time.Millisecond)).
		SetServerAPIOptions(serverAPI).
		SetMonitor(monitor))
	if err != nil {
		panic(err)
	}
	// disconnect MongoDB when the function returns
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// ping MongoDB
	if err := client.Ping(ctx, &readpref.ReadPref{}); err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	// show databases
	showDBs, err := client.ListDatabaseNames(ctx, bson.D{}, options.ListDatabases().SetNameOnly(true))
	if err != nil {
		panic(err)
	}
	log.Println("MONGODB Databases: ", showDBs)

	userRepository := mongodb.NewUserRepository(client)

	// Find All set skip & limit projection
	resultFindAllUsers, err := userRepository.GetAllUsers(ctx)
	if err != nil {
		panic(err)
	}
	for _, user := range resultFindAllUsers {
		fmt.Println("user:", user)
	}

	// Find One by objectID
	userId := "6945141cb0361299495b7078"
	resultFindOneUser, err := userRepository.GetUserById(ctx, userId)
	if err != nil {
		panic(err)
	}
	log.Printf("Find One User: %+v\n", resultFindOneUser)

	// Update One by objectID
	updateResult, err := userRepository.UpdateUserPhoneOnById(ctx, userId, "0994443331")
	if err != nil {
		panic(err)
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Aggregate user
	resultAggregateUser, err := userRepository.AggregateUsers(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Aggregate Users:", resultAggregateUser)
}
