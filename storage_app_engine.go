package main

import (
	"appengine"
	"appengine/datastore"
	"fmt"
)

type SurveyID string
type Question struct {
	Survey SurveyID `datastore:"survey"`
	Choices []string `datastore:"choices"`
	Precision int `datastore:"precision"`
	Seen bool `datastore:"seen"`
	Response []int `datastore:"response"`
}

func NextQuestion(c appengine.Context, survey SurveyID) (int64, *Question, error) {
	q := datastore.NewQuery("Question").
	        Filter("survey =", survey).
	        Filter("seen =", false).
	        Limit(1).
	        KeysOnly()
	
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("GetAll failed: %s", err)
	}
	if len(keys) == 0 {
		return 0, nil, nil
	}
	
	var question Question
	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		if err := datastore.Get(c, keys[0], &question); err != nil {
			return fmt.Errorf("Get failed: %s", err)
		}
		question.Seen = true
		_, err := datastore.Put(c, keys[0], &question)
		return err
	}, nil)
	if err != nil {
		return 0, nil, err
	}
	
	return keys[0].IntID(), &question, nil
}

func AnswerQuestion(c appengine.Context, id int64, response []int) error {
	var question Question
	key := datastore.NewKey(c, "Question", "", id, nil)
	if err := datastore.Get(c, key, &question); err != nil {
		return err
	}
	
	if len(question.Response) != 0 {
		return nil
	}
	
	if len(question.Choices) != len(response) {
		return nil
	}
	
	question.Response = response
	_, err := datastore.Put(c, key, &question)
	return err
}

func AllQuestions(c appengine.Context, survey SurveyID) (<-chan Question) {
	iterator := datastore.NewQuery("Question").Filter("survey =", survey).Run(c)
	result := make(chan Question)
	
	go func() {
		for {
			var question Question
			if key, err := iterator.Next(&question); err != nil || key == nil {
				break
			}
			result <- question
		}
		close(result)
	}()
	
	return result
}

func AddQuestions(c appengine.Context, questions []Question) error {
	keys := make([]*datastore.Key, len(questions))
	for i := range keys {
		keys[i] = datastore.NewIncompleteKey(c, "Question",nil)
	}
	_, err := datastore.PutMulti(c, keys, questions)
	return err
}