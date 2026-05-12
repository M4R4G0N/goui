package components

import (
	"fmt"
)

type WatchAction string

const (
	WatchText        WatchAction = "text"
	WatchTag         WatchAction = "tag"
	WatchColor       WatchAction = "color"
	WatchSize        WatchAction = "size"
	WatchWeight      WatchAction = "weight"
	WatchClass       WatchAction = "class"
	WatchPlaceholder WatchAction = "placeholder"
	WatchValue       WatchAction = "value"
	WatchType        WatchAction = "type"
)

type WatchOption struct {
	SourceID string
	Action   WatchAction
}

func Watch(source interface{ GetID() string }, action WatchAction) WatchOption {
	return WatchOption{SourceID: source.GetID(), Action: action}
}

type BindOption struct {
	Template string
	Sources  map[string]string
}

func Bind(template string, sources map[string]string) BindOption {
	return BindOption{Template: template, Sources: sources}
}

// SyncText sincroniza o valor de um input com o texto de um elemento.
func SyncText(sourceID, targetID, fallback string) Component {
	script := fmt.Sprintf(`<script>
		(function() {
			var src = document.getElementById('%s');
			var target = document.getElementById('%s');
			if (src && target) {
				var update = function() {
					target.innerText = src.value || '%s';
				};
				['input', 'change'].forEach(ev => src.addEventListener(ev, update));
				update();
			}
		})();
	</script>`, sourceID, targetID, fallback)
	return HTML(script)
}

// SyncRange sincroniza dois inputs (início/fim) com um label.
func SyncRange(startID, endID, labelID, separator string) Component {
	script := fmt.Sprintf(`<script>
		(function() {
			var s = document.getElementById('%s');
			var e = document.getElementById('%s');
			var l = document.getElementById('%s');
			if (s && e && l) {
				var up = function() { l.innerText = s.value + ' %s ' + e.value; };
				s.addEventListener('input', up);
				e.addEventListener('input', up);
				up();
			}
		})();
	</script>`, startID, endID, labelID, separator)
	return HTML(script)
}

// SyncCSV sincroniza o upload de um CSV com uma DataTable (Client-side).
func SyncCSV(uploaderID, tableID string) Component {
	script := fmt.Sprintf(`<script>
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
						thead.innerHTML = ''; tbody.innerHTML = '';
						var sep = lines[0].indexOf(';') > -1 ? ';' : ',';
						var headers = lines[0].split(sep);
						headers.forEach(h => { var th = document.createElement('th'); th.innerText = h.trim(); thead.appendChild(th); });
						for (var i = 1; i < lines.length; i++) {
							if (lines[i].trim() === "") continue;
							var tr = document.createElement('tr');
							var cells = lines[i].split(sep);
							cells.forEach(c => { var td = document.createElement('td'); td.innerText = c.trim(); tr.appendChild(td); });
							tbody.appendChild(tr);
						}
					};
					reader.readAsText(file);
				});
			}
		})();
	</script>`, uploaderID, tableID)
	return HTML(script)
}

// SyncServer envia um arquivo para o backend e atualiza uma tabela com o resultado.
func SyncServer(uploaderID, tableID, apiPath string) Component {
	script := fmt.Sprintf(`<script>
		(function() {
			var up = document.getElementById('%s');
			var table = document.getElementById('%s');
			if (up && table) {
				up.addEventListener('change', function(e) {
					var file = e.target.files[0];
					if (!file) return;
					var formData = new FormData();
					formData.append('file', file);
					fetch('%s', { method: 'POST', body: formData })
					.then(r => r.text())
					.then(text => {
						var lines = text.split(/\r?\n/);
						var thead = table.querySelector('thead tr');
						var tbody = table.querySelector('tbody');
						thead.innerHTML = ''; tbody.innerHTML = '';
						var sep = lines[0].indexOf(';') > -1 ? ';' : ',';
						lines[0].split(sep).forEach(h => { var th = document.createElement('th'); th.innerText = h.trim(); thead.appendChild(th); });
						for (var i = 1; i < lines.length; i++) {
							if (lines[i].trim() === "") continue;
							var tr = document.createElement('tr');
							lines[i].split(sep).forEach(c => { var td = document.createElement('td'); td.innerText = c.trim(); tr.appendChild(td); });
							tbody.appendChild(tr);
						}
					});
				});
			}
		})();
	</script>`, uploaderID, tableID, apiPath)
	return HTML(script)
}

// SyncSubmit envia dados de múltiplos inputs para uma API.
func SyncSubmit(sourceIDs []string, apiPath, targetID string) Component {
	idsJSON := "["
	for i, id := range sourceIDs {
		if i > 0 {
			idsJSON += ","
		}
		idsJSON += fmt.Sprintf("'%s'", id)
	}
	idsJSON += "]"
	script := fmt.Sprintf(`<script>
		(function() {
			window.gouiSubmit = function() {
				var ids = %s;
				var data = {};
				ids.forEach(id => { var el = document.getElementById(id); if (el) data[id] = el.value; });
				fetch('%s', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data) })
				.then(r => r.text())
				.then(res => { var t = document.getElementById('%s'); if (t) t.innerText = res; });
			};
		})();
	</script>`, idsJSON, apiPath, targetID)
	return HTML(script)
}

// TableEditSync liga um Toggle ao modo de edição de uma tabela.
// Equivalente à lógica que costumava ficar como <script> nos comp_ de demonstração.
func TableEditSync(toggleID, tableID string) Component {
	return HTML(fmt.Sprintf(`<script>
(function(){
  var sw  = document.getElementById('%s');
  var tbl = document.getElementById('%s');
  if (!sw || !tbl) return;
  sw.addEventListener('change', function(e) {
    tbl.querySelectorAll('td').forEach(function(c) {
      c.contentEditable = e.target.checked;
    });
  });
})();
</script>`, toggleID, tableID))
}

// FetchInto registra a função JS gouiFetchInto(url, targetID) que faz fetch de uma
// URL e injeta o HTML resultante no elemento alvo — substituindo o conteúdo anterior.
func FetchInto() Component {
	return HTML(`<script>
if (typeof window.gouiFetchInto !== 'function') {
  window.gouiFetchInto = function(url, targetID) {
    fetch(url).then(function(r){ return r.text(); }).then(function(html){
      var el = document.getElementById(targetID);
      if (el) el.innerHTML = html;
    });
  };
}
</script>`)
}

// CalendarActionScript registra uma função JS que lê os valores de um CalendarRange
// e de um Calendar simples, e os envia via goui.action para o backend.
func CalendarActionScript(fnName, actionID, actionKey, fieldName, resultID, rangeCalID, singleCalID string) Component {
	return HTML(fmt.Sprintf(`<script>
window.%s = function() {
  var res = document.getElementById('%s');
  if (res) res.innerText = 'Go processando action...';
  var params = { action: '%s' };
  var sStart = document.getElementById('%s-start') ? document.getElementById('%s-start').value : '';
  var sEnd   = document.getElementById('%s-end')   ? document.getElementById('%s-end').value   : '';
  if (sStart && sEnd) {
    params['%s'] = [sStart, sEnd];
  } else {
    var single = document.getElementById('%s');
    params['%s'] = single ? single.value : '';
  }
  goui.action('%s', params).then(function(txt) {
    if (res) res.innerText = txt;
  });
};
</script>`, fnName, resultID, actionKey,
		rangeCalID, rangeCalID, rangeCalID, rangeCalID,
		fieldName, singleCalID, fieldName, actionID))
}
