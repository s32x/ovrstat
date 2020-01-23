package ovrstat

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://playoverwatch.com/en-us/career"

	apiURL = "https://playoverwatch.com/en-us/search/account-by-name/"

	// PlatformXBL is platform : XBOX
	PlatformXBL = "xbl"

	// PlatformPSN is the platform : Playstation Network
	PlatformPSN = "psn"

	// PlatformPC is the platform : PC
	PlatformPC = "pc"

	PlatformNS = "nintendo-switch"
)

var (
	// ErrPlayerNotFound is thrown when a player doesn't exist
	ErrPlayerNotFound = errors.New("Player not found")

	// ErrInvalidPlatform is thrown when the passed params are incorrect
	ErrInvalidPlatform = errors.New("Invalid platform")
)

// Stats retrieves player stats
// Universal method if you don't need to differentiate it
func Stats(platform, tag string) (*PlayerStats, error) {
	switch platform {
	case PlatformPC:
		return PCStats(tag) // Perform a stats lookup for PC
	case PlatformPSN, PlatformXBL, PlatformNS:
		return ConsoleStats(platform, tag) // Perform a stats lookup for Console
	default:
		return nil, ErrInvalidPlatform
	}
}

// ConsoleStats retrieves player stats for Console
func ConsoleStats(platform, tag string) (*PlayerStats, error) {
	return playerStats(fmt.Sprintf("/%s/%s", platform, tag), platform)
}

// PCStats retrieves player stats for PC
func PCStats(tag string) (*PlayerStats, error) {
	return playerStats(fmt.Sprintf("/pc/%s", tag), "pc")
}

// playerStats retrieves all Overwatch statistics for a given player
func playerStats(profilePath string, platform string) (*PlayerStats, error) {
	// Create the profile url for scraping
	url := baseURL + profilePath

	// Perform the stats request and decode the response
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve profile")
	}
	defer res.Body.Close()

	// Parses the stats request into a goquery document
	pd, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create goquery document")
	}

	// Checks if profile not found, site still returns 200 in this case
	if pd.Find("h1.u-align-center").First().Text() == "Profile Not Found" {
		return nil, ErrPlayerNotFound
	}

	// Scrapes all stats for the passed user and sets struct member data
	ps := parseGeneralInfo(pd.Find("div.masthead").First())

	// Perform api request
	type Platform struct {
		Platform    string `json:"platform"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		URLName     string `json:"urlName"`
		PlayerLevel int    `json:"playerLevel"`
		Portrait    string `json:"portrait"`
		IsPublic    bool   `json:"isPublic"`
	}
	var platforms []Platform
	apires, err := http.Get(apiURL + strings.Replace(profilePath[strings.LastIndex(profilePath, "/")+1:], "-", "%23", -1))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to perform platform API request")
	}
	defer apires.Body.Close()

	// Decode received JSON
	if err := json.NewDecoder(apires.Body).Decode(&platforms); err != nil {
		return nil, errors.Wrap(err, "Failed to decode platform API response")
	}

	for _, p := range platforms {
		if p.Platform == platform {
			ps.Name = p.Name
			ps.Prestige = int(math.Floor(float64(p.PlayerLevel) / 100))
		}
	}

	if pd.Find("p.masthead-permission-level-text").First().Text() == "Private Profile" {
		ps.Private = true
		return &ps, nil
	}

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
	ps.Level, _ = strconv.Atoi(s.Find("div.player-level div.u-vertical-center").First().Text())
	ps.LevelIcon, _ = s.Find("div.player-level").Attr("style")
	ps.LevelIcon = strings.Replace(ps.LevelIcon, "background-image:url(", "", -1)
	ps.LevelIcon = strings.Replace(ps.LevelIcon, ")", "", -1)
	ps.LevelIcon = strings.TrimSpace(ps.LevelIcon)
	ps.PrestigeIcon, _ = s.Find("div.player-rank").Attr("style")
	ps.PrestigeIcon = strings.Replace(ps.PrestigeIcon, "background-image:url(", "", -1)
	ps.PrestigeIcon = strings.Replace(ps.PrestigeIcon, ")", "", -1)
	ps.PrestigeIcon = strings.TrimSpace(ps.PrestigeIcon)
	ps.Endorsement, _ = strconv.Atoi(s.Find("div.EndorsementIcon-tooltip div.u-center").First().Text())
	ps.EndorsementIcon, _ = s.Find("div.EndorsementIcon").Attr("style")
	ps.EndorsementIcon = strings.Replace(ps.EndorsementIcon, "background-image:url(", "", -1)
	ps.EndorsementIcon = strings.Replace(ps.EndorsementIcon, ")", "", -1)

	// Ratings.
	s.Find("div.show-for-lg div.competitive-rank div.competitive-rank-role").Each(func(i int, rankSel *goquery.Selection) {
		// Rank selections.
		sel := rankSel.Find("div.competitive-rank-section")

		role, _ := sel.Find("div.competitive-rank-tier.competitive-rank-tier-tooltip").Attr("data-ow-tooltip-text")
		roleIcon, _ := sel.Find("img.competitive-rank-role-icon").Attr("src")
		rankIcon, _ := sel.Find("div.competitive-rank-tier.competitive-rank-tier-tooltip img.competitive-rank-tier-icon").Attr("src")
		level, _ := strconv.Atoi(sel.Find("div.competitive-rank-level").Text())
		formatedRole := strings.TrimSuffix(role, " Skill Rating")

		ps.Ratings = append(ps.Ratings, Rating{
			Level:    level,
			Role:     strings.ToLower(formatedRole),
			RoleIcon: roleIcon,
			RankIcon: rankIcon,
		})
	})

	ps.GamesWon, _ = strconv.Atoi(strings.Replace(s.Find("div.masthead p.masthead-detail.h4 span").Text(), " games won", "", -1))

	return ps
}

// parseDetailedStats populates the passed stats collection with detailed statistics
func parseDetailedStats(playModeSelector *goquery.Selection) statsCollection {
	var sc statsCollection
	sc.TopHeroes = parseHeroStats(playModeSelector.Find("div.progress-category").Parent())
	sc.CareerStats = parseCareerStats(playModeSelector.Find("div.js-stats").Parent())
	return sc
}

// parseHeroStats : Parses stats for each individual hero and returns a map
func parseHeroStats(heroStatsSelector *goquery.Selection) map[string]*topHeroStats {
	bhsMap := make(map[string]*topHeroStats)

	heroStatsSelector.Find("div.progress-category").Each(func(i int, heroGroupSel *goquery.Selection) {
		categoryID, _ := heroGroupSel.Attr("data-category-id")
		categoryID = strings.Replace(categoryID, "0x0860000000000", "", -1)
		heroGroupSel.Find("div.ProgressBar").Each(func(i2 int, statSel *goquery.Selection) {
			heroName := cleanJSONKey(statSel.Find("div.ProgressBar-title").Text())
			statVal := statSel.Find("div.ProgressBar-description").Text()

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
			case "31C":
				bhsMap[heroName].ObjectiveKills, _ = strconv.ParseFloat(statVal, 64)
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
	careerStatsSelector.Find("select option").Each(func(i int, heroSel *goquery.Selection) {
		heroVal, _ := heroSel.Attr("value")
		heroMap[heroVal] = heroSel.Text()
	})

	// Iterates over every hero div
	careerStatsSelector.Find("div.row.js-stats").Each(func(i int, heroStatsSel *goquery.Selection) {
		currentHero, _ := heroStatsSel.Attr("data-category-id")
		currentHero = cleanJSONKey(heroMap[currentHero])

		// Iterates over every stat box
		heroStatsSel.Find("div.card-stat-block-container").Each(func(i2 int, statBoxSel *goquery.Selection) {
			statType := statBoxSel.Find(".stat-title").Text()
			statType = cleanJSONKey(statType)

			// Iterates over stat row
			statBoxSel.Find("table.DataTable tbody tr").Each(func(i3 int, statSel *goquery.Selection) {

				// Iterates over every stat td
				statKey := ""
				statVal := ""
				statSel.Find("td").Each(func(i4 int, statKV *goquery.Selection) {
					switch i4 {
					case 0:
						statKey = transformKey(cleanJSONKey(statKV.Text()))
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
								csMap[currentHero].Assists = make(map[string]interface{})
							}
							csMap[currentHero].Assists[statKey] = parseType(statVal)
						case "average":
							if csMap[currentHero].Average == nil {
								csMap[currentHero].Average = make(map[string]interface{})
							}
							csMap[currentHero].Average[statKey] = parseType(statVal)
						case "best":
							if csMap[currentHero].Best == nil {
								csMap[currentHero].Best = make(map[string]interface{})
							}
							csMap[currentHero].Best[statKey] = parseType(statVal)
						case "combat":
							if csMap[currentHero].Combat == nil {
								csMap[currentHero].Combat = make(map[string]interface{})
							}
							csMap[currentHero].Combat[statKey] = parseType(statVal)
						case "deaths":
							if csMap[currentHero].Deaths == nil {
								csMap[currentHero].Deaths = make(map[string]interface{})
							}
							csMap[currentHero].Deaths[statKey] = parseType(statVal)
						case "heroSpecific":
							if csMap[currentHero].HeroSpecific == nil {
								csMap[currentHero].HeroSpecific = make(map[string]interface{})
							}
							csMap[currentHero].HeroSpecific[statKey] = parseType(statVal)
						case "game":
							if csMap[currentHero].Game == nil {
								csMap[currentHero].Game = make(map[string]interface{})
							}
							csMap[currentHero].Game[statKey] = parseType(statVal)
						case "matchAwards":
							if csMap[currentHero].MatchAwards == nil {
								csMap[currentHero].MatchAwards = make(map[string]interface{})
							}
							csMap[currentHero].MatchAwards[statKey] = parseType(statVal)
						case "miscellaneous":
							if csMap[currentHero].Miscellaneous == nil {
								csMap[currentHero].Miscellaneous = make(map[string]interface{})
							}
							csMap[currentHero].Miscellaneous[statKey] = parseType(statVal)
						}
					}
				})
			})
		})
	})
	return csMap
}

func parseType(val string) interface{} {
	i, err := strconv.Atoi(val)
	if err == nil {
		return i
	}
	f, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return f
	}
	return val
}

var (
	keyReplacer = strings.NewReplacer("-", " ", ".", " ", ":", " ", "'", "", "ú", "u", "ö", "o")
)

// cleanJSONKey
func cleanJSONKey(str string) string {
	// Removes localization rubish
	if strings.Contains(str, "} other {") {
		re := regexp.MustCompile("{count, plural, one {.+} other {(.+)}}")
		if len(re.FindStringSubmatch(str)) == 2 {
			otherForm := re.FindStringSubmatch(str)[1]
			str = re.ReplaceAllString(str, otherForm)
		}
	}

	str = keyReplacer.Replace(str) // Removes all dashes, dots, and colons from titles
	str = strings.ToLower(str)
	str = strings.Title(str)                // Uppercases lowercase leading characters
	str = strings.Replace(str, " ", "", -1) // Removes Spaces
	for i, v := range str {                 // Lowercases initial character
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
