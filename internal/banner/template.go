package banner

// BannerData holds the data to be injected into the banner template
type BannerData struct {
	RepoName        string
	RepoDescription string
	Owner           string
	Language        string
	Stars           int
	Forks           int
}

// Default dimensions for the banner
const (
	BannerWidth  = 1280
	BannerHeight = 640
)
