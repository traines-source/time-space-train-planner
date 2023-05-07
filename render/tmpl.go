package render

import (
	"fmt"
	"io"
	"text/template"
	"time"

	"traines.eu/time-space-train-planner/internal"
)

func (c *container) render(wr io.Writer) {
	t := template.Must(template.New("time-space.tmpl.svg").ParseFiles("./render/time-space.tmpl.svg"))
	err := t.Execute(wr, c)
	if err != nil {
		panic(err)
	}
}

func (c *container) X(coord Coord) int {
	if coord.SpaceAxis == "" {
		return 0
	}
	return int(float32(c.Stations[coord.SpaceAxis].SpaceAxis)/float32(c.MaxSpace)*float32(c.SpaceAxisSize-100) + 50.0)
}

func (c *container) Y(coord Coord) int {
	if coord.TimeAxis.IsZero() {
		return 50 + c.Stations[coord.SpaceAxis].SpaceAxisHeap*20
	}
	delta := float32(coord.TimeAxis.Unix() - c.MinTime.Unix())
	return int(delta/c.TimeAxisDistance*float32(c.TimeAxisSize-100) + 100.0)
}

func (e *EdgePath) Label() string {
	if e.Line == nil {
		return ""
	}
	label := e.Line.Name
	if e.Message != "" {
		label += " (" + substr(e.Message, 0, 30) + "...)"
	}
	if e.Line.Type == "Foot" {
		return "🚶 " + label
	}
	return label
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func (e *EdgePath) Type() string {
	if e.Line == nil {
		return ""
	}
	return e.Line.Type
}

func (p *EdgePath) Departure() string {
	return p.time(func(stop internal.StopInfo) time.Time { return stop.Departure }, func(stop internal.StopInfo) string { return stop.DepartureTrack })
}

func (p *EdgePath) Arrival() string {
	return p.time(func(stop internal.StopInfo) time.Time { return stop.Arrival }, func(stop internal.StopInfo) string { return stop.ArrivalTrack })
}

func (p *EdgePath) LiveDataDeparture() string {
	return p.liveDataClass(func(stop internal.StopInfo) time.Time { return stop.Departure })
}

func (p *EdgePath) LiveDataArrival() string {
	return p.liveDataClass(func(stop internal.StopInfo) time.Time { return stop.Arrival })
}

func (e *EdgePath) time(timeResolver func(internal.StopInfo) time.Time, trackResolver func(internal.StopInfo) string) string {
	if e.Line == nil {
		return ""
	}
	label := simpleTime(timeResolver(e.Actual)) + delay(timeResolver(e.Current), timeResolver(e.Planned))
	if trackResolver(e.Planned) != "" {
		label += "Pl." + trackResolver(e.Planned)
	}
	return label
}

func (e *EdgePath) liveDataClass(timeResolver func(internal.StopInfo) time.Time) string {
	if e.Line == nil {
		return ""
	}
	current := timeResolver(e.Current)
	if current.IsZero() {
		return ""
	}
	if current.Sub(timeResolver(e.Planned)).Minutes() > 5 {
		return "live-red"
	}
	return "live-green"
}

func (c *container) Minutes(time time.Time) string {
	return fmt.Sprintf("%.0f", time.Sub(c.MinTime).Minutes())
}

func simpleTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d ", t.Hour(), t.Minute())
}

func (c *container) SimpleTime(t time.Time) string {
	return simpleTime(t)
}

func delay(current time.Time, planned time.Time) string {
	if !current.IsZero() {
		return " (" + fmt.Sprintf("%+.0f", current.Sub(planned).Minutes()) + ") "
	}
	return ""
}
