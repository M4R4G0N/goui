package components

// Calendar creates a date picker input component.
// It is a specialized Input with type="date".
//
// Examples:
//	Calendar(Name("birthday"))
//	Calendar("id=start-date")
func Calendar(opts ...any) *InputComponent {
	// We just reuse the Input component with the "date" type forced
	args := append([]any{"date", Class("goui-input")}, opts...)
	c := Input(args...)

	// Add auto-show behavior: clicking anywhere opens the calendar
	if c.Attrs == nil {
		c.Attrs = make(Attr)
	}
	c.Attrs["onclick"] = "if(this.showPicker)this.showPicker()"
	
	return c
}

func (c *InputComponent) AddAttr(key, value string) {
	if c.Attrs == nil {
		c.Attrs = make(Attr)
	}
	c.Attrs[key] = value
}

