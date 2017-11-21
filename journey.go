package renfeparser

import (
	"fmt"
	"time"
)

type Journey struct {
	Origin, Destiny, Train_type string
	Departure, Arrival          time.Time
	Prices_by_class             map[string]float64
}

func (j Journey) String() string {
	format := "%s -> %s (from %s to %s)"
	journey_text := fmt.Sprintf(format, j.Origin, j.Destiny, j.Departure.Format("2006-01-02 15:04"), j.Arrival.Format("15:04"))
	for class, price := range j.Prices_by_class {
		journey_text = journey_text + fmt.Sprintf("\n\t%s: %v euros", class, price)
	}
	return journey_text
}
