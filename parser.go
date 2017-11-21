package renfeparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

func createJourneyURL(origin string, destiny string, search_date time.Time) (url string) {
	year, month_name, day := search_date.Date()
	month := int(month_name)
	log.Print(fmt.Sprintf("Parsing journey  %s -> %s at %v-%v-%v", origin, destiny, year, month, day))
	url = fmt.Sprintf("http://horarios.renfe.com/HIRRenfeWeb/buscar.do?O=%s&D=%s&AF=%v&MF=%v&DF=%v", origin, destiny, year, month, day)
	log.Print(url)
	return
}

func ParseJourneysForDay(origin string, destiny string, search_date time.Time) (journeys []Journey) {
	url := createJourneyURL(origin, destiny, search_date)
	doc, err := goquery.NewDocument(url)
	CheckError(err)

	doc.Find("#row tbody .odd,.even").Each(func(row_number int, row_journey *goquery.Selection) {
		journey := parseJourneyRow(row_journey, origin, destiny, search_date)
		journeys = append(journeys, journey)
	})

	return
}

func ParseJourneysForPeriod(origin string, destiny string, start_date time.Time, end_date time.Time) (journeys []Journey) {
    search_date := start_date
    for search_date.Before(end_date) {
        days_journeys := ParseJourneysForDay(origin, destiny, search_date)
        journeys = append(journeys, days_journeys...)
        search_date = search_date.AddDate(0,0,1)
    }
    return
}

func parseJourneyRow(tr *goquery.Selection, origin string, destiny string, search_date time.Time) Journey {
	journey := Journey{Origin: origin, Destiny: destiny}
	tr.Find("td").Each(func(column_number int, column_content *goquery.Selection) {
		switch column_number {
		case 0:
			journey.Train_type = standardizeSpaces(column_content.Text())
		case 1: //departure
			journey.Departure = parse_time(search_date, column_content.Text())
		case 2: //arrival
			journey.Arrival = parse_time(search_date, column_content.Text())
		case 4:
			journey.Prices_by_class = parseJourneyPrices(column_content)
		}

	})
	return journey
}

func parse_time(date time.Time, str_time string) time.Time {
	hour_and_minute := strings.Split(str_time, ".")
	hour, err := strconv.Atoi(hour_and_minute[0])
	CheckError(err)
	minute, err := strconv.Atoi(hour_and_minute[1])
	CheckError(err)
	parsed_date := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())

	return parsed_date
}

func parseJourneyPrices(tr *goquery.Selection) map[string]float64 {
	prices_by_class := make(map[string]float64)
	tr.Find("tr #divcont").Each(func(column_number int, column_content *goquery.Selection) {
		content := standardizeSpaces(column_content.Text())

		class_and_price := strings.Split(content, ": ")
		str_price := strings.Replace(class_and_price[1], ",", ".", 1)
		class := class_and_price[0]
		price, err := strconv.ParseFloat(str_price, 2)
		CheckError(err)
		prices_by_class[class] = price

		//log.Print(fmt.Sprintf("\tTarifa %s: %v",class, price))

	})
	return prices_by_class
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
