package components

import (
	"fmt"
)

type BadgeVariant string

const (
	DefaultBadge BadgeVariant = "default"
	SuccessBadge BadgeVariant = "success"
	ErrorBadge   BadgeVariant = "error"
	WarningBadge BadgeVariant = "warning"
	InfoBadge    BadgeVariant = "info"
)

type BadgeComponent struct {
	Content string
	Variant BadgeVariant
	Class   string
}

// Badge renders a small pill label.
// Example: Badge("Active", SuccessBadge)
func Badge(content string, opts ...any) *BadgeComponent {
	b := &BadgeComponent{
		Content: content,
		Variant: DefaultBadge,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case BadgeVariant:
			b.Variant = v
		case string:
			b.Class = v
		case Class:
			b.Class = string(v)
		}
	}
	return b
}

func (b *BadgeComponent) Render() string {
	return fmt.Sprintf(`<span class="goui-badge goui-badge-%s %s">%s</span>`, b.Variant, b.Class, b.Content)
}
