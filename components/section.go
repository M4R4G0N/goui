package components

import (
	"fmt"
)

type SectionComponent struct {
	Title string
	Items []Component
}

// Section creates a grouped area with a title and child components.
func Section(title string, items ...Component) *SectionComponent {
	return &SectionComponent{
		Title: title,
		Items: items,
	}
}

func (s *SectionComponent) Render() string {
	var body string
	for _, item := range s.Items {
		body += item.Render()
	}

	return fmt.Sprintf(`
		<section class="goui-section-group">
			<h2 class="goui-section-title">%s</h2>
			<div class="goui-section-content">
				%s
			</div>
		</section>
	`, s.Title, body)
}
