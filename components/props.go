package components

import (
	"fmt"
	"strings"
)

// Props define the HTML attributes and inline styling for a component
type Props struct {
	Tag   string // Used by Text to override default tags (e.g. "h2", "p")
	Style Style
	CSS   CSS
	Class string
	ID    string
}

// Render returns the HTML attributes string representing these properties.
func (p Props) Render() string {
	var attrs []string

	if p.ID != "" {
		attrs = append(attrs, fmt.Sprintf(`id="%s"`, p.ID))
	}

	if p.Class != "" {
		attrs = append(attrs, fmt.Sprintf(`class="%s"`, p.Class))
	}

	var styles []string
	if p.Style != nil {
		styles = append(styles, p.Style.Render())
	}
	if p.CSS != "" {
		styles = append(styles, string(p.CSS))
	}

	if len(styles) > 0 {
		var validStyles []string
		for _, s := range styles {
			if s != "" {
				validStyles = append(validStyles, s)
			}
		}
		if len(validStyles) > 0 {
			attrs = append(attrs, fmt.Sprintf(`style="%s"`, strings.Join(validStyles, "; ")))
		}
	}

	if len(attrs) > 0 {
		return " " + strings.Join(attrs, " ")
	}
	return ""
}
