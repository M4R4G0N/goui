package components

// Card is a semantic wrapper for a Div with the goui-card class.
// It accepts the same arguments as a Div (Children, Class, Style, etc).
func Card(args ...any) *DivComponent {
	// Iniciamos com a classe goui-card por padrão
	d := Div(Class("goui-card"))
	
	// Processamos os argumentos extras
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			ParseStringAttr(v, &d.Class, &d.ID, &d.Attrs)
		case Style:
			d.Style = v
		case Class:
			if d.Class != "" {
				d.Class += " " + string(v)
			} else {
				d.Class = string(v)
			}
		case ID:
			d.ID = string(v)
		case Attr:
			d.Attrs = v
		case Component:
			d.Children = append(d.Children, v)
		}
	}
	return d
}
