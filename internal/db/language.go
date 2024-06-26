package db

import (
	"GameDB/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveLanguage(language *model.Language) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"id": language.ID}
	update := bson.M{"$set": language}
	opts := options.Update().SetUpsert(true)
	_, err := LanguageCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetLanguages(ids []int) ([]*model.Language, error) {
	var languages []*model.Language
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var filter interface{}
	if len(ids) == 1 {
		filter = bson.M{"id": ids[0]}
	} else {
		filter = bson.M{"id": bson.M{"$in": ids}}
	}
	cursor, err := LanguageCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var language model.Language
		if err := cursor.Decode(&language); err != nil {
			return nil, err
		}
		languages = append(languages, &language)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return languages, nil
}
