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
	TopHeros          map[string]*HeroStats
}

type HeroStats struct {
	TimePlayed          string
	GamesWon            int
	WinPercentage       int
	WeaponAccuracy      int
	EliminationsPerLife float64
	MultiKillBest       int
	ObjectiveKillsAvg   float64
}

// GetPlayerStats : Gets all stats available for a player
func GetPlayerStats(platform, region, tag string) (*PlayerStats, error) {
	var playerStats PlayerStats

	// Performs the http request on the Overwatch website to retrieve all player info
	playerDoc, err := goquery.NewDocument("https://playoverwatch.com" +
		"/en-us/career/" + platform + "/" + region + "/" + tag)

	// Scrapes general stats info for player
	if err = populateGeneralInfo(
		playerDoc.Find("div.masthead").First(),
		&playerStats); err != nil {
		return nil, err
	}

	// Scrapes all Quickplay stats for player
	if err = populateDetailedStats(
		playerDoc.Find("div#quick-play").First(),
		&playerStats.QuickPlayStats); err != nil {
		return nil, err
	}

	// Scrapes all Competitive stats for player
	if err = populateDetailedStats(
		playerDoc.Find("div#competitive-play").First(),
		&playerStats.CompetitiveStats); err != nil {
		return nil, err
	}

	return &playerStats, nil
}

func populateGeneralInfo(generalSelector *goquery.Selection, playerStats *PlayerStats) error {
	// Populates all general basic stats for the player
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

func populateDetailedStats(playModeSelector *goquery.Selection, statsColl *StatsCollection) error {
	// Populates all detailed basic stats for the player
	playModeSelector.Find("li.column").Each(func(i int, statSel *goquery.Selection) {
		statType := strings.ToLower(statSel.Find("p.card-copy").First().Text())
		statVal := strings.Replace(statSel.Find("h3.card-heading").Text(), ",", "", -1)
		if strings.Contains(statType, "eliminations") {
			statsColl.EliminationsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "damage done") {
			statsColl.DamageDoneAvg, _ = strconv.ParseInt(statVal, 10, 64)
		}
		if strings.Contains(statType, "deaths") {
			statsColl.DeathsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "deaths") {
			statsColl.DeathsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "final blows") {
			statsColl.FinalBlowsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "healing done") {
			statsColl.HealingDoneAvg, _ = strconv.ParseInt(statVal, 10, 64)
		}
		if strings.Contains(statType, "objective kills") {
			statsColl.ObjectiveKillsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
		if strings.Contains(statType, "objective time") {
			statsColl.ObjectiveTimeAvg = statVal
		}
		if strings.Contains(statType, "solo kills") {
			statsColl.SoloKillsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
	})

	// Parses out top hero stats and assigns it to our parent struct
	statsColl.TopHeros = parseHeroStats(playModeSelector.Find("section.hero-comparison-section").First())

	return nil
}

func parseHeroStats(heroStatsSelector *goquery.Selection) map[string]*HeroStats {
	tempHeroStatMap := make(map[string]*HeroStats)

	heroStatsSelector.Find("div.progress-category").Each(func(i int, heroGroupSel *goquery.Selection) {
		categoryID, _ := heroGroupSel.Attr("data-category-id")
		categoryID = strings.Replace(categoryID, "overwatch.guid.0x0860000000000", "", -1)
		heroGroupSel.Find("div.progress-2").Each(func(i2 int, statSel *goquery.Selection) {
			heroName := statSel.Find("div.title").Text()
			statVal := statSel.Find("div.description").Text()

			// Creates hero map if it doesn't exist
			if tempHeroStatMap[heroName] == nil {
				tempHeroStatMap[heroName] = new(HeroStats)
			}

			// Sets hero stats
			if categoryID == "021" {
				tempHeroStatMap[heroName].TimePlayed = statVal
			} else if categoryID == "039" {
				tempHeroStatMap[heroName].GamesWon, _ = strconv.Atoi(statVal)
			} else if categoryID == "3D1" {
				tempHeroStatMap[heroName].WinPercentage, _ = strconv.Atoi(strings.Replace(statVal, "%", "", -1))
			} else if categoryID == "02F" {
				tempHeroStatMap[heroName].WeaponAccuracy, _ = strconv.Atoi(strings.Replace(statVal, "%", "", -1))
			} else if categoryID == "3D2" {
				tempHeroStatMap[heroName].EliminationsPerLife, _ = strconv.ParseFloat(statVal, 64)
			} else if categoryID == "346" {
				tempHeroStatMap[heroName].MultiKillBest, _ = strconv.Atoi(statVal)
			} else if categoryID == "39C" {
				tempHeroStatMap[heroName].ObjectiveKillsAvg, _ = strconv.ParseFloat(statVal, 64)
			}
		})
	})

	return tempHeroStatMap
}
