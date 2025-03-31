package ui

import (
	"net/http"
	"os"
)

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load HTML", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}
