
 # renfeparser
    import "github.com/diegosolor/renfeparser"


Library to parse train's schedule from spanish renfe operator

Parsing a period:
``` go

package main

import (
    "log"
    "time"
    "github.com/diegosolor/renfeparser"
)

func main() {
    tz_madrid, err := time.LoadLocation("Europe/Madrid")
    renfeparser.CheckError(err)
    start_date := time.Date(2017,12,16,0, 0, 0, 0, tz_madrid)
    end_date := time.Date(2017,12,28,0, 0, 0, 0, tz_madrid)
    journeys := renfeparser.ParseJourneysForPeriod("MADRI","SANTA", start_date, end_date)
    for _, journey := range(journeys) {
        log.Print(journey)
    }

    renfeparser.ExportToCSV(journeys, "result.csv", nil)

}

```
