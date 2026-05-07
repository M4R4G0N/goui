package components

import "strings"

// NavItem is implemented by anything that can appear in the sidebar nav.
type NavItem interface {
	renderNavItem() string
}

// Link is a plain sidebar navigation link.
type Link struct {
	Href string
	Text string
}

func (l Link) renderNavItem() string {
	return `<a class="goui-sidebar-link" href="` + l.Href + `">` + l.Text + `</a>`
}

// NavGroupComponent is a collapsible group of NavItems in the sidebar.
// It supports nesting: pass another NavGroup as an item to create sub-groups.
type NavGroupComponent struct {
	Label string
	Items []NavItem
	Open  bool
}

// NavGroup creates a collapsible group. Items can be Link or another NavGroup.
func NavGroup(label string, items ...NavItem) NavGroupComponent {
	return NavGroupComponent{Label: label, Items: items}
}

func (g NavGroupComponent) renderNavItem() string {
	var b strings.Builder

	expanded := "false"
	openClass := ""
	if g.Open {
		expanded = "true"
		openClass = " open"
	}

	b.WriteString(`<div class="goui-nav-group">`)
	b.WriteString(`<button class="goui-nav-group-btn" onclick="gouiToggleGroup(this)" aria-expanded="`)
	b.WriteString(expanded)
	b.WriteString(`"><span>`)
	b.WriteString(g.Label)
	b.WriteString(`</span>`)
	b.WriteString(`<svg class="goui-nav-group-arrow" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>`)
	b.WriteString(`</button>`)

	b.WriteString(`<div class="goui-nav-group-items` + openClass + `">`)
	b.WriteString(`<div class="goui-nav-group-inner">`)
	for _, item := range g.Items {
		switch v := item.(type) {
		case Link:
			// Child link — use child style
			b.WriteString(`<a class="goui-sidebar-link goui-sidebar-link-child" href="` + v.Href + `">` + v.Text + `</a>`)
		default:
			// Nested NavGroup — render recursively
			b.WriteString(item.renderNavItem())
		}
	}
	b.WriteString(`</div>`)
	b.WriteString(`</div>`)
	b.WriteString(`</div>`)
	return b.String()
}

// ── Navbar ────────────────────────────────────────────────────────────────────

type NavbarComponent struct {
	Logo  string
	Items []NavItem
}

func Navbar(logo string, items ...NavItem) *NavbarComponent {
	return &NavbarComponent{Logo: logo, Items: items}
}

const iconSun = `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>`
const iconMoon = `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>`

func (n *NavbarComponent) Render() string {
	var b strings.Builder

	b.WriteString(`<nav class="goui-sidebar-nav">`)
	if n.Logo != "" {
		b.WriteString(`<div class="goui-sidebar-title">` + n.Logo + `</div>`)
	}
	for _, item := range n.Items {
		b.WriteString(item.renderNavItem())
	}
	b.WriteString(`</nav>`)

	b.WriteString(`<div class="goui-sidebar-footer">`)
	b.WriteString(`<button class="goui-theme-toggle" onclick="gouiToggleTheme()" title="Toggle theme" aria-label="Toggle theme">`)
	b.WriteString(`<span class="goui-icon-sun">` + iconSun + `</span>`)
	b.WriteString(`<span class="goui-icon-moon">` + iconMoon + `</span>`)
	b.WriteString(`<span class="goui-theme-label-sun">Light mode</span>`)
	b.WriteString(`<span class="goui-theme-label-moon">Dark mode</span>`)
	b.WriteString(`</button>`)
	b.WriteString(`</div>`)

	return b.String()
}
