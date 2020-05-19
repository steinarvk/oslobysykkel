package oslobysykkel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// StaticSource is a DataSource reading from static files; meant for testing.
type StaticSource map[string]*Station

func (s StaticSource) GetAllStations(ctx context.Context) (map[string]*Station, error) {
	return s, nil
}

func readStationStatus(filename string) (*StationStatusResponse, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rv StationStatusResponse
	if err := json.Unmarshal(data, &rv); err != nil {
		return nil, err
	}
	return &rv, nil
}

func readStationInfo(filename string) (*StationInformationResponse, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rv StationInformationResponse
	if err := json.Unmarshal(data, &rv); err != nil {
		return nil, err
	}
	return &rv, nil
}

// NewStaticDataSource creates a new static data source.
func NewStaticDataSource(statusFilename, infoFilename string) (StaticSource, error) {
	status, err := readStationStatus(statusFilename)
	if err != nil {
		return nil, fmt.Errorf("error reading station status: %w", err)
	}

	info, err := readStationInfo(infoFilename)
	if err != nil {
		return nil, fmt.Errorf("error reading station information: %w", err)
	}

	rv := map[string]*Station{}

	getStation := func(k string) *Station {
		st, ok := rv[k]
		if ok {
			return st
		}
		st = &Station{}
		rv[k] = st
		return st
	}

	for _, stationStatus := range status.Data.Stations {
		status := stationStatus
		getStation(stationStatus.StationId).Status = &status
	}

	for _, stationInfo := range info.Data.Stations {
		info := stationInfo
		getStation(stationInfo.StationId).Info = &info
	}

	for stationId, station := range rv {
		if station.Status == nil {
			return nil, fmt.Errorf("station %q is missing Status", stationId)
		}
		if station.Info == nil {
			return nil, fmt.Errorf("station %q is missing Info", stationId)
		}
	}

	return StaticSource(rv), nil
}
