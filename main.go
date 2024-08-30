package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	//creamos un context con 10 segundos de timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// conectamos a la base de mongo
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// ping para ver que ande todo
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	type ejemplo struct {
		Saturno string `bson:"saturno"`
		Jupiter string `bson:"jupiterrrrrr"`
	}

	e := ejemplo{
		"planeta con anillos",
		"planeta grade",
	}

	// nos conectamos a una base y una collection
	collection := client.Database("testing").Collection("colltest")
	result, err := collection.InsertOne(ctx, e)
	if err != nil {
		log.Fatal(err)
	}

	// nos busca la "tabla" entera
	cur, err := collection.Find(ctx, e)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	// iteramos sobre las "filas" / documento
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("id:")
	fmt.Println(result.InsertedID)
}
