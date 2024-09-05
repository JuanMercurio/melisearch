package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) NewSynonyms(ctx context.Context, searchTerm string, newSynonyms []string) error {

	col := r.db.Collection("synonyms")

	filter := bson.M{"word": searchTerm}
	update := bson.M{
		"$addToSet": bson.M{"synonyms": bson.M{"$each": newSynonyms}},
	}
	options := options.Update().SetUpsert(true)

	_, err := col.UpdateMany(ctx, filter, update, options)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetWordFromSynonym(ctx context.Context, synonym string) (string, error) {
	col := r.db.Collection("synonyms")
	filter := bson.M{"synonyms": synonym}

	var result struct {
		Word string `json:"word"`
	}

	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return "", err
	}

	return result.Word, nil
}

func (r *Repository) List(ctx context.Context, word string) ([]string, error) {
	col := r.db.Collection("synonyms")

	filter := bson.M{"word": word}
	onlySynonyms := bson.D{{Key: "synonyms", Value: 1}}

	res := col.FindOne(ctx, filter, options.FindOne().SetProjection(onlySynonyms))
	if res.Err() != nil {
		return nil, res.Err()
	}

	var result struct {
		Synonyms []string `bson:"synonyms"`
	}
	res.Decode(&result)

	return result.Synonyms, nil
}
