package version

// Version information set via ldflags at build time
var (
	// Version is the semantic version (e.g., "1.2.3")
	Version = "dev"

	// Commit is the git commit hash
	Commit = "unknown"
)
