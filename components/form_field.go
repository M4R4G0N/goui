package components

// FormField creates a standard form group with a label, an input component, and an optional help text.
func FormField(label string, input Component, helpText string) *DivComponent {
	items := []Component{
		Text(label, "label", Class("goui-label")),
		input,
	}

	if helpText != "" {
		items = append(items, Text(helpText, "p", Class("goui-text-muted goui-mt-1"), Style{"font-size": "0.8rem"}))
	}

	return Div(
		Class("goui-form-field goui-mb-20"),
		items,
	)
}
