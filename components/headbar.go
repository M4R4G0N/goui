package components

type HeadbarComponent struct {
	Title string
}

func Headbar(title string) *HeadbarComponent {
	return &HeadbarComponent{Title: title}
}

// RenderHead returns the content injected inside <head>: title, theme script and CSS.
func (h *HeadbarComponent) RenderHead() string {
	return `<title>` + h.Title + `</title>
<meta name="color-scheme" content="light dark">` + ThemeScript + Theme
}

// Render returns the visible top header bar rendered inside <body>.
func (h *HeadbarComponent) Render() string {
	return `<header class="goui-header"><span class="goui-header-title">` + h.Title + `</span></header>`
}
