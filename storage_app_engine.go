// +build appengine

package rankingsurvey

import (
	"appengine"
	"appengine/datastore"
	"github.com/fatlotus/collaborativepermute"
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"encoding/gob"
	"bytes"
)

func init() {
	http.Handle("/", MakeHandler())
}

type State struct {
	Data []byte
	Options []string
}

func (s *State) Engine() (e *collaborativepermute.Engine) {
	buf := bytes.NewBuffer(s.Data)
	err := gob.NewDecoder(buf).Decode(&e)
	if err != nil {
		panic(err)
	}
	return
}

func (s *State) SetEngine(e *collaborativepermute.Engine) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(e); err != nil {
		panic(err)
	}
	s.Data = buf.Bytes()
}

func NextQuestion(r *http.Request, sv SurveyID) (string, *Question, error){
	c := appengine.NewContext(r)
	k := datastore.NewKey(c, "State", "state", 0, nil)
	var s State

	if err := datastore.Get(c, k, &s); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return "", nil, nil
		} else {
			return "", nil, err
		}
	}

	q := s.Engine().Generate(0)
	labels := make([]string, len(q.Choices))
	for i := range labels {
		labels[i] = s.Options[q.Choices[i]]
	}
	
	return fmt.Sprintf("%d-%d", q.Choices[0], q.Choices[1]), &Question{
		Survey: sv,
		Choices: labels,
		Precision: 0,
	}, nil
}

func AnswerQuestion(r *http.Request, key string, response []int) error {
	c := appengine.NewContext(r)
	k := datastore.NewKey(c, "State", "state", 0, nil)

	frags := strings.Split(key, "-")
	if len(frags) != 2 {
		return fmt.Errorf("invalid key %#v", key)
	}

	a, err := strconv.ParseInt(frags[0], 32, 10)
	if err != nil {
		return err
	}

	b, err := strconv.ParseInt(frags[1], 32, 10)
	if err != nil {
		return err
	}

	if response[0] == 1 {
		a, b = b, a
	}

	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		var s State
		if err := datastore.Get(c, k, &s); err != nil {
			return err
		}

		e := s.Engine()
		e.Respond(collaborativepermute.Query{
			User: 0,
			Choices: []int{int(a), int(b)},
		})
		s.SetEngine(e)

		_, err := datastore.Put(c, k, &s)
		return err
	}, nil)

}

func AllQuestions(r *http.Request, survey SurveyID) (<-chan Question) {
	qs := make(chan Question)
	close(qs)
	return qs
}

func AddQuestions(r *http.Request, questions []Question) error {
	c := appengine.NewContext(r)
	k := datastore.NewKey(c, "State", "state", 0, nil)

	items := make(map[string]int)

	for _, question := range questions {
		for _, choice := range question.Choices {
			items[choice] = 0
		}
	}

	s := new(State)
	s.SetEngine(collaborativepermute.NewEngine(len(items), len(items)))
	s.Options = make([]string, 0)

	for option := range items {
		s.Options = append(s.Options, option)
	}

	_, err := datastore.Put(c, k, s)
	return err
}