package rankingsurvey

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go static/ templates/

func MakeHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", survey)
	mux.Handle("/static/", staticHandler())
	return mux
}

var indexTmpl = template.Must(template.New("index.html").Funcs(
	template.FuncMap{
		"prog": func(a, b int) int {
			return 100 * a / b
		},
	}).Parse(string(MustAsset("templates/index.html"))))

var surveyTmpl = template.Must(template.New("survey.html").Funcs(
	template.FuncMap{
		"loop": func(n int) []string {
			return make([]string, n)
		},
		"div": func(a, b int) int {
			return a / b
		},
		"add1": func(n int) int {
			return n + 1
		},
		"asHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"pairwise": func(q Question) bool {
			return q.Exclusive && len(q.Choices) == 2 && q.Precision == 2
		},
		"contpairwise": func(q Question) bool {
			return q.Exclusive && len(q.Choices) == 2 && q.Precision > 10
		},
		"justone": func(q Question) bool {
			return len(q.Choices) == 1
		},
		"binary": func(q Question) bool {
			return len(q.Choices) == 1 && q.Precision == 2
		},
		"align": func(index, total int) string {
			if total == 1 {
				return "center"
			} else if index == 0 {
				return "left"
			} else if index == total-1 {
				return "right"
			} else {
				return "center"
			}
		},
	}).Parse(string(MustAsset("templates/survey.html"))))

func homePage(w http.ResponseWriter, r *http.Request) {
	surveys := []Survey(nil)
	if IsAdmin(r) {
		var err error
		surveys, err = AllSurveys(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if r.Method == "POST" {
		options := make([]string, 0)
		for i := 0; i < 100; i++ {
			// this is stupid, but it is faster than scanning the
			// database each time
			options = append(options, fmt.Sprintf("experiment/%d/2", i))
			options = append(options, fmt.Sprintf("experiment/%d/5", i))
			options = append(options, fmt.Sprintf("experiment/%d/100", i))
			options = append(options, fmt.Sprintf("experiment/%d/cmp", i))
		}

		for try := len(options); try >= 0; try -= 10 {
			selected := rand.Intn(len(options))

			free, err := IsFree(r, SurveyID(options[selected]))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if free {
				http.Redirect(w, r, "/"+options[selected]+
					"?email="+url.QueryEscape(r.FormValue("email")), 303)
				return
			}
		}

		http.Error(w, "no surveys left?", 500)
		return
	}

	w.Header().Set("Content-type", "text/html")
	err := indexTmpl.Execute(w, struct {
		Surveys []Survey
		Admin   bool
	}{surveys, IsAdmin(r)})

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func survey(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".json") {
		result(w, r)
		return
	} else if r.URL.Path == "/" {
		homePage(w, r)
		return
	}

	survey := SurveyID(r.URL.Path[1:])
	if survey == SurveyID("") {
		survey = SurveyID("survey")
	}

	// Process existing response, if one is given.
	question := r.FormValue("question")
	email := r.FormValue("email")
	if question != "" {
		values := make([]int, 0)

		for i := 0; ; i++ {
			sval := r.FormValue(fmt.Sprintf("response-%d", i))
			val, err := strconv.ParseInt(sval, 10, 64)
			if err != nil {
				break
			}
			values = append(values, int(val))
		}
		if err := AnswerQuestion(r, question, email, values); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	// Display a form to respond to a new question.
	question, q, err, complete, total := NextQuestion(r, survey)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-type", "text/html")

	err = surveyTmpl.Execute(w, struct {
		Survey   SurveyID
		Question *Question
		ID       string
		Admin    bool
		Number   int
		Total    int
	}{survey, q, question, IsAdmin(r), total - complete + 1, total})

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
