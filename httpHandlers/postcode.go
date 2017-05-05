package httpHandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"regexp"

	"github.com/wpferg/house-prices/store"
)

var postcodeRegex = regexp.MustCompile("/postcode/(?P<postcode>.+)")

func PostcodeSearch(w http.ResponseWriter, r *http.Request) {
	matchedDetails := postcodeRegex.FindStringSubmatch(r.URL.Path)

	if len(matchedDetails) < 2 {
		log.Println("matched details", matchedDetails)
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	postcode := matchedDetails[1]
	log.Println("Searching postcodes:", postcode)

	results := store.SearchPostcode(postcode)
	log.Println("Postcodes matching", postcode, len(results))

	headers := w.Header()
	headers.Add("Content-Type", "application/json")
	headers.Add("Access-Control-Allow-Origin", "*")

	w.WriteHeader(200)
	jsonData, _ := json.Marshal(results)
	w.Write(jsonData)
}
