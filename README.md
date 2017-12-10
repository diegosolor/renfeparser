
 # renfeparser
    import "github.com/diegosolor/renfeparser"


Library to parse train's schedule from spanish renfe operator

## Usage

    $ go run main/main_parser.go --origin SANTA --destiny MADRI 

by default it searches the two next months, you can change this behaviur wiyh the flags start_date and end_date
The result will be a csv with the cheapest journey for each train

## TODO

dictionary for station's names
allow to store the data in a data base
paralelize requests

## Disclaimer

This is a personal project with two objetives: learn some Go language and enter in renfe's web as few time as possible. Any feedback is more than wellcome
