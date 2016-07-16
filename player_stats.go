package goow

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// PlayerStats : Holds all stats on a specified Overwatch player
type PlayerStats struct {
	Icon             string
	Name             string
	Level            int
	LevelIcon        string
	Rating           string
	RatingIcon       string
	GamesPlayed      int
	QuickPlayStats   StatsCollection
	CompetitiveStats StatsCollection
}

type StatsCollection struct {
	EliminationsAvg   float64
	DamageDoneAvg     int64
	DeathsAvg         float64
	FinalBlowsAvg     float64
	HealingDoneAvg    int64
	ObjectiveKillsAvg float64
	ObjectiveTimeAvg  string
	SoloKillsAvg      float64
	//TopHeros          []HeroStats
	//CareerStats     CareerStats
}

// type HeroStats struct {
// 	Name                string
// 	TimePlayed          string
// 	GamesWon            int
// 	WinPercentage       int
// 	WeaponAccuracy      int
// 	EliminationsPerLife float32
// 	MultiKillBest       int
// 	ObjectiveKillsAvg   float32
// }
//
// type CareerStats struct {
// 	Combat        CombatStats
// 	Death         DeathStats
// 	Game          GameStats
// 	Assists       AssistsStats
// 	Average       AverageStats
// 	Miscellaneous MiscellaneousStats
// 	Best          BestStats
// 	MatchAwards   MatchAwardsStats
// }

// GetPlayerStats : Gets all stats available for a player
func GetPlayerStats(platform, region, tag string) (*PlayerStats, error) {
	var playerStats PlayerStats

	// Performs the http request on the Overwatch website to retrieve all player info
	playerDoc, err := goquery.NewDocument("https://playoverwatch.com" +
		"/en-us/career/" + platform + "/" + region + "/" + tag)

	// Scrapes general stats info for player
	if err = unmarshalGeneralStats(
		playerDoc.Find("div.masthead").First(),
		&playerStats); err != nil {
		return nil, err
	}

	// Scrapes all Quickplay stats for player
	if err = unmarshalDetailedStats(
		playerDoc.Find("div#quick-play div.row ul.row").First(),
		&playerStats.QuickPlayStats); err != nil {
		return nil, err
	}

	// Scrapes all Competitive stats for player
	if err = unmarshalDetailedStats(
		playerDoc.Find("div#competitive-play div.row ul.row").First(),
		&playerStats.CompetitiveStats); err != nil {
		return nil, err
	}

	return &playerStats, nil
}

func unmarshalGeneralStats(generalSelector *goquery.Selection, playerStats *PlayerStats) error {
	playerStats.Icon, _ = generalSelector.Find("img.player-portrait").Attr("src")
	playerStats.Name = generalSelector.Find("h1.header-masthead").Text()
	playerStats.Level, _ = strconv.Atoi(generalSelector.Find("div.player-level div.u-vertical-center").Text())
	playerStats.LevelIcon, _ = generalSelector.Find("div.player-level").Attr("style")
	playerStats.LevelIcon = strings.Replace(playerStats.LevelIcon, "background-image:url(", "", -1)
	playerStats.LevelIcon = strings.Replace(playerStats.LevelIcon, ")", "", -1)
	playerStats.Rating = generalSelector.Find("div.competitive-rank div.u-align-center").Text()
	playerStats.RatingIcon, _ = generalSelector.Find("div.competitive-rank img").Attr("src")
	playerStats.GamesPlayed, _ = strconv.Atoi(strings.Replace(generalSelector.Find("div.masthead-player p.masthead-detail.h4 span").Text(), " games won", "", -1))
	return nil
}

func unmarshalDetailedStats(quickPlaySelector *goquery.Selection, quickPlayStats *StatsCollection) error {
	quickPlaySelector.Find("li.column").Each(func(i int, stat *goquery.Selection) {
		statType := strings.ToLower(stat.Find("p.card-copy").First().Text())
		statVal := strings.Replace(stat.Find("h3.card-heading").Text(), ",", "", -1)
		if strings.Contains(statType, "eliminations") {
			quickPlayStats.EliminationsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "damage done") {
			quickPlayStats.DamageDoneAvg, _ = strconv.ParseInt(statVal, 10, 64)
		}
		if strings.Contains(statType, "deaths") {
			quickPlayStats.DeathsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "deaths") {
			quickPlayStats.DeathsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "final blows") {
			quickPlayStats.FinalBlowsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "healing done") {
			quickPlayStats.HealingDoneAvg, _ = strconv.ParseInt(statVal, 10, 64)
		}
		if strings.Contains(statType, "objective kills") {
			quickPlayStats.ObjectiveKillsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "objective time") {
			quickPlayStats.ObjectiveTimeAvg = statVal
		}
		if strings.Contains(statType, "solo kills") {
			quickPlayStats.SoloKillsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
	})
	return nil
}
