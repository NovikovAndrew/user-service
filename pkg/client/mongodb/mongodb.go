package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var optinDB options.Credential
	var mongoURI string
	var isAuth bool

	if authDB == "" {
		authDB = database
	}

	if username == "" && password == "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
		optinDB = options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		}
		isAuth = true
	}

	clientOptins := options.Client().ApplyURI(mongoURI)

	if isAuth {
		clientOptins.SetAuth(optinDB)
	}

	client, err := mongo.Connect(ctx)

	if err != nil {
		return nil, fmt.Errorf("Can not to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Can not connect to ping MongoDB: %v", err)
	}

	return client.Database(database), nil
}
