// +build appengine

package rankingsurvey

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

func init() {
	http.Handle("/", MakeHandler())
}

func IsAdmin(r *http.Request) bool {
	return user.IsAdmin(appengine.NewContext(r))
}

func AllSurveys(r *http.Request) ([]Survey, error) {
	c := appengine.NewContext(r)

	questions := make([]Question, 0)
	q := datastore.NewQuery("Question").Project("survey", "seen")
	_, err := q.GetAll(c, &questions)
	if err != nil {
		return nil, err
	}

	seen := map[SurveyID]int{}
	total := map[SurveyID]int{}

	for _, q := range questions {
		if !q.Seen.IsZero() {
			seen[q.Survey] += 1
		}
		total[q.Survey] += 1
	}

	surveys := []Survey{}
	for id, total := range total {
		surveys = append(surveys, Survey{
			Survey: id,
			Seen:   seen[id],
			Total:  total,
		})
	}

	sort.Sort(ByID(surveys))

	return surveys, nil
}

func NextQuestion(r *http.Request, s SurveyID) (i string, q *Question, e error, f int, t int) {
	c := appengine.NewContext(r)

	// Fetch the next key outside of a transaction.
	query := datastore.NewQuery("Question").
		Filter("survey =", s).
		Filter("seen <=", time.Unix(0, 0)).
		KeysOnly()

	keys, err := query.GetAll(c, nil)
	if err != nil {
		return
	}

	f = len(keys)

	t, err = datastore.NewQuery("Question").Filter("survey =", s).Count(c)
	if err != nil {
		return
	}

	if len(keys) == 0 {
		return
	}
	// Select a random key from the set.
	idx := rand.Intn(len(keys))
	key := keys[idx]

	// Retreive and update the given question.
	var question Question
	err = datastore.RunInTransaction(c, func(c context.Context) error {
		if err := datastore.Get(c, key, &question); err != nil {
			return err
		}
		question.Seen = time.Now()
		_, err := datastore.Put(c, key, &question)
		return err
	}, nil)
	if err != nil {
		return
	}

	i = key.StringID()
	q = &question
	return
}

func AnswerQuestion(r *http.Request, key, email string, response []int) error {
	c := appengine.NewContext(r)

	var question Question
	dbkey := datastore.NewKey(c, "Question", key, 0, nil)
	if err := datastore.Get(c, dbkey, &question); err != nil {
		return err
	}

	if len(question.Response) != 0 {
		// panic("wrong length")
		return nil
	}

	if len(question.Choices) != len(response) {
		// panic("mismatch")
		return nil
	}

	question.Responded = time.Now()
	question.Response = response
	question.Email = email
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

	for j := 0; j < len(keys); j += 500 {
		k := j + 500
		if k >= len(keys) {
			k = len(keys) - 1
		}
		if _, err := datastore.PutMulti(c, keys[j:k], questions[j:k]); err != nil {
			return err
		}
	}

	// Remove old entities from the datastore.
	for key, count := range counts {
		name := fmt.Sprintf("%s-%010d", key, count)
		first := datastore.NewKey(c, "Question", name, 0, nil)
		name = fmt.Sprintf("%s-%010d", key, 9999999999)
		last := datastore.NewKey(c, "Question", name, 0, nil)

		q := datastore.NewQuery("Question").Filter("__key__ >=", first)
		q = q.Filter("__key__ <=", last).KeysOnly()
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
