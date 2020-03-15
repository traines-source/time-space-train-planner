package main

import (
	"traines.eu/time-space-train-planner/internal"
)

func main() {
	c := &internal.Consumer{}
	c.CallProviders()
}
