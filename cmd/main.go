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
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("./res"))))
	http.HandleFunc("/tstp", renderTimeSpace)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func renderTimeSpace(w http.ResponseWriter, r *http.Request) {
	var from = queryIntList(r.URL.Query()["from"])
	var to = queryIntList(r.URL.Query()["to"])
	var vias = queryIntList(r.URL.Query()["vias"])

	var form = r.URL.Query()["form"]
	var datetime = r.URL.Query().Get("datetime")

	if len(from) > 0 && len(to) > 0 {
		stations, lines := internal.ObtainData(from[0], to[0], vias, datetime)
		if len(vias) > 0 && len(form) == 0 {
			log.Print("Request:", r.URL.RawQuery)
			w.Header().Set("Content-Type", "image/svg+xml")
			render.TimeSpace(stations, lines, w, r.URL.RawQuery)
			return
		}
		log.Print(to[0], stations)
		render.Vias(stations, from[0], to[0], datetime, w)
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
