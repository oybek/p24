package rest

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

func (rest *Rest) TripFind(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.Has("city_a") || !query.Has("city_b") || !query.Has("date") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cityA := query.Get("city_a")
	cityB := query.Get("city_b")
	dateRaw := query.Get("date")
	date, err := time.Parse(time.RFC3339, dateRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	trips, err := rest.mc.TripFind(cityA, cityB, date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trips)
}
