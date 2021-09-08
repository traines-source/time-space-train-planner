package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"traines.eu/time-space-train-planner/internal"
	"traines.eu/time-space-train-planner/render"
)

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/tstp", renderTimeSpace)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func renderTimeSpace(w http.ResponseWriter, r *http.Request) {
	var from = queryIntList(r.URL.Query()["from"])
	var to = queryIntList(r.URL.Query()["to"])
	var vias = queryIntList(r.URL.Query()["vias"])

	var form = r.URL.Query()["form"]

	if len(from) > 0 && len(to) > 0 {
		stations, lines := internal.ObtainData(from[0], to[0], vias)
		if len(vias) > 0 && len(form) == 0 {
			w.Header().Set("Content-Type", "image/svg+xml")
			render.TimeSpace(stations, lines, w)
			return
		}
		log.Print(to[0], stations)
		render.Vias(stations, from[0], to[0], w)
		return
	}
	render.Index(w)
}

func queryIntList(params []string) []int {
	var ints = []int{}
	for _, i := range params {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		if j == 0 {
			continue
		}
		ints = append(ints, j)
	}
	return ints
}
