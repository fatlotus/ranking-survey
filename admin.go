package rankingsurvey

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var nonalnum = regexp.MustCompile("[^a-zA-Z0-9]+")

type byTime []Question

func (b byTime) Less(i, j int) bool {
	if b[i].Responded.IsZero() {
		return false
	}
	if b[j].Responded.IsZero() {
		return true
	}
	return b[i].Responded.Before(b[i].Responded)
}

func (b byTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byTime) Len() int      { return len(b) }

func result(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Error(w, "forbidden", 403)
		return
	}

	survey := SurveyID(r.URL.Path[1 : len(r.URL.Path)-5])

	// Allow file uploads of new questions.
	if r.Method != "HEAD" && r.Method != "GET" {
		upload, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		decoder := json.NewDecoder(upload)
		questions := make([]Question, 0)
		for {
			var question Question
			err := decoder.Decode(&question)
			if err == io.EOF {
				break
			} else if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			questions = append(questions, question)
		}

		if err := AddQuestions(r, questions); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/", 303)
		return
	}

	// Otherwise, serve the current question set.
	w.Header().Add("Content-type", "text/plain")

	if r.FormValue("download") != "" {
		filename := nonalnum.ReplaceAllString(string(survey), "-")
		w.Header().Add("Content-disposition",
			fmt.Sprintf("attachment; filename=\"%s.json\"", filename))
	}

	questions, err := AllQuestions(r, survey)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	sort.Sort(byTime(questions))

	for i, question := range questions {
		ids := []int(nil)

		for j, choice := range question.Choices {
			if len(choice) > 4 && choice[:4] == "<!--" {
				if len(ids) == 0 {
					ids = make([]int, len(question.Choices))
				}
				fragments := strings.Split(choice[4:], "-->")
				id, _ := strconv.ParseInt(fragments[0], 10, 64)
				ids[j] = int(id)
			}
		}

		if len(ids) > 0 {
			questions[i].ChoiceIDs = ids
		}
	}

	encoder := json.NewEncoder(w)
	for _, question := range questions {
		err := encoder.Encode(&question)
		if err != nil {
			panic(err)
		}
	}
}
