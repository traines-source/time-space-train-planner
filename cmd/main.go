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
	main := http.NewServeMux()
	main.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("./res"))))
	main.HandleFunc("/tstp", renderTimeSpace)

	api := http.NewServeMux()
	main.Handle("/api/v1/", http.StripPrefix("/api/v1", api))
	api.HandleFunc("/vias", apiVias)
	api.HandleFunc("/timespace", apiTimespace)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), main))
}

func renderTimeSpace(w http.ResponseWriter, r *http.Request) {
	var from = queryIntList(r.URL.Query()["from"])
	var to = queryIntList(r.URL.Query()["to"])
	var vias = queryIntList(r.URL.Query()["vias"])

	var form = r.URL.Query()["form"]
	var datetime = r.URL.Query().Get("datetime")

	if len(from) > 0 && len(to) > 0 {
		stations, lines, err := internal.ObtainData(from[0], to[0], vias, datetime)
		if err == nil && len(vias) > 0 && len(form) == 0 {
			log.Print("Request:", r.URL.RawQuery)
			w.Header().Set("Content-Type", "image/svg+xml")
			render.TimeSpace(stations, lines, w, r.URL.RawQuery)
			return
		}
		log.Print(to[0], stations)
		if err != nil {
			//v, ok := interface{}(err).(internal.ErrorWithCode)
			w.WriteHeader(err.ErrorCode())
		}
		render.Vias(stations, from[0], to[0], datetime, w, err)
		return
	}
	render.Index(w)
}

func apiVias(w http.ResponseWriter, r *http.Request) {
	var from = queryIntList(r.URL.Query()["from"])
	var to = queryIntList(r.URL.Query()["to"])
	var vias = queryIntList(r.URL.Query()["vias"])

	var datetime = r.URL.Query().Get("datetime")

	w.Header().Set("Content-Type", "application/json")
	if len(from) > 0 && len(to) > 0 {
		stations, _, err := internal.ObtainData(from[0], to[0], vias, datetime)
		if err != nil {
			w.WriteHeader(err.ErrorCode())
		}
		render.ViasApi(stations, from[0], to[0], datetime, w, err)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func apiTimespace(w http.ResponseWriter, r *http.Request) {
	var from = queryIntList(r.URL.Query()["from"])
	var to = queryIntList(r.URL.Query()["to"])
	var vias = queryIntList(r.URL.Query()["vias"])

	var datetime = r.URL.Query().Get("datetime")

	w.Header().Set("Content-Type", "application/json")
	if len(from) > 0 && len(to) > 0 && len(vias) > 0 {
		stations, lines, err := internal.ObtainData(from[0], to[0], vias, datetime)
		if err != nil {
			w.WriteHeader(err.ErrorCode())
		}
		render.TimeSpaceApi(stations, lines, w, r.URL.RawQuery)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func queryIntList(params []string) []int {
	var ints = []int{}
	for _, i := range params {
		if i == "" {
			continue
		}
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
