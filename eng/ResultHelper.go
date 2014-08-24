package eng

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	rp "resultprovider"
	"strconv"
	"time"
)

var downloadUrl string = "http://www.football-data.co.uk/mmz4281/1415/"

func getInternalSeasonId(externalId string) string {

	if externalId == "2015" {
		return "E0.csv"
	} else {
		panic("Unknown season id")
	}
}

func GetAllResults(id string) ([]rp.Result, error) {
	results := make([]rp.Result, 0)
	internalSeasonId := getInternalSeasonId(id)
	resp, err := http.Get(downloadUrl + internalSeasonId)
	if nil != err {
		return results, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	re := regexp.MustCompile(",")
	for scanner.Scan() {
		line := scanner.Text()
		match := re.Split(line, 24)
		if match[2] == "HomeTeam" {
			continue
		}
		if match[4] == "" {
			continue
		}

		results = append(results,
			rp.Result{Id: "1",
				Date:                getDate(match[1]),
				Round:               1,
				HomeTeamName:        match[2],
				AwayTeamName:        match[3],
				HomeGoals:           stringToUInt8(match[4]),
				AwayGoals:           stringToUInt8(match[5]),
				HomeGoalsAtHalfTime: stringToUInt8(match[7]),
				AwayGoalsAtHalfTime: stringToUInt8(match[8]),
				Cards:               rp.CardInfo{HomeTeamNumberOfYellowCards: stringToUInt8(match[19]), HomeTeamNumberOfRedCards: stringToUInt8(match[21]), AwayTeamNumberOfYellowCards: stringToUInt8(match[20]), AwayTeamNumberOfRedCards: stringToUInt8(match[22])}})
	}
	return results, nil
}

func getDate(dateString string) time.Time {
	formatString := ""
	if 8 == len(dateString) {
		formatString = "02/01/06"
	} else if 10 == len(dateString) {
		formatString = "02/01/2006"
	}

	resultTime, err := time.Parse(formatString, dateString)
	if nil != err {
		fmt.Printf("Could not parse time %v, using now!", dateString)
	}
	return resultTime
}
func stringToUInt8(str string) uint8 {
	value, err := strconv.Atoi(str)
	if nil != err {
		fmt.Printf("Failed parsing string to int. %v", err)
	}
	return uint8(value)
}
