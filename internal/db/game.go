package db

import (
	"GameDB/internal/cache"
	"GameDB/internal/config"
	"GameDB/internal/log"
	"GameDB/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	removeDelimiter            = regexp.MustCompile(`[:\-\+]`)
	removeRepeatingSpacesRegex = regexp.MustCompile(`\s+`)
)

func GetAllGameDownloadsWithAuthor(regex string) ([]*model.GameDownload, error) {
	var items []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "author", Value: primitive.Regex{Pattern: regex, Options: "i"}}}
	cursor, err := GameDownloadCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameDownload
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		items = append(items, &game)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	return items, err
}

func IsGameCrawled(flag string, author string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "author", Value: primitive.Regex{Pattern: author, Options: "i"}},
		{Key: "update_flag", Value: flag},
	}
	var game model.GameDownload
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return false
		}
		log.Logger.Error("Failed to find game", zap.Error(err))
		return false
	}
	return true
}

func IsGameCrawledByURL(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "url", Value: url},
	}
	var game model.GameDownload
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return false
		}
		log.Logger.Error("Failed to find game", zap.Error(err))
		return false
	}
	return true
}

func SaveGameDownload(item *model.GameDownload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()
	item.Size = strings.Replace(item.Size, "gb", "GB", -1)
	item.Size = strings.Replace(item.Size, "mb", "MB", -1)
	filter := bson.M{"_id": item.ID}
	update := bson.M{"$set": item}
	opts := options.Update().SetUpsert(true)
	_, err := GameDownloadCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func SaveGameInfo(item *model.GameInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()
	filter := bson.M{"_id": item.ID}
	update := bson.M{"$set": item}
	opts := options.Update().SetUpsert(true)
	_, err := GameInfoCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func SaveGameDownloads(items []*model.GameDownload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	operations := []mongo.WriteModel{}
	for _, item := range items {
		if item.ID.IsZero() {
			item.ID = primitive.NewObjectID()
		}
		if item.CreatedAt.IsZero() {
			item.CreatedAt = time.Now()
		}
		item.UpdatedAt = time.Now()
		filter := bson.M{"_id": item.ID}
		update := bson.M{"$set": item}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, model)
	}
	_, err := GameDownloadCollection.BulkWrite(ctx, operations)
	if err != nil {
		return err
	}
	return nil
}

func GetAllGameDownloads() ([]*model.GameDownload, error) {
	var items []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := GameDownloadCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameDownload
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		items = append(items, &game)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	return items, err
}

func GetGameDownloadByUrl(url string) (*model.GameDownload, error) {
	var item model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"url": url}
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return &model.GameDownload{}, nil
		}
		return nil, err
	}
	return &item, nil
}

func GetGameDownloadByID(id primitive.ObjectID) (*model.GameDownload, error) {
	var item model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	err := GameDownloadCollection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return &model.GameDownload{}, nil
		}
		return nil, err
	}
	return &item, nil
}

func GetGameDownloadsByIDs(ids []primitive.ObjectID) ([]*model.GameDownload, error) {
	var items []*model.GameDownload
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := GameDownloadCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameDownload
		if err = cursor.Decode(&game); err != nil {
			return nil, err
		}
		items = append(items, &game)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	return items, err
}

func SearchGameInfos(name string, page int, pageSize int) ([]*model.GameInfo, int, error) {
	var items []*model.GameInfo
	name = removeDelimiter.ReplaceAllString(name, " ")
	name = removeRepeatingSpacesRegex.ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)
	name = strings.Replace(name, " ", ".*", -1)
	name = fmt.Sprintf("%s.*", name)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"$or": []interface{}{
		bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}},
		bson.M{"aliases": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}},
	}}
	totalCount, err := GameInfoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)
	findOpts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := GameInfoCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var game model.GameInfo
		if err = cursor.Decode(&game); err != nil {
			return nil, 0, err
		}
		game.Games, err = GetGameDownloadsByIDs(game.GameIDs)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, &game)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	return items, int(totalPages), nil
}

func SearchGameInfosCache(name string, page int, pageSize int) ([]*model.GameInfo, int, error) {
	type res struct {
		Items     []*model.GameInfo
		TotalPage int
	}
	if config.Config.RedisAvaliable {
		key := fmt.Sprintf("searchGameDetails:%s:%d:%d", name, page, pageSize)
		val, exist := cache.Redis.Get(key)
		if exist {
			var data res
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				return nil, 0, err
			}
			return data.Items, data.TotalPage, nil
		} else {
			data, totalPage, err := SearchGameInfos(name, page, pageSize)
			if err != nil {
				return nil, 0, err
			}
			dataBytes, err := json.Marshal(res{Items: data, TotalPage: totalPage})
			if err != nil {
				return nil, 0, err
			}
			err = cache.Redis.AddWithExpire(key, string(dataBytes), 10*time.Minute)
			if err != nil {
				log.Logger.Warn("Failed to add cache", zap.Error(err))
			}
			return data, totalPage, nil
		}
	} else {
		return SearchGameInfos(name, page, pageSize)
	}
}

func GetGameInfoByPlatformID(idtype string, id int) (*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var filter interface{}
	switch idtype {
	case "steam":
		filter = bson.M{"steam_id": id}
	case "gog":
		filter = bson.M{"gog_id": id}
	case "igdb":
		filter = bson.M{"igdb_id": id}
	}
	var game model.GameInfo
	err := GameInfoCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func IsGameInfoExist(idtype string, id int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var filter interface{}
	switch idtype {
	case "steam":
		filter = bson.M{"steam_id": id}
	case "gog":
		filter = bson.M{"gog_id": id}
	case "igdb":
		filter = bson.M{"igdb_id": id}
	}
	var game model.GameInfo
	err := GameInfoCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return false
		}
		log.Logger.Error("Failed to find game", zap.Error(err))
		return false
	}
	return true
}

func GetGameDownloadsNotInGameInfos(num int) ([]*model.GameDownload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var gamesNotInDetails []*model.GameDownload
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "game_infos"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "games"},
			{Key: "as", Value: "gameDetail"},
		}}},
	}
	if num != -1 && num > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: num}})
	}
	pipeline = append(pipeline,
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "gameDetail", Value: bson.D{{Key: "$size", Value: 0}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "name", Value: 1}}}},
	)

	cursor, err := GameDownloadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var game model.GameDownload
		if err := cursor.Decode(&game); err != nil {
			return nil, err
		}
		gamesNotInDetails = append(gamesNotInDetails, &game)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gamesNotInDetails, nil
}

func GetGameInfoByID(id primitive.ObjectID) (*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var game model.GameInfo
	err := GameInfoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func DeduplicateGames() error {
	type queryRes struct {
		ID    string               `bson:"_id"`
		Total int                  `bson:"total"`
		IDs   []primitive.ObjectID `bson:"ids"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res []queryRes
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$magnet"},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ids", Value: bson.D{{Key: "$push", Value: "$_id"}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "total", Value: bson.D{{Key: "$gt", Value: 1}}},
		}}},
	}
	cursor, err := GameDownloadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return err
	}
	for _, item := range res {
		idsToDelete := item.IDs[1:]
		_, err = GameDownloadCollection.DeleteMany(ctx, bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: idsToDelete}}}})
		if err != nil {
			return err
		}
		log.Logger.Info("Removed duplicates", zap.Any("ids", idsToDelete))
		cursor, err := GameInfoCollection.Find(ctx, bson.M{"games": bson.M{"$in": idsToDelete}})
		if err != nil {
			return err
		}
		var infos []*model.GameInfo
		if err := cursor.All(ctx, &infos); err != nil {
			return err
		}
		for _, info := range infos {
			newGames := make([]primitive.ObjectID, 0, len(info.GameIDs))
			for _, id := range info.GameIDs {
				if !slices.Contains(idsToDelete, id) {
					newGames = append(newGames, id)
				}
			}
			info.GameIDs = newGames
			if err := SaveGameInfo(info); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetGameInfosByName(name string) ([]*model.GameInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	name = strings.TrimSpace(name)
	name = fmt.Sprintf("^%s$", name)
	filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}}
	cursor, err := GameInfoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var games []*model.GameInfo
	if err = cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}
