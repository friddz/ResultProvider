package ice

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	rp "resultprovider"
	"resultprovider/ice/leikuratburdir"
	"resultprovider/ice/motleikir"
	"sort"
	"strconv"
	"time"
	//"strings"
)

const goalId = 11
const goalFromPenaltyId = 17
const ownGoalId = 12

const confirmedResult = "S"

func getInternalSeasonId(externalId string) string {

	if externalId == "2015" {
		return "33503"
	} else if externalId == "2014" {
		return "32265"
	} else if externalId == "2013" {
		return "29465"
	} else if externalId == "2012" {
		return "26463"
	} else if externalId == "2011" {
		return "23425"
	} else if externalId == "2010" {
		return "19847"
	} else if externalId == "2009" {
		return "17548"
	} else if externalId == "2008" {
		return "16788"
	} else if externalId == "2007" {
		return "14843"
	} else if externalId == "2006" {
		return "12486"
	} else if externalId == "2005" {
		return "8585"
	} else {
		panic("Unknown season id")
	}
}

//Usage: results := getAllResults(id)
//Pre: id is of type string and is an unique id of a season in one of the icelandic divisions (see: http://www.ksi.is/mot/XML/)
//Post: results is an array containing all the results for the given season.
func getAllResults(seasonId string) ([]rp.Result, error) {

	results := make([]rp.Result, 0)
	motLeikir := &motleikir.RequestMotLeikir{MotNumer: seasonId, Xmlns: "http://www2.ksi.is/vefthjonustur/mot/"}
	envelope := motleikir.RequestEnvelope{Xmlns: "http://schemas.xmlsoap.org/soap/envelope/"}
	soapBody := new(motleikir.RequestSoapBody)
	soapBody.Content = motLeikir
	envelope.Soap = soapBody
	buf, _ := xml.Marshal(envelope)
	body := bytes.NewBuffer([]byte(buf))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://www2.ksi.is/vefthjonustur/mot.asmx?op=LeikurAtburdir", body)
	req.Header.Add("Content-Type", "text/xml")
	r, _ := client.Do(req)
	responseBody, _ := ioutil.ReadAll(r.Body)
	v := motleikir.Response{}
	err := xml.Unmarshal([]byte(responseBody), &v)
	if nil != err {
		return results, err
	}

	for _, r := range v.Body.MotLeikirResponse.MotLeikirSvar.ArrayMotLeikir.MotLeikur {
		if confirmedResult == r.SkyrslaStada {
			results = append(results,
				rp.Result{Id: r.LeikurNumer,
					Date:                parseToTime(r.LeikDagur),
					Round:               parseToUint8(r.UmferdNumer),
					HomeTeamName:        r.FelagHeimaNafn,
					AwayTeamName:        r.FelagUtiNafn,
					HomeGoals:           parseToUint8(r.UrslitHeima),
					AwayGoals:           parseToUint8(r.UrslitUti),
					HomeGoalsAtHalfTime: parseToUint8(r.StadaFyrriHalfleikHeima),
					AwayGoalsAtHalfTime: parseToUint8(r.StadaFyrriHalfleikUti)})
		}
	}

	return results, nil
}

//Usage: 	goals,err := getResultGoals(id)
//Pre:		id is of type string and is an unique id of a result in one of the icelancdic divisions (see: http://www.ksi.is/mot/XML/)
//Post:		goals contains all the goals scored in the game with the given id. If err is non nil, an error occured.
func getResultGoals(resultId string) ([]rp.Goal, error) {
	goals := make([]rp.Goal, 0)
	leikurAtburdir := &leikuratburdir.RequestLeikurAtburdir{LeikurNumer: resultId, Xmlns: "http://www2.ksi.is/vefthjonustur/mot/"}
	envelope := leikuratburdir.RequestEnvelope{Xmlns: "http://schemas.xmlsoap.org/soap/envelope/"}
	soapBody := &leikuratburdir.RequestSoapBody{Xmlns: "http://www.w3.org/2001/XMLSchema-instance"}
	soapBody.Content = leikurAtburdir
	envelope.Soap = soapBody
	buf, _ := xml.Marshal(envelope)
	body := bytes.NewBuffer([]byte(buf))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://www2.ksi.is/vefthjonustur/mot.asmx?", body)
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	r, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	responseBody, _ := ioutil.ReadAll(r.Body)
	v := leikuratburdir.Response{}

	err = xml.Unmarshal([]byte(responseBody), &v)
	if nil != err {
		return goals, err
	}

	for _, r := range v.Body.LeikurAtburdirResponse.LeikurAtburdirSvar.ArrayLeikurAtburdir.LeikurAtburdir {
		if goalId == r.AtburdurNumer || goalFromPenaltyId == r.AtburdurNumer || ownGoalId == r.AtburdurNumer {
			goals = append(goals, rp.Goal{GoalScorerName: r.LeikmadurNafn, Minute: r.AtburdurMinuta, TeamName: r.FelagNafn, Type: getGoalType(r.AtburdurNumer)})
		}
	}

	return goals, nil
}

func getGoalType(eventId uint8) rp.GoalType {
	goalType := rp.RegularGoal
	if goalFromPenaltyId == eventId {
		goalType = rp.GoalFromPenalty
	} else if ownGoalId == eventId {
		goalType = rp.OwnGoal
	}

	return goalType
}
func parseToTime(s string) time.Time {
	value, _ := time.Parse("2006-01-02T15:04:05", s)
	return value
}

func parseToUint8(s string) uint8 {
	value, _ := strconv.ParseUint(s, 10, 8)
	return uint8(value)
}

type WorkRequest struct {
	id    string
	index int
}

type WorkResponse struct {
	index int
	goals []rp.Goal
}

func Worker(in <-chan *WorkRequest, out chan<- *WorkResponse) {

	for w := range in {
		goals, _ := getResultGoals(w.id)
		out <- &WorkResponse{index: w.index, goals: goals}
	}
}

func GetAllResultsForMultipleSeasons(ids []string) ([]rp.Season, error) {
	fileName := getFileName(ids)
	if existsInCache(fileName) {
		return getResultsFromCache(fileName)
	} else {
		seasons, err := getAllResultsForMultipleSeasons(ids)
		if err != nil {
			return nil, err
		}
		err = writeToCache(seasons, ids)
		return seasons, err
	}
}

func getAllResultsForMultipleSeasons(ids []string) ([]rp.Season, error) {
	seasons := make([]rp.Season, 0)
	results := make([]rp.Result, 0)
	for _, id := range ids {
		results, _ = GetAllResults(getInternalSeasonId(id))
		seasons = append(seasons, rp.Season{Id: id, Results: results})
	}
	return seasons, nil

}

func getResultsFromCache(fileName string) ([]rp.Season, error) {
	seasons := []rp.Season{}
	data, err := ioutil.ReadFile(fileName)
	if nil != err {
		return seasons, err
	}
	err = json.Unmarshal(data, &seasons)
	return seasons, err
}

func existsInCache(fileName string) bool {
	_, err := os.Stat(fileName)
	return nil == err
}
func writeToCache(seasons []rp.Season, ids []string) error {
	blob, err := json.Marshal(seasons)
	if err != nil {
		return err
	}
	fileName := getFileName(ids)
	err = ioutil.WriteFile(fileName, blob, 0644)
	return err
}

func getFileName(ids []string) string {
	return "C:/tmp/" + generateHash(ids) + ".json"
}

func generateHash(ids []string) string {
	hash := 23
	sort.Strings(ids)
	for _, i := range ids {
		val, err := strconv.Atoi(i)
		if nil != err {
			panic(err)
		}
		hash = hash*31 + val
	}
	//return strings.Replace(strconv.Itoa(hash), "-", "_", 1)
	return strconv.Itoa(hash)
}

func GetAllResults(id string) ([]rp.Result, error) {
	results, err := getAllResults(id)
	if nil != err {
		return nil, err
	}

	in := make(chan *WorkRequest, len(results))
	out := make(chan *WorkResponse, len(results))

	for i := 0; i < 40; i++ {
		go Worker(in, out)
	}

	for i, r := range results {
		in <- &WorkRequest{id: r.Id, index: i}
	}

	close(in)

	for i := 0; i < len(results); i++ {
		resp := <-out
		results[resp.index].Goals = resp.goals
	}

	return results, nil
}
