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
	var httpCacheDurationFlag string
	var apiCacheDurationFlag string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.StringVar(&httpCacheDurationFlag, "http-cache", "", "HTTP cache duration (e.g., '1h', '30m', '300s')")
	flag.StringVar(&apiCacheDurationFlag, "api-cache", "", "API cache duration (e.g., '1h', '30m', '300s')")
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

	// Override cache durations from command-line flags if provided
	if httpCacheDurationFlag != "" {
		appConfig.Cache.HTTPCacheDuration = httpCacheDurationFlag
	}
	if apiCacheDurationFlag != "" {
		appConfig.Cache.APICacheDuration = apiCacheDurationFlag
	}

	// Parse cache durations
	httpCacheDuration, err := time.ParseDuration(appConfig.Cache.HTTPCacheDuration)
	if err != nil {
		log.Printf("Invalid HTTP cache duration '%s', using default 1h: %v", appConfig.Cache.HTTPCacheDuration, err)
		httpCacheDuration = 1 * time.Hour
	}
	cfg.HTTPCacheDuration = httpCacheDuration

	apiCacheDuration, err := time.ParseDuration(appConfig.Cache.APICacheDuration)
	if err != nil {
		log.Printf("Invalid API cache duration '%s', using default 1h: %v", appConfig.Cache.APICacheDuration, err)
		apiCacheDuration = 1 * time.Hour
	}
	cfg.APICacheDuration = apiCacheDuration

	log.Printf("Cache configuration: HTTP=%v, API=%v", httpCacheDuration, apiCacheDuration)

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

	githubClient := github.NewClient(appConfig.GitHub.Token, cfg.APICacheDuration)

	// Create handler
	handler := api.NewHandler(svgBuilder, githubClient, cfg)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	r.HandleFunc("/banner/{owner}/{repo}.svg", handler.GenerateBanner).Methods("GET")
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
