package crawler

import (
	"GameDB/internal/db"
	"GameDB/internal/model"
	"GameDB/internal/utils"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateGameInfo(idtype string, id int) (*model.GameInfo, error) {
	switch idtype {
	case "steam":
		return GenerateSteamGameInfo(id)
	case "gog":
		return GenerateGOGGameInfo(id)
	case "igdb":
		return GenerateIGDBGameInfo(id)
	default:
		return nil, errors.New("Invalid ID type")
	}
}

func AddGameInfoManually(gameID primitive.ObjectID, idtype string, id int) error {
	info, err := GenerateGameInfo(idtype, id)
	if err != nil {
		return err
	}
	info.GameIDs = append(info.GameIDs, gameID)
	info.GameIDs = utils.Unique(info.GameIDs)
	return db.SaveGameInfo(info)
}
