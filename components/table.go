package components

import (
	"fmt"
)

type DataTableComponent struct {
	ID       string
	Headers  []string
	Rows     [][]string
	Class    string
	Style    Style
	Attrs    Attr
	SyncWith *FileUploaderComponent
}

func (t *DataTableComponent) GetID() string {
	if t.ID == "" {
		t.ID = AutoID()
	}
	return t.ID
}

// DataTable creates an interactive table that can be populated from CSV.
func DataTable(headers []string, rows [][]string, opts ...any) *DataTableComponent {
	t := &DataTableComponent{
		Headers: headers,
		Rows:    rows,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case *FileUploaderComponent:
			t.SyncWith = v
		case string:
			ParseStringAttr(v, &t.Class, &t.ID, &t.Attrs)
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

func (t *DataTableComponent) Render() string {
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

	return fmt.Sprintf(`
		<div class="goui-table-container %s" id="container-%s">
			<div class="goui-table-actions">
				<button class="goui-btn goui-btn-primary" onclick="exportTableToCSV('%s')">
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right:5px"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
					Exportar CSV
				</button>
			</div>
			<div style="overflow-x: auto;">
				<table class="goui-table" id="%s">
					<thead><tr>%s</tr></thead>
					<tbody>%s</tbody>
				</table>
			</div>
			%s
			<script>
				window.exportTableToCSV = function(tableID) {
					var table = document.getElementById(tableID);
					var csv = [];
					for (var i = 0; i < table.rows.length; i++) {
						var row = [], cols = table.rows[i].cells;
						for (var j = 0; j < cols.length; j++) row.push('"' + cols[j].innerText + '"');
						csv.push(row.join(","));
					}
					var csvFile = new Blob([csv.join("\n")], {type: "text/csv"});
					var downloadLink = document.createElement("a");
					downloadLink.download = "export_goui.csv";
					downloadLink.href = window.URL.createObjectURL(csvFile);
					downloadLink.style.display = "none";
					document.body.appendChild(downloadLink);
					downloadLink.click();
				};
			</script>
		</div>`, t.Class, id, id, id, headerHTML, rowsHTML, syncScript)
}
