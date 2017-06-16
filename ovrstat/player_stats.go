package ovrstat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://playoverwatch.com/en-us/career"
)

// PCStats retrieves player stats for PC
func PCStats(region, tag string) (*PlayerStats, error) {
	return playerStats(fmt.Sprintf("/pc/%s/%s", region, tag))
}

// ConsoleStats retrieves player stats for Console
func ConsoleStats(platform, tag string) (*PlayerStats, error) {
	return playerStats(fmt.Sprintf("/%s/%s", platform, tag))
}

// playerStats retrieves all Overwatch statistics for a given player
func playerStats(profilePath string) (*PlayerStats, error) {
	// Performs the stats request
	pd, err := goquery.NewDocument(baseURL + profilePath)
	if err != nil {
		return nil, err
	}

	// Scrapes all stats for the passed user and sets struct member data
	ps := parseGeneralInfo(pd.Find("div.masthead").First())
	ps.QuickPlayStats = parseDetailedStats(pd.Find("div#quickplay").First())
	ps.CompetitiveStats = parseDetailedStats(pd.Find("div#competitive").First())
	return &ps, nil
}

// populateGeneralInfo extracts the users general info and returns it in a
// PlayerStats struct
func parseGeneralInfo(s *goquery.Selection) PlayerStats {
	var ps PlayerStats

	// Populates all general player information
	ps.Icon, _ = s.Find("img.player-portrait").Attr("src")
	ps.Name = s.Find("h1.header-masthead").Text()
	ps.Level, _ = strconv.Atoi(s.Find("div.player-level div.u-vertical-center").First().Text())
	ps.LevelIcon, _ = s.Find("div.player-level").Attr("style")
	ps.LevelIcon = strings.Replace(ps.LevelIcon, "background-image:url(", "", -1)
	ps.LevelIcon = strings.Replace(ps.LevelIcon, ")", "", -1)
	ps.Prestige = getPrestigeByIcon(ps.LevelIcon)
	ps.PrestigeIcon, _ = s.Find("div.player-rank").Attr("style")
	ps.PrestigeIcon = strings.Replace(ps.PrestigeIcon, "background-image:url(", "", -1)
	ps.PrestigeIcon = strings.Replace(ps.PrestigeIcon, ")", "", -1)
	ps.Rating, _ = strconv.Atoi(s.Find("div.competitive-rank div.u-align-center").First().Text())
	ps.RatingIcon, _ = s.Find("div.competitive-rank img").Attr("src")
	ps.GamesWon, _ = strconv.Atoi(strings.Replace(s.Find("div.masthead p.masthead-detail.h4 span").Text(), " games won", "", -1))

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

	sc.TopHeros = parseHeroStats(playModeSelector.Find("section.hero-comparison-section").First())
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
	r, _ := regexp.Compile(`0x0250000000000(.+?)_Border`)
	iconID := r.FindSubmatch([]byte(levelIcon))
	if len(iconID) != 2 {
		return 0
	}
	return rankMap[string(iconID[1])]
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
