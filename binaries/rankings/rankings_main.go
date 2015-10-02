// +build !appengine

package main

import (
	"github.com/fatlotus/rankings"
	"net/http"
)

func main() {
	panic(http.ListenAndServe(":8080", rankings.MakeHandler()))
}