package components

import (
	"fmt"
	"strings"
)

type SliderComponent struct {
	ID    string
	Min   int
	Max   int
	Step  int
	Value int
	Class string
	Style Style
	Attrs Attr
}

func (s *SliderComponent) GetID() string {
	if s.ID == "" {
		s.ID = AutoID()
	}
	return s.ID
}

// Min sets the minimum value for a slider.
type Min int

// Max sets the maximum value for a slider.
type Max int

// Step sets the increment step for a slider.
type Step int

// Slider creates a range input component.
func Slider(opts ...any) *SliderComponent {
	s := &SliderComponent{Min: 0, Max: 100, Step: 1}
	for _, opt := range opts {
		switch v := opt.(type) {
		case int:
			s.Value = v
		case Min:
			s.Min = int(v)
		case Max:
			s.Max = int(v)
		case Step:
			s.Step = int(v)
		case string:
			if !ParseStringAttr(v, &s.Class, &s.ID, &s.Attrs) {
				// Could handle placeholder or other strings if needed
			}
		case Style:
			s.Style = v
		case Class:
			s.Class = string(v)
		case ID:
			s.ID = string(v)
		case Attr:
			s.Attrs = v
		}
	}
	return s
}

func (s *SliderComponent) Render() string {
	styleStr := ""
	if len(s.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, strings.Join(s.Style.entries(), ";"))
	}

	return fmt.Sprintf(`
		<div class="goui-slider-container %s"%s>
			<input type="range" class="goui-slider" id="%s" 
				min="%d" max="%d" step="%d" value="%d"%s>
		</div>`, s.Class, styleStr, s.GetID(), s.Min, s.Max, s.Step, s.Value, renderAttrs(s.Attrs))
}
