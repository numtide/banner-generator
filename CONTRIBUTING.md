# Contributing to Banner Generator

Thank you for your interest in contributing to Banner Generator!

## Development Setup

```bash
git clone https://github.com/numtide/banner-generator.git
cd banner-generator
nix develop
```

## Project Structure

```
banner-generator/
├── cmd/                    # Application entry points
│   ├── banner-api/        # HTTP API server
│   └── banner-cli/        # Command-line tool
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and middleware
│   ├── banner/           # Banner generation logic
│   ├── config/           # Configuration management
│   ├── converter/        # SVG to PNG conversion
│   ├── fonts/            # Font management and resolution
│   ├── github/           # GitHub API client
│   └── utils/            # Shared utilities
├── deploy/               # Deployment configuration and assets
│   ├── fonts/           # Font files and configuration
│   ├── templates/       # SVG templates
│   └── *.toml           # Configuration files
└── docs/                 # Documentation

```

## Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Follow Go conventions and idioms
   - Keep functions small and focused
   - Add comments for exported functions
   - Update documentation as needed

3. **Format your code**:
   ```bash
   make fmt
   ```

4. **Run linters**:
   ```bash
   make lint
   ```

5. **Test your changes**:
   ```bash
   make test
   ```


## Configuration

The configuration file requires a `template_path` field:

```toml
# Required - path to your SVG template
template_path = "templates/banner.svg"

[fonts]
enable_web_fonts = false  # Embed fonts as base64 when false
default_family = "GT Pressura"
```

## Adding New Fonts

1. Add font files to the `deploy/fonts/` directory (WOFF format preferred)
2. Update `deploy/fonts/fonts.toml` with the new font configuration:
   ```toml
   [[fonts]]
   family = "your-font-family"
   name = "Your Font Name"
   aliases = ["YourFont", "your font"]
   # WOFF format is preferred over WOFF2
   variants = { ttf = "your-font.ttf", woff = "web/your-font.woff" }
   ```
3. The font will be automatically detected from `font-family` attributes in templates

## Creating New Templates

Templates are pure SVG files - no templating engine required:

1. Create a new `.svg` file in `deploy/templates/`
2. Add IDs to elements that should be dynamic:
   - `id="repo-name"` - Repository name text element
   - `id="description"` - Description text element (with tspan children for multi-line)
   - `id="stats-stars"` - Stars count text element
   - `id="stats-forks"` - Forks count text element
   - `id="stats-language"` - Language text element
   - `id="stats-group"` - Group to hide if no stats
   - `id="font-css"` - Style element where font CSS will be injected
3. Use `font-family` attributes on text elements - fonts will be automatically detected and embedded

### How It Works

- The system uses regex-based string replacement to update SVG elements by ID
- Text content is wrapped in `<tspan>` elements to preserve structure
- Multi-line descriptions maintain their tspan layout with proper x/dy attributes
- Fonts are embedded as base64 data URIs or referenced as URLs based on configuration
- No XML parsing means no namespace duplication issues

## Submitting Changes

1. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   Follow conventional commit format:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `refactor:` for code refactoring
   - `test:` for test additions/changes
   - `chore:` for maintenance tasks

2. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**:
   - Fill in the PR template
   - Link any related issues
   - Ensure CI passes

## Code Style

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Keep line length reasonable (80-120 characters)
- Use meaningful variable and function names
- Avoid deep nesting

## Questions?

Feel free to open an issue for any questions or discussions!