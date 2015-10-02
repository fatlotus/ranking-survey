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
	http.HandleFunc("/result", result)
	http.HandleFunc("/fixture", fixture)
}

var tmpl = template.Must(template.ParseFiles("templates/survey.html"))

func survey(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	// Process existing response.
	response, _ := strconv.ParseInt(r.FormValue("question"), 10, 64)
	if response != 0 {
		values := make([]int, 0)
		
		for i := 0;; i++ {
			sval := r.FormValue(fmt.Sprintf("response-%d", i))
			val, err := strconv.ParseInt(sval, 10, 64)
			if err != nil {
				break
			}
			values = append(values, int(val))
		}
		AnswerQuestion(c, response, values)
	}
	
	// Display a form to respond to a new question.
	qid, question, err := NextQuestion(c, "survey")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Header().Set("Content-type", "text/html")
	
	err = tmpl.Execute(w, struct {
		Question *Question
		ID int64
	} { question, qid })
	if err != nil {
		panic(err)
	}
}
