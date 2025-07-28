# Banner Generator

A Go-based tool that generates beautiful SVG and PNG banners for GitHub repositories with automatic font detection, customizable templates, and GitHub integration.

## Features

- **Dynamic SVG/PNG Generation**: Create banners in both SVG (for web embedding) and PNG (for social media) formats
- **GitHub Integration**: Automatically fetches repository metadata including stars, forks, and language
- **Smart Font Management**: Automatic font detection from templates with configurable font registry
- **Light/Dark Mode Support**: Banners automatically adapt to the viewer's color scheme preference
- **Mustache Templates**: Fully customizable banner designs using Mustache templating
- **Web API & CLI**: Available as both HTTP API server and command-line tool
- **Access Control**: Optional allowlist for restricting banner generation to specific repositories
- **Performance**: Built-in caching and rate limiting for efficient operation
- **Multiple Deployment Options**: Run standalone, with Docker, or via Nix

## Quick Start

### Embed in Your README

Add this to your README.md:

```markdown
![Banner](https://your-banner-server.com/banner/owner/repo.svg)
```

### CLI Examples

Generate a PNG banner:

```bash
banner-cli generate owner/repo -o banner.png
```

Generate and upload to GitHub:

```bash
banner-cli generate-upload owner/repo --token=$GITHUB_TOKEN
```

## Installation

### Using Go Install

```bash
go install github.com/numtide/banner-generator/cmd/banner-api@latest
go install github.com/numtide/banner-generator/cmd/banner-cli@latest
```

### Using Nix

```bash
nix run github:numtide/banner-generator#banner-api
nix run github:numtide/banner-generator#banner-cli
```

### From Source

```bash
# Clone the repository
git clone https://github.com/numtide/banner-generator
cd banner-generator

# Install dependencies
go mod download

# Build binaries
make build

# Binaries will be in ./bin/
```

### Prerequisites

For PNG generation, install one of these tools:
- `rsvg-convert` (recommended) - Part of librsvg
- ImageMagick (`convert` command)
- Inkscape

## Usage

### API Server

Start the server with default settings:

```bash
banner-api
```

Start with custom configuration:

```bash
banner-api \
  -port 8080 \
  -font-dir ./fonts \
  -template ./templates/banner.svg.mustache \
  -github-token $GITHUB_TOKEN \
  -allow "numtide,nixos" \
  -web-fonts \
  -base-url https://your-domain.com
```

#### API Endpoints

- `GET /banner/{owner}/{repo}.svg` - Generate SVG banner
- `GET /banner/{owner}/{repo}.png` - Generate PNG banner  
- `GET /health` - Health check endpoint
- `GET /` - Landing page with usage examples
- `GET /fonts/*` - Serve web font files (when using web fonts)

### CLI Tool

The CLI provides three main commands:

#### Generate PNG

```bash
# Generate banner.png in current directory
banner-cli generate owner/repo

# Specify output path
banner-cli generate owner/repo -o /path/to/banner.png

# Use custom font
banner-cli generate owner/repo --font ./custom-font.ttf
```

#### Upload to GitHub

```bash
# Upload existing PNG to repository
banner-cli upload owner/repo ./banner.png --token $GITHUB_TOKEN
```

#### Generate and Upload

```bash
# Generate and upload with defaults
banner-cli generate-upload owner/repo --token $GITHUB_TOKEN

# Customize upload settings
banner-cli generate-upload owner/repo \
  --token $GITHUB_TOKEN \
  --branch main \
  --path assets/banner.png \
  --message "Update repository banner"
```

## Configuration

### Environment Variables

You can configure the banner generator using environment variables. Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Available variables:

- `GITHUB_TOKEN` - GitHub API token for authenticated requests (increases rate limits)
- `PORT` - API server port (default: 8080)
- `FONT_DIR` - Directory containing font files (default: fonts)
- `WEB_FONTS` - Use web fonts instead of embedded fonts (true/false)
- `BASE_URL` - Base URL for web fonts when WEB_FONTS=true
- `ALLOW_LIST` - Comma-separated allowlist of orgs/repos

### API Server Options

| Flag | Description | Default |
|------|-------------|---------|
| `-port` | Port to listen on | 8080 |
| `-font-dir` | Directory containing font files | fonts |
| `-template` | Path to SVG mustache template | banner.svg.mustache |
| `-github-token` | GitHub API token | $GITHUB_TOKEN |
| `-allow` | Comma-separated allowlist | (none) |
| `-web-fonts` | Use web fonts instead of embedded | false |
| `-base-url` | Base URL for web fonts | http://localhost:PORT |

### CLI Options

| Flag | Description | Default |
|------|-------------|---------|
| `--token` | GitHub API token | $GITHUB_TOKEN |
| `--font` | Path to font file | fonts/gt-pressura-regular.ttf |
| `--output`, `-o` | Output path for generated files | banner.png |
| `--branch` | Target branch (generate-upload) | main |
| `--path` | Repository path (generate-upload) | assets/banner.png |
| `--message` | Commit message (generate-upload) | Update repository banner |

### Access Control

Restrict access to specific organizations or repositories:

```bash
# Allow only specific organizations
banner-api -allow "myorg,mycompany"

# Allow specific repos and orgs
banner-api -allow "myorg,torvalds/linux,microsoft/vscode"

# Mix of users and repos
banner-api -allow "johndoe,janedoe/project,acme-corp"
```

Allowlist formats:
- `org` or `username` - Allow all repositories from that organization/user
- `owner/repo` - Allow only that specific repository

## Customization

### Templates

Create custom banner designs using Mustache templates:

```bash
banner-api -template /path/to/custom-template.svg.mustache
```

Templates receive these variables:

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

Example template snippet:

```svg
<text font-family="gt-pressura" font-size="{{repoNameFontSize}}">
  {{repoName}}
</text>
{{#hasStats}}
<text>⭐ {{stars}} 🍴 {{forks}}</text>
{{/hasStats}}
```

### Font Management

Fonts are configured in `fonts/fonts.toml`:

```toml
[[fonts]]
family = "gt-pressura"
name = "GT Pressura Regular"
aliases = ["GT Pressura", "gt pressura", "GTpressura"]
variants = { 
  ttf = "gt-pressura-regular.ttf", 
  woff = "web/gt-pressura-regular.woff", 
  woff2 = "web/gt-pressura-regular.woff2" 
}
```

The system automatically:
- Detects fonts used in templates via `font-family` attributes
- Resolves fonts using family names and aliases
- Generates @font-face CSS
- Embeds font data or serves via web fonts

## Development

### Setup

```bash
# Clone repository
git clone https://github.com/numtide/banner-generator
cd banner-generator

# Enter development environment (with Nix)
nix develop

# Install dependencies
go mod download
```

### Common Tasks

```bash
# Run API server with hot reload
make dev

# Build binaries
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linters
make lint

# Build release binaries
make release
```

### Project Structure

```
banner-generator/
├── cmd/                    # Application entry points
│   ├── banner-api/        # HTTP API server
│   └── banner-cli/        # Command-line tool
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and middleware
│   ├── banner/           # Banner generation logic
│   ├── cli/              # CLI commands implementation
│   ├── config/           # Configuration management
│   ├── converter/        # SVG to PNG conversion
│   ├── fonts/            # Font management and registry
│   ├── github/           # GitHub API client
│   └── utils/            # Shared utilities
├── templates/            # SVG Mustache templates
├── fonts/                # Font files and configuration
│   ├── fonts.toml       # Font registry configuration
│   └── web/             # Web font variants
└── docs/                 # Additional documentation
```

## GitHub API and Rate Limits

### Authentication

While the banner generator works without authentication for public repositories, providing a GitHub token is recommended:

- **Higher rate limits**: 60 requests/hour (unauthenticated) vs 5,000 requests/hour (authenticated)
- **Access to private repositories**: Generate banners for your private repos
- **More reliable service**: Avoid rate limit errors during high usage

Get a token from [GitHub Settings → Personal Access Tokens](https://github.com/settings/tokens).

### Caching

The API server implements intelligent caching to minimize GitHub API calls:

- Repository metadata is cached for optimal performance
- Cache automatically invalidates when repository data changes
- Reduces load on GitHub API and improves response times

## Deployment

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o banner-api ./cmd/banner-api

FROM alpine:latest
RUN apk add --no-cache librsvg
COPY --from=builder /app/banner-api /banner-api
COPY fonts /fonts
COPY templates /templates
EXPOSE 8080
CMD ["/banner-api"]
```

### Systemd Service

```ini
[Unit]
Description=Banner Generator API
After=network.target

[Service]
Type=simple
User=banner
ExecStart=/usr/local/bin/banner-api -port 8080
Restart=on-failure
Environment="GITHUB_TOKEN=your-token-here"

[Install]
WantedBy=multi-user.target
```

### Reverse Proxy (Nginx)

```nginx
server {
    listen 443 ssl http2;
    server_name banners.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_cache_valid 200 1h;
    }
}
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Development setup
- Code style and conventions
- Adding new fonts and templates
- Submitting pull requests

## License

MIT