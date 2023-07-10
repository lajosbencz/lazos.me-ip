package main

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"strings"
)

//go:embed favicon.ico
var faviconFile embed.FS

func getClientIP(r *http.Request) string {
	// Get the client IP address from the X-Forwarded-For header if available
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		// Return the first IP address from the list
		return strings.TrimSpace(ips[0])
	}

	// Get the remote IP address from the request
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	return remoteIP
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	w.Header().Add("Content-Type", "text/plain")
	response := fmt.Sprintf("%s", ip)
	fmt.Fprintln(w, response)
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	faviconData, err := faviconFile.ReadFile("favicon.ico")
	if err != nil {
		http.Error(w, "Favicon not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	_, _ = w.Write(faviconData)
}

func main() {
	port := 8080
	http.HandleFunc("/", handlerRoot)
	http.HandleFunc("/favicon.ico", handlerFavicon)
	fmt.Printf("Starting server on http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
