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
}
