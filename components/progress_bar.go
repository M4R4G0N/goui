package components

import (
	"fmt"
	"strings"
)

type ProgressVariant string

const (
	ProgressDefault ProgressVariant = ""
	ProgressSuccess ProgressVariant = "success"
	ProgressError   ProgressVariant = "error"
	ProgressWarning ProgressVariant = "warning"
	ProgressInfo    ProgressVariant = "info"
)

// Stream é a opção que ativa SSE em um ProgressBar.
// Sem ela, a barra é estática — ideal para valores pré-computados.
// Com ela, cada Add() faz broadcast automático e o Render() abre EventSource.
//
//	bar := ProgressBar(ProgressSuccess, Stream) // SSE ativado
//	bar := ProgressBar(ProgressSuccess)         // estático
type Stream struct{}

type ProgressBarComponent struct {
	ID         string
	Total      int
	Current    int
	Variant    ProgressVariant
	ShowLabel  bool
	sseEnabled bool
	Class      string
	Style      Style
	Attrs      Attr
}

func (p *ProgressBarComponent) GetID() string {
	if p.ID == "" {
		p.ID = AutoID()
	}
	return p.ID
}

// ProgressBar cria uma barra de progresso.
// Passe Stream para ativar SSE — cada Add() empurra o estado ao cliente automaticamente.
func ProgressBar(opts ...any) *ProgressBarComponent {
	pb := &ProgressBarComponent{ShowLabel: true}
	for _, opt := range opts {
		switch v := opt.(type) {
		case ProgressVariant:
			pb.Variant = v
		case Stream:
			pb.sseEnabled = true
		case bool:
			pb.ShowLabel = v
		case string:
			ParseStringAttr(v, &pb.Class, &pb.ID, &pb.Attrs)
		case Class:
			pb.Class = string(v)
		case ID:
			pb.ID = string(v)
		case Style:
			pb.Style = v
		case Attr:
			pb.Attrs = v
		}
	}
	return pb
}

// SetTotal define o valor máximo da barra.
func (p *ProgressBarComponent) SetTotal(n int) *ProgressBarComponent {
	p.Total = n
	return p
}

// SetVariant muda a variante de cor da barra em tempo de execução.
// Se SSE estiver ativo, a mudança é propagada ao frontend automaticamente no próximo Add().
func (p *ProgressBarComponent) SetVariant(v ProgressVariant) *ProgressBarComponent {
	p.Variant = v
	if p.sseEnabled {
		pct := 0
		if p.Total > 0 {
			pct = (p.Current * 100) / p.Total
			if pct > 100 {
				pct = 100
			}
		}
		SSEBroadcast(p.GetID(), fmt.Sprintf(
			`{"current":%d,"total":%d,"pct":%d,"variant":"%s"}`,
			p.Current, p.Total, pct, string(v),
		))
	}
	return p
}

// Add incrementa o progresso. Passo padrão: 1.
// O passo deve ser um inteiro positivo (> 0); panics caso contrário.
// Se a opção Stream foi passada, faz broadcast automático via SSE.
func (p *ProgressBarComponent) Add(n ...int) *ProgressBarComponent {
	step := 1
	if len(n) > 0 {
		if n[0] <= 0 {
			panic(fmt.Sprintf("goui: ProgressBar.Add requires a positive value (> 0), got %d", n[0]))
		}
		step = n[0]
	}
	p.Current += step
	if p.Total > 0 && p.Current > p.Total {
		p.Current = p.Total
	}

	if p.sseEnabled {
		pct := 0
		if p.Total > 0 {
			pct = (p.Current * 100) / p.Total
			if pct > 100 {
				pct = 100
			}
		}
		SSEBroadcast(p.GetID(), fmt.Sprintf(
			`{"current":%d,"total":%d,"pct":%d,"variant":"%s"}`,
			p.Current, p.Total, pct, string(p.Variant),
		))
	}

	return p
}

// ProgressBarOnDone retorna um script que reage ao evento "done" do ProgressBar via SSE.
// Mostra msgID por 2s depois reseta btnID — padrão idêntico a SyncText / SyncRange.
func ProgressBarOnDone(barID, btnID, msgID string) Component {
	return HTML(fmt.Sprintf(`<script>
(function(){
  document.addEventListener('goui:sse:done:%s', function() {
    var msg = document.getElementById('%s');
    var btn = document.getElementById('%s');
    if (msg) { msg.style.display = 'flex'; }
    setTimeout(function() {
      if (msg) msg.style.display = 'none';
      if (btn) { btn.disabled = false; btn.innerText = 'Reiniciar'; }
    }, 2000);
  });
})();
</script>`, barID, msgID, btnID))
}

func (p *ProgressBarComponent) Render() string {
	pct := 0
	if p.Total > 0 {
		pct = (p.Current * 100) / p.Total
		if pct > 100 {
			pct = 100
		}
	}

	fillClass := "goui-progress-fill"
	if p.Variant != "" {
		fillClass += " goui-progress-fill-" + string(p.Variant)
	}

	styleStr := ""
	if len(p.Style) > 0 {
		styleStr = fmt.Sprintf(` style="%s"`, strings.Join(p.Style.entries(), ";"))
	}

	attrs := renderAttrs(p.Attrs)
	id := p.GetID()
	fillID := id + "-fill"
	labelID := id + "-label"

	labelHTML := ""
	if p.ShowLabel {
		labelHTML = fmt.Sprintf(
			`<div class="goui-progress-label" id="%s"><span>%d / %d</span><span>%d%%</span></div>`,
			labelID, p.Current, p.Total, pct,
		)
	}

	// SSE script: emitido apenas quando Stream foi passado no construtor.
	// Abre um único EventSource para este componente. Eventos nomeados
	// são despachados como CustomEvent no document para outros scripts reagirem.
	sseScript := ""
	if p.sseEnabled {
		sseScript = fmt.Sprintf(`<script>
(function(){
  var fill  = document.getElementById('%s');
  var label = document.getElementById('%s');
  if (!fill) return;
  var es = new EventSource('/api/goui/stream?id=%s');
  es.onmessage = function(e) {
    try {
      var d = JSON.parse(e.data);
      fill.style.width = d.pct + '%%';
      if (d.variant !== undefined) {
        fill.className = 'goui-progress-fill' + (d.variant ? ' goui-progress-fill-' + d.variant : '');
      }
      if (label) label.innerHTML =
        '<span>' + d.current + ' / ' + d.total + '</span>' +
        '<span>' + d.pct + '%%</span>';
    } catch(_) {}
  };
  es.addEventListener('done', function(e) {
    document.dispatchEvent(new CustomEvent('goui:sse:done:%s', {detail: e.data}));
  });
  es.onerror = function(e) {
    if (es.readyState === EventSource.CLOSED) es.close();
  };
})();
</script>`, fillID, labelID, id, id)
	}

	return fmt.Sprintf(
		`<div class="goui-progress-container %s" data-goui-id="%s"%s%s>`+
			`<div class="goui-progress">`+
			`<div class="%s" id="%s" style="width:%d%%"></div>`+
			`</div>`+
			`%s`+
			`</div>`+
			`%s`,
		p.Class, id, styleStr, attrs,
		fillClass, fillID, pct,
		labelHTML,
		sseScript,
	)
}
