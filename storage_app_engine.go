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

func NextQuestion(c appengine.Context, s SurveyID) (i string, q *Question, e error) {
	// Fetch the next key outside of a transaction.
	query := datastore.NewQuery("Question").
	            Filter("survey =", s).
	            Filter("seen =", false).
	            Limit(1).
	            KeysOnly()

	keys, err := query.GetAll(c, nil)
	if err != nil {
		return
	}
	if len(keys) == 0 {
		return
	}

	// Retreive and update the given question.
	var question Question
	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		if err := datastore.Get(c, keys[0], &question); err != nil {
			return err
		}
		question.Seen = true
		_, err := datastore.Put(c, keys[0], &question)
		return err
	}, nil)
	if err != nil {
		return
	}

	i = keys[0].StringID()
	q = &question
	return
}

func AnswerQuestion(c appengine.Context, key string, response []int) error {
	var question Question
	dbkey := datastore.NewKey(c, "Question", key, 0, nil)
	if err := datastore.Get(c, dbkey, &question); err != nil {
		return err
	}
	
	if len(question.Response) != 0 {
		return nil
	}
	
	if len(question.Choices) != len(response) {
		return nil
	}
	
	question.Response = response
	_, err := datastore.Put(c, dbkey, &question)
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
	// Put new versions of added queries.
	keys := make([]*datastore.Key, len(questions))
	counts := make(map[SurveyID]int, 0)

	for i := range keys {
		survey := questions[i].Survey
		name := fmt.Sprintf("%s-%010d", survey, counts[survey])
		keys[i] = datastore.NewKey(c, "Question", name, 0, nil)
		counts[survey] += 1
	}

	if _, err := datastore.PutMulti(c, keys, questions); err != nil {
		return err
	}

	// Remove old entities from the datastore.
	for key, count := range counts {
		name := fmt.Sprintf("%s-%010d", key, count)
		key := datastore.NewKey(c, "Question", name, 0, nil)
		q := datastore.NewQuery("Question").Filter("__key__ >=", key).KeysOnly()
		keys, err := q.GetAll(c, nil)
		if err != nil {
			return err
		}
		if err := datastore.DeleteMulti(c, keys); err != nil {
			return err
		}
	}

	return nil
}