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
	return `<header class="goui-header">
		<button class="goui-mobile-menu-btn" onclick="gouiToggleSidebar()">
			<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="18" x2="21" y2="18"></line></svg>
		</button>
		<span class="goui-header-title">` + h.Title + `</span>
	</header>`
}
