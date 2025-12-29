package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
)

//go:embed static/*
var staticFiles embed.FS

func (s *Server) registerGUIRoutes(mux *http.ServeMux) {
	// Serve static files
	staticFS, _ := fs.Sub(staticFiles, "static")
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/", fileServer)

	// API endpoints
	mux.HandleFunc("/api/config", s.handleConfig)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	// Security check
	currentConfig := loadConfig()
	isLocked := currentConfig.Token != ""

	// If locked, require valid token in query or header
	if isLocked {
		authHeader := r.Header.Get("Authorization")
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		if token != currentConfig.Token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	switch r.Method {
	case "GET":
		// Mask secrets if needed, but for local config GUI we might want to show them or leave empty
		// For now, let's return them so the user can see what's set
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentConfig)

	case "POST":
		var config Config
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Security: If locked, users can only update if they passed the Auth check above.
		// (which is already enforced by the middleware logic at the start of the function)

		if err := saveConfig(config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		go func() {
			os.Exit(0)
		}()

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
