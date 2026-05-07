package components

import (
	"fmt"
	"strings"
)

// TextareaComponent renders a <textarea> with auto-resize.
type TextareaComponent struct {
	Name        string
	ID          string
	Placeholder string
	Value       string
	Rows        int
	Class       string
	Style       Style
	Attrs       Attr
	Validation  *Validation
}

func (t *TextareaComponent) GetID() string {
	if t.ID == "" {
		t.ID = AutoID()
	}
	return t.ID
}

// Rows sets the number of visible text lines in a Textarea.
type Rows int

// Textarea creates a multi-line text area that auto-resizes as the user types.
//
// Example:
//
//	Textarea(Placeholder("Escreva aqui..."), Name("body"), Rows(5))
//	Textarea(Validation{Required: true, MinLen: 10, RequiredMsg: "Campo obrigatório."})
func Textarea(opts ...any) *TextareaComponent {
	t := &TextareaComponent{Rows: 3}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if !ParseStringAttr(v, &t.Class, &t.ID, &t.Attrs) {
				if t.Placeholder == "" {
					t.Placeholder = v
				}
			}
		case Placeholder:
			t.Placeholder = string(v)
		case Value:
			t.Value = string(v)
		case Name:
			t.Name = string(v)
		case ID:
			t.ID = string(v)
		case Rows:
			t.Rows = int(v)
		case Style:
			t.Style = v
		case Class:
			t.Class = string(v)
		case Attr:
			t.Attrs = v
		case Validation:
			t.Validation = &v
		}
	}
	return t
}

func (t *TextareaComponent) Render() string {
	id := t.GetID()
	var attrs []string
	if t.Name != "" {
		attrs = append(attrs, fmt.Sprintf(`name="%s"`, t.Name))
	}
	attrs = append(attrs, fmt.Sprintf(`id="%s"`, id))
	attrs = append(attrs, fmt.Sprintf(`rows="%d"`, t.Rows))
	if t.Placeholder != "" {
		attrs = append(attrs, fmt.Sprintf(`placeholder="%s"`, t.Placeholder))
	}
	cls := "goui-input goui-textarea"
	if t.Class != "" {
		cls += " " + t.Class
	}
	attrs = append(attrs, fmt.Sprintf(`class="%s"`, cls))
	if len(t.Style) > 0 {
		attrs = append(attrs, fmt.Sprintf(`style="%s"`, t.Style.Render()))
	}
	for k, v := range t.Attrs {
		attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
	}
	if t.Validation != nil {
		attrs = append(attrs, renderValidationAttrs(t.Validation)...)
	}

	validationScript := ""
	if t.Validation != nil {
		validationScript = renderValidationScript(id, t.Validation)
	}

	autoResize := fmt.Sprintf(
		`<script>(function(){var el=document.getElementById('%s');if(!el)return;function resize(){el.style.height='auto';el.style.height=el.scrollHeight+'px';}el.addEventListener('input',resize);resize();})();</script>`,
		id,
	)

	return fmt.Sprintf(`<textarea %s>%s</textarea>%s%s`,
		strings.Join(attrs, " "), t.Value, autoResize, validationScript)
}
