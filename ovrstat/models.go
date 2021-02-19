package ovrstat /* import "s32x.com/ovrstat/ovrstat" */

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
	Ratings          []Rating        `json:"ratings"`
	GamesWon         int             `json:"gamesWon"`
	QuickPlayStats   statsCollection `json:"quickPlayStats"`
	CompetitiveStats statsCollection `json:"competitiveStats"`
	Private          bool            `json:"private"`
}

type Rating struct {
	Level    int    `json:"level"`
	Role     string `json:"role"`
	RoleIcon string `json:"roleIcon"`
	RankIcon string `json:"rankIcon"`
}

// statsCollection holds a collection of stats for a particular player
type statsCollection struct {
	TopHeroes   map[string]*topHeroStats `json:"topHeroes"`
	CareerStats map[string]*careerStats  `json:"careerStats"`
}

// topHeroStats holds basic stats for each hero
type topHeroStats struct {
	TimePlayed          string  `json:"timePlayed"`
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

// Platform represents a response from the search-by-name api request
type Platform struct {
	Platform    string `json:"platform"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URLName     string `json:"urlName"`
	PlayerLevel int    `json:"playerLevel"`
	Portrait    string `json:"portrait"`
	IsPublic    bool   `json:"isPublic"`
}
