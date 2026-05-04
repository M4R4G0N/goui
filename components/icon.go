package components

import (
	"fmt"
)

type IconComponent struct {
	Name  string
	Size  int
	Class string
}

// Icon renders a Lucide icon by name.
// Example: Icon("search"), Icon("user", 24), Icon("box", "goui-text-primary")
func Icon(name string, opts ...any) *IconComponent {
	i := &IconComponent{
		Name: name,
		Size: 18,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case int:
			i.Size = v
		case string:
			i.Class = v
		case Class:
			i.Class = string(v)
		}
	}
	return i
}

func (i *IconComponent) Render() string {
	return fmt.Sprintf(`
		<i data-lucide="%s" class="%s" style="width: %dpx; height: %dpx;"></i>
		<script src="https://unpkg.com/lucide@latest"></script>
		<script>if(window.lucide) lucide.createIcons();</script>
	`, i.Name, i.Class, i.Size, i.Size)
}
