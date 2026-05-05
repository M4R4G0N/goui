package components

import (
	"fmt"
)

type SnippetComponent struct {
	Title       string
	Code        string
	ID          string
	Collapsible bool
	IsCollapsed bool
}

// Snippet renders a premium code block with syntax highlighting and accordion effect.
func Snippet(title, code string, opts ...any) *SnippetComponent {
	s := &SnippetComponent{
		Title:       title,
		Code:        code,
		ID:          AutoID(),
		Collapsible: true, // Padrão agora é ser colapsável
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case ID:
			s.ID = string(v)
		case bool: // Se passar um bool falso, desativa o colapso
			s.Collapsible = v
		}
	}
	return s
}

func (s *SnippetComponent) Render() string {
	codeID := s.ID + "-code"
	bodyID := s.ID + "-body"
	arrowID := s.ID + "-arrow"

	return fmt.Sprintf(`
		<div class="goui-code-section">
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css">
			<div class="goui-code-header" style="cursor: pointer" onclick="toggleSnippet('%s', '%s')">
				<div class="goui-flex goui-items-center goui-gap-10">
					<span id="%s" style="transition: transform 0.2s; display: inline-block; color: #e5e5e5 !important;">▼</span>
					<span class="goui-code-header-label">%s</span>
				</div>
				<button class="goui-btn goui-btn-sm goui-btn-ghost" onclick="event.stopPropagation(); copySnippet('%s')">Copy</button>
			</div>
			<div id="%s" class="goui-code-block" style="padding: 0 !important; display: block;">
				<pre style="margin: 0; background: transparent;"><code id="%s" class="language-go" style="text-align: left; display: block; white-space: pre;">%s</code></pre>
			</div>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-go.min.js"></script>
			<script>
				if (typeof window.copySnippet === 'undefined') {
					window.copySnippet = function(id) {
						var code = document.getElementById(id).innerText;
						navigator.clipboard.writeText(code).then(function() {
							var btn = event.target;
							var old = btn.innerText;
							btn.innerText = 'Copied!';
							setTimeout(function() { btn.innerText = old; }, 2000);
						});
					};
				}

				if (typeof window.toggleSnippet === 'undefined') {
					window.toggleSnippet = function(bodyId, arrowId) {
						var body = document.getElementById(bodyId);
						var arrow = document.getElementById(arrowId);
						if (body.style.display === 'none') {
							body.style.display = 'block';
							arrow.style.transform = 'rotate(0deg)';
						} else {
							body.style.display = 'none';
							arrow.style.transform = 'rotate(-90deg)';
						}
					};
				}

				if (window.Prism) Prism.highlightAll();
			</script>
		</div>
	`, bodyID, arrowID, arrowID, s.Title, codeID, bodyID, codeID, s.Code)
}
