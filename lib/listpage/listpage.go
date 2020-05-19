// Package listpage serves generated HTML for the root endpoint of the bysykkel app.
package listpage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	"kaldager.com/oslobysykkel/lib/oslobysykkel"
)

type Handler struct {
	Source oslobysykkel.DataSource
}

func (h Handler) prepareParams(ctx context.Context, req *http.Request) (*pageParams, error) {
	allStations, err := h.Source.GetAllStations(ctx)
	if err != nil {
		return nil, err
	}

	var stations []*oslobysykkel.Station
	for _, station := range allStations {
		stations = append(stations, station)
	}
	sort.Slice(stations, func(i, j int) bool {
		return stations[i].Info.Name < stations[j].Info.Name
	})

	return &pageParams{Stations: stations}, nil
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.serveHTTP(w, req)
}

func (h Handler) serveHTTP(w http.ResponseWriter, req *http.Request) error {
	if req.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed")
	}
	params, err := h.prepareParams(req.Context(), req)
	if err != nil {
		params = &pageParams{Error: err}
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	if err := listpageTmpl.Execute(w, params); err != nil {
		log.Printf("warning: write failed: %v", err)
		return err
	}

	return nil
}
