package goui

import (
	"fmt"
	"net/http"

	"github.com/M4R4G0N/goUI/components"
)

// App configures the main application
type App struct {
	Routes map[string]components.Component
}

// NewApp creates a new goui application
func NewApp() *App {
	return &App{
		Routes: make(map[string]components.Component),
	}
}

// RegisterRoute adds a manual route.
// In a real file-based routing, this is populated automatically by the router/CLI.
func (a *App) RegisterRoute(path string, c components.Component) {
	a.Routes[path] = c
}

// RegisterHandler allows registering a custom http.HandlerFunc for dynamic endpoints
// (e.g. HTMX API routes that return partial HTML).
func (a *App) RegisterHandler(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
}

// Start boots the web server
func (a *App) Start(ip, port string) error {
	for path, comp := range a.Routes {
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, comp.Render())
		})
	}

	addr := ip + ":" + port
	displayAddr := addr
	if ip == "" {
		displayAddr = "localhost:" + port
	}

	fmt.Printf("\n  %sgoUI v0.1.0%s\n", "\033[1;34m", "\033[0m")
	fmt.Printf("  %s➜%s  %sLocal:%s   %shttp://%s/%s\n", "\033[34m", "\033[0m", "\033[1m", "\033[0m", "\033[36m", displayAddr, "\033[0m")
	fmt.Printf("  %s➜%s  %sNetwork:%s %suse --host to expose%s\n\n", "\033[34m", "\033[0m", "\033[1m", "\033[0m", "\033[90m", "\033[0m")

	return http.ListenAndServe(addr, nil)
}
