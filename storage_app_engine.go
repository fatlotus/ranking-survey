// +build appengine

package rankingsurvey

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"net/http"
	"time"
)

func init() {
	http.Handle("/", MakeHandler())
}

func IsAdmin(r *http.Request) bool {
	return user.IsAdmin(appengine.NewContext(r))
}

func NextQuestion(r *http.Request, s SurveyID) (i string, q *Question, e error) {
	c := appengine.NewContext(r)

	// Fetch the next key outside of a transaction.
	query := datastore.NewQuery("Question").
		Filter("survey =", s).
		Filter("seen <=", time.Unix(0, 0)).
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
	err = datastore.RunInTransaction(c, func(c context.Context) error {
		if err := datastore.Get(c, keys[0], &question); err != nil {
			return err
		}
		question.Seen = time.Now()
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

func AnswerQuestion(r *http.Request, key string, response []int) error {
	c := appengine.NewContext(r)

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

	question.Responded = time.Now()
	question.Response = response
	_, err := datastore.Put(c, dbkey, &question)
	return err
}

func AllQuestions(r *http.Request, survey SurveyID) <-chan Question {
	c := appengine.NewContext(r)

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

func AddQuestions(r *http.Request, questions []Question) error {
	c := appengine.NewContext(r)

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
