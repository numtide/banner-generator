package converter

import (
	"bytes"
	"fmt"
	"os/exec"
)

// SVGToPNG converts SVG data to PNG format
// This uses rsvg-convert if available, or falls back to other methods
func SVGToPNG(svgData []byte) ([]byte, error) {
	// Try rsvg-convert first (most reliable for complex SVGs)
	if pngData, err := convertWithRSVG(svgData); err == nil {
		return pngData, nil
	}

	// Try ImageMagick convert
	if pngData, err := convertWithImageMagick(svgData); err == nil {
		return pngData, nil
	}

	// Try Inkscape
	if pngData, err := convertWithInkscape(svgData); err == nil {
		return pngData, nil
	}

	return nil, fmt.Errorf("no SVG to PNG converter available. Please install rsvg-convert, ImageMagick, or Inkscape")
}

func convertWithRSVG(svgData []byte) ([]byte, error) {
	cmd := exec.Command("rsvg-convert", "-f", "png", "-w", "1280", "-h", "640")
	cmd.Stdin = bytes.NewReader(svgData)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("rsvg-convert failed: %w", err)
	}

	return out.Bytes(), nil
}

func convertWithImageMagick(svgData []byte) ([]byte, error) {
	cmd := exec.Command("convert", "-background", "none", "-density", "96", "-resize", "1280x640", "svg:-", "png:-")
	cmd.Stdin = bytes.NewReader(svgData)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ImageMagick convert failed: %w", err)
	}

	return out.Bytes(), nil
}

func convertWithInkscape(svgData []byte) ([]byte, error) {
	// Inkscape requires file input, so we'll use process substitution
	cmd := exec.Command("inkscape", "--pipe", "--export-type=png", "--export-width=1280", "--export-height=640")
	cmd.Stdin = bytes.NewReader(svgData)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("inkscape failed: %w", err)
	}

	return out.Bytes(), nil
}
