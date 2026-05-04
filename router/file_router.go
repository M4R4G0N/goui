package router

import (
	goui "github.com/M4R4G0N/goUI"
	"github.com/M4R4G0N/goUI/components"
)

// PageHandler represents a function that returns a UI component,
// receiving its own title and path.
type PageHandler func(title, path string) components.Component

type RouteDef struct {
	Path    string
	Title   string
	Handler PageHandler
}

// globalRoutes holds the registered routes from pages
var globalRoutes []RouteDef

// RegisterPage registers a route globally.
// It should be called from an init() function within each page file.
func RegisterPage(path, title string, handler PageHandler) {
	globalRoutes = append(globalRoutes, RouteDef{
		Path:    path,
		Title:   title,
		Handler: handler,
	})
}

// LayoutFunc defines a wrapper around a page's content
type LayoutFunc func(title, path string, body components.Component) components.Component

// InjectRoutes registers all accumulated global routes into the application.
// You can optionally pass a LayoutFunc to wrap all pages (e.g. with a Navbar).
func InjectRoutes(app *goui.App, layouts ...LayoutFunc) {
	for _, route := range globalRoutes {
		// Evaluates the handler injecting its own title and path
		comp := route.Handler(route.Title, route.Path)

		// Se tem layout configurado, embrulha o componente nele
		if len(layouts) > 0 && layouts[0] != nil {
			comp = layouts[0](route.Title, route.Path, comp)
		}

		app.RegisterRoute(route.Path, comp)
	}
}
