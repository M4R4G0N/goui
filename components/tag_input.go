package components

import (
	"fmt"
	"strings"
)

// TagInputComponent renders a chip/tag input field.
type TagInputComponent struct {
	ID          string
	Name        string
	Placeholder string
	Tags        []string
	Class       string
	Attrs       Attr
}

func (t *TagInputComponent) GetID() string {
	if t.ID == "" {
		t.ID = AutoID()
	}
	return t.ID
}

// TagInput creates a chip/tag input where users type a value and press Enter
// (or comma) to add a tag. Backspace on an empty field removes the last tag.
// The hidden input stores the tags as a comma-separated string.
//
// Example:
//
//	TagInput(Name("skills"), Placeholder("Adicionar habilidade..."), []string{"Go", "Docker"})
func TagInput(opts ...any) *TagInputComponent {
	t := &TagInputComponent{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			if !ParseStringAttr(v, &t.Class, &t.ID, &t.Attrs) {
				if t.Placeholder == "" {
					t.Placeholder = v
				}
			}
		case Placeholder:
			t.Placeholder = string(v)
		case Name:
			t.Name = string(v)
		case ID:
			t.ID = string(v)
		case []string:
			t.Tags = append(t.Tags, v...)
		case Class:
			t.Class = string(v)
		case Attr:
			t.Attrs = v
		}
	}
	return t
}

func (t *TagInputComponent) Render() string {
	id := t.GetID()
	wrapID := id + "-wrap"
	inputID := id + "-field"
	hiddenID := id + "-hidden"

	placeholder := t.Placeholder
	if placeholder == "" {
		placeholder = "Adicionar tag..."
	}

	var quotedTags []string
	for _, tag := range t.Tags {
		quotedTags = append(quotedTags, fmt.Sprintf(`"%s"`, strings.ReplaceAll(tag, `"`, `\"`)))
	}
	initialJSON := "[" + strings.Join(quotedTags, ",") + "]"

	cls := "goui-taginput"
	if t.Class != "" {
		cls += " " + t.Class
	}

	return fmt.Sprintf(`
<div class="%s" id="%s">
  <div class="goui-taginput-chips" id="%s"></div>
  <input type="text" class="goui-taginput-field" id="%s" placeholder="%s" autocomplete="off">
  <input type="hidden" id="%s" name="%s">
</div>
<script>
(function(){
  var wrap  = document.getElementById('%s');
  var field = document.getElementById('%s');
  var hid   = document.getElementById('%s');
  var tags  = %s;

  function render() {
    wrap.innerHTML = '';
    tags.forEach(function(tag, i) {
      var chip = document.createElement('span');
      chip.className = 'goui-taginput-chip';
      var txt = document.createTextNode(tag + ' ');
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'goui-taginput-remove';
      btn.textContent = '×';
      btn.addEventListener('click', (function(idx){ return function(){ tags.splice(idx,1); render(); }; })(i));
      chip.appendChild(txt);
      chip.appendChild(btn);
      wrap.appendChild(chip);
    });
    hid.value = tags.join(',');
  }

  field.addEventListener('keydown', function(e) {
    if ((e.key === 'Enter' || e.key === ',') && field.value.trim() !== '') {
      e.preventDefault();
      var val = field.value.trim().replace(/,$/, '');
      if (val && tags.indexOf(val) === -1) {
        tags.push(val);
        field.value = '';
        render();
      }
    } else if (e.key === 'Backspace' && field.value === '' && tags.length > 0) {
      tags.pop();
      render();
    }
  });

  render();
})();
</script>`,
		cls, id,
		wrapID, inputID, placeholder, hiddenID, t.Name,
		wrapID, inputID, hiddenID, initialJSON,
	)
}
