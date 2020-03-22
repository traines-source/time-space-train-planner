package main

import (
	"log"
	"net/http"
	"os"

	"traines.eu/time-space-train-planner/internal"
	"traines.eu/time-space-train-planner/render"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	stations, lines := internal.ObtainData()
	render.TimeSpace(stations, lines, w)
}
