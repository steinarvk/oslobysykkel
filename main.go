// Binary bysykkel serves real-time information about Oslo city bikes.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"kaldager.com/oslobysykkel/lib/apihandler"
	"kaldager.com/oslobysykkel/lib/listpage"
	"kaldager.com/oslobysykkel/lib/oslobysykkel"
)

func mainCore() error {
	ctx := context.Background()

	// Create a data source from which API information will be pulled when necessary.
	// (We don't need to pull once for every request, only often enough to keep
	// the freshness generally lower than 10s.)
	src, err := oslobysykkel.NewDynamicDataSource("https://gbfs.urbansharing.com/oslobysykkel.no")
	if err != nil {
		return err
	}

	// Pull once from the data source at startup to fail-fast in case of errors.
	// We don't want to begin serving if we can't reach the data source.
	if err := selfcheck(ctx, src); err != nil {
		return fmt.Errorf("startup self-check failed: %w", err)
	}

	// Main handler, to serve generated HTML intended for humans.
	h := listpage.Handler{Source: src}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			// The default mux treats Handle("/", ...) as a path prefix.
			// Reject non-exact matches.
			http.NotFound(w, req)
			return
		}
		h.ServeHTTP(w, req)
	})

	// API handler; this serves data intended for machines.
	http.Handle("/api/get-all-stations", apihandler.GetAllHandler{Source: src})
	http.Handle("/api/get-distances", apihandler.GetDistancesHandler{Source: src})

	// Serve static files required for the HTML page.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("static/fonts/"))))

	// Begin serving. This should normally never return.
	return serveForever()
}

func serveForever() error {
	portstring := os.Getenv("PORT")
	if portstring == "" {
		portstring = "8080"
	}
	addr := ":" + portstring

	log.Printf("serving on: %s", addr)

	return http.ListenAndServe(addr, nil)
}

func selfcheck(ctx context.Context, src oslobysykkel.DataSource) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := src.GetAllStations(ctx)
	return err
}

func main() {
	if err := mainCore(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}
