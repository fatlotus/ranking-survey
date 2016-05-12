// +build !appengine

package rankingsurvey

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type sharedState struct {
	Questions []Question
	sync.Mutex
}

func IsAdmin(r *http.Request) bool {
	return true
}

var globalState sharedState

func AllSurveys(r *http.Request) ([]Survey, error) {
	globalState.Lock()
	defer globalState.Unlock()

	seen := 0
	total := 0
	for _, q := range globalState.Questions {
		if !q.Seen.IsZero() {
			seen += 1
		}
		total += 1
	}

	return Survey{
		Survey: "survey",
		Seen:   seen,
		Total:  total,
	}
}

func NextQuestion(r *http.Request, s SurveyID) (string, *Question, error, int, int) {
	globalState.Lock()
	defer globalState.Unlock()

	length := len(globalState.Questions)

	for index, question := range globalState.Questions {
		if question.Seen.IsZero() {
			globalState.Questions[index].Seen = time.Now()
			return strconv.FormatInt(int64(index), 10), &question, nil, index, length
		}
	}

	return "", nil, nil, length, length
}

func AnswerQuestion(r *http.Request, key, email string, response []int) error {
	globalState.Lock()
	defer globalState.Unlock()

	index, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		panic(err)
		return err
	}

	if index >= int64(len(globalState.Questions)) || index < 0 {
		panic(err)
		return errors.New("out of range")
	}

	globalState.Questions[int(index)].Email = email
	globalState.Questions[int(index)].Response = response
	return nil
}

func AllQuestions(r *http.Request, survey SurveyID) <-chan Question {
	channel := make(chan Question)
	go func() {
		globalState.Lock()
		defer globalState.Unlock()

		for _, q := range globalState.Questions {
			channel <- q
		}
		close(channel)
	}()
	return channel
}

func AddQuestions(r *http.Request, questions []Question) error {
	globalState.Lock()
	defer globalState.Unlock()

	globalState.Questions = questions
	return nil
}
