# Font Registry Architecture

The font registry provides a unified system for managing fonts in the banner generator, supporting both web serving and SVG-to-PNG conversion.

## Key Components

### 1. Font Registry (`internal/fonts/registry.go`)
- Central registry for all fonts
- Maps font families to file paths
- Serves fonts via HTTP
- Generates @font-face CSS
- Loads font data for embedding

### 2. Auto Font Builder (`internal/banner/auto_font_builder.go`)
- Automatically detects fonts from SVG templates
- Uses font resolver to extract and resolve font requirements
- Supports two modes:
  - **Embedded fonts**: Base64-encoded in SVG (default)
  - **Web fonts**: External references with @font-face

### 3. Font Resolver (`internal/fonts/resolver.go`)
- Extracts font requirements from SVG content
- Resolves fonts using the registry
- Generates appropriate CSS for both embedded and web modes

## Usage

### API Server

```bash
# Default: Embedded fonts (self-contained SVGs)
./bin/banner-api

# Web fonts mode (smaller SVGs, requires font hosting)
./bin/banner-api -web-fonts

# Custom base URL for fonts
./bin/banner-api -web-fonts -base-url https://cdn.example.com

# Custom font directory
./bin/banner-api -font-dir /path/to/fonts
```

### Font Directory Structure

```
fonts/
├── gt-pressura-regular.ttf      # Main TTF file
└── web/                          # Web-optimized formats
    ├── gt-pressura-regular.woff2
    └── gt-pressura-regular.woff
```

### Font Serving

The registry serves fonts at `/fonts/{family}/{family}.{format}`:
- `/fonts/gt-pressura/gt-pressura.woff2`
- `/fonts/gt-pressura/gt-pressura.woff`
- `/fonts/gt-pressura/gt-pressura.ttf`

## Benefits

1. **Flexibility**: Switch between embedded and web fonts with a flag
2. **Performance**: Web fonts mode reduces SVG size from ~113KB to ~4KB
3. **Caching**: Browser caches fonts, improving load times
4. **Maintainability**: Centralized font management
5. **Extensibility**: Easy to add new fonts to the registry

## Adding New Fonts

To add a new font:

1. Place font files in the fonts directory with appropriate subdirectories:
   - Main TTF file in root
   - Web formats in `web/` subdirectory

2. Create or update `fonts/fonts.toml`:

```toml
[[fonts]]
family = "new-font"
name = "New Font Name"
aliases = ["New Font", "newfont"]

[fonts.variants]
ttf = "new-font.ttf"
woff = "web/new-font.woff"
woff2 = "web/new-font.woff2"
```

3. The font will be automatically available in templates

## File Size Comparison

| Mode | SVG Size | First Load | Subsequent Loads |
|------|----------|------------|------------------|
| Embedded | ~113KB | 113KB | 113KB each |
| Web Fonts | ~4KB | 19KB (4KB + 15KB font) | 4KB (font cached) |
