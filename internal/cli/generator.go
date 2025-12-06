package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/numtide/banner-generator/internal/banner"
	"github.com/numtide/banner-generator/internal/config"
	"github.com/numtide/banner-generator/internal/converter"
	"github.com/numtide/banner-generator/internal/fonts"
	"github.com/numtide/banner-generator/internal/github"
)

// Generator handles PNG banner generation
type Generator struct {
	svgBuilder   banner.Builder
	githubClient *github.Client
}

// NewGeneratorWithConfig creates a new PNG generator with provided config
func NewGeneratorWithConfig(appConfig *config.AppConfig) (*Generator, error) {
	// Create font manager from config
	fontManager := fonts.NewManager(appConfig.Fonts.FontsDir)

	// Use template path from config
	templatePath := appConfig.TemplatePath

	svgBuilder := banner.NewSimpleSVGBuilder(
		fontManager,
		templatePath,
		appConfig.Fonts.EnableWebFonts,
		appConfig.Fonts.WebFontsBaseURL,
	)

	return &Generator{
		svgBuilder:   svgBuilder,
		githubClient: github.NewClient(appConfig.GitHub.Token, 1*time.Hour), // Use 1 hour cache for CLI
	}, nil
}

// GeneratePNG generates a PNG banner for the specified repository
func (g *Generator) GeneratePNG(repoPath, outputPath string, noStats bool) error {
	// Parse owner/repo format
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format, expected owner/repo")
	}
	owner, repo := parts[0], parts[1]

	// Set default output path
	if outputPath == "" {
		outputPath = "banner.png"
	}

	fmt.Printf("Fetching repository data for %s/%s...\n", owner, repo)

	// Fetch repository data
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repoData, err := g.githubClient.GetRepositoryData(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to fetch repository data: %w", err)
	}

	// Clear stats if --no-stats flag is set
	if noStats {
		repoData.StargazersCount = 0
		repoData.ForksCount = 0
		repoData.Language = ""
	}

	fmt.Printf("Generating banner for: %s\n", repoData.Name)
	if repoData.Description != "" {
		fmt.Printf("Description: %s\n", repoData.Description)
	}

	// Generate SVG
	svg, err := g.svgBuilder.BuildBanner(repoData)
	if err != nil {
		return fmt.Errorf("failed to generate SVG: %w", err)
	}

	// Convert SVG to PNG
	fmt.Println("Converting SVG to PNG...")
	pngData, err := converter.SVGToPNG([]byte(svg))
	if err != nil {
		return fmt.Errorf("failed to convert SVG to PNG: %w", err)
	}

	// Save PNG file
	if err := os.WriteFile(outputPath, pngData, 0644); err != nil {
		return fmt.Errorf("failed to save PNG file: %w", err)
	}

	fmt.Printf("Banner saved to: %s\n", outputPath)
	return nil
}
