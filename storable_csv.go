package renfeparser

import (
	"encoding/csv"
	"fmt"
	"os"
)

type JourneysToCSV struct {
	Headers             []string
	File_name           string
	Only_cheapest_class bool
    Journeys            []Journey
}

func defaultHeaders() []string {

	headers := []string{"origin", "destiny", "departure_date", "arrival_date", "class", "price"}
	return headers
}

func ExportToCSV(journeys []Journey, file_name string, headers []string) {

	if headers == nil {
		headers = defaultHeaders()
	}

	to_csv := JourneysToCSV{
		Headers:             headers,
		File_name:           file_name,
		Only_cheapest_class: true,
        Journeys : journeys}

	to_csv.WriteFile()

}

func (to_csv JourneysToCSV) WriteFile() {

	file, err := os.Create(to_csv.File_name)
	CheckError(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(to_csv.Headers)

	for _, journey := range to_csv.Journeys {
		to_csv.writeJourney(journey, writer)
	}
}

func (to_csv JourneysToCSV) writeJourney(journey Journey, csv_writer *csv.Writer) {
	classes := journey.Classes()
	if to_csv.Only_cheapest_class {
		classes = []string{journey.CheapestClass()}
	}
	for _, class := range classes {
		journey_values := to_csv.journeyValues(journey, class)
		err := csv_writer.Write(journey_values)
		CheckError(err)
	}
}

func (to_csv JourneysToCSV) journeyValues(journey Journey, class string) []string {
	journey_values := make([]string, 0)
	for _, column := range to_csv.Headers {
		var value string
		switch column {
		case "origin":
			value = journey.Origin
		case "destiny":
			value = journey.Destiny
		case "departure_date":
			value = journey.Departure.Format("2006-01-02 15:04")
		case "arrival_date":
			value = journey.Arrival.Format("2006-01-02 15:04")
		case "class":
			value = class
		case "price":
			value = fmt.Sprintf("%v", journey.Price(class))
		}
		journey_values = append(journey_values, value)
	}

	return journey_values
}
