package rankings

type SurveyID string
type Question struct {
	Survey SurveyID `datastore:"survey"`
	Choices []string `datastore:"choices"`
	Precision int `datastore:"precision"`
	Seen bool `datastore:"seen"`
	Response []int `datastore:"response"`
}
