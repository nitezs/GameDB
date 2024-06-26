package model

type IGDBGameDetail struct {
	ID                    int     `json:"id,omitempty"`
	ParentGame            int     `json:"parent_game,omitempty"`
	AgeRatings            []int   `json:"age_ratings,omitempty"`
	AlternativeNames      []int   `json:"alternative_names,omitempty"`
	Category              int     `json:"category,omitempty"`
	Cover                 int     `json:"cover,omitempty"`
	CreatedAt             int     `json:"created_at,omitempty"`
	ExternalGames         []int   `json:"external_games,omitempty"`
	FirstReleaseDate      int     `json:"first_release_date,omitempty"`
	Franchises            []int   `json:"franchises,omitempty"`
	GameModes             []int   `json:"game_modes,omitempty"`
	Genres                []int   `json:"genres,omitempty"`
	InvolvedCompanies     []int   `json:"involved_companies,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Platforms             []int   `json:"platforms,omitempty"`
	PlayerPerspectives    []int   `json:"player_perspectives,omitempty"`
	Rating                float64 `json:"rating,omitempty"`
	RatingCount           int     `json:"rating_count,omitempty"`
	ReleaseDates          []int   `json:"release_dates,omitempty"`
	Screenshots           []int   `json:"screenshots,omitempty"`
	SimilarGames          []int   `json:"similar_games,omitempty"`
	Slug                  string  `json:"slug,omitempty"`
	Summary               string  `json:"summary,omitempty"`
	Tags                  []int   `json:"tags,omitempty"`
	Themes                []int   `json:"themes,omitempty"`
	TotalRating           float64 `json:"total_rating,omitempty"`
	TotalRatingCount      int     `json:"total_rating_count,omitempty"`
	UpdatedAt             int     `json:"updated_at,omitempty"`
	URL                   string  `json:"url,omitempty"`
	VersionParent         int     `json:"version_parent,omitempty"`
	VersionTitle          string  `json:"version_title,omitempty"`
	Checksum              string  `json:"checksum,omitempty"`
	Websites              []int   `json:"websites,omitempty"`
	GameLocalizations     []int   `json:"game_localizations,omitempty"`
	AggregatedRating      float64 `json:"aggregated_rating,omitempty"`
	AggregatedRatingCount int     `json:"aggregated_rating_count,omitempty"`
	Artworks              []int   `json:"artworks,omitempty"`
	Bundles               []int   `json:"bundles,omitempty"`
	Collection            int     `json:"collection,omitempty"`
	GameEngines           []int   `json:"game_engines,omitempty"`
	Keywords              []int   `json:"keywords,omitempty"`
	MultiplayerModes      []int   `json:"multiplayer_modes,omitempty"`
	StandaloneExpansions  []int   `json:"standalone_expansions,omitempty"`
	Storyline             string  `json:"storyline,omitempty"`
	Videos                []int   `json:"videos,omitempty"`
	LanguageSupports      []int   `json:"language_supports,omitempty"`
	Collections           []int   `json:"collections,omitempty"`
}

type IGDBGameDetails []*IGDBGameDetail

type IGDBScreenshot struct {
	ID           int    `json:"id,omitempty"`
	AlphaChannel bool   `json:"alpha_channel,omitempty"`
	Animated     bool   `json:"animated,omitempty"`
	Game         int    `json:"game,omitempty"`
	Height       int    `json:"height,omitempty"`
	ImageID      string `json:"image_id,omitempty"`
	URL          string `json:"url,omitempty"`
	Width        int    `json:"width,omitempty"`
	Checksum     string `json:"checksum,omitempty"`
}

type IGDBScreenshots []*IGDBScreenshot

type IGDBCover struct {
	ID           int    `json:"id,omitempty"`
	AlphaChannel bool   `json:"alpha_channel,omitempty"`
	Animated     bool   `json:"animated,omitempty"`
	Game         int    `json:"game,omitempty"`
	Height       int    `json:"height,omitempty"`
	ImageID      string `json:"image_id,omitempty"`
	URL          string `json:"url,omitempty"`
	Width        int    `json:"width,omitempty"`
	Checksum     string `json:"checksum,omitempty"`
}

type IGDBCovers []*IGDBCover

type IGDBAlternativeName struct {
	ID       int    `json:"id,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Game     int    `json:"game,omitempty"`
	Name     string `json:"name,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}

type IGDBAlternativeNames []*IGDBAlternativeName

type IGDBLanguageSupport struct {
	ID                  int    `json:"id,omitempty"`
	Game                int    `json:"game,omitempty"`
	Language            int    `json:"language,omitempty"`
	LanguageSupportType int    `json:"language_support_type,omitempty"`
	CreatedAt           int    `json:"created_at,omitempty"`
	UpdatedAt           int    `json:"updated_at,omitempty"`
	Checksum            string `json:"checksum,omitempty"`
}

type IGDBLanguageSupports []*IGDBLanguageSupport

type IGDBLanguage struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	NativeName string `json:"native_name,omitempty"`
	Locale     string `json:"locale,omitempty"`
	CreatedAt  int    `json:"created_at,omitempty"`
	UpdatedAt  int    `json:"updated_at,omitempty"`
	Checksum   string `json:"checksum,omitempty"`
}

type IGDBLanguages []*IGDBLanguage

type IGDBInvolvedCompany struct {
	ID         int    `json:"id"`
	Company    int    `json:"company"`
	CreatedAt  int    `json:"created_at"`
	Developer  bool   `json:"developer"`
	Game       int    `json:"game"`
	Porting    bool   `json:"porting"`
	Publisher  bool   `json:"publisher"`
	Supporting bool   `json:"supporting"`
	UpdatedAt  int    `json:"updated_at"`
	Checksum   string `json:"checksum"`
}

type IGDBInvolvedCompanies []*IGDBInvolvedCompany

type IGDBCompany struct {
	ID                 int    `json:"id"`
	ChangeDateCategory int    `json:"change_date_category"`
	Country            int    `json:"country"`
	CreatedAt          int    `json:"created_at"`
	Description        string `json:"description"`
	Developed          []int  `json:"developed"`
	Logo               int    `json:"logo"`
	Name               string `json:"name"`
	Parent             int    `json:"parent"`
	Published          []int  `json:"published"`
	Slug               string `json:"slug"`
	StartDate          int    `json:"start_date"`
	StartDateCategory  int    `json:"start_date_category"`
	UpdatedAt          int    `json:"updated_at"`
	URL                string `json:"url"`
	Websites           []int  `json:"websites"`
	Checksum           string `json:"checksum"`
}

type IGDBCompanies []*IGDBCompany
