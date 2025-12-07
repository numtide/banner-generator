package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/numtide/banner-generator/internal/banner"
	"github.com/numtide/banner-generator/internal/config"
	"github.com/numtide/banner-generator/internal/github"
	"github.com/numtide/banner-generator/internal/version"
)

//go:embed index.html
var indexHTML []byte

// Handler handles HTTP requests
type Handler struct {
	svgBuilder   banner.Builder
	githubClient *github.Client
	config       *config.Config
}

// NewHandler creates a new API handler
func NewHandler(svgBuilder banner.Builder, githubClient *github.Client, cfg *config.Config) *Handler {
	return &Handler{
		svgBuilder:   svgBuilder,
		githubClient: githubClient,
		config:       cfg,
	}
}

// HealthCheck returns the health status of the service
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"version": version.Version,
		"commit":  version.Commit,
		"time":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but don't write to response as headers are already sent
		log.Printf("Failed to encode health check response: %v", err)
	}
}

// GenerateBanner generates an SVG banner for a GitHub repository
func (h *Handler) GenerateBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]

	if owner == "" || repo == "" {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}

	// Check access control
	if !h.config.IsAllowed(owner, repo) {
		http.Error(w, "Access denied: This repository is not allowed", http.StatusForbidden)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Fetch repository data
	repoData, err := h.githubClient.GetRepositoryData(ctx, owner, repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch repository data: %v", err), http.StatusNotFound)
		return
	}

	// Generate SVG
	svg, err := h.svgBuilder.BuildBanner(repoData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate banner: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(h.config.HTTPCacheDuration.Seconds())))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write SVG
	if _, err := w.Write([]byte(svg)); err != nil {
		// Log error but can't send error response as headers are already sent
		log.Printf("Failed to write SVG response: %v", err)
	}
}

// Index returns a simple landing page
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(indexHTML); err != nil {
		log.Printf("Failed to write HTML response: %v", err)
	}
}
