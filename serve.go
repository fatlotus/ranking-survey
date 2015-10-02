package main

import (
	"net/http"
	"appengine"
	"fmt"
	"strconv"
	"html/template"
)

func init() {
	http.HandleFunc("/", survey)
	http.Handle("/static", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/questions.csv", result)
}

var tmpl = template.Must(template.New("survey.html").Funcs(
	template.FuncMap{
		"loop": func(n int) []string {
			return make([]string, n)
		},
	}).ParseFiles("templates/survey.html"))

func survey(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	// Process existing response, if one is given.
	question := r.FormValue("question")
	if question != "" {
		values := make([]int, 0)
		
		for i := 0;; i++ {
			sval := r.FormValue(fmt.Sprintf("response-%d", i))
			val, err := strconv.ParseInt(sval, 10, 64)
			if err != nil {
				break
			}
			values = append(values, int(val))
		}
		if err := AnswerQuestion(c, question, values); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	
	// Display a form to respond to a new question.
	question, q, err := NextQuestion(c, "survey")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Header().Set("Content-type", "text/html")
	
	err = tmpl.Execute(w, struct {
		Question *Question
		ID string
	} { q, question })
	
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
