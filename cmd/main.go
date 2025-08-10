package main

import (
	"log"
	"net/http"
	"os"

	"traines.eu/time-space-train-planner/internal"
	"traines.eu/time-space-train-planner/render"
)

func main() {
	log.Print("Starting...")
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
	var from = queryStationList(r.URL.Query()["from"])
	var to = queryStationList(r.URL.Query()["to"])
	var vias = queryStationList(r.URL.Query()["vias"])

	var form = r.URL.Query()["form"]
	var datetime = r.URL.Query().Get("datetime")
	var regionly = r.URL.Query().Get("regionly") == "true"

	if len(from) > 0 && len(to) > 0 {
		if len(vias) > 0 && len(form) == 0 {
			stations, lines, err := internal.ObtainData("de_db", from[0], to[0], vias, datetime, regionly)
			if err == nil {
				log.Print("Request:", r.URL.RawQuery)
				w.Header().Set("Content-Type", "image/svg+xml")
				render.TimeSpace(stations, lines, w, r.URL.RawQuery)
				return
			} else {
				w.WriteHeader(err.ErrorCode())
			}
		}
		stations, err := internal.ObtainVias("de_db", from[0], to[0], vias, datetime, regionly)
		if err != nil {
			w.WriteHeader(err.ErrorCode())
		}
		render.Vias(stations, from[0], to[0], datetime, w, err)
		return
	}
	render.Index(w)
}

func apiVias(w http.ResponseWriter, r *http.Request) {
	var system = r.URL.Query().Get("system")
	var from = queryStationList(r.URL.Query()["from"])
	var to = queryStationList(r.URL.Query()["to"])
	var vias = queryStationList(r.URL.Query()["vias"])

	var datetime = r.URL.Query().Get("datetime")
	var regionly = r.URL.Query().Get("regionly") == "true"

	w.Header().Set("Content-Type", "application/json")
	if len(from) > 0 && len(to) > 0 {
		stations, err := internal.ObtainVias(system, from[0], to[0], vias, datetime, regionly)
		if err != nil {
			w.WriteHeader(err.ErrorCode())
		}
		render.ViasApi(stations, from[0], to[0], datetime, w, err)
		return
	}
	log.Print("Bad Request, from/to unset")
	w.WriteHeader(http.StatusBadRequest)
}

func apiTimespace(w http.ResponseWriter, r *http.Request) {
	var system = r.URL.Query().Get("system")
	var from = queryStationList(r.URL.Query()["from"])
	var to = queryStationList(r.URL.Query()["to"])
	var vias = queryStationList(r.URL.Query()["vias"])

	var datetime = r.URL.Query().Get("datetime")
	var regionly = r.URL.Query().Get("regionly") == "true"

	w.Header().Set("Content-Type", "application/json")
	if len(from) > 0 && len(to) > 0 && len(vias) > 0 {
		stations, lines, err := internal.ObtainData(system, from[0], to[0], vias, datetime, regionly)
		if err != nil {
			w.WriteHeader(err.ErrorCode())
			return
		}
		render.TimeSpaceApi(stations, lines, w, r.URL.RawQuery)
		return
	}
	log.Print("Bad Request, from/to/vias unset")
	w.WriteHeader(http.StatusBadRequest)
}

func queryStationList(params []string) []string {
	var ids = []string{}
	for _, id := range params {
		if id == "" || id == "0" {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
