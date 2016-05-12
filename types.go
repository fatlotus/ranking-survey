package rankingsurvey

import (
	"time"
)

type SurveyID string
type Question struct {
	Survey    SurveyID  `datastore:"survey" json:"survey"`
	Choices   []string  `datastore:"choices" json:"choices"`
	Precision int       `datastore:"precision" json:"precision"`
	Exclusive bool      `datastore:"exclusive" json:"exclusive"`
	Seen      time.Time `datastore:"seen" json:"seenTime"`
	Responded time.Time `datastore:"seenTime" json:"respondedTime"`
	Response  []int     `datastore:"response" json:"response"`
	Email     string    `datastore:"email" json:"email"`
}

type Survey struct {
	Survey      SurveyID
	Seen, Total int
}

type ByID []Survey

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].Survey < a[j].Survey }
