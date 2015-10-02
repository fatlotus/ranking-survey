package main

import (
	"net/http"
	"appengine"
	"fmt"
	"strings"
	"encoding/csv"
	"strconv"
)

func result(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	if r.Method != "HEAD" && r.Method != "GET" {
		upload, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		reader := csv.NewReader(upload)
		rows, err := reader.ReadAll()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		questions := make([]Question, len(rows) - 1)
		for idx, row := range rows[1:] {
			precision, _ := strconv.ParseInt(row[1], 10, 64)
			questions[idx] = Question{
				Survey: "survey",
				Choices: strings.Split(row[0], ";"),
				Precision: int(precision),
			}
		}
		
		AddQuestions(c, questions)
		return
	}
	
	w.Header().Add("Content-type", "text/plain")

	writer := csv.NewWriter(w)
	writer.Write([]string {"choices", "range"})

	for question := range AllQuestions(c, "survey") {
		
		response := make([]string, len(question.Choices))
		for i, choice := range question.Choices {
			if len(question.Response) == len(question.Choices) {
				response[i] = fmt.Sprintf("%s:%d", choice, question.Response[i])
			} else if question.Seen {
				response[i] = fmt.Sprintf("%s:?", choice)
			} else {
				response[i] = choice
			}
		}
		
		writer.Write([]string {
			strings.Join(response, ";"),
			strconv.FormatInt(int64(question.Precision), 10),
		})
	}
	
	writer.Flush()
}

func fixture(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := AddQuestions(c, []Question{
		// Question{"sa", []string{"Cats"}, 5, false, []int{}},
		// Question{"sa", []string{"Dogs"}, 5, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		Question{"sa", []string{"Cats", "Dogs"}, 0, false, []int{}},
		// Question{"sa", []string{"Cats", "Dogs"}, 2, false, []int{}},
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	http.Redirect(w, r, "/", 303)
}