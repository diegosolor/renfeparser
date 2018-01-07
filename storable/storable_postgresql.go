package storable

import (
	"database/sql"
	"fmt"
	"github.com/diegosolor/renfeparser/common"
	_ "github.com/lib/pq"
)

type JourneysToPostgreSQL struct {
	Journeys            []common.Journey
	Only_cheapest_class bool
	data_base           *sql.DB
}

const (
	DB_USER      = "renfe"
	DB_PASSWORD  = "renfe"
	DB_NAME      = "renfe"
	QUERY_INSERT = `INSERT INTO public.journeys (origin, destiny, departure_date, arrival_date, class, price) 
                    VALUES($1,$2,$3,$4,$5,$6) 
                    on conflict on constraint journeys_pkey do update set check_date = EXCLUDED.check_date, price = EXCLUDED.price;`
)

func ExportToPostgreSQL(journeys []common.Journey) {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	common.CheckError(err)

	to_pg := JourneysToPostgreSQL{
		Journeys:            journeys,
		Only_cheapest_class: true,
		data_base:           db}
	to_pg.Write()
}

func (to_pg JourneysToPostgreSQL) Write() {

	for _, journey := range to_pg.Journeys {
		to_pg.writeJourney(journey)
	}
}

func (to_pg JourneysToPostgreSQL) writeJourney(journey common.Journey) {
	classes := journey.Classes()
	if to_pg.Only_cheapest_class {
		classes = []string{journey.CheapestClass()}
	}
	for _, class := range classes {
		_, err := to_pg.data_base.Exec(
			QUERY_INSERT,
			journey.Origin,
			journey.Destiny,
			journey.Departure.Format("2006-01-02 15:04"),
			journey.Arrival.Format("2006-01-02 15:04"),
			class,
			fmt.Sprintf("%v", journey.Price(class)))
		common.CheckError(err)
	}
}

func (to_pg JourneysToPostgreSQL) journeyValues(journey common.Journey, class string) []string {
	journey_values := []string{
		journey.Origin,
		journey.Destiny,
		journey.Departure.Format("2006-01-02 15:04"),
		journey.Arrival.Format("2006-01-02 15:04"),
		class,
		fmt.Sprintf("%v", journey.Price(class))}
	return journey_values
}
