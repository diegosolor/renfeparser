package main

import (
    "log"
    "time"
    "flag"
	"fmt"
    "github.com/diegosolor/renfeparser/parser"
    "github.com/diegosolor/renfeparser/common"
    "github.com/diegosolor/renfeparser/storable"
)

func main() {
    origin := flag.String("origin","MADRI","Origin station")
    destiny := flag.String("destiny","SANTA","Destiny station")
    str_start_date := flag.String("start_date",defaultStartDate(),"Start search date, in format YYYYMMDD")
    str_end_date := flag.String("end_date",defaultEndDate(),"End search date, in format YYYYMMDD")
    flag.Parse()

    start_date := parseDate(*str_start_date)
    end_date := parseDate(*str_end_date)
    journeys := parser.ParseJourneysForPeriod(*origin, *destiny, start_date, end_date)
    for _, journey := range(journeys) {
        log.Print(journey)
    }

    csv_name := fmt.Sprintf("journeys_%v_%v_from_%v_to_%v.csv", *origin, *destiny, *str_start_date, *str_end_date)
    storable.ExportToCSV(journeys, csv_name, nil)
    storable.ExportToPostgreSQL(journeys)

}

func parseDate(str_date string) (date time.Time) {
    tz_madrid, err := time.LoadLocation("Europe/Madrid")
    date, err = time.ParseInLocation(dateFormat(), str_date, tz_madrid)
	common.CheckError(err)
    return
}

func defaultEndDate() (string) {
    now := time.Now()
    two_months_from_now := now.AddDate(0,2,0)
	return two_months_from_now.Format(dateFormat())
}

func defaultStartDate() (string) {
    now := time.Now()
	return now.Format(dateFormat())
}

func dateFormat() (string) {
    const date_format = "20060102"
    return date_format
}
