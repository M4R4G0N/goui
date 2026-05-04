package components

import (
	"fmt"
	"strings"
)

// Style represents inline CSS properties as a key-value map.
type Style map[string]string

func (s Style) Render() string {
	return strings.Join(s.entries(), ";")
}

func (s Style) entries() []string {
	out := make([]string, 0, len(s))
	for k, v := range s {
		out = append(out, fmt.Sprintf("%s:%s", k, v))
	}
	return out
}

// CSS represents a raw CSS string for inline styles.
type CSS string

func (c CSS) Render() string { return string(c) }

// Class is a CSS class name passed as an option to Text, Div, etc.
type Class string

// ID sets the HTML id attribute.
type ID string

// Attr represents arbitrary HTML attributes (e.g. data-*, onchange).
type Attr map[string]string

// ParseStringAttr checks if a string is an attribute like 'class="my-class"'.
// It extracts the key and value, applying it to class, id, or attrs appropriately.
func ParseStringAttr(v string, class *string, id *string, attrs *Attr) bool {
	if !strings.Contains(v, "=") {
		return false
	}
	parts := strings.SplitN(v, "=", 2)
	key := strings.TrimSpace(parts[0])
	val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

	if key == "class" {
		if *class != "" {
			*class += " " + val
		} else {
			*class = val
		}
	} else if key == "id" {
		*id = val
	} else {
		if *attrs == nil {
			*attrs = make(Attr)
		}
		(*attrs)[key] = val
	}
	return true
}
