package components

type IconComponent struct {
	Name string // E.g., "home", "user", "settings"
	Size int
}

func Icon(name string) *IconComponent {
	return &IconComponent{Name: name, Size: 24}
}

// Render returns a generic SVG placeholder for demo purposes, 
// but in a real library it would map 'name' to an actual SVG path.
func (i *IconComponent) Render() string {
	// A simple SVG star as a placeholder for any icon
	return `<svg class="goui-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
		<polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
		<text x="12" y="24" font-size="6" text-anchor="middle" fill="black" stroke="none">` + i.Name + `</text>
	</svg>`
}
