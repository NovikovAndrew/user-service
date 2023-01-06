package db

import (
	"context"
	"errors"
	"fmt"
	"rest-api/cmd/internal/user"
	"rest-api/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collectionName string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collectionName),
		logger:     logger,
	}
}

func (db *db) FindAll(ctx context.Context) ([]user.User, error) {
	var users []user.User
	result, err := db.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, fmt.Errorf("Fauler to find users: %v", err)
	}

	if err := result.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("Failed to deocde users, error: %v", err)
	}

	return users, nil
}

func (db *db) Create(ctx context.Context, user user.User) (string, error) {
	db.logger.Debugf("Create user: %v", user)

	result, err := db.collection.InsertOne(ctx, user)

	if err != nil {
		return "", fmt.Errorf("Can not create user %v: %v", user, err)
	}

	db.logger.Debug("Convert InsertedID to ObjectID")
	if res, ok := result.InsertedID.(primitive.ObjectID); ok {
		return res.Hex(), nil
	}

	db.logger.Trace(user)
	return "", fmt.Errorf("Failed to convert from object ID to hex")
}

func (db *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	objctID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return u, fmt.Errorf("Can not convert hex to object ID %v: %v", id, err)
	}

	filter := bson.M{"_id": objctID}
	result := db.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// TODO: error notify not found
			return u, err
		}

		return u, fmt.Errorf("Failed to find user by id %s, error %v", id, result.Err())
	}

	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("Failed to decode user by id %s, error %v", id, err)
	}

	return u, nil
}

func (db *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)

	if err != nil {
		return fmt.Errorf("Failed to get object ID from user %s, error %v", user, err)
	}

	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)

	if err != nil {
		return fmt.Errorf("Failed to marshal user %s, error %v", user, err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)

	if err != nil {
		return fmt.Errorf("Failed to unmarshal userbytes %s, error %v", user)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := db.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("Failed to update user query %s, error %v", user, err)
	}

	if result.MatchedCount == 0 {
		// TODO: error entity not found
		return fmt.Errorf("NOT FOUND")
	}

	db.logger.Trace("Macted %d documents, modified %d", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (db *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("Failed to get object ID %s, error %v", id, err)
	}

	filter := bson.M{"_id": objectID}

	result, err := db.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("Failed to delete user by id %s, error %v", objectID, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("NOT FOUND")
	}

	db.logger.Trace("Deleted %d documents", result.DeletedCount)

	return nil
}
