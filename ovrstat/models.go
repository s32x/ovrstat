package ovrstat

// PlayerStats holds all stats on a specified Overwatch player
type PlayerStats struct {
	Icon             string          `json:"icon"`
	Name             string          `json:"name"`
	Level            int             `json:"level"`
	LevelIcon        string          `json:"levelIcon"`
	Endorsement      int             `json:"endorsement"`
	EndorsementIcon  string          `json:"endorsementIcon"`
	Prestige         int             `json:"prestige"`
	PrestigeIcon     string          `json:"prestigeIcon"`
	Rating           int             `json:"rating"`
	RatingIcon       string          `json:"ratingIcon"`
	GamesWon         int             `json:"gamesWon"`
	QuickPlayStats   statsCollection `json:"quickPlayStats"`
	CompetitiveStats statsCollection `json:"competitiveStats"`
	Private          bool            `json:"private"`
}

// statsCollection holds a collection of stats for a particular player
type statsCollection struct {
	TopHeroes   map[string]*topHeroStats `json:"topHeroes"`
	CareerStats map[string]*careerStats  `json:"careerStats"`
}

// topHeroStats holds basic stats for each hero
type topHeroStats struct {
	TimePlayed          string  `json:"timePlayed"`
	TimePlayedInSeconds int     `json:"timePlayedInSeconds"`
	GamesWon            int     `json:"gamesWon"`
	WinPercentage       int     `json:"winPercentage"`
	WeaponAccuracy      int     `json:"weaponAccuracy"`
	EliminationsPerLife float64 `json:"eliminationsPerLife"`
	MultiKillBest       int     `json:"multiKillBest"`
	ObjectiveKills      float64 `json:"objectiveKills"`
}

// careerStats holds very detailed stats for each hero
type careerStats struct {
	Assists       map[string]interface{} `json:"assists"`
	Average       map[string]interface{} `json:"average"`
	Best          map[string]interface{} `json:"best"`
	Combat        map[string]interface{} `json:"combat"`
	Deaths        map[string]interface{} `json:"deaths"`
	HeroSpecific  map[string]interface{} `json:"heroSpecific"`
	Game          map[string]interface{} `json:"game"`
	MatchAwards   map[string]interface{} `json:"matchAwards"`
	Miscellaneous map[string]interface{} `json:"miscellaneous"`
}
