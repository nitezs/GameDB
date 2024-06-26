package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameInfo struct {
	ID              primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string               `json:"name,omitempty" bson:"name,omitempty"`
	Description     string               `json:"description,omitempty" bson:"description,omitempty"`
	Aliases         []string             `json:"aliases,omitempty" bson:"aliases,omitempty"`
	Developers      []string             `json:"developers,omitempty" bson:"developers,omitempty"`
	Publishers      []string             `json:"publishers,omitempty" bson:"publishers,omitempty"`
	IGDBID          int                  `json:"-" bson:"igdb_id,omitempty"`
	SteamID         int                  `json:"-" bson:"steam_id,omitempty"`
	GOGID           int                  `json:"-" bson:"gog_id,omitempty"`
	HowLongToBeatID int                  `json:"-" bson:"how_long_to_beat_id,omitempty"`
	Cover           string               `json:"cover,omitempty" bson:"cover,omitempty"`
	Languages       []string             `json:"languages,omitempty" bson:"languages,omitempty"`
	Screenshots     []string             `json:"screenshots,omitempty" bson:"screenshots,omitempty"`
	GameIDs         []primitive.ObjectID `json:"game_ids,omitempty" bson:"games,omitempty"`
	Games           []*GameDownload      `json:"game_downloads,omitempty" bson:"-"`
	CreatedAt       time.Time            `json:"-" bson:"created_at,omitempty"`
	UpdatedAt       time.Time            `json:"-" bson:"updated_at,omitempty"`
}

func (g *GameInfo) MarshalJSON() ([]byte, error) {
	type Alias GameInfo
	aux := &struct {
		ID string `json:"id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if !g.ID.IsZero() {
		aux.ID = g.ID.Hex()
	}

	return json.Marshal(aux)
}

type GameDownload struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"-" bson:"name,omitempty"`
	RawName    string             `json:"raw_name,omitempty" bson:"raw_name,omitempty"`
	Magnet     string             `json:"download_link,omitempty" bson:"magnet,omitempty"`
	Size       string             `json:"size,omitempty" bson:"size,omitempty"`
	Url        string             `json:"url" bson:"url,omitempty"`
	Author     string             `json:"author,omitempty" bson:"author,omitempty"`
	UpdateFlag string             `json:"-" bson:"update_flag,omitempty"`
	CreatedAt  time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"-" bson:"updated_at,omitempty"`
}

type Language struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	LID        int                `bson:"id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	NativeName string             `bson:"native_name,omitempty"`
}
