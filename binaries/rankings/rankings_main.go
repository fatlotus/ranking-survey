// +build !appengine

package main

import (
	"github.com/fatlotus/rankingsurvey"
	"net/http"
)

func main() {
	panic(http.ListenAndServe(":8080", rankingsurvey.MakeHandler()))
}