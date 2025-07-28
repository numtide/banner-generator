![Numtide Banner Generator](https://banner.numtide.com/banner/numtide/banner-generator.svg)

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
# Generate PNG
banner-cli generate owner/repo -o banner.png

# Generate and upload to GitHub
banner-cli generate-upload owner/repo --token $GITHUB_TOKEN
```

## Configuration

See `deploy/banner-generator.toml` for configuration options.

## Template Variables

Templates receive the following variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{fontCSS}}` | Auto-generated @font-face CSS | (CSS string) |
| `{{repoName}}` | Repository name | "banner-generator" |
| `{{repoNameFontSize}}` | Calculated font size | "32" |
| `{{owner}}` | Repository owner | "numtide" |
| `{{description}}` | Repository description | "Generate banners..." |
| `{{descriptionLines}}` | Multi-line description array | [{text: "...", y: 100}] |
| `{{language}}` | Primary language | "Go" |
| `{{stars}}` | Formatted star count | "1.2k" |
| `{{forks}}` | Formatted fork count | "45" |
| `{{hasDescription}}` | Has description? | true/false |
| `{{hasLanguage}}` | Has language? | true/false |
| `{{hasStats}}` | Has stats? | true/false |


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
