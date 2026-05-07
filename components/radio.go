package components

import (
	"fmt"
	"strings"
)

// RadioItem represents one option in a RadioGroup.
type RadioItem struct {
	Label    string
	Value    string
	Checked  bool
	Disabled bool
}

// RadioGroupComponent renders a labeled group of mutually-exclusive radio buttons.
type RadioGroupComponent struct {
	Label string
	Name  string
	Items []RadioItem
	Class string
	ID    string
}

func (r *RadioGroupComponent) GetID() string {
	if r.ID == "" {
		r.ID = AutoID()
	}
	return r.ID
}

// RadioGroup creates a labeled group of radio buttons.
//
// Example:
//
//	RadioGroup("Plano", Name("plan"),
//	    RadioItem{Label: "Gratuito", Value: "free", Checked: true},
//	    RadioItem{Label: "Pro", Value: "pro"},
//	)
func RadioGroup(label string, opts ...any) *RadioGroupComponent {
	g := &RadioGroupComponent{Label: label}
	for _, opt := range opts {
		switch v := opt.(type) {
		case RadioItem:
			g.Items = append(g.Items, v)
		case []RadioItem:
			g.Items = append(g.Items, v...)
		case Name:
			g.Name = string(v)
		case ID:
			g.ID = string(v)
		case Class:
			g.Class = string(v)
		}
	}
	return g
}

func (r *RadioGroupComponent) Render() string {
	var b strings.Builder
	cls := "goui-radio-group"
	if r.Class != "" {
		cls += " " + r.Class
	}
	fmt.Fprintf(&b, `<div class="%s" id="%s">`, cls, r.GetID())
	if r.Label != "" {
		fmt.Fprintf(&b, `<span class="goui-label">%s</span>`, r.Label)
	}
	b.WriteString(`<div class="goui-radio-group-items">`)
	for _, item := range r.Items {
		id := AutoID()
		var attrs []string
		if r.Name != "" {
			attrs = append(attrs, fmt.Sprintf(`name="%s"`, r.Name))
		}
		attrs = append(attrs, fmt.Sprintf(`value="%s"`, item.Value))
		if item.Checked {
			attrs = append(attrs, "checked")
		}
		if item.Disabled {
			attrs = append(attrs, "disabled")
		}
		fmt.Fprintf(&b,
			`<label class="goui-radio-label"><input type="radio" id="%s" class="goui-radio-input" %s><span class="goui-radio-circle"></span><span class="goui-radio-text">%s</span></label>`,
			id, strings.Join(attrs, " "), item.Label,
		)
	}
	b.WriteString(`</div></div>`)
	return b.String()
}
