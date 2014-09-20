package resultprovider

import "time"

type Result struct {
	Id                  string
	Date                time.Time
	Round               uint8
	HomeTeamName        string
	AwayTeamName        string
	HomeGoals           uint8
	AwayGoals           uint8
	HomeGoalsAtHalfTime uint8
	AwayGoalsAtHalfTime uint8
	Goals               []Goal
	Cards               CardInfo
}

type Goal struct {
	GoalScorerName string
	Minute         uint8
	TeamName       string
	Type           GoalType
}

type Season struct {
	Id      string
	Results []Result
}

type CardInfo struct {
	HomeTeamNumberOfYellowCards uint8
	HomeTeamNumberOfRedCards    uint8
	AwayTeamNumberOfYellowCards uint8
	AwayTeamNumberOfRedCards    uint8
}

type GoalType uint8

const (
	RegularGoal     GoalType = 1
	OwnGoal                  = 2
	GoalFromPenalty          = 3
)

type Fixture struct {
	Date time.Time
	HomeTeamName string
	AwayTeamName string
}