package components

import (
	"fmt"
)

type DownloadButtonComponent struct {
	ID       string
	Label    string
	URL      string
	Filename string
	Class    string
	Style    Style
	Attrs    Attr
	SyncWith *FileUploaderComponent
}

func (d *DownloadButtonComponent) GetID() string {
	if d.ID == "" {
		d.ID = AutoID()
	}
	return d.ID
}

// Filename sets the suggested name for the downloaded file.
type Filename string

// DownloadButton creates a button that triggers a file download.
func DownloadButton(label, url string, opts ...any) *DownloadButtonComponent {
	d := &DownloadButtonComponent{
		Label: label,
		URL:   url,
		Class: "goui-btn goui-btn-primary",
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case *FileUploaderComponent:
			d.SyncWith = v
		case Filename:
			d.Filename = string(v)
		case string:
			if !ParseStringAttr(v, &d.Class, &idDummy, &d.Attrs) {
				d.Label = v
			}
		case Style:
			d.Style = v
		case Class:
			d.Class = string(v)
		case Attr:
			d.Attrs = v
		case ID:
			d.ID = string(v)
		}
	}
	return d
}

var idDummy string // Used to discard ID from ParseStringAttr in DownloadButton

func (d *DownloadButtonComponent) Render() string {
	downloadAttr := "download"
	if d.Filename != "" {
		downloadAttr = fmt.Sprintf(`download="%s"`, d.Filename)
	}

	styleStr := ""
	if len(d.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, d.Style.Render())
	}

	id := d.GetID()
	syncScript := ""
	if d.SyncWith != nil {
		// Injeta a lógica do SyncFile automaticamente
		syncScript = fmt.Sprintf(`
			<script>
				(function() {
					var up = document.getElementById('%s');
					var dl = document.getElementById('%s');
					if (up && dl) {
						up.addEventListener('change', function(e) {
							var file = e.target.files[0];
							if (file) {
								var url = URL.createObjectURL(file);
								dl.href = url;
								dl.setAttribute('download', file.name);
								dl.innerText = ' Baixar ' + file.name;
							}
						});
					}
				})();
			</script>
		`, d.SyncWith.GetID(), id)
	}

	return fmt.Sprintf(`
		<div style="display: inline-block;">
			<a href="%s" %s id="%s" class="%s"%s %s>
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right:8px"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
				%s
			</a>
			%s
		</div>`, d.URL, downloadAttr, id, d.Class, styleStr, renderAttrs(d.Attrs), d.Label, syncScript)
}
