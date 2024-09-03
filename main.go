package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

	repo := NewRepository(client.Database("testing"))

	router := gin.Default()
	router.GET("/word", func(ctx *gin.Context) {
		synonym := ctx.Query("synonym")
		word, err := repo.GetWordFromSynonym(ctx, synonym)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"word": word})

	})
	router.GET("/synonyms/:word", func(ctx *gin.Context) {
		word := ctx.Param("word")
		synonyms, err := repo.List(ctx, word)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"synonyms": synonyms})
	})

	router.POST("/synonyms/:word", func(ctx *gin.Context) {
		word := ctx.Param("word")
		synonyms := strings.Split(ctx.Query("synonyms"), ",")

		for i, s := range synonyms {
			synonyms[i] = strings.TrimSpace(s)
		}
		for i := len(synonyms) - 1; i >= 0; i-- {
			if synonyms[i] == "" {
				synonyms = append(synonyms[:i], synonyms[i+1:]...)
			}
		}

		if err := repo.NewSynonyms(ctx, word, synonyms); err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"status": "ok"})

	})

	router.Run(":8080")
	return
}
