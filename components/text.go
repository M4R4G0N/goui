package components

import (
	"fmt"
	"strings"
)

type TextComponent struct {
	Content      string
	Tag          string
	Class        string
	Style        Style
	ID           string
	Attrs        Attr
	SyncSource   SyncSource
	SyncFallback string
	Watchers     []WatchOption
	Binding      *BindOption
}

// GetID returns the ID of the text element, automatically generating one if empty.
func (t *TextComponent) GetID() string {
	if t.ID == "" {
		t.ID = AutoID()
	}
	return t.ID
}

// Text renders a typed text element.
func Text(contentOrSource any, opts ...any) *TextComponent {
	t := &TextComponent{Tag: "p"}

	switch v := contentOrSource.(type) {
	case string:
		t.Content = v
	case SyncSource:
		t.SyncSource = v
		t.SyncFallback = "..."
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if !ParseStringAttr(v, &t.Class, &t.ID, &t.Attrs) {
				t.Tag = v
			}
		case Style:
			t.Style = v
		case Class:
			t.Class = string(v)
		case ID:
			t.ID = string(v)
		case Attr:
			t.Attrs = v
		case WatchOption:
			t.Watchers = append(t.Watchers, v)
		case BindOption:
			t.Binding = &v
		}
	}
	return t
}

func (t *TextComponent) Render() string {
	id := t.GetID()
	styleStr := ""
	if len(t.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, t.Style.Render())
	}
	attrs := renderAttrs(t.Attrs)

	// Injeta scripts de Watch/Bind se necessário
	scripts := ""
	for _, w := range t.Watchers {
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
							if (action === 'text') target.innerText = val;
							if (action === 'color') target.style.color = val;
							if (action === 'size') target.style.fontSize = val + 'px';
							if (action === 'weight') target.style.fontWeight = val;
							if (action === 'class') target.className = 'goui-pg-preview-box ' + val;
							if (action === 'tag') {
								var newEl = document.createElement(val);
								newEl.id = target.id;
								newEl.className = target.className;
								newEl.style.cssText = target.style.cssText;
								newEl.innerText = target.innerText;
								target.parentNode.replaceChild(newEl, target);
							}
						});
					}
				})();
			</script>
		`, w.SourceID, w.Action, id)
	}

	if t.Binding != nil {
		sourcesJSON := "{"
		first := true
		for alias, sid := range t.Binding.Sources {
			if !first {
				sourcesJSON += ","
			}
			sourcesJSON += fmt.Sprintf("'%s': '%s'", alias, sid)
			first = false
		}
		sourcesJSON += "}"

		// Escapa quebras de linha e aspas para o JS
		jsTemplate := strings.ReplaceAll(t.Binding.Template, "\n", "\\n")
		jsTemplate = strings.ReplaceAll(jsTemplate, "\r", "")
		jsTemplate = strings.ReplaceAll(jsTemplate, "\"", "\\\"")

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
							var val = el ? el.value : '';
							content = content.replace('${' + alias + '}', val);
						}
						target.innerText = content;
					}
					for (var alias in sources) {
						var el = document.getElementById(sources[alias]);
						if (el) el.addEventListener('input', update);
					}
					update();
				})();
			</script>
		`, sourcesJSON, id, jsTemplate)
	}

	if t.SyncSource != nil {
		scripts += SyncText(t.SyncSource.GetID(), id, t.SyncFallback).Render()
	}

	return fmt.Sprintf(`<%s id="%s" class="%s"%s %s>%s</%s>%s`, 
		t.Tag, id, t.Class, styleStr, attrs, t.Content, t.Tag, scripts)
}
