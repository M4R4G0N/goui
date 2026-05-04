package components

import (
	"fmt"
	"strings"
)

type ToggleComponent struct {
	ID      string
	Checked bool
	Class   string
	Style   Style
	Attrs   Attr
}

func (t *ToggleComponent) GetID() string {
	if t.ID == "" {
		t.ID = AutoID()
	}
	return t.ID
}

// Toggle creates a switch-like checkbox.
func Toggle(opts ...any) *ToggleComponent {
	t := &ToggleComponent{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case bool:
			t.Checked = v
		case string:
			if !ParseStringAttr(v, &t.Class, &t.ID, &t.Attrs) {
				if v == "checked" {
					t.Checked = true
				}
			}
		case Style:
			t.Style = v
		case Class:
			t.Class = string(v)
		case ID:
			t.ID = string(v)
		case Attr:
			t.Attrs = v
		}
	}
	return t
}

func (t *ToggleComponent) Render() string {
	checkedAttr := ""
	if t.Checked {
		checkedAttr = " checked"
	}

	styleStr := ""
	if len(t.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, strings.Join(t.Style.entries(), ";"))
	}

	// Toggle Switch styling using standard CSS classes
	return fmt.Sprintf(`
		<label class="goui-toggle %s"%s>
			<input type="checkbox"%s id="%s"%s>
			<span class="goui-toggle-slider"></span>
		</label>`, t.Class, styleStr, checkedAttr, t.GetID(), renderAttrs(t.Attrs))
}

func renderAttrs(attrs Attr) string {
	if len(attrs) == 0 {
		return ""
	}
	var res []string
	for k, v := range attrs {
		res = append(res, fmt.Sprintf(` %s="%s"`, k, v))
	}
	return strings.Join(res, "")
}
