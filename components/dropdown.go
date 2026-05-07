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
		case []string:
			for _, s := range v {
				d.Options = append(d.Options, Option{Value: s, Label: s})
			}
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
		case Value:
			val := string(v)
			for i, opt := range d.Options {
				if opt.Value == val {
					d.Options[i].Selected = true
				}
			}
		case Attr:
			d.Attrs = v
		}
	}
	return d
}

func (d *DropdownComponent) Render() string {
	id := d.GetID()
	var b strings.Builder

	// Wrapper seguindo sua sugestão: position: relative
	wrapperClass := "goui-dropdown"
	if d.Class != "" {
		wrapperClass += " " + d.Class
	}
	
	fmt.Fprintf(&b, `<div class="%s" id="wrapper-%s" style="position: relative;">`, wrapperClass, id)

	// Botão que ativa o dropdown
	selectedLabel := "Selecione..."
	if len(d.Options) > 0 {
		for _, opt := range d.Options {
			if opt.Selected {
				selectedLabel = opt.Label
				break
			}
		}
	}

	fmt.Fprintf(&b, `
		<button type="button" class="goui-select" onclick="gouiToggleCustomDropdown('%s')" id="btn-%s">
			<span class="selected-text">%s</span>
			<span class="goui-dropdown-arrow">
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
			</span>
		</button>`, id, id, selectedLabel)

	// Menu de Opções (Oculto por padrão, position: absolute)
	fmt.Fprintf(&b, `
		<ul class="goui-dropdown-menu" id="menu-%s" style="display: none; position: absolute; top: 100%%; left: 0; z-index: 1000; margin-top: 4px;">`, id)
	
	for _, opt := range d.Options {
		activeClass := ""
		if opt.Selected { activeClass = "active" }
		fmt.Fprintf(&b, `
			<li class="goui-dropdown-item %s" onclick="gouiSelectOption('%s', '%s', '%s')">%s</li>`, 
			activeClass, id, opt.Value, opt.Label, opt.Label)
	}
	b.WriteString(`</ul>`)

	// Input oculto para manter compatibilidade com forms
	fmt.Fprintf(&b, `<input type="hidden" name="%s" id="%s" value="%s">`, d.Name, id, func() string {
		for _, opt := range d.Options { if opt.Selected { return opt.Value } }
		return ""
	}())

	b.WriteString(`</div>`)

	// Lógica JS para o Custom Dropdown
	b.WriteString(`
		<script>
			if (typeof window.gouiToggleCustomDropdown !== 'function') {
				window.gouiToggleCustomDropdown = function(id) {
					var menu = document.getElementById('menu-' + id);
					var allMenus = document.querySelectorAll('.goui-dropdown-menu');
					allMenus.forEach(m => { if (m.id !== 'menu-' + id) m.style.display = 'none'; });
					menu.style.display = menu.style.display === 'none' ? 'block' : 'none';
				};
				window.gouiSelectOption = function(id, val, label) {
					document.getElementById(id).value = val;
					document.getElementById('btn-' + id).querySelector('.selected-text').innerText = label;
					document.getElementById('menu-' + id).style.display = 'none';
					// Trigger change event for SyncWithForm
					var input = document.getElementById(id);
					var event = new Event('change', { bubbles: true });
					input.dispatchEvent(event);
				};
				window.addEventListener('click', function(e) {
					if (!e.target.closest('.goui-dropdown')) {
						document.querySelectorAll('.goui-dropdown-menu').forEach(m => m.style.display = 'none');
					}
				});
			}
		</script>
	`)

	return b.String()
}

// Name sets the HTML name attribute (for form submission).
type Name string

// Size sets the number of visible rows in a multiselect.
type Size int
