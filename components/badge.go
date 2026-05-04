package components

type BadgeComponent struct {
	Content string
}

// Badge renders a small pill label using the goui-badge style.
func Badge(content string) *BadgeComponent {
	return &BadgeComponent{Content: content}
}

func (b *BadgeComponent) Render() string {
	return `<span class="goui-badge">` + b.Content + `</span>`
}
