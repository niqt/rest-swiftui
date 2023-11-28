package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
)

var booksDb *mongo.Database

type key string

var ctx context.Context

const (
	hostKey         = key("hostKey")
	BooksCollection = "books"
	UserCollection  = "users"
	DBName          = "library"
	mongoKey        = "LIBRARY_MONGO"
)

// Initialize the database

func init() {
	mongoServer := os.Getenv(mongoKey)
	if mongoServer == "" {
		mongoServer = "localhost"
	}
	ctx = context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, hostKey, mongoServer)

	var err error
	booksDb, err = configDB(ctx)

	if err != nil {
		fmt.Printf("db configuration %s\n", err)
	}
	fmt.Printf("db connected ")
}

// Configure the database
func configDB(ctx context.Context) (*mongo.Database, error) {
	uri := fmt.Sprintf(`mongodb://%s`,
		ctx.Value(hostKey),
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("todo: couldn't connect to mongo: %v", err)
	}

	db := client.Database(DBName)
	return db, nil
}

// return the reference to the database
func GetDB() *mongo.Database {
	return booksDb
}

// return the context
func GetCtx() context.Context {
	return ctx
}
