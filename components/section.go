package components

import (
	"fmt"
	"strings"
)

type SectionComponent struct {
	Title       string
	Items       []Component
	Class       string
	ID          string
	Style       Style
	Attrs       Attr
	Collapsible bool
	Open        bool
}

// Section creates a grouped area with a title and child components.
func Section(title string, args ...any) *SectionComponent {
	s := &SectionComponent{
		Title: title,
		Items: []Component{},
		Open:  true, // Padrão é começar aberto
	}
	for _, arg := range args {
		if comp, ok := arg.(Component); ok {
			s.Items = append(s.Items, comp)
			continue
		}

		switch v := arg.(type) {
		case bool: // Se passar true, ativa o colapsável
			s.Collapsible = v
		case string:
			ParseStringAttr(v, &s.Class, &s.ID, &s.Attrs)
		case Class:
			s.Class = string(v)
		case ID:
			s.ID = string(v)
		case Style:
			s.Style = v
		case Attr:
			s.Attrs = v
		}
	}
	return s
}

func (s *SectionComponent) Render() string {
	id := s.ID
	if id == "" {
		id = AutoID()
	}

	var body strings.Builder
	for _, item := range s.Items {
		body.WriteString(item.Render())
	}

	styleStr := ""
	if len(s.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, s.Style.Render())
	}

	// Lógica de Colapsável
	cursor := ""
	onclick := ""
	chevron := ""
	contentStyle := ""

	if s.Collapsible {
		cursor = "cursor: pointer; user-select: none;"
		onclick = fmt.Sprintf(` onclick="gouiToggleSection('%s')"`, id)
		chevron = `<span class="goui-section-chevron" id="chevron-` + id + `">▼</span>`
	}

	return fmt.Sprintf(`
		<section class="goui-section-group %s" id="%s"%s>
			<div class="goui-section-header" style="%s"%s>
				<h2 class="goui-section-title" style="margin-bottom:0; border-bottom:none;">%s</h2>
				%s
			</div>
			<div class="goui-section-content" id="content-%s" style="%s">
				%s
			</div>
		</section>
		<script>
			if (typeof window.gouiToggleSection === 'undefined') {
				window.gouiToggleSection = function(id) {
					var content = document.getElementById('content-' + id);
					var chevron = document.getElementById('chevron-' + id);
					if (content.style.display === 'none') {
						content.style.display = 'flex';
						chevron.style.transform = 'rotate(0deg)';
					} else {
						content.style.display = 'none';
						chevron.style.transform = 'rotate(-90deg)';
					}
				};
			}
		</script>
	`, s.Class, id, styleStr, cursor, onclick, s.Title, chevron, id, contentStyle, body.String())
}
