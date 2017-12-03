package renfeparser

import (
    "fmt"
    "time"
    "log"
)

type Journey struct {
    Origin, Destiny, Train_type string
    Departure, Arrival time.Time
    Prices_by_class map[string]float64 
    Train_id string
}

func (j Journey) String() string {
    format := "%s %s -> %s (from %s to %s)"
    journey_text := fmt.Sprintf(format, j.Train_type, j.Origin, j.Destiny, j.Departure.Format("2006-01-02 15:04"), j.Arrival.Format("15:04"))
    for class, price := range j.Prices_by_class {
        journey_text = journey_text + fmt.Sprintf("\n\t%s: %v euros",class, price)
    }
    return journey_text
}

func (j *Journey) AddPrices(prices map[string]float64) {
    if  j.Prices_by_class == nil {
        j.Prices_by_class = make(map[string]float64)
    }

    for class, price := range prices {
        j.Prices_by_class[class] = price
    }
}

func (j *Journey) Price(class string) float64 {
    price, class_exists := j.Prices_by_class[class]
    if (!class_exists) {
        log.Fatal(fmt.Sprintf("Class %v not found in journey %v", class, j))
    }
    return price

}
