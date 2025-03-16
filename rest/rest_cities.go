package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

func (rest *Rest) Cities(w http.ResponseWriter, r *http.Request) {
	cityNames, err := rest.mc.CityNamesGet()
	if err != nil {
		log.Printf("failed to load cityNames: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cities := []City{}
	for k, v := range cityNames {
		if k == "_id" {
			continue
		}
		cities = append(cities, City{k, v})
	}

	json, err := json.Marshal(cities)
	if err != nil {
		log.Printf("failed to marshall cities: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
