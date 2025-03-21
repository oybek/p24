package rest

import (
	"encoding/json"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (rest *Rest) TripFind(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.Has("user_type") || !query.Has("city_a") || !query.Has("city_b") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userType := query.Get("user_type")
	cityA := query.Get("city_a")
	cityB := query.Get("city_b")

	trips, err := rest.mc.TripFind(userType, cityA, cityB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trips)
}

func (rest *Rest) TripCard(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !query.Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(query.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	trip, err := rest.mc.TripGetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := rest.mc.UserGetByChatID(trip.ChatID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := rest.bot.DrawCard(trip, user)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(bytes)
}
