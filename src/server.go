package src

import (
	"net/http"
	"time"
)

// NewMux creates a new ServeMux and registers the app routes.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static assets like background images.
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Register the home route.
	mux.HandleFunc("/", HomeHandler)

	// Register the JSON API route.
	mux.HandleFunc("/api/ascii-art", ApiAsciiArtHandler)

	return mux
}

// StartServer creates the router and starts the HTTP server.
func StartServer(address string) error {
	mux := NewMux()

	// Wrap the mux with rate limiting: 5 requests per second
	limiter := NewRateLimiter(5, 1*time.Second)
	limitedMux := RateLimitMiddleware(mux, limiter)

	server := &http.Server{
		Addr:         address,
		Handler:      limitedMux, // use the wrapped mux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}
