package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

// NewGenerator creates a new PNG generator with default config
func NewGenerator(githubToken string) (*Generator, error) {
	// Create font registry with default font directory
	fontRegistry := fonts.DefaultRegistryWithConfig(config.DefaultFontDir)

	// Use default template with auto font builder
	templatePath := banner.LocateTemplate("banner.svg.mustache")
	svgBuilder := banner.NewAutoFontBuilder(fontRegistry, templatePath, "")

	return &Generator{
		svgBuilder:   svgBuilder,
		githubClient: github.NewClient(githubToken),
	}, nil
}

// NewGeneratorWithConfig creates a new PNG generator with provided config
func NewGeneratorWithConfig(appConfig *config.AppConfig) (*Generator, error) {
	// Create font registry from config
	fontRegistry := createFontRegistry(appConfig)

	// Determine template path
	templateName := appConfig.Templates.DefaultTemplate
	if template, ok := appConfig.Templates.Templates[templateName]; ok {
		templatePath := filepath.Join(appConfig.Templates.TemplatesDir, template.Path)

		// Determine base URL for web fonts
		fontBaseURL := ""
		if appConfig.Fonts.EnableWebFonts {
			fontBaseURL = appConfig.Fonts.WebFontsBaseURL
		}

		svgBuilder := banner.NewAutoFontBuilder(fontRegistry, templatePath, fontBaseURL)

		return &Generator{
			svgBuilder:   svgBuilder,
			githubClient: github.NewClient(appConfig.GitHub.Token),
		}, nil
	}

	return nil, fmt.Errorf("template '%s' not found in configuration", templateName)
}

// createFontRegistry creates a font registry from the configuration
func createFontRegistry(appConfig *config.AppConfig) *fonts.Registry {
	// Always use the existing font system which loads from fonts.toml
	// The config loader already merges fonts.toml into the app config
	return fonts.DefaultRegistryWithConfig(appConfig.Fonts.FontsDir)
}

// GeneratePNG generates a PNG banner for the specified repository
func (g *Generator) GeneratePNG(repoPath, outputPath string) error {
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

	fmt.Printf("Generating banner for: %s\n", repoData.RepoName)
	if repoData.RepoDescription != "" {
		fmt.Printf("Description: %s\n", repoData.RepoDescription)
	}

	// Generate SVG
	svg, err := g.svgBuilder.BuildSVG(repoData)
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
