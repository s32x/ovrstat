package goow

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// PlayerStats holds all stats on a specified Overwatch player
type PlayerStats struct {
	Icon             string          `json:"icon"`
	Name             string          `json:"name"`
	Level            int             `json:"level"`
	LevelIcon        string          `json:"levelIcon"`
	Prestige         int             `json:"prestige"`
	PrestigeIcon     string          `json:"prestigeIcon"`
	Rating           string          `json:"rating"`
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

// GetPlayerStats gets all stats available for a player
func GetPlayerStats(platform, region, tag string) (*PlayerStats, error) {
	// Creates the profile url page based on platform
	url := "https://playoverwatch.com/en-us/career/" + platform + "/" + region + "/" + tag
	if platform != "pc" {
		url = "https://playoverwatch.com/en-us/career/" + platform + "/" + tag
	}

	// Performs the http request on the Overwatch website to retrieve all player info
	playerDoc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	// Scrapes all stats for passed player and sets struct member data
	ps := parseGeneralInfo(playerDoc.Find("div.masthead").First())
	ps.QuickPlayStats = parseDetailedStats(playerDoc.Find("div#quickplay").First())
	ps.CompetitiveStats = parseDetailedStats(playerDoc.Find("div#competitive").First())

	return &ps, nil
}

// populateGeneralInfo populates the passed playerStats with generic play stats
func parseGeneralInfo(generalSelector *goquery.Selection) PlayerStats {
	var ps PlayerStats

	// Populates all general basic stats for the player
	ps.Icon, _ = generalSelector.Find("img.player-portrait").Attr("src")
	ps.Name = generalSelector.Find("h1.header-masthead").Text()
	ps.Level, _ = strconv.Atoi(generalSelector.Find("div.player-level div.u-vertical-center").First().Text())
	ps.LevelIcon, _ = generalSelector.Find("div.player-level").Attr("style")
	ps.LevelIcon = strings.Replace(ps.LevelIcon, "background-image:url(", "", -1)
	ps.LevelIcon = strings.Replace(ps.LevelIcon, ")", "", -1)
	ps.Prestige = getPrestigeByIcon(ps.LevelIcon)
	ps.PrestigeIcon, _ = generalSelector.Find("div.player-rank").Attr("style")
	ps.PrestigeIcon = strings.Replace(ps.PrestigeIcon, "background-image:url(", "", -1)
	ps.PrestigeIcon = strings.Replace(ps.PrestigeIcon, ")", "", -1)
	ps.Rating = generalSelector.Find("div.competitive-rank div.u-align-center").First().Text()
	ps.RatingIcon, _ = generalSelector.Find("div.competitive-rank img").Attr("src")
	ps.GamesWon, _ = strconv.Atoi(strings.Replace(generalSelector.Find("div.masthead p.masthead-detail.h4 span").Text(), " games won", "", -1))

	return ps
}

// parseDetailedStats populates the passed stats collection with detailed statistics
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

// parseCareerStats
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

func getPrestigeByIcon(levelIcon string) int {
	levelIconID := strings.Replace(levelIcon, "https://blzgdapipro-a.akamaihd.net/game/playerlevelrewards/", "", -1)
	levelIconID = strings.Replace(levelIconID, "_Border.png", "", -1)

	rankMap := map[string]int{
		//# Bronze 0 - 5
		"0x0250000000000918": 0,
		"0x0250000000000919": 0,
		"0x025000000000091A": 0,
		"0x025000000000091B": 0,
		"0x025000000000091C": 0,
		"0x025000000000091D": 0,
		"0x025000000000091E": 0,
		"0x025000000000091F": 0,
		"0x0250000000000920": 0,
		"0x0250000000000921": 0,
		"0x0250000000000922": 1,
		"0x0250000000000924": 1,
		"0x0250000000000925": 1,
		"0x0250000000000926": 1,
		"0x025000000000094C": 1,
		"0x0250000000000927": 1,
		"0x0250000000000928": 1,
		"0x0250000000000929": 1,
		"0x025000000000092B": 1,
		"0x0250000000000950": 1,
		"0x025000000000092A": 2,
		"0x025000000000092C": 2,
		"0x0250000000000937": 2,
		"0x025000000000093B": 2,
		"0x0250000000000933": 2,
		"0x0250000000000923": 2,
		"0x0250000000000944": 2,
		"0x0250000000000948": 2,
		"0x025000000000093F": 2,
		"0x0250000000000951": 2,
		"0x025000000000092D": 3,
		"0x0250000000000930": 3,
		"0x0250000000000934": 3,
		"0x0250000000000938": 3,
		"0x0250000000000940": 3,
		"0x0250000000000949": 3,
		"0x0250000000000952": 3,
		"0x025000000000094D": 3,
		"0x0250000000000945": 3,
		"0x025000000000093C": 3,
		"0x025000000000092E": 4,
		"0x0250000000000931": 4,
		"0x0250000000000935": 4,
		"0x025000000000093D": 4,
		"0x0250000000000946": 4,
		"0x025000000000094A": 4,
		"0x0250000000000953": 4,
		"0x025000000000094E": 4,
		"0x0250000000000939": 4,
		"0x0250000000000941": 4,
		"0x025000000000092F": 5,
		"0x0250000000000932": 5,
		"0x025000000000093E": 5,
		"0x0250000000000936": 5,
		"0x025000000000093A": 5,
		"0x0250000000000942": 5,
		"0x0250000000000947": 5,
		"0x025000000000094F": 5,
		"0x025000000000094B": 5,
		"0x0250000000000954": 5,
		//# Silver 6 - 11
		"0x0250000000000956": 6,
		"0x025000000000095C": 6,
		"0x025000000000095D": 6,
		"0x025000000000095E": 6,
		"0x025000000000095F": 6,
		"0x0250000000000960": 6,
		"0x0250000000000961": 6,
		"0x0250000000000962": 6,
		"0x0250000000000963": 6,
		"0x0250000000000964": 6,
		"0x0250000000000957": 7,
		"0x0250000000000965": 7,
		"0x0250000000000966": 7,
		"0x0250000000000967": 7,
		"0x0250000000000968": 7,
		"0x0250000000000969": 7,
		"0x025000000000096A": 7,
		"0x025000000000096B": 7,
		"0x025000000000096C": 7,
		"0x025000000000096D": 7,
		"0x0250000000000958": 8,
		"0x025000000000096E": 8,
		"0x025000000000096F": 8,
		"0x0250000000000970": 8,
		"0x0250000000000971": 8,
		"0x0250000000000972": 8,
		"0x0250000000000973": 8,
		"0x0250000000000974": 8,
		"0x0250000000000975": 8,
		"0x0250000000000976": 8,
		"0x0250000000000959": 9,
		"0x0250000000000977": 9,
		"0x0250000000000978": 9,
		"0x0250000000000979": 9,
		"0x025000000000097A": 9,
		"0x025000000000097B": 9,
		"0x025000000000097C": 9,
		"0x025000000000097D": 9,
		"0x025000000000097E": 9,
		"0x025000000000097F": 9,
		"0x025000000000095A": 10,
		"0x0250000000000980": 10,
		"0x0250000000000981": 10,
		"0x0250000000000982": 10,
		"0x0250000000000983": 10,
		"0x0250000000000984": 10,
		"0x0250000000000985": 10,
		"0x0250000000000986": 10,
		"0x0250000000000987": 10,
		"0x0250000000000988": 10,
		"0x025000000000095B": 11,
		"0x0250000000000989": 11,
		"0x025000000000098A": 11,
		"0x025000000000098B": 11,
		"0x025000000000098C": 11,
		"0x025000000000098D": 11,
		"0x025000000000098E": 11,
		"0x025000000000098F": 11,
		"0x0250000000000991": 11,
		"0x0250000000000990": 11,
		//# Gold 12 - 17
		"0x0250000000000992": 12,
		"0x0250000000000993": 12,
		"0x0250000000000994": 12,
		"0x0250000000000995": 12,
		"0x0250000000000996": 12,
		"0x0250000000000997": 12,
		"0x0250000000000998": 12,
		"0x0250000000000999": 12,
		"0x025000000000099A": 12,
		"0x025000000000099B": 12,
		"0x025000000000099C": 13,
		"0x025000000000099D": 13,
		"0x025000000000099E": 13,
		"0x025000000000099F": 13,
		"0x02500000000009A0": 13,
		"0x02500000000009A1": 13,
		"0x02500000000009A2": 13,
		"0x02500000000009A3": 13,
		"0x02500000000009A4": 13,
		"0x02500000000009A5": 13,
		"0x02500000000009A6": 14,
		"0x02500000000009A7": 14,
		"0x02500000000009A8": 14,
		"0x02500000000009A9": 14,
		"0x02500000000009AA": 14,
		"0x02500000000009AB": 14,
		"0x02500000000009AC": 14,
		"0x02500000000009AD": 14,
		"0x02500000000009AE": 14,
		"0x02500000000009AF": 14,
		"0x02500000000009B0": 15,
		"0x02500000000009B1": 15,
		"0x02500000000009B2": 15,
		"0x02500000000009B3": 15,
		"0x02500000000009B4": 15,
		"0x02500000000009B5": 15,
		"0x02500000000009B6": 15,
		"0x02500000000009B7": 15,
		"0x02500000000009B8": 15,
		"0x02500000000009B9": 15,
		"0x02500000000009BA": 16,
		"0x02500000000009BB": 16,
		"0x02500000000009BC": 16,
		"0x02500000000009BD": 16,
		"0x02500000000009BE": 16,
		"0x02500000000009BF": 16,
		"0x02500000000009C0": 16,
		"0x02500000000009C1": 16,
		"0x02500000000009C2": 16,
		"0x02500000000009C3": 16,
		"0x02500000000009C4": 17,
		"0x02500000000009C5": 17,
		"0x02500000000009C6": 17,
		"0x02500000000009C7": 17,
		"0x02500000000009C8": 17,
		"0x02500000000009C9": 17,
		"0x02500000000009CA": 17,
		"0x02500000000009CB": 17,
		"0x02500000000009CC": 17,
		"0x02500000000009CD": 17,
	}
	return rankMap[levelIconID]
}

// cleanJSONKey
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
