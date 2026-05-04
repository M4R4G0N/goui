package components

import "strings"

// DocTableRow describes one row in a documentation classes table.
type DocTableRow struct {
	Class   string    // CSS class name (e.g. "goui-text-muted")
	Effect  string    // Human-readable description
	Preview Component // Live preview component
}

// DocTable renders a styled table for documenting CSS classes.
func DocTable(rows []DocTableRow) *HTML {
	var b strings.Builder

	b.WriteString(`<table class="goui-classes-table">
<thead>
<tr>
<th>Class</th>
<th>Effect</th>
<th>Preview</th>
</tr>
</thead>
<tbody>`)

	for _, r := range rows {
		preview := ""
		if r.Preview != nil {
			preview = r.Preview.Render()
		}
		b.WriteString(`<tr>
<td><code class="goui-class-tag">` + r.Class + `</code></td>
<td>` + r.Effect + `</td>
<td>` + preview + `</td>
</tr>`)
	}

	b.WriteString(`</tbody>
</table>`)

	h := HTML(b.String())
	return &h
}
