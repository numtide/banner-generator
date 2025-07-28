package banner

// Builder is the interface for SVG banner builders
type Builder interface {
	BuildSVG(data *BannerData) (string, error)
}
