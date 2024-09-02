package main

import (
	"context"
	"fmt"
	"log"

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

func (r *Repository) NewSynonyms(ctx context.Context, word string, newSynonyms []string) error {

	col := r.db.Collection("synonyms")

	filter := bson.M{"word": word}
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

	cur, err := col.Find(ctx, filter)
	if err != nil {
		return "", err
	}

	var word struct {
		Word string `json:"word"`
	}

	for cur.Next(ctx) {
		err := cur.Decode(&word)
		if err != nil {
			return "", err
		}
	}

	return word.Word, nil
}

func (r *Repository) List(ctx context.Context, word string) {
	col := r.db.Collection("synonyms")

	filter := bson.M{"word": word}

	cur, err := col.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var synonyms struct {
		Synonyms []string `json:"synonyms"`
	}

	// esto solofunciona si tiene 1 doc de esa palabra
	for cur.Next(ctx) {
		err := cur.Decode(&synonyms)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(synonyms.Synonyms)

}
