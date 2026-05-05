package components

import (
	"fmt"
)

type TableComponent struct {
	ID        string
	Headers   []string
	Rows      [][]string
	Class     string
	Style     Style
	Attrs     Attr
	SyncWith       *FileUploaderComponent
	MaxHeight      string // E.g., "600px"
	ExportFileName string
}

// MaxHeight sets the maximum height of the table with a scrollbar.
type MaxHeight string

// WithExport enables the CSV export button with a specific filename.
type WithExport string

func (t *TableComponent) GetID() string {
	if t.ID == "" {
		t.ID = AutoID()
	}
	return t.ID
}

// Table creates an interactive table that can be populated from CSV.
// It supports headers, rows, and optional MaxHeight via components.Style or raw string.
func Table(headers []string, rows [][]string, opts ...any) *TableComponent {
	t := &TableComponent{
		Headers: headers,
		Rows:    rows,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case *FileUploaderComponent:
			t.SyncWith = v
		case MaxHeight:
			t.MaxHeight = string(v)
		case WithExport:
			t.ExportFileName = string(v)
		case string:
			if !ParseStringAttr(v, &t.Class, &t.ID, &t.Attrs) {
				t.MaxHeight = v
			}
		case Style:
			t.Style = v
		case Class:
			t.Class = string(v)
		case Attr:
			t.Attrs = v
		case ID:
			t.ID = string(v)
		}
	}
	return t
}

func (t *TableComponent) Render() string {
	id := t.GetID()
	
	syncScript := ""
	if t.SyncWith != nil {
		syncScript = fmt.Sprintf(`
			<script>
				(function() {
					var up = document.getElementById('%s');
					var table = document.getElementById('%s');
					if (up && table) {
						up.addEventListener('change', function(e) {
							var file = e.target.files[0];
							if (!file) return;
							var reader = new FileReader();
							reader.onload = function(event) {
								var contents = event.target.result;
								var lines = contents.split(/\r?\n/);
								if (lines.length === 0) return;
								
								var thead = table.querySelector('thead tr') || table.createTHead().insertRow();
								var tbody = table.querySelector('tbody') || table.createTBody();
								thead.innerHTML = '';
								tbody.innerHTML = '';
								
								var separator = lines[0].indexOf(';') > -1 ? ';' : ',';
								var headers = lines[0].split(separator);
								headers.forEach(function(h) {
									var th = document.createElement('th');
									th.innerText = h.replace(/"/g, '').trim();
									thead.appendChild(th);
								});
								
								for (var i = 1; i < lines.length; i++) {
									if (lines[i].trim() === "") continue;
									var tr = document.createElement('tr');
									var cells = lines[i].split(separator);
									cells.forEach(function(c) {
										var td = document.createElement('td');
										td.contentEditable = "true";
										td.innerText = c.replace(/"/g, '').trim();
										tr.appendChild(td);
									});
									tbody.appendChild(tr);
								}
							};
							reader.readAsText(file);
						});
					}
				})();
			</script>
		`, t.SyncWith.GetID(), id)
	}

	headerHTML := ""
	for _, h := range t.Headers {
		headerHTML += fmt.Sprintf("<th>%s</th>", h)
	}

	rowsHTML := ""
	for _, row := range t.Rows {
		rowsHTML += "<tr>"
		for _, cell := range row {
			rowsHTML += fmt.Sprintf(`<td contenteditable="true">%s</td>`, cell)
		}
		rowsHTML += "</tr>"
	}

	containerStyle := ""
	if t.MaxHeight != "" {
		containerStyle = fmt.Sprintf("max-height: %s; overflow-y: auto;", t.MaxHeight)
	}

	exportButton := ""
	if t.ExportFileName != "" {
		exportButton = fmt.Sprintf(`
			<div class="goui-table-actions">
				<button class="goui-btn goui-btn-sm goui-btn-secondary" onclick="exportTableToCSV('%s', '%s')" title="Exportar CSV">
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
				</button>
			</div>`, id, t.ExportFileName)
	}

	return fmt.Sprintf(`
		<div class="goui-table-container %s" id="container-%s">
			%s
			<div style="overflow-x: auto; %s">
				<table class="goui-table" id="%s">
					<thead><tr style="position: sticky; top: 0; z-index: 10; background: var(--goui-surface-2);">%s</tr></thead>
					<tbody>%s</tbody>
				</table>
			</div>
			%s
			<script>
				if (typeof window.exportTableToCSV !== 'function') {
					window.exportTableToCSV = function(tableID, filename) {
						var table = document.getElementById(tableID);
						var csv = [];
						for (var i = 0; i < table.rows.length; i++) {
							var row = [], cols = table.rows[i].cells;
							for (var j = 0; j < cols.length; j++) row.push('"' + cols[j].innerText + '"');
							csv.push(row.join(","));
						}
						var csvFile = new Blob([csv.join("\n")], {type: "text/csv"});
						var downloadLink = document.createElement("a");
						downloadLink.download = filename || "export.csv";
						downloadLink.href = window.URL.createObjectURL(csvFile);
						downloadLink.style.display = "none";
						document.body.appendChild(downloadLink);
						downloadLink.click();
						document.body.removeChild(downloadLink);
					};
				}
			</script>
		</div>`, t.Class, id, exportButton, containerStyle, id, headerHTML, rowsHTML, syncScript)
}
