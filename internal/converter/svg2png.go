package converter

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// findChromePath searches for a Chrome/Chromium executable
func findChromePath() string {
	// Check CHROME_PATH environment variable first
	if path := os.Getenv("CHROME_PATH"); path != "" {
		return path
	}

	// List of common Chrome/Chromium executable names
	candidates := []string{
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"chrome",
	}

	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}

	return ""
}

// ColorScheme represents the preferred color scheme for rendering
type ColorScheme string

const (
	ColorSchemeLight ColorScheme = "light"
	ColorSchemeDark  ColorScheme = "dark"
)

// SVGToPNG converts SVG data to PNG format using headless Chrome with light mode
func SVGToPNG(svgData []byte) ([]byte, error) {
	return SVGToPNGWithColorScheme(svgData, ColorSchemeLight)
}

// SVGToPNGWithColorScheme converts SVG data to PNG format with specified color scheme
func SVGToPNGWithColorScheme(svgData []byte, colorScheme ColorScheme) ([]byte, error) {
	log.Printf("Starting SVG to PNG conversion (color scheme: %s)", colorScheme)
	log.Printf("SVG data size: %d bytes", len(svgData))

	// Find Chrome/Chromium executable
	chromePath := findChromePath()
	if chromePath == "" {
		return nil, fmt.Errorf("no Chrome/Chromium executable found. Install chromium or set CHROME_PATH")
	}
	log.Printf("Using browser: %s", chromePath)

	// Write SVG to a temporary file (data URIs have size limits in Chrome)
	tmpFile, err := os.CreateTemp("", "banner-*.svg")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpFileName := tmpFile.Name()
	defer func() { _ = os.Remove(tmpFileName) }()

	if _, err := tmpFile.Write(svgData); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("failed to write SVG to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	fileURL := "file://" + tmpFileName
	log.Printf("Using temp file: %s", tmpFileName)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create chrome instance with logging
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	// Enable verbose logging if DEBUG env var is set
	var ctxOpts []chromedp.ContextOption
	if os.Getenv("DEBUG") != "" {
		ctxOpts = append(ctxOpts, chromedp.WithDebugf(log.Printf))
	}

	ctx, cancel = chromedp.NewContext(allocCtx, ctxOpts...)
	defer cancel()

	log.Println("Chrome instance created")

	var pngData []byte

	// Navigate to SVG and take screenshot
	log.Println("Setting viewport to 1280x640")
	err = chromedp.Run(ctx, chromedp.EmulateViewport(1280, 640))
	if err != nil {
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	log.Printf("Setting color scheme to: %s", colorScheme)
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return emulation.SetEmulatedMedia().
			WithFeatures([]*emulation.MediaFeature{
				{Name: "prefers-color-scheme", Value: string(colorScheme)},
			}).
			Do(ctx)
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to set color scheme: %w", err)
	}

	log.Println("Navigating to SVG file...")
	err = chromedp.Run(ctx, chromedp.Navigate(fileURL))
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to SVG: %w", err)
	}

	log.Println("Taking screenshot...")
	err = chromedp.Run(ctx, chromedp.FullScreenshot(&pngData, 100))
	if err != nil {
		return nil, fmt.Errorf("failed to take screenshot: %w", err)
	}

	log.Printf("Screenshot captured: %d bytes", len(pngData))
	return pngData, nil
}
