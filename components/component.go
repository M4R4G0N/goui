package components

import (
	"fmt"
	"sync/atomic"
)

var globalIDCounter int64

// AutoID generates a unique HTML ID string for a component if it doesn't have one.
func AutoID() string {
	return fmt.Sprintf("goui-comp-%d", atomic.AddInt64(&globalIDCounter, 1))
}

// SyncSource is implemented by any component that can act as a reactive source for SyncWith.
type SyncSource interface {
	GetID() string
}

// Component represents any UI element that can be rendered to HTML.
type Component interface {
	Render() string
}

// headMetaRenderer is satisfied by HeadbarComponent, which needs to inject
// both <head> meta content and a visible top bar.
type headMetaRenderer interface {
	RenderHead() string
}

// HTML is a generic wrapper for raw HTML strings.
type HTML string

func (h HTML) Render() string { return string(h) }

// PageLayout controls how wide the main content renders.
type PageLayout int

const (
	LayoutFull     PageLayout = iota // fills available width (default)
	LayoutCentered                   // constrained to max-width and centered
)

// Page represents a full HTML document with a top header, sidebar, and main content.
type Page struct {
	Head   Component
	Nav    Component
	Body   Component
	Layout PageLayout
}

func NewPage(head, nav, body Component, opts ...any) *Page {
	p := &Page{Head: head, Nav: nav, Body: body}
	for _, opt := range opts {
		if l, ok := opt.(PageLayout); ok {
			p.Layout = l
		}
	}
	return p
}

func (p *Page) Render() string {
	// Separate <head> meta from the visible header bar.
	headMeta := p.Head.Render()
	headerBar := ""
	if hr, ok := p.Head.(headMetaRenderer); ok {
		headMeta = hr.RenderHead()
		headerBar = p.Head.Render()
	}

	bodyHTML := p.Body.Render()
	if p.Layout == LayoutCentered {
		bodyHTML = `<div class="goui-body-centered">` + bodyHTML + `</div>`
	}

	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    ` + headMeta + `
</head>
<body>
    ` + headerBar + `
    <div class="goui-layout">
        <aside class="goui-sidebar">` + p.Nav.Render() + `</aside>
        <main class="goui-main">` + bodyHTML + `</main>
    </div>
</body>
</html>`
}
