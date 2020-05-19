// Package oslobysykkel defines JSON-compatible structs for the oslobysykkel API and data sources to pull from it.
package oslobysykkel

import "context"

type StationStatus struct {
	StationId               string `json:"station_id"`
	IsInstalled             int    `json:"is_installed"`
	IsRenting               int    `json:"is_renting"`
	IsReturning             int    `json:"is_returning"`
	LastReportedUnixSeconds int64  `json:"last_reported"`
	NumBikesAvailable       int    `json:"num_bikes_available"`
	NumDocksAvailable       int    `json:"num_docks_available"`
}

type StationInformation struct {
	Address   string  `json:"address"`
	StationId string  `json:"station_id"`
	Name      string  `json:"name"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Capacity  int     `json:"capacity"`
}

type StationStatusResponse struct {
	Data struct {
		Stations []StationStatus `json:"stations"`
	} `json:"data"`
	LastUpdated int64 `json:"last_updated"`
	TTL         int   `json:"ttl"`
}

type StationInformationResponse struct {
	Data struct {
		Stations []StationInformation `json:"stations"`
	} `json:"data"`
	LastUpdated int64 `json:"last_updated"`
	TTL         int   `json:"ttl"`
}

type Station struct {
	Status *StationStatus      `json:"status"`
	Info   *StationInformation `json:"info"`
}

type DataSource interface {
	GetAllStations(context.Context) (map[string]*Station, error)
}
