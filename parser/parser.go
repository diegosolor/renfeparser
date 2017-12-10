package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
    "github.com/diegosolor/renfeparser/common"
)

func createJourneyURL(origin string, destiny string, search_date time.Time) (url string) {
	year, month_name, day := search_date.Date()
	month := int(month_name)
	log.Print(fmt.Sprintf("Parsing journey  %s -> %s at %v-%v-%v", origin, destiny, year, month, day))
	url = fmt.Sprintf("http://horarios.renfe.com/HIRRenfeWeb/buscar.do?O=%s&D=%s&AF=%v&MF=%v&DF=%v", origin, destiny, year, month, day)
	log.Print(url)
	return
}

func ParseJourneysForDay(origin string, destiny string, search_date time.Time) (journeys []common.Journey) {
	url := createJourneyURL(origin, destiny, search_date)
	doc, err := goquery.NewDocument(url)
	common.CheckError(err)

	doc.Find("#row tbody .odd,.even").Each(func(row_number int, row_journey *goquery.Selection) {
		journey := parseJourneyRow(row_journey, origin, destiny, search_date)
        journey = CheckPosibleOffer(journey)
		journeys = append(journeys, journey)
	})
	return
}

func ParseJourneysForPeriod(origin string, destiny string, start_date time.Time, end_date time.Time) (journeys []common.Journey) {
	search_date := start_date
	for search_date.Before(end_date) {
		days_journeys := ParseJourneysForDay(origin, destiny, search_date)
		if len(days_journeys) == 0 {
			//if none journeys are found we suppose this is the last day renfe operator pusblished journeys
			log.Print(fmt.Sprintf("No journeys are found for date %v, we stop looking forward", search_date))
			break
		}
		journeys = append(journeys, days_journeys...)
		search_date = search_date.AddDate(0, 0, 1)
	}
	return
}

func parseJourneyRow(tr *goquery.Selection, origin string, destiny string, search_date time.Time) common.Journey {
	journey := common.Journey{Origin: origin, Destiny: destiny}
	tr.Find("td").Each(func(column_number int, column_content *goquery.Selection) {
		switch column_number {
		case 0: //id and train type
                journey.Train_type, journey.Train_id = parseTrainTypeId(column_content.Text())
		case 1: //departure
			journey.Departure = parse_time(search_date, column_content.Text())
		case 2: //arrival
			journey.Arrival = parse_time(search_date, column_content.Text())
		case 4:
            prices_by_class := parseJourneyPrices(column_content)
            journey.AddPrices(prices_by_class)
		}
	})
	return journey
}

func parseTrainTypeId(train_type_and_id string) (train_type string, id string) {
    train_type_and_id = standardizeSpaces(train_type_and_id)
    type_and_id := strings.Split(train_type_and_id, " ")
    train_type = type_and_id[1]
    id = type_and_id[0]
    return
}

func parse_time(date time.Time, str_time string) time.Time {
	hour_and_minute := strings.Split(str_time, ".")
	hour, err := strconv.Atoi(hour_and_minute[0])
	common.CheckError(err)
	minute, err := strconv.Atoi(hour_and_minute[1])
	common.CheckError(err)
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
		common.CheckError(err)
		prices_by_class[class] = price
	})
	return prices_by_class
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
