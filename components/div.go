package components

import (
	"fmt"
	"strings"
)

type DivComponent struct {
	Children []Component
	Class    string
	Style    Style
	ID       string
	Attrs    Attr
}

// Div is a container that groups multiple components together.
// You can pass components.Class, components.Style, etc., alongside children components.
func Div(args ...any) *DivComponent {
	d := &DivComponent{}
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			if !ParseStringAttr(v, &d.Class, &d.ID, &d.Attrs) {
				// If it's a raw string that's not an attribute, we could treat it as Text Component, but Div usually takes Components.
				d.Children = append(d.Children, HTML(v))
			}
		case Style:
			d.Style = v
		case Class:
			d.Class = string(v)
		case ID:
			d.ID = string(v)
		case Attr:
			d.Attrs = v
		case Component:
			d.Children = append(d.Children, v)
		case []Component:
			d.Children = append(d.Children, v...)
		case []any:
			for _, item := range v {
				if comp, ok := item.(Component); ok {
					d.Children = append(d.Children, comp)
				}
			}
		}
	}
	return d
}

func (d *DivComponent) Render() string {
	var parts []string
	if d.Class != "" {
		parts = append(parts, fmt.Sprintf(`class="%s"`, d.Class))
	}
	if len(d.Style) > 0 {
		parts = append(parts, fmt.Sprintf(`style="%s"`, strings.Join(d.Style.entries(), ";")))
	}
	if d.ID != "" {
		parts = append(parts, fmt.Sprintf(`id="%s"`, d.ID))
	}
	if len(d.Attrs) > 0 {
		for k, v := range d.Attrs {
			parts = append(parts, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}
	attrs := ""
	if len(parts) > 0 {
		attrs = " " + strings.Join(parts, " ")
	}

	var builder strings.Builder
	builder.WriteString("<div" + attrs + ">\n")
	for _, child := range d.Children {
		builder.WriteString("\t" + child.Render() + "\n")
	}
	builder.WriteString("</div>")

	return builder.String()
}
