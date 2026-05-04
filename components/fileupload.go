package components

import (
	"fmt"
	"strings"
)

type FileUploaderComponent struct {
	ID       string
	Label    string
	Accept   string
	MaxSize  int64 // in bytes
	Multiple bool
	Class    string
	Style    Style
	Attrs    Attr
}

func (f *FileUploaderComponent) GetID() string {
	if f.ID == "" {
		f.ID = AutoID()
	}
	return f.ID
}

// Accept sets the allowed file extensions (e.g., ".pdf,.csv").
type Accept string

// MaxSize sets the maximum file size in bytes.
type MaxSize int64

// Multiple enables multiple file selection.
type Multiple bool

// FileUploader creates a drag-and-drop file upload area.
func FileUploader(label string, opts ...any) *FileUploaderComponent {
	f := &FileUploaderComponent{
		Label: label,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case Accept:
			f.Accept = string(v)
		case MaxSize:
			f.MaxSize = int64(v)
		case Multiple:
			f.Multiple = bool(v)
		case string:
			ParseStringAttr(v, &f.Class, &f.ID, &f.Attrs)
		case Style:
			f.Style = v
		case Class:
			f.Class = string(v)
		case Attr:
			f.Attrs = v
		}
	}
	return f
}

func (f *FileUploaderComponent) Render() string {
	id := f.GetID()
	acceptAttr := ""
	if f.Accept != "" {
		acceptAttr = fmt.Sprintf(` accept="%s"`, f.Accept)
	}
	multiAttr := ""
	if f.Multiple {
		multiAttr = " multiple"
	}

	styleStr := ""
	if len(f.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, strings.Join(f.Style.entries(), ";"))
	}

	label := f.Label
	if f.MaxSize > 0 {
		label = fmt.Sprintf("%s (%dMB)", label, f.MaxSize/(1024*1024))
	}

	return fmt.Sprintf(`
		<div class="goui-uploader %s" id="container-%s"%s>
			<div class="goui-uploader-icon">
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
			</div>
			<div class="goui-uploader-text">%s</div>
			<input type="file" class="goui-uploader-input" id="%s"%s%s%s>
			<script>
				(function() {
					var container = document.getElementById('container-%s');
					var input = document.getElementById('%s');
					var accept = '%s';
					var maxSize = %d;
					
					container.addEventListener('dragover', function(e) {
						e.preventDefault();
						container.classList.add('dragover');
					});
					container.addEventListener('dragleave', function() {
						container.classList.remove('dragover');
					});
					container.addEventListener('drop', function(e) {
						e.preventDefault();
						container.classList.remove('dragover');
						
						var files = e.dataTransfer.files;
						if (files.length > 0) {
							console.log('Drop detectado no container, arquivos:', files.length);
							for (var i=0; i<files.length; i++) {
								if (maxSize > 0 && files[i].size > maxSize) {
									alert('Arquivo muito grande: ' + files[i].name + '. O limite é de ' + (maxSize / 1024 / 1024) + 'MB.');
									return;
								}
								if (accept) {
									var allowed = accept.split(',').map(function(ext) { return ext.trim().toLowerCase(); });
									var ext = '.' + files[i].name.split('.').pop().toLowerCase();
									if (allowed.indexOf(ext) === -1) {
										alert('Tipo de arquivo inválido: ' + files[i].name + '. Apenas ' + accept + ' são permitidos.');
										return;
									}
								}
							}
							input.files = files;
							console.log('Disparando evento change para o input...');
							input.dispatchEvent(new Event('change', { bubbles: true }));
						}
					});

					input.addEventListener('change', function() {
						console.log('Input file mudou!', this.files.length > 0 ? this.files[0].name : 'vazio');
						var files = this.files;
						for (var i=0; i<files.length; i++) {
							if (maxSize > 0 && files[i].size > maxSize) {
								alert('Arquivo muito grande: ' + files[i].name + '. O limite é de ' + (maxSize / 1024 / 1024) + 'MB.');
								this.value = ''; // Limpa o input
								return;
							}
						}
					});
				})();
			</script>
		</div>`, f.Class, id, styleStr, label, id, acceptAttr, multiAttr, renderAttrs(f.Attrs), id, id, f.Accept, f.MaxSize)
}
