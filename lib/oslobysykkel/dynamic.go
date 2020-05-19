package oslobysykkel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	// StalenessThreshold is the maximum acceptable staleness.
	// If the data is already fresher than this, no new data is pulled from the source.
	// If the data is staler than this but fresher than SyncStalenessThreshold, a new
	// pull from the data source is initiated but the current requested is served using
	// the cached values.
	StalenessThreshold = 10 * time.Second

	// SyncStalenessThreshold is the maximum tolerable staleness.
	// If the data is staler than this, we will _wait_ for a refresh to complete before
	// serving the user's request.
	SyncStalenessThreshold = 1 * time.Minute
)

// DynamicSource is a live data source for the oslobysykkel API, pulling data over HTTP/HTTPS.
type DynamicSource struct {
	mu   sync.Mutex
	cond *sync.Cond

	apiPrefix string

	stations map[string]*Station

	lastCompletedUpdate time.Time
	updateOngoing       bool
}

// NewDynamicDataSource creates a new dynamic data source.
// apiPrefix should usually be equal to "https://gbfs.urbansharing.com/oslobysykkel.no".
func NewDynamicDataSource(apiPrefix string) (*DynamicSource, error) {
	rv := &DynamicSource{
		apiPrefix: apiPrefix,
		stations:  map[string]*Station{},
	}
	rv.cond = sync.NewCond(&rv.mu)
	return rv, nil
}

var (
	httpClient = &http.Client{
		Timeout: time.Second * 30,
	}
)

// getJSON makes one HTTP(S) GET request and unmarshals the response as JSON.
func getJSON(ctx context.Context, dest interface{}, uri string) error {
	t0 := time.Now()
	log.Printf("fetching %q", uri)
	defer func() {
		log.Printf("done fetching %q after %v", uri, time.Since(t0))
	}()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Client-Identifier", "steinarkaldager-oslobysykkkel.app.kaldager.com")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dest)
}

// update synchronously performs an update from the source API, integrating the results into the stored values.
func (s *DynamicSource) update(ctx context.Context, apiPrefix string) error {
	wg := sync.WaitGroup{}

	wg.Add(2)

	var statusResponse StationStatusResponse
	var statusErr error

	statusURL := apiPrefix + "/station_status.json"
	infoURL := apiPrefix + "/station_information.json"

	var infoResponse StationInformationResponse
	var infoErr error

	go func() {
		defer wg.Done()
		statusErr = getJSON(ctx, &statusResponse, statusURL)
	}()

	go func() {
		defer wg.Done()
		infoErr = getJSON(ctx, &infoResponse, infoURL)
	}()

	wg.Wait()

	if statusErr != nil {
		return fmt.Errorf("error fetching status: %w", statusErr)
	}
	if infoErr != nil {
		return fmt.Errorf("error fetching info: %w", infoErr)
	}

	if err := s.integrateUpdates(&statusResponse, &infoResponse); err != nil {
		return fmt.Errorf("error integrating updates: %w", err)
	}

	return nil
}

// integrateUpdates integrates refreshed responses into the stored values.
func (s *DynamicSource) integrateUpdates(statusResponse *StationStatusResponse, infoResponse *StationInformationResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	getStation := func(k string) *Station {
		st, ok := s.stations[k]
		if ok {
			return st
		}
		st = &Station{}
		s.stations[k] = st
		return st
	}

	for _, stationStatus := range statusResponse.Data.Stations {
		status := stationStatus
		getStation(stationStatus.StationId).Status = &status
	}

	for _, stationInfo := range infoResponse.Data.Stations {
		info := stationInfo
		getStation(stationInfo.StationId).Info = &info
	}

	s.lastCompletedUpdate = time.Now()

	s.cond.Broadcast()

	return nil
}

// holdingLockMaybeStartUpdate spawns (in the background) a new update, unless one is already going on.
// As the name indicates, the lock should be held when calling this.
func (s *DynamicSource) holdingLockMaybeStartUpdate(ctx context.Context) {
	if s.updateOngoing {
		return
	}

	apiPrefix := s.apiPrefix

	s.updateOngoing = true
	go func() {

		defer func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.updateOngoing = false
		}()

		if err := s.update(ctx, apiPrefix); err != nil {
			log.Printf("error updating data: %v", err)
		}
	}()
}

// maybeRefreshData refreshes data if there's a need for it, depending on staleness.
// It may spawn an update in the background, and it may or may not wait for the next
// update to complete before returning.
func (s *DynamicSource) maybeRefreshData(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	staleness := time.Since(s.lastCompletedUpdate)

	shouldRequestUpdate := staleness > StalenessThreshold
	shouldWaitForUpdate := staleness > SyncStalenessThreshold

	if shouldRequestUpdate {
		log.Printf("requesting update (due to staleness %v)", staleness)
		s.holdingLockMaybeStartUpdate(ctx)
	}

	if shouldWaitForUpdate {
		for staleness > SyncStalenessThreshold {
			log.Printf("waiting synchronously for update (staleness %v)", staleness)
			s.holdingLockMaybeStartUpdate(ctx)
			if err := ctx.Err(); err != nil {
				return err
			}
			s.cond.Wait()
			staleness = time.Since(s.lastCompletedUpdate)
		}
	}

	return nil
}

// GetAllStations returns data from the dynamic data source.
func (s *DynamicSource) GetAllStations(ctx context.Context) (map[string]*Station, error) {
	t0 := time.Now()
	log.Printf("received GetAllStations request")
	defer func() {
		log.Printf("done processing GetAllStations request after %v", time.Since(t0))
	}()
	if err := s.maybeRefreshData(ctx); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Take a shallow copy of the stations map.
	rv := map[string]*Station{}
	for k, v := range s.stations {
		if v.Info == nil || v.Status == nil {
			// Skip any entries with partial data.
			log.Printf("warning: station %q is missing either Info or Status", k)
			continue
		}
		rv[k] = v
	}

	return rv, nil
}
