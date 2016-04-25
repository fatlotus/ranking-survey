package rankingsurvey

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go static/ templates/

func MakeHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", survey)
	mux.Handle("/static/", http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	return mux
}

var tmpl = template.Must(template.New("survey").Funcs(
	template.FuncMap{
		"loop": func(n int) []string {
			return make([]string, n)
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
		"binary": func(q Question) bool {
			return len(q.Choices) == 1 && q.Precision == 2
		},
	}).Parse(string(MustAsset("templates/survey.html"))))

func survey(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".json") {
		result(w, r)
		return
	}

	survey := SurveyID(r.URL.Path[1:])
	if survey == SurveyID("") {
		survey = SurveyID("survey")
	}

	// Process existing response, if one is given.
	question := r.FormValue("question")
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
		if err := AnswerQuestion(r, question, values); err != nil {
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

	err = tmpl.Execute(w, struct {
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
