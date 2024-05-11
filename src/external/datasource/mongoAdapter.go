package mongodb

import (
	"context"
	"errors"
	"strings"
	"time"

	coreErrors "github.com/hcsouza/fiap-tech-fast-food/src/common/errors"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const MongoDuplicateKeyErrorCode = 11000

type mongoAdapter[T any] struct {
	client         mongo.Client
	collectionName string
	collection     mongo.Collection
	domain         T
}

func NewMongoAdapter[T any](client mongo.Client, databaseName, collectionName string) interfaces.DatabaseSource {
	collection := client.Database(databaseName).Collection(collectionName)
	return &mongoAdapter[T]{
		collectionName: collectionName,
		client:         client,
		collection:     *collection,
	}
}

func (ad *mongoAdapter[T]) FindAll(fieldName, fieldValue string) ([]interface{}, error) {
	ctx := context.TODO()
	param := bson.D{}
	var results []T

	if fieldName != "" && fieldValue != "" {
		param = bson.D{{Key: fieldName, Value: fieldValue}}
	}

	cursor, err := ad.collection.Find(ctx, param)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var interfaceResults []interface{}
	for _, result := range results {
		interfaceResults = append(interfaceResults, result)
	}

	return interfaceResults, nil
}

func (ad *mongoAdapter[T]) FindOne(key, value string) (interface{}, error) {
	ctx := context.TODO()
	var result T

	err := ad.collection.FindOne(
		ctx,
		bson.D{{Key: key, Value: value}},
	).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("record not found")
		} else {
			return nil, err
		}
	}
	return &result, err
}

func (ad *mongoAdapter[T]) Save(data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ad.collection.InsertOne(ctx, data)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			return nil, coreErrors.ErrDuplicatedKey
		}
		return nil, err
	}

	return res.InsertedID, err
}

func (ad *mongoAdapter[T]) Update(identifier string, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ad.collection.UpdateOne(ctx, bson.M{"_id": identifier}, bson.D{{"$set", data}})

	if err != nil {
		return nil, err
	}

	return res.UpsertedID, err
}

func (ad *mongoAdapter[T]) Delete(identifier string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := ad.collection.DeleteOne(ctx, bson.M{"_id": identifier})
	return res.DeletedCount, err
}
