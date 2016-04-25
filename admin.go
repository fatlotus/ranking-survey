package rankingsurvey

import (
	"encoding/json"
	"io"
	"net/http"
)

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

	encoder := json.NewEncoder(w)
	for question := range AllQuestions(r, survey) {
		err := encoder.Encode(&question)
		if err != nil {
			panic(err)
		}
	}
}
