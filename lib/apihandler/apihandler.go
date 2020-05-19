// Package apihandler contains API handlers for the bysykkel app.
package apihandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/umahmood/haversine"
	"kaldager.com/oslobysykkel/lib/oslobysykkel"
)

type GetAllHandler struct {
	Source oslobysykkel.DataSource
}

func (h GetAllHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allStations, err := h.Source.GetAllStations(req.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(allStations)
}

type GetDistancesHandler struct {
	Source oslobysykkel.DataSource
}

func (h GetDistancesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allStations, err := h.Source.GetAllStations(req.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}
	origins, ok := req.URL.Query()["origin"]
	if !ok || len(origins) != 1 {
		http.Error(w, fmt.Sprintf("error: expected exactly 1 origin, got %d", len(origins)), http.StatusBadRequest)
		return
	}
	originStationID := origins[0]
	originStation, ok := allStations[originStationID]
	if !ok {
		http.Error(w, fmt.Sprintf("error: no such station %q", originStationID), http.StatusBadRequest)
		return
	}

	rv := map[string]interface{}{}
	c0 := haversine.Coord{Lat: originStation.Info.Lat, Lon: originStation.Info.Lon}
	for stationID, station := range allStations {
		c1 := haversine.Coord{Lat: station.Info.Lat, Lon: station.Info.Lon}
		_, km := haversine.Distance(c0, c1)
		rv[stationID] = km
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(rv)
}
