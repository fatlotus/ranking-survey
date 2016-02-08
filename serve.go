package rankingsurvey

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"html/template"
	"net/http"
	"strconv"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go static/ templates/

func MakeHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", survey)
	mux.Handle("/static/", http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	mux.HandleFunc("/questions.csv", result)
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
	}).Parse(string(MustAsset("templates/survey.html"))))

func survey(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
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
	question, q, err := NextQuestion(r, "survey")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-type", "text/html")

	err = tmpl.Execute(w, struct {
		Question *Question
		ID       string
		Admin    bool
	}{q, question, IsAdmin(r)})

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
