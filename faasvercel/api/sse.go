package api

import (
	"fmt"
	"net/http"
	"time"
)

func Sse(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a ticker to send events periodically
	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()
	cnt := 0
	for {
		select {
		case <-r.Context().Done():
			// Client closed the connection
			return
		case t := <-ticker.C:
			cnt++
			if cnt > 20 {
				return
			}
			// Send an event to the client
			fmt.Fprintf(w, "data: %s\n\n", t.String())
			// Flush the response to ensure the event is sent
			flusher, ok := w.(http.Flusher)
			if ok {
				flusher.Flush()
			}
		}
	}
}
