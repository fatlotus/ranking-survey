// +build !appengine

package rankings

import (
	"errors"
	"sync"
	"net/http"
	"strconv"
)

type sharedState struct {
	Questions []Question
	sync.Mutex
}

var globalState sharedState

func NextQuestion(r *http.Request, s SurveyID) (string, *Question, error) {
	globalState.Lock()
	defer globalState.Unlock()
	
	for index, question := range globalState.Questions {
		if !question.Seen {
			globalState.Questions[index].Seen = true
			return strconv.FormatInt(int64(index), 10), &question, nil 
		}
	}
	
	return "", nil, nil
}

func AnswerQuestion(r *http.Request, key string, response []int) error {
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
	
	globalState.Questions[int(index)].Response = response
	return nil
}

func AllQuestions(r *http.Request, survey SurveyID) (<-chan Question) {
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
