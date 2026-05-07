package components

import (
	"fmt"
)

// ColorPickerComponent renders a styled color picker with a swatch and hex display.
type ColorPickerComponent struct {
	ID    string
	Name  string
	Value string
	Class string
	Attrs Attr
}

func (c *ColorPickerComponent) GetID() string {
	if c.ID == "" {
		c.ID = AutoID()
	}
	return c.ID
}

// ColorPicker creates a styled color picker with a clickable swatch and hex label.
// The default color is #7c3aed (goUI primary).
//
// Example:
//
//	ColorPicker(Name("brand_color"), Value("#ff5733"))
func ColorPicker(opts ...any) *ColorPickerComponent {
	c := &ColorPickerComponent{Value: "#7c3aed"}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if !ParseStringAttr(v, &c.Class, &c.ID, &c.Attrs) {
				// Treat bare hex strings as Value
				if len(v) > 0 && v[0] == '#' {
					c.Value = v
				}
			}
		case Value:
			c.Value = string(v)
		case Name:
			c.Name = string(v)
		case ID:
			c.ID = string(v)
		case Class:
			c.Class = string(v)
		case Attr:
			c.Attrs = v
		}
	}
	return c
}

func (c *ColorPickerComponent) Render() string {
	id := c.GetID()
	swatchID := id + "-swatch"
	hexID := id + "-hex"

	nameAttr := ""
	if c.Name != "" {
		nameAttr = fmt.Sprintf(` name="%s"`, c.Name)
	}

	cls := "goui-colorpicker"
	if c.Class != "" {
		cls += " " + c.Class
	}

	attrStr := renderAttrs(c.Attrs)

	return fmt.Sprintf(`
<div class="%s" id="%s">
  <div class="goui-colorpicker-swatch" id="%s" style="background:%s;" onclick="document.getElementById('%s').click()"></div>
  <input type="color" id="%s"%s value="%s" class="goui-colorpicker-input"%s>
  <span class="goui-colorpicker-hex" id="%s">%s</span>
</div>
<script>
(function(){
  var picker = document.getElementById('%s');
  var swatch = document.getElementById('%s');
  var hex    = document.getElementById('%s');
  if(!picker) return;
  picker.addEventListener('input', function() {
    swatch.style.background = picker.value;
    hex.textContent = picker.value;
  });
})();
</script>`,
		cls, id,
		swatchID, c.Value, id,
		id, nameAttr, c.Value, attrStr,
		hexID, c.Value,
		id, swatchID, hexID,
	)
}
