package ovrstat

// PlayerStats holds all stats on a specified Overwatch player
type PlayerStats struct {
	Icon             string          `json:"icon"`
	Name             string          `json:"name"`
	Level            int             `json:"level"`
	LevelIcon        string          `json:"levelIcon"`
	Prestige         int             `json:"prestige"`
	PrestigeIcon     string          `json:"prestigeIcon"`
	Rating           int             `json:"rating"`
	RatingIcon       string          `json:"ratingIcon"`
	GamesWon         int             `json:"gamesWon"`
	QuickPlayStats   statsCollection `json:"quickPlayStats"`
	CompetitiveStats statsCollection `json:"competitiveStats"`
}

// statsCollection holds a collection of stats for a particular player
type statsCollection struct {
	EliminationsAvg   float64                  `json:"eliminationsAvg"`
	DamageDoneAvg     int64                    `json:"damageDoneAvg"`
	DeathsAvg         float64                  `json:"deathsAvg"`
	FinalBlowsAvg     float64                  `json:"finalBlowsAvg"`
	HealingDoneAvg    int64                    `json:"healingDoneAvg"`
	ObjectiveKillsAvg float64                  `json:"objectiveKillsAvg"`
	ObjectiveTimeAvg  string                   `json:"objectiveTimeAvg"`
	SoloKillsAvg      float64                  `json:"soloKillsAvg"`
	TopHeros          map[string]*topHeroStats `json:"topHeros"`
	CareerStats       map[string]*careerStats  `json:"careerStats"`
}

// topHeroStats holds basic stats for each hero
type topHeroStats struct {
	TimePlayed          string  `json:"timePlayed"`
	GamesWon            int     `json:"gamesWon"`
	WinPercentage       int     `json:"winPercentage"`
	WeaponAccuracy      int     `json:"weaponAccuracy"`
	EliminationsPerLife float64 `json:"eliminationsPerLife"`
	MultiKillBest       int     `json:"multiKillBest"`
	ObjectiveKillsAvg   float64 `json:"objectiveKillsAvg"`
}

// careerStats holds very detailed stats for each hero
type careerStats struct {
	Assists       map[string]string `json:"assists,omitempty"`
	Average       map[string]string `json:"average,omitempty"`
	Best          map[string]string `json:"best,omitempty"`
	Combat        map[string]string `json:"combat,omitempty"`
	Deaths        map[string]string `json:"deaths,omitempty"`
	HeroSpecific  map[string]string `json:"heroSpecific,omitempty"`
	Game          map[string]string `json:"game,omitempty"`
	MatchAwards   map[string]string `json:"matchAwards,omitempty"`
	Miscellaneous map[string]string `json:"miscellaneous,omitempty"`
}
