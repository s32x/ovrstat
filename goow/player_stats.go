package goow

import (
	"strconv"
	"strings"
	"unicode"

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

// topHeroStats : Holds basic stats for each hero
type topHeroStats struct {
	TimePlayed          string  `json:"timePlayed"`
	GamesWon            int     `json:"gamesWon"`
	WinPercentage       int     `json:"winPercentage"`
	WeaponAccuracy      int     `json:"weaponAccuracy"`
	EliminationsPerLife float64 `json:"eliminationsPerLife"`
	MultiKillBest       int     `json:"multiKillBest"`
	ObjectiveKillsAvg   float64 `json:"objectiveKillsAvg"`
}

// careerStats : Holds very detailed stats for each hero
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

// GetPlayerStats : Gets all stats available for a player
func GetPlayerStats(platform, region, tag string) (PlayerStats, error) {
	// Creates the profile url page based on platform
	url := "https://playoverwatch.com" + "/en-us/career/" + platform + "/" + region + "/" + tag
	if platform != "pc" {
		url = "https://playoverwatch.com" + "/en-us/career/" + platform + "/" + tag
	}

	// Performs the http request on the Overwatch website to retrieve all player info
	playerDoc, err := goquery.NewDocument(url)
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

	// Parses out career stats and assigns it to our parent struct
	sc.CareerStats = parseCareerStats(playModeSelector.Find("section.career-stats-section").First())

	return sc
}

// parseHeroStats : Parses stats for each individual hero and returns a map
func parseHeroStats(heroStatsSelector *goquery.Selection) map[string]*topHeroStats {
	bhsMap := make(map[string]*topHeroStats)

	heroStatsSelector.Find("div.progress-category").Each(func(i int, heroGroupSel *goquery.Selection) {
		categoryID, _ := heroGroupSel.Attr("data-category-id")
		categoryID = strings.Replace(categoryID, "overwatch.guid.0x0860000000000", "", -1)
		heroGroupSel.Find("div.progress-2").Each(func(i2 int, statSel *goquery.Selection) {
			heroName := cleanJSONKey(statSel.Find("div.title").Text())
			statVal := statSel.Find("div.description").Text()

			// Creates hero map if it doesn't exist
			if bhsMap[heroName] == nil {
				bhsMap[heroName] = new(topHeroStats)
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

func parseCareerStats(careerStatsSelector *goquery.Selection) map[string]*careerStats {
	csMap := make(map[string]*careerStats)

	heroMap := make(map[string]string)
	// Populates tempHeroMap to match hero ID to name in second scrape
	careerStatsSelector.Find("select.js-career-select option").Each(func(i int, heroSel *goquery.Selection) {
		heroVal, _ := heroSel.Attr("value")
		heroMap[heroVal] = heroSel.Text()
	})

	// Iterates over every hero div
	careerStatsSelector.Find("div.row div.js-stats").Each(func(i int, heroStatsSel *goquery.Selection) {
		currentHero, _ := heroStatsSel.Attr("data-category-id")
		currentHero = cleanJSONKey(heroMap[currentHero])

		// Iterates over every stat box
		heroStatsSel.Find("div.column.xs-12").Each(func(i2 int, statBoxSel *goquery.Selection) {
			statType := statBoxSel.Find("span.stat-title").Text()
			statType = cleanJSONKey(statType)

			// Iterates over stat row
			statBoxSel.Find("table.data-table tbody tr").Each(func(i3 int, statSel *goquery.Selection) {

				// Iterates over every stat td
				statKey := ""
				statVal := ""
				statSel.Find("td").Each(func(i4 int, statKV *goquery.Selection) {
					switch i4 {
					case 0:
						statKey = cleanJSONKey(statKV.Text())
					case 1:
						statVal = strings.Replace(statKV.Text(), ",", "", -1) // Removes commas from 1k+ values

						// Creates stat map if it doesn't exist
						if csMap[currentHero] == nil {
							csMap[currentHero] = new(careerStats)
						}

						// Switches on type, creating category stat maps if exists (will omitempty on json marshal)
						switch statType {
						case "assists":
							if csMap[currentHero].Assists == nil {
								csMap[currentHero].Assists = make(map[string]string)
							}
							csMap[currentHero].Assists[statKey] = statVal
						case "average":
							if csMap[currentHero].Average == nil {
								csMap[currentHero].Average = make(map[string]string)
							}
							csMap[currentHero].Average[statKey] = statVal
						case "best":
							if csMap[currentHero].Best == nil {
								csMap[currentHero].Best = make(map[string]string)
							}
							csMap[currentHero].Best[statKey] = statVal
						case "combat":
							if csMap[currentHero].Combat == nil {
								csMap[currentHero].Combat = make(map[string]string)
							}
							csMap[currentHero].Combat[statKey] = statVal
						case "deaths":
							if csMap[currentHero].Deaths == nil {
								csMap[currentHero].Deaths = make(map[string]string)
							}
							csMap[currentHero].Deaths[statKey] = statVal
						case "hero specific":
							if csMap[currentHero].HeroSpecific == nil {
								csMap[currentHero].HeroSpecific = make(map[string]string)
							}
							csMap[currentHero].HeroSpecific[statKey] = statVal
						case "game":
							if csMap[currentHero].Game == nil {
								csMap[currentHero].Game = make(map[string]string)
							}
							csMap[currentHero].Game[statKey] = statVal
						case "match awards":
							if csMap[currentHero].MatchAwards == nil {
								csMap[currentHero].MatchAwards = make(map[string]string)
							}
							csMap[currentHero].MatchAwards[statKey] = statVal
						case "miscellaneous":
							if csMap[currentHero].Miscellaneous == nil {
								csMap[currentHero].Miscellaneous = make(map[string]string)
							}
							csMap[currentHero].Miscellaneous[statKey] = statVal
						}
					}
				})
			})
		})
	})

	return csMap
}

func cleanJSONKey(str string) string {
	str = strings.Replace(str, "-", " ", -1) // Removes all dashes from titles
	str = strings.ToLower(str)
	str = strings.Title(str)                // Uppercases lowercase leading characters
	str = strings.Replace(str, " ", "", -1) // Removes Spaces
	for i, v := range str {                 // Lowercases initial character
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
