package components

import (
	"fmt"
	"strings"
)

// Option represents a single item in a Dropdown.
type Option struct {
	Value    string
	Label    string
	Selected bool
	Disabled bool
}

// multi is an unexported sentinel type so Multi cannot be confused with other values.
type multi bool

// Multi enables multiselect mode when passed to Dropdown.
func Multi(v bool) multi { return multi(v) }

type DropdownComponent struct {
	Options  []Option
	Multi    bool
	Class    string
	Style    Style
	Name     string
	Size     int // visible rows for multiselect (0 = browser default)
	ID       string
	Attrs    Attr
}

// GetID returns the ID of the dropdown, automatically generating one if empty.
func (d *DropdownComponent) GetID() string {
	if d.ID == "" {
		d.ID = AutoID()
	}
	return d.ID
}

// Dropdown renders a styled select element.
// Pass Multi to allow multiple selections.
// Pass Style{} or Class("") to customise appearance.
// Pass Name("fieldname") to set the form field name.
//
// Examples:
//
//	Dropdown(components.Option{Value: "go", Label: "Go"}, components.Option{Value: "rs", Label: "Rust"})
//	Dropdown(components.Option{...}, Multi)
//	Dropdown(components.Option{...}, Class("my-select"), Name("language"))
func Dropdown(opts ...any) *DropdownComponent {
	d := &DropdownComponent{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			ParseStringAttr(v, &d.Class, &d.ID, &d.Attrs)
		case Option:
			d.Options = append(d.Options, v)
		case []Option:
			d.Options = append(d.Options, v...)
		case multi:
			d.Multi = bool(v)
		case Style:
			d.Style = v
		case Class:
			d.Class = string(v)
		case Name:
			d.Name = string(v)
		case Size:
			d.Size = int(v)
		case ID:
			d.ID = string(v)
		case Attr:
			d.Attrs = v
		}
	}
	return d
}

func (d *DropdownComponent) Render() string {
	var b strings.Builder

	// Wrapper
	wrapperClass := "goui-dropdown"
	if d.Multi {
		wrapperClass += " goui-dropdown-multi"
	}
	fmt.Fprintf(&b, `<div class="%s">`, wrapperClass)

	// <select> attributes
	selectAttrs := `class="goui-select"`
	if d.Name != "" {
		selectAttrs += fmt.Sprintf(` name="%s"`, d.Name)
	}
	if d.Multi {
		selectAttrs += ` multiple`
		if d.Size > 0 {
			selectAttrs += fmt.Sprintf(` size="%d"`, d.Size)
		}
	}
	if d.Class != "" {
		selectAttrs += fmt.Sprintf(` class="%s"`, d.Class) // NOTE: fixed data-class to class for select
	}
	if d.GetID() != "" {
		selectAttrs += fmt.Sprintf(` id="%s"`, d.ID)
	}
	if len(d.Attrs) > 0 {
		for k, v := range d.Attrs {
			selectAttrs += fmt.Sprintf(` %s="%s"`, k, v)
		}
	}

	var styleStr string
	if len(d.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, strings.Join(d.Style.entries(), ";"))
	}

	fmt.Fprintf(&b, `<select %s%s>`, selectAttrs, styleStr)

	// Options
	for _, opt := range d.Options {
		attrs := fmt.Sprintf(`value="%s"`, opt.Value)
		if opt.Selected {
			attrs += ` selected`
		}
		if opt.Disabled {
			attrs += ` disabled`
		}
		fmt.Fprintf(&b, `<option %s>%s</option>`, attrs, opt.Label)
	}

	b.WriteString(`</select>`)

	// Custom arrow icon for single select (hidden for multi)
	if !d.Multi {
		b.WriteString(`<span class="goui-dropdown-arrow" aria-hidden="true">`)
		b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>`)
		b.WriteString(`</span>`)
	}

	b.WriteString(`</div>`)
	return b.String()
}

// Name sets the HTML name attribute (for form submission).
type Name string

// Size sets the number of visible rows in a multiselect.
type Size int
