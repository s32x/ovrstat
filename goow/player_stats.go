package goow

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// PlayerStats : Holds all stats on a specified Overwatch player
type PlayerStats struct {
	Icon             string          `json:"icon"`
	Name             string          `json:"name"`
	Level            int             `json:"level"`
	LevelIcon        string          `json:"levelIcon"`
	Rating           string          `json:"rating"`
	RatingIcon       string          `json:"ratingIcon"`
	GamesPlayed      int             `json:"gamesPlayed"`
	QuickPlayStats   statsCollection `json:"quickPlayStats"`
	CompetitiveStats statsCollection `json:"competitiveStats"`
}

// statsCollection : Holds a collection of stats for a particular player
type statsCollection struct {
	EliminationsAvg   float64                    `json:"eliminationsAvg"`
	DamageDoneAvg     int64                      `json:"damageDoneAvg"`
	DeathsAvg         float64                    `json:"deathsAvg"`
	FinalBlowsAvg     float64                    `json:"finalBlowsAvg"`
	HealingDoneAvg    int64                      `json:"healingDoneAvg"`
	ObjectiveKillsAvg float64                    `json:"objectiveKillsAvg"`
	ObjectiveTimeAvg  string                     `json:"objectiveTimeAvg"`
	SoloKillsAvg      float64                    `json:"soloKillsAvg"`
	TopHeros          map[string]*basicHeroStats `json:"topHeros"`
}

// basicHeroStats : Holds specific stats for each hero
type basicHeroStats struct {
	TimePlayed          string  `json:"timePlayed"`
	GamesWon            int     `json:"gamesWon"`
	WinPercentage       int     `json:"winPercentage"`
	WeaponAccuracy      int     `json:"weaponAccuracy"`
	EliminationsPerLife float64 `json:"eliminationsPerLife"`
	MultiKillBest       int     `json:"multiKillBest"`
	ObjectiveKillsAvg   float64 `json:"objectiveKillsAvg"`
}

// GetPlayerStats : Gets all stats available for a player
func GetPlayerStats(platform, region, tag string) (PlayerStats, error) {
	// Performs the http request on the Overwatch website to retrieve all player info
	playerDoc, err := goquery.NewDocument("https://playoverwatch.com" +
		"/en-us/career/" + platform + "/" + region + "/" + tag)
	if err != nil {
		return PlayerStats{}, err
	}

	// Scrapes all stats for passed player and sets struct member data
	ps := parseGeneralInfo(playerDoc.Find("div.masthead").First())
	ps.QuickPlayStats = parseDetailedStats(playerDoc.Find("div#quick-play").First())
	ps.CompetitiveStats = parseDetailedStats(playerDoc.Find("div#competitive-play").First())

	return ps, nil
}

// populateGeneralInfo : Populates the passed playerStats with generic play stats
func parseGeneralInfo(generalSelector *goquery.Selection) PlayerStats {
	var ps PlayerStats

	// Populates all general basic stats for the player
	ps.Icon, _ = generalSelector.Find("img.player-portrait").Attr("src")
	ps.Name = generalSelector.Find("h1.header-masthead").Text()
	ps.Level, _ = strconv.Atoi(generalSelector.Find("div.player-level div.u-vertical-center").Text())
	ps.LevelIcon, _ = generalSelector.Find("div.player-level").Attr("style")
	ps.LevelIcon = strings.Replace(ps.LevelIcon, "background-image:url(", "", -1)
	ps.LevelIcon = strings.Replace(ps.LevelIcon, ")", "", -1)
	ps.Rating = generalSelector.Find("div.competitive-rank div.u-align-center").Text()
	ps.RatingIcon, _ = generalSelector.Find("div.competitive-rank img").Attr("src")
	ps.GamesPlayed, _ = strconv.Atoi(strings.Replace(generalSelector.Find("div.masthead-player p.masthead-detail.h4 span").Text(), " games won", "", -1))

	return ps
}

// parseDetailedStats : Populates the passed stats collection with detailed statistics
func parseDetailedStats(playModeSelector *goquery.Selection) statsCollection {
	var sc statsCollection

	// Populates all detailed basic stats for the player
	playModeSelector.Find("li.column").Each(func(i int, statSel *goquery.Selection) {
		statType := statSel.Find("p.card-copy").First().Text()
		statType = strings.Replace(strings.ToLower(statType), " - average", "", -1)
		statVal := strings.Replace(statSel.Find("h3.card-heading").Text(), ",", "", -1)

		switch statType {
		case "eliminations":
			sc.EliminationsAvg, _ = strconv.ParseFloat(statVal, 64)
		case "damage done":
			sc.DamageDoneAvg, _ = strconv.ParseInt(statVal, 10, 64)
		case "deaths":
			sc.DeathsAvg, _ = strconv.ParseFloat(statVal, 64)
		case "final blows":
			sc.FinalBlowsAvg, _ = strconv.ParseFloat(statVal, 64)
		case "healing done":
			sc.HealingDoneAvg, _ = strconv.ParseInt(statVal, 10, 64)
		case "objective kills":
			sc.ObjectiveKillsAvg, _ = strconv.ParseFloat(statVal, 64)
		case "objective time":
			sc.ObjectiveTimeAvg = statVal
		case "solo kills":
			sc.SoloKillsAvg, _ = strconv.ParseFloat(statVal, 64)
		}
	})

	// Parses out top hero stats and assigns it to our parent struct
	sc.TopHeros = parseHeroStats(playModeSelector.Find("section.hero-comparison-section").First())

	return sc
}

// parseHeroStats : Parses stats for each individual hero and returns a map
func parseHeroStats(heroStatsSelector *goquery.Selection) map[string]*basicHeroStats {
	bhsMap := make(map[string]*basicHeroStats)

	heroStatsSelector.Find("div.progress-category").Each(func(i int, heroGroupSel *goquery.Selection) {
		categoryID, _ := heroGroupSel.Attr("data-category-id")
		categoryID = strings.Replace(categoryID, "overwatch.guid.0x0860000000000", "", -1)
		heroGroupSel.Find("div.progress-2").Each(func(i2 int, statSel *goquery.Selection) {
			heroName := statSel.Find("div.title").Text()
			statVal := statSel.Find("div.description").Text()

			// Creates hero map if it doesn't exist
			if bhsMap[heroName] == nil {
				bhsMap[heroName] = new(basicHeroStats)
			}

			// Sets hero stats based on stat category type
			switch categoryID {
			case "021":
				bhsMap[heroName].TimePlayed = statVal
			case "039":
				bhsMap[heroName].GamesWon, _ = strconv.Atoi(statVal)
			case "3D1":
				bhsMap[heroName].WinPercentage, _ = strconv.Atoi(strings.Replace(statVal, "%", "", -1))
			case "02F":
				bhsMap[heroName].WeaponAccuracy, _ = strconv.Atoi(strings.Replace(statVal, "%", "", -1))
			case "3D2":
				bhsMap[heroName].EliminationsPerLife, _ = strconv.ParseFloat(statVal, 64)
			case "346":
				bhsMap[heroName].MultiKillBest, _ = strconv.Atoi(statVal)
			case "39C":
				bhsMap[heroName].ObjectiveKillsAvg, _ = strconv.ParseFloat(statVal, 64)
			}
		})
	})

	return bhsMap
}
