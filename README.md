![Numtide Banner Generator](https://banner.numtide.com/banner/numtide/banner-generator.svg?1)

Generates SVG and PNG banners for Numtide's GitHub repositories with automatic metadata fetching.

## Quick Start

### Embed in Your README

Add this to any Numtide repository README:

```markdown
![Banner](https://banner.numtide.com/banner/numtide/your-repo.svg)
```

## Usage

```bash
nix run github:numtide/banner-generator#banner-api
nix run github:numtide/banner-generator#banner-cli
```

## API Endpoints

- `GET /banner/{owner}/{repo}.svg` - Generate SVG banner
- `GET /banner/{owner}/{repo}.png` - Generate PNG banner

## CLI Usage

```bash
# Generate PNG banner
banner-cli generate owner/repo -o banner.png

# Generate without stats (for social preview that won't get stale)
banner-cli generate owner/repo --no-stats -o banner.png
```

After generating, upload the PNG as social preview via:
Repository Settings > Social preview > Edit

## Configuration

See `deploy/banner-generator.toml` for configuration options.

## Template Structure

Templates are pure SVG files with specific element IDs that get replaced dynamically:

| Element ID | Description | Example Content |
|------------|-------------|-----------------|
| `repo-name` | Repository name text | "banner-generator" |
| `description` | Description text (can contain tspan elements) | "Generate banners..." |
| `stats-stars` | Stars count text | "⭐ 1.2k" |
| `stats-forks` | Forks count text | "🍴 45" |
| `stats-language` | Primary language text | "Go" |
| `stats-group` | Stats container (hidden if no data) | - |
| `font-css` | Style element for font injection | - |


## Development

```bash
nix develop

make dev     # Run with hot reload
make build   # Build binaries
make test    # Run tests
make lint    # Run linters
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT
