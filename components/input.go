package components

import (
	"fmt"
	"strings"
)

// Input renders a bare <input> element.
// The caller composes labels, wrappers, etc. as needed.
type InputComponent struct {
	Name        string
	ID          string
	Type        string
	Value       string
	Placeholder string
	Class       string
	Style       Style
	Attrs       Attr
	Watchers    []WatchOption
	Binding     *BindOption
	Validation  *Validation
}

// GetID returns the ID of the input, automatically generating one if empty.
func (i *InputComponent) GetID() string {
	if i.ID == "" {
		i.ID = AutoID()
	}
	return i.ID
}

// Placeholder sets the input placeholder text.
type Placeholder string

// Type represents the HTML input type (text, color, number, etc).
type Type string

// Value sets the initial value of an input or dropdown.
type Value string

// Input creates a new input component with variadic options.
func Input(opts ...any) *InputComponent {
	i := &InputComponent{Type: "text"}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if !ParseStringAttr(v, &i.Class, &i.ID, &i.Attrs) {
				types := map[string]bool{
					"text": true, "password": true, "color": true, "date": true,
					"email": true, "number": true, "checkbox": true, "radio": true,
				}
				if types[v] {
					i.Type = v
				} else if i.Placeholder == "" {
					i.Placeholder = v
				} else if i.Name == "" {
					i.Name = v
				} else {
					i.Value = v
				}
			}
		case Placeholder:
			i.Placeholder = string(v)
		case Type:
			i.Type = string(v)
		case Value:
			i.Value = string(v)
		case Name:
			i.Name = string(v)
		case ID:
			i.ID = string(v)
		case Style:
			i.Style = v
		case Class:
			i.Class = string(v)
		case Attr:
			i.Attrs = v
		case WatchOption:
			i.Watchers = append(i.Watchers, v)
		case BindOption:
			i.Binding = &v
		case Validation:
			i.Validation = &v
		}
	}
	return i
}

func (i *InputComponent) Render() string {
	var attrs []string

	if i.Name != "" {
		attrs = append(attrs, fmt.Sprintf(`name="%s"`, i.Name))
	}
	if i.GetID() != "" {
		attrs = append(attrs, fmt.Sprintf(`id="%s"`, i.ID))
	}

	attrs = append(attrs, fmt.Sprintf(`type="%s"`, i.Type))

	if i.Value != "" {
		attrs = append(attrs, fmt.Sprintf(`value="%s"`, i.Value))
	}
	if i.Placeholder != "" {
		attrs = append(attrs, fmt.Sprintf(`placeholder="%s"`, i.Placeholder))
	}
	fullClass := "goui-input"
	if i.Class != "" {
		fullClass += " " + i.Class
	}
	attrs = append(attrs, fmt.Sprintf(`class="%s"`, fullClass))
	if len(i.Style) > 0 {
		attrs = append(attrs, fmt.Sprintf(`style="%s"`, i.Style.Render()))
	}
	if len(i.Attrs) > 0 {
		for k, v := range i.Attrs {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}
	if i.Validation != nil {
		attrs = append(attrs, renderValidationAttrs(i.Validation)...)
	}

	id := i.GetID()
	scripts := ""
	if i.Validation != nil {
		scripts += renderValidationScript(id, i.Validation)
	}
	for _, w := range i.Watchers {
		scripts += fmt.Sprintf(`
			<script>
				(function() {
					var src = document.getElementById('%s');
					if (src) {
						var action = '%s';
						src.addEventListener('input', function() {
							var target = document.getElementById('%s');
							if (!target) return;
							var val = src.value;
							if (action === 'value') target.value = val;
							if (action === 'placeholder') target.placeholder = val;
							if (action === 'type') target.type = val;
							if (action === 'class') target.className = val;
						});
					}
				})();
			</script>
		`, w.SourceID, w.Action, id)
	}

	if i.Binding != nil {
		sourcesJSON := "{"
		first := true
		for alias, sid := range i.Binding.Sources {
			if !first {
				sourcesJSON += ","
			}
			sourcesJSON += fmt.Sprintf("'%s': '%s'", alias, sid)
			first = false
		}
		sourcesJSON += "}"

		scripts += fmt.Sprintf(`
			<script>
				(function() {
					var sources = %s;
					var target = document.getElementById('%s');
					var template = "%s";
					function update() {
						var content = template;
						for (var alias in sources) {
							var el = document.getElementById(sources[alias]);
							content = content.replace('${' + alias + '}', el ? el.value : '');
						}
						target.value = content;
					}
					for (var alias in sources) {
						var el = document.getElementById(sources[alias]);
						if (el) el.addEventListener('input', update);
					}
					update();
				})();
			</script>
		`, sourcesJSON, id, strings.ReplaceAll(i.Binding.Template, "\"", "\\\""))
	}

	return "<input " + joinAttrs(attrs) + ">" + scripts
}

func joinAttrs(attrs []string) string {
	var result string
	for i, a := range attrs {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}
