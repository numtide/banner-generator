package banner

import "github.com/numtide/banner-generator/internal/github"

// Builder is the interface for SVG banner builders
type Builder interface {
	BuildBanner(repo *github.Repository) (string, error)
}
