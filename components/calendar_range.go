package components

import (
	"fmt"
)

type DateRangeComponent struct {
	IDStart string
	IDEnd   string
	Class   string
	Style   Style
	Attrs   Attr
}

// GetID returns the ID of the start input (default for SyncSource)
func (d *DateRangeComponent) GetID() string {
	if d.IDStart == "" {
		d.IDStart = AutoID()
	}
	return d.IDStart
}

// GetEndID returns the ID of the end input
func (d *DateRangeComponent) GetEndID() string {
	if d.IDEnd == "" {
		d.IDEnd = AutoID()
	}
	return d.IDEnd
}

// DateRange creates a component with two linked date inputs.
func DateRange(opts ...any) *DateRangeComponent {
	d := &DateRangeComponent{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			ParseStringAttr(v, &d.Class, &d.IDStart, &d.Attrs)
		case Style:
			d.Style = v
		case Class:
			d.Class = string(v)
		case Attr:
			d.Attrs = v
		}
	}
	return d
}

func (d *DateRangeComponent) Render() string {
	start := Calendar(ID(d.GetID()), Class(d.Class)).Render()
	end := Calendar(ID(d.GetEndID()), Class(d.Class)).Render()

	return fmt.Sprintf(`
		<div class="goui-date-range" style="display: flex; align-items: center; gap: 10px">
			<div class="goui-date-range-start">%s</div>
			<div class="goui-date-range-separator">até</div>
			<div class="goui-date-range-end">%s</div>
		</div>`, start, end)
}
