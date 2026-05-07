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
	Editable       bool
	ShowEditToggle bool
	FormSync       *FormSyncOption
}

type FormSyncOption struct {
	ButtonID string
	Inputs   []SyncSource
}

// MaxHeight sets the maximum height of the table with a scrollbar.
type MaxHeight string

// WithExport enables the CSV export button with a specific filename.
type WithExport string

// Editable enables or disables cell editing in the table.
type Editable bool

// ShowEditToggle adds a toggle to the table header to enable/disable editing.
type ShowEditToggle bool

// SyncWithForm links a button and a set of inputs to the table for automatic data entry.
func SyncWithForm(btn SyncSource, inputs ...SyncSource) *FormSyncOption {
	return &FormSyncOption{ButtonID: btn.GetID(), Inputs: inputs}
}

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
		case Editable:
			t.Editable = bool(v)
		case ShowEditToggle:
			t.ShowEditToggle = bool(v)
		case *FormSyncOption:
			t.FormSync = v
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
			editable := ""
			if t.Editable {
				editable = ` contenteditable="true"`
			}
			rowsHTML += fmt.Sprintf(`<td%s>%s</td>`, editable, cell)
		}
		rowsHTML += "</tr>"
	}

	formSyncScript := ""
	if t.FormSync != nil {
		inputIDs := []string{}
		for _, in := range t.FormSync.Inputs {
			inputIDs = append(inputIDs, in.GetID())
		}
		
		idsList := ""
		for i, id := range inputIDs {
			if i > 0 { idsList += "," }
			idsList += fmt.Sprintf("'%s'", id)
		}

		formSyncScript = fmt.Sprintf(`
			<script>
				(function() {
					var btn = document.getElementById('%s');
					var tbl = document.getElementById('%s').querySelector('tbody');
					var inputIDs = [%s];
					
					if (btn && tbl) {
						btn.addEventListener('click', function() {
							var row = document.createElement('tr');
							var hasContent = false;
							
							inputIDs.forEach(function(id) {
								var el = document.getElementById(id);
								var td = document.createElement('td');
								if (!el) {
									td.innerText = '-';
								} else {
									hasContent = true;
									if (el.tagName === 'SELECT') {
										// Pega o wrapper (.goui-dropdown) para manter o estilo e a seta
										var wrapper = el.closest('.goui-dropdown') || el;
										td.innerHTML = wrapper.outerHTML;
										var newSel = td.querySelector('select');
										if (newSel) {
											newSel.value = el.value;
											newSel.id = ''; // Evita IDs duplicados na tabela
										}
									} else if (el.type === 'checkbox') {
										// Clona o toggle
										var wrapper = el.closest('.goui-toggle') || el;
										td.innerHTML = wrapper.outerHTML;
										var chk = td.querySelector('input');
										chk.checked = el.checked;
										chk.id = '';
									} else {
										td.innerText = el.value;
									}
								}
								row.appendChild(td);
							});
							
							if (hasContent) {
								// Sincroniza com modo edição mestre
								var masterToggle = document.getElementById('toggle-edit-%s');
								var canEdit = masterToggle ? masterToggle.checked : false;
								row.querySelectorAll('td').forEach(c => c.contentEditable = canEdit);
								row.querySelectorAll('input, select').forEach(i => i.disabled = !canEdit);
								
								tbl.appendChild(row);
								// Reset inputs de texto
								inputIDs.forEach(function(id) {
									var el = document.getElementById(id);
									if (el && el.type === 'text') el.value = '';
								});
							}
						});
					}
				})();
			</script>
		`, t.FormSync.ButtonID, id, idsList, id)
	}

	containerStyle := ""
	if t.MaxHeight != "" {
		containerStyle = fmt.Sprintf("max-height: %s; overflow-y: auto;", t.MaxHeight)
	}

	tableActions := ""
	if t.ExportFileName != "" || t.ShowEditToggle {
		tableActions = `<div class="goui-table-actions" style="display: flex; gap: 10px; align-items: center;">`
		
		if t.ShowEditToggle {
			tableActions += fmt.Sprintf(`
				<div style="display: flex; align-items: center; gap: 8px; margin-right: auto; padding-left: 5px;">
					<label class="goui-toggle goui-toggle-sm">
						<input type="checkbox" id="toggle-edit-%s" %s 
							onclick="toggleTableEdit('%s', this.checked)">
						<span class="goui-toggle-slider"></span>
					</label>
					<span style="font-size: 0.75rem; font-weight: bold; color: var(--goui-text-muted); text-transform: uppercase;">Editar</span>
				</div>`, id, func() string { if t.Editable { return "checked" }; return "" }(), id)
		}

		if t.ExportFileName != "" {
			tableActions += fmt.Sprintf(`
				<button class="goui-btn goui-btn-sm goui-btn-secondary" onclick="exportTableToCSV('%s', '%s')" title="Exportar CSV">
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
				</button>`, id, t.ExportFileName)
		}
		tableActions += `</div>`
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
							for (var j = 0; j < cols.length; j++) {
								var cell = cols[j];
								var val = "";
								
								// Se tiver um select (Dropdown), pega o valor selecionado
								var sel = cell.querySelector('select');
								if (sel) {
									val = sel.value;
								} 
								// Se tiver um checkbox (Toggle), pega o status
								else {
									var chk = cell.querySelector('input[type="checkbox"]');
									if (chk) {
										val = chk.checked ? "Ativo" : "Inativo";
									} else {
										val = cell.innerText;
									}
								}
								
								row.push('"' + val.trim() + '"');
							}
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
				if (typeof window.toggleTableEdit !== 'function') {
					window.toggleTableEdit = function(tableID, canEdit) {
						var table = document.getElementById(tableID);
						var cells = table.querySelectorAll('td');
						cells.forEach(function(c) {
							c.contentEditable = canEdit;
						});
						// Também desabilita/habilita inputs e selects dentro da tabela
						var controls = table.querySelectorAll('input, select');
						controls.forEach(function(ctrl) {
							ctrl.disabled = !canEdit;
						});
					};
				}
			</script>
		</div>` + formSyncScript, t.Class, id, tableActions, containerStyle, id, headerHTML, rowsHTML, syncScript)
}
