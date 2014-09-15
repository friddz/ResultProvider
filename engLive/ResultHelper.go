package englive

import (
	"github.com/PuerkitoBio/goquery"
	rp "github.com/friddz/ResultProvider"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"sort"
)

type ByDateAsc []rp.Result

func (p ByDateAsc) Len() int {
	return len(p)
}

func (p ByDateAsc) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByDateAsc) Less(i int, j int) bool {
	return p[i].Date.Before(p[j].Date)
}


func GetAllResults(id string) ([]rp.Result, error) {
	results := make([]rp.Result, 0)

	doc, err := goquery.NewDocument("http://www.livescore.com/soccer/england/premier-league/results/all/")
	if err != nil {
		log.Fatal(err)
	}

	date := time.Now()
	doc.Find(".league-table tr").Each(func(i int, s *goquery.Selection) {
		dateString := strings.TrimSpace(s.Find(".date").Text())
		if (len(dateString)) > 0 {
			date, _ = time.Parse("2006 January 02", "2014 "+dateString)
		} else {
			link, _ := (s.Find(".fs a").Attr("href"))
			link = strings.TrimSpace(link)
			results = append(results, resultDetails("http://livescore.com/"+link, date))
		}
	})

	sort.Sort(ByDateAsc(results))
	for i,_ := range results {
		results[i].Round = uint8(i / 10 + 1)
	}
	return results, nil
}

func resultDetails(url string, date time.Time) rp.Result {
	result := rp.Result{Goals: []rp.Goal{}, Cards: rp.CardInfo{}}
	result.Date = date
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	r, _ := regexp.Compile("[\\d]")
	doc.Find(".match-details tr").Each(func(i int, s *goquery.Selection) {
		if 0 == i {
			result.HomeTeamName = strings.TrimSpace(s.Find(".home").Text())
			result.AwayTeamName = strings.TrimSpace(s.Find(".awy").Text())
			scoreString := strings.TrimSpace(s.Find(".sco").Text())
			goals := r.FindAllString(scoreString, -1)
			result.HomeGoals = parseToUint8(goals[0])
			result.AwayGoals = parseToUint8(goals[1])
		} else if 1 == i {
			scoreString := strings.TrimSpace(s.Find(".sco").Text())
			goals := r.FindAllString(scoreString, -1)
			homeGoals := uint8(0)
			awayGoals := uint8(0)
			if len(goals) > 1 {
				homeGoals = parseToUint8(goals[0])
				awayGoals = parseToUint8(goals[1])
			}
			result.HomeGoalsAtHalfTime = homeGoals
			result.AwayGoalsAtHalfTime = awayGoals
		}
		if len(s.Find(".goal").Nodes) > 0 {
			min := getMinute(s)
			isHomeTeam := isHomeTeam(s, ".goal")
			name := s.Find(".goal").Parent().Find(".name").Text()
			teamName := ""
			if isHomeTeam {
				teamName = result.HomeTeamName
			} else {
				teamName = result.AwayTeamName
			}
			result.Goals = append(result.Goals, rp.Goal{GoalScorerName: name, Minute: min, TeamName: teamName})
		}
		if len(s.Find(".yellowcard").Nodes) > 0 {
			isHomeTeam := isHomeTeam(s, ".yellowcard")
			if isHomeTeam {
				result.Cards.HomeTeamNumberOfYellowCards = result.Cards.HomeTeamNumberOfYellowCards + 1
			} else {
				result.Cards.AwayTeamNumberOfYellowCards = result.Cards.AwayTeamNumberOfYellowCards + 1
			}
		}
		if len(s.Find("span.redcard").Nodes) > 0 {
			isHomeTeam := isHomeTeam(s, "span.redcard")
			if isHomeTeam {
				result.Cards.HomeTeamNumberOfRedCards = result.Cards.HomeTeamNumberOfRedCards + 1
			} else {
				result.Cards.AwayTeamNumberOfRedCards = result.Cards.AwayTeamNumberOfRedCards + 1
			}
		}
		if (len(s.Find("span.redyellowcard").Nodes) > 0) {
			isHomeTeam := isHomeTeam(s, "span.redyellowcard")
			if isHomeTeam {
				result.Cards.HomeTeamNumberOfRedCards = result.Cards.HomeTeamNumberOfRedCards + 1
			} else {
				result.Cards.AwayTeamNumberOfRedCards = result.Cards.AwayTeamNumberOfRedCards + 1
			}
	 	}
	 })
	return result
}

func parseToUint8(s string) uint8 {
	value, _ := strconv.ParseUint(s, 10, 8)
	return uint8(value)
}

func isHomeTeam(s *goquery.Selection, id string) bool {
	classAttr, exists := s.Find(id).Parent().Find(".name").Attr("class")
	isHomeTeam := true
	if exists {
		if strings.Contains(classAttr, "left") {
			isHomeTeam = false
		}
	}
	return isHomeTeam
}
func getMinute(s *goquery.Selection) uint8 {
	txt := strings.TrimSpace(s.Find(".min").Text())
	txt = strings.Replace(txt, "'", "", -1)
	i, err := strconv.ParseUint(txt, 10, 8)
	if nil != err {
		panic(err)
	}
	return uint8(i)
}