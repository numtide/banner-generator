package banner

import (
	"github.com/numtide/banner-generator/internal/resources"
)

// LocateTemplate finds a template file using the resource locator
func LocateTemplate(filename string) string {
	locator := resources.NewResourceLocator()
	return locator.FindTemplate(filename)
}
