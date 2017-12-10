package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
    "github.com/diegosolor/renfeparser/common"
)

func CheckPosibleOffer(journey common.Journey) common.Journey {

	template_post_values := `
callCount=1
scriptSessionId=99999
c0-param4=number:99999
c0-id=99999
c0-scriptName=preciosManager
c0-methodName=calcularPreciosAJAX
c0-param0=string:%v
c0-param1=string:%v
c0-param2=string:%v
c0-param3=string:%v
batchId=1`

	post_values := fmt.Sprintf(template_post_values,
		journey.Train_id,
		journey.Departure.Format("02-01-2006"),
		journey.Origin,
		journey.Destiny)

	post_body_values := bytes.NewBufferString(post_values)

	url_offers := "http://horarios.renfe.com/HIRRenfeWeb/dwr/call/plaincall/preciosManager.calcularPreciosAJAX.dwr"

	resp, err_post := http.Post(url_offers, " text/plain", post_body_values)
	common.CheckError(err_post)

	defer resp.Body.Close()
	body, err_read := ioutil.ReadAll(resp.Body)
	common.CheckError(err_read)
	js_response := string(body[:])

	prices_by_class := parseJSResponse(js_response)
	journey.AddPrices(prices_by_class)

	return journey
}

func parseJSResponse(js_response string) map[string]float64 {
	parse_js := regexp.MustCompile(`clase="([A-Z]{1,2})".*precio="([0-9]+,[0-9]{2})"`)
	matches := parse_js.FindAllStringSubmatch(js_response, -1)
	prices_by_class := make(map[string]float64)

	for _, offer_match := range matches {
		str_price := strings.Replace(offer_match[2], ",", ".", 1)
		price, err := strconv.ParseFloat(str_price, 2)
		common.CheckError(err)
		class := offer_match[1]
		prices_by_class[class] = price

	}

	return prices_by_class

}
