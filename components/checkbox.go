package components

import (
	"fmt"
	"strings"
)

// CheckboxItem represents one option in a CheckboxGroup.
type CheckboxItem struct {
	Label    string
	Value    string
	Checked  bool
	Disabled bool
}

// CheckboxComponent renders a single styled checkbox with a label.
type CheckboxComponent struct {
	Label   string
	Name    string
	ID      string
	Value   string
	Checked bool
	Class   string
	Attrs   Attr
}

func (c *CheckboxComponent) GetID() string {
	if c.ID == "" {
		c.ID = AutoID()
	}
	return c.ID
}

// Checkbox creates a single styled checkbox with a label.
//
// Pass true as an option to pre-check it. Pass Disabled to disable it.
//
// Example:
//
//	Checkbox("Aceito os termos", true, Name("terms"), Value("1"))
func Checkbox(label string, opts ...any) *CheckboxComponent {
	c := &CheckboxComponent{Label: label}
	for _, opt := range opts {
		switch v := opt.(type) {
		case bool:
			c.Checked = v
		case string:
			ParseStringAttr(v, &c.Class, &c.ID, &c.Attrs)
		case Name:
			c.Name = string(v)
		case ID:
			c.ID = string(v)
		case Value:
			c.Value = string(v)
		case Class:
			c.Class = string(v)
		case Attr:
			if c.Attrs == nil {
				c.Attrs = make(Attr)
			}
			for k, val := range v {
				c.Attrs[k] = val
			}
		}
	}
	return c
}

func (c *CheckboxComponent) Render() string {
	id := c.GetID()
	var attrs []string
	if c.Name != "" {
		attrs = append(attrs, fmt.Sprintf(`name="%s"`, c.Name))
	}
	if c.Value != "" {
		attrs = append(attrs, fmt.Sprintf(`value="%s"`, c.Value))
	}
	if c.Checked {
		attrs = append(attrs, "checked")
	}
	for k, v := range c.Attrs {
		attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
	}
	extra := ""
	if len(attrs) > 0 {
		extra = " " + strings.Join(attrs, " ")
	}
	labelClass := "goui-checkbox-label"
	if c.Class != "" {
		labelClass += " " + c.Class
	}
	return fmt.Sprintf(
		`<label class="%s"><input type="checkbox" id="%s" class="goui-checkbox-input"%s><span class="goui-checkbox-box"></span><span class="goui-checkbox-text">%s</span></label>`,
		labelClass, id, extra, c.Label,
	)
}

// CheckboxGroupComponent renders a labeled group of checkboxes.
type CheckboxGroupComponent struct {
	Label string
	Name  string
	Items []CheckboxItem
	Class string
	ID    string
}

func (g *CheckboxGroupComponent) GetID() string {
	if g.ID == "" {
		g.ID = AutoID()
	}
	return g.ID
}

// CheckboxGroup creates a group of checkboxes under a shared label.
//
// Example:
//
//	CheckboxGroup("Tecnologias", Name("tech"),
//	    CheckboxItem{Label: "Go", Value: "go", Checked: true},
//	    CheckboxItem{Label: "Rust", Value: "rust"},
//	)
func CheckboxGroup(label string, opts ...any) *CheckboxGroupComponent {
	g := &CheckboxGroupComponent{Label: label}
	for _, opt := range opts {
		switch v := opt.(type) {
		case CheckboxItem:
			g.Items = append(g.Items, v)
		case []CheckboxItem:
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

func (g *CheckboxGroupComponent) Render() string {
	var b strings.Builder
	cls := "goui-checkbox-group"
	if g.Class != "" {
		cls += " " + g.Class
	}
	fmt.Fprintf(&b, `<div class="%s" id="%s">`, cls, g.GetID())
	if g.Label != "" {
		fmt.Fprintf(&b, `<span class="goui-label">%s</span>`, g.Label)
	}
	b.WriteString(`<div class="goui-checkbox-group-items">`)
	for _, item := range g.Items {
		id := AutoID()
		var attrs []string
		if g.Name != "" {
			attrs = append(attrs, fmt.Sprintf(`name="%s"`, g.Name))
		}
		attrs = append(attrs, fmt.Sprintf(`value="%s"`, item.Value))
		if item.Checked {
			attrs = append(attrs, "checked")
		}
		if item.Disabled {
			attrs = append(attrs, "disabled")
		}
		fmt.Fprintf(&b,
			`<label class="goui-checkbox-label"><input type="checkbox" id="%s" class="goui-checkbox-input" %s><span class="goui-checkbox-box"></span><span class="goui-checkbox-text">%s</span></label>`,
			id, strings.Join(attrs, " "), item.Label,
		)
	}
	b.WriteString(`</div></div>`)
	return b.String()
}
