package components

import "fmt"

// ButtonVariant controls the visual style of a Button.
type ButtonVariant string

const (
	Primary   ButtonVariant = "primary"
	Ghost     ButtonVariant = "ghost"
	Secondary ButtonVariant = "secondary"
)

// SubmitAction represents a server-side call configuration.
type SubmitAction struct {
	Path    string
	Target  string
	Sources []string
}

// OnSubmit creates a configuration for a button to send data to the server.
func OnSubmit(path, target string, sources ...string) SubmitAction {
	return SubmitAction{Path: path, Target: target, Sources: sources}
}

type ButtonComponent struct {
	ID      string
	Label   string
	Variant ButtonVariant
	Href    string
	Class   string
	Style   Style
	Attrs   Attr
	Action  SubmitAction
}


func (b *ButtonComponent) GetID() string {
	if b.ID == "" {
		b.ID = AutoID()
	}
	return b.ID
}

// Button renders a styled action element.
// If href is provided it renders as <a>, otherwise as <button>.
func Button(label string, opts ...any) *ButtonComponent {
	b := &ButtonComponent{
		Label:   label,
		Variant: Primary,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case ButtonVariant:
			b.Variant = v
		case SubmitAction:
			b.Action = v
		case string:
			if !ParseStringAttr(v, &b.Class, &b.ID, &b.Attrs) {
				if b.Href == "" {
					b.Href = v
				}
			}
		case Style:
			b.Style = v
		case Class:
			b.Class = string(v)
		case ID:
			b.ID = string(v)
		case Attr:
			b.Attrs = v
		}
	}
	return b
}

func (b *ButtonComponent) Render() string {
	id := b.GetID()
	class := fmt.Sprintf("goui-btn goui-btn-%s %s", b.Variant, b.Class)
	styleStr := ""
	if len(b.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, b.Style.Render())
	}
	
	// Se tiver uma ação, injetamos o onclick e o script
	actionAttr := ""
	actionScript := ""
	if b.Action.Path != "" {
		actionAttr = ` onclick="gouiSubmit(this)"`
		
		sourcesJSON := "["
		for i, s := range b.Action.Sources {
			if i > 0 { sourcesJSON += "," }
			sourcesJSON += fmt.Sprintf("'%s'", s)
		}
		sourcesJSON += "]"

		actionScript = fmt.Sprintf(`
			<script>
				function gouiSubmit(btn) {
					var data = {};
					var sources = %s;
					sources.forEach(function(s) {
						var el = document.getElementById(s);
						if (el) data[s] = el.value;
					});
					fetch('%s', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify(data)
					})
					.then(r => r.text())
					.then(res => {
						var target = document.getElementById('%s');
						if (target) target.innerText = res;
					});
				}
			</script>
		`, sourcesJSON, b.Action.Path, b.Action.Target)
	}

	attrs := renderAttrs(b.Attrs)

	if b.Href != "" {
		return fmt.Sprintf(`<a id="%s" class="%s" href="%s"%s %s%s>%s</a>%s`, id, class, b.Href, styleStr, attrs, actionAttr, b.Label, actionScript)
	}
	return fmt.Sprintf(`<button id="%s" class="%s"%s %s%s>%s</button>%s`, id, class, styleStr, attrs, actionAttr, b.Label, actionScript)
}
