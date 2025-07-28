# Contributing to Banner Generator

Thank you for your interest in contributing to Banner Generator!

## Development Setup

1. **Prerequisites**:
   - Go 1.21 or later
   - Git
   - One of: rsvg-convert, ImageMagick, or Inkscape (for PNG conversion)

2. **Clone the repository**:
   ```bash
   git clone https://github.com/numtide/banner-generator.git
   cd banner-generator
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Build the project**:
   ```bash
   make build
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
├── templates/            # SVG Mustache templates
├── fonts/                # Font files and configuration
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
   ./example.sh  # Run the example script
   ```

## Adding New Fonts

1. Add font files to the `fonts/` directory
2. Update `fonts/fonts.toml` with the new font configuration:
   ```toml
   [[fonts]]
   family = "your-font-family"
   name = "Your Font Name"
   aliases = ["YourFont", "your font"]
   variants = { ttf = "your-font.ttf", woff = "web/your-font.woff", woff2 = "web/your-font.woff2" }
   ```
3. Test that the font works in templates

## Creating New Templates

1. Create a new `.svg.mustache` file in `templates/`
2. Use `font-family` attributes on text elements - fonts will be automatically detected
3. Available template variables:
   - `{{repoName}}` - Repository name (without owner)
   - `{{repoNameFontSize}}` - Calculated font size
   - `{{description}}` - Repository description
   - `{{descriptionLines}}` - Multi-line description array
   - `{{language}}` - Primary programming language
   - `{{hasDescription}}` - Boolean for conditionals
   - `{{hasLanguage}}` - Boolean for conditionals

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