package main

import (
	"net/http"

	"github.com/arishlabroo/traintomo/internal"
)

func main() {
	http.ListenAndServe(":80", internal.NewServer(internal.NewInMemoryDB()))
}
