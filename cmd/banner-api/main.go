package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/numtide/banner-generator/internal/api"
	"github.com/numtide/banner-generator/internal/banner"
	"github.com/numtide/banner-generator/internal/config"
	"github.com/numtide/banner-generator/internal/fonts"
	"github.com/numtide/banner-generator/internal/github"
)

func main() {
	// Parse command-line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	loader := config.NewConfigLoader()
	appConfig, err := loader.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create configuration for access control
	var allowedEntries []string
	if appConfig.AccessControl.Enabled {
		// Combine orgs and users into allowlist format
		allowedEntries = append(allowedEntries, appConfig.AccessControl.AllowedOrgs...)
		allowedEntries = append(allowedEntries, appConfig.AccessControl.AllowedUsers...)
		log.Printf("Access control enabled with allowlist: %v", allowedEntries)
	}
	cfg := config.NewConfig(allowedEntries)

	// Create font manager from config
	fontManager := fonts.NewManager(appConfig.Fonts.FontsDir)
	log.Printf("Font directory: %s", appConfig.Fonts.FontsDir)

	// Use template path from config
	templatePath := appConfig.TemplatePath
	log.Printf("Using template: %s", templatePath)

	// Determine base URL for web fonts
	fontBaseURL := ""
	if appConfig.Fonts.EnableWebFonts {
		if appConfig.Fonts.WebFontsBaseURL != "" {
			fontBaseURL = appConfig.Fonts.WebFontsBaseURL
		} else {
			// Use local server URL
			fontBaseURL = fmt.Sprintf("http://%s:%d", appConfig.Server.Host, appConfig.Server.Port)
		}
		log.Printf("Web fonts enabled, base URL: %s", fontBaseURL)
	}

	// Initialize components
	svgBuilder := banner.NewSimpleSVGBuilder(fontManager, templatePath, appConfig.Fonts.EnableWebFonts, fontBaseURL)
	log.Printf("Using simple SVG-based banner generation")

	githubClient := github.NewClient(appConfig.GitHub.Token)

	// Create handler
	handler := api.NewHandler(svgBuilder, githubClient, cfg)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	r.HandleFunc("/banner/{owner}/{repo}.svg", handler.GenerateBanner).Methods("GET")
	r.HandleFunc("/banner/{owner}/{repo}.png", handler.GeneratePNGBanner).Methods("GET")
	r.HandleFunc("/", handler.Index).Methods("GET")

	// Serve font files using font manager
	r.PathPrefix("/fonts/").Handler(fontManager)

	// Setup middleware
	r.Use(api.LoggingMiddleware)

	// Create server
	readTimeout, _ := time.ParseDuration(appConfig.Server.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(appConfig.Server.WriteTimeout)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Starting server on %s:%d", appConfig.Server.Host, appConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
