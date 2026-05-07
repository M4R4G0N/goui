package components

import (
	"fmt"
	"net/http"
	"reflect"
	"sync/atomic"
)

var (
	globalIDCounter int64
	// ComponentRegistry armazena instâncias vivas de componentes para busca posterior.
	ComponentRegistry = make(map[string]interface{})
	// AllIDs mantém a lista de todos os IDs gerados na sessão atual.
	AllIDs []string
	// ActionHandlers associa um ID a uma função de resposta (endpoint aleatório).
	ActionHandlers = make(map[string]func(r *http.Request) string)
	// GlobalActionHandler permite centralizar todas as ações em um único seletor no backend.
	GlobalActionHandler func(id, action string, r *http.Request) string
)

func init() {
	// Despachante genérico com Auto-Sync de Estado
	http.HandleFunc("/api/goui/action", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		action := r.FormValue("action")

		// ─────────────────────────────────────────────────────────────────────────
		// 🔴 TODO: EVOLUIR - FOCO EM RESPOSTA PARA BACKEND APENAS VIA FORMULÁRIO
		// ─────────────────────────────────────────────────────────────────────────

		// Sincronização Automática: Se o componente está no registro, atualizamos o Value dele
		if comp, ok := ComponentRegistry[id]; ok {
			// Usamos reflexão para ver se o componente tem um campo "Value" e atualizá-lo
			v := reflect.ValueOf(comp)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() == reflect.Struct {
				f := v.FieldByName("Value")
				if f.IsValid() && f.CanSet() {
					// Tentamos pegar o valor do campo "value" ou "campo_venda" ou o nome do componente
					newVal := r.FormValue("value")
					if newVal == "" {
						newVal = r.FormValue("q") // fallback
					}
					if newVal != "" {
						f.SetString(newVal)
					}
				}
			}
		}

		// 1. Tenta o handler individual do componente
		if handler, ok := ActionHandlers[id]; ok {
			w.Write([]byte(handler(r)))
			return
		}

		// 2. Tenta o seletor global (Controller Style)
		if GlobalActionHandler != nil {
			w.Write([]byte(GlobalActionHandler(id, action, r)))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Nenhum handler para ID %s ou Action %s", id, action)
	})
}

// Register adiciona um componente ao registro para introspecção.
func Register(id string, comp interface{}) {
	ComponentRegistry[id] = comp
}

// AutoID generates a unique HTML ID string for a component.
func AutoID() string {
	id := fmt.Sprintf("goui-comp-%d", atomic.AddInt64(&globalIDCounter, 1))
	AllIDs = append(AllIDs, id)
	return id
}

// RegisterAction vincula um comportamento dinâmico a um componente específico.
func RegisterAction(id string, handler func(r *http.Request) string) {
	ActionHandlers[id] = handler
}

// SyncSource is implemented by any component that can act as a reactive source for SyncWith.
type SyncSource interface {
	GetID() string
}

// Component represents any UI element that can be rendered to HTML.
type Component interface {
	Render() string
}

// SyncWith is a helper that implements SyncSource for a raw ID string.
type SyncWith string

func (s SyncWith) GetID() string { return string(s) }

// headMetaRenderer is satisfied by HeadbarComponent, which needs to inject
// both <head> meta content and a visible top bar.
type headMetaRenderer interface {
	RenderHead() string
}

// HTML is a generic wrapper for raw HTML strings.
type HTML string

func (h HTML) Render() string { return string(h) }

// PageLayout controls how wide the main content renders.
type PageLayout int

const (
	LayoutNarrow   PageLayout = iota // constrained to ~730px, Streamlit-style — default
	LayoutFull                       // fills available width
	LayoutCentered                   // constrained to ~900px and centered
)

// Page represents a full HTML document with a top header, sidebar, and main content.
type Page struct {
	Head    Component
	Nav     Component
	Body    Component
	Layout  PageLayout
	Palette Component
}

func NewPage(head, nav, body Component, opts ...any) *Page {
	p := &Page{Head: head, Nav: nav, Body: body}
	for _, opt := range opts {
		switch v := opt.(type) {
		case PageLayout:
			p.Layout = v
		case *CommandPaletteComponent:
			p.Palette = v
		}
	}
	return p
}

// LayoutToggle renders a button that toggles LayoutCentered on/off at runtime.
// When active it wraps the main content in goui-body-centered; when inactive reverts to goui-body-narrow.
// Accepts optional label strings: first = label when centered, second = label when expanded.
func LayoutToggle(opts ...any) Component {
	id := AutoID()
	labelOn := "Centralizar"
	labelOff := "Expandir"

	labels := []string{}
	for _, opt := range opts {
		if s, ok := opt.(string); ok {
			labels = append(labels, s)
		}
	}
	if len(labels) >= 1 {
		labelOn = labels[0]
	}
	if len(labels) >= 2 {
		labelOff = labels[1]
	}

	return HTML(fmt.Sprintf(`<button id="%s" class="goui-btn goui-btn-secondary"></button>
<script>(function(){
	var btn = document.getElementById('%s');
	var getWrapper = function(){ return document.querySelector('.goui-main > div'); };
	var centered = !!(getWrapper() && getWrapper().classList.contains('goui-body-centered'));
	var render = function(){
		btn.innerText = centered ? '%s' : '%s';
		btn.classList.toggle('goui-btn-primary', centered);
		btn.classList.toggle('goui-btn-secondary', !centered);
	};
	btn.addEventListener('click', function(){
		var w = getWrapper();
		if (!w) return;
		centered = !centered;
		w.classList.toggle('goui-body-centered', centered);
		w.classList.toggle('goui-body-narrow', !centered);
		render();
	});
	render();
})();</script>`, id, id, labelOn, labelOff))
}

func (p *Page) Render() string {
	// Separate <head> meta from the visible header bar.
	headMeta := p.Head.Render()
	headerBar := ""
	if hr, ok := p.Head.(headMetaRenderer); ok {
		headMeta = hr.RenderHead()
		headerBar = p.Head.Render()
	}

	bodyHTML := p.Body.Render()
	switch p.Layout {
	case LayoutCentered:
		bodyHTML = `<div class="goui-body-centered">` + bodyHTML + `</div>`
	case LayoutNarrow:
		bodyHTML = `<div class="goui-body-narrow">` + bodyHTML + `</div>`
	}

	paletteHTML := ""
	if p.Palette != nil {
		paletteHTML = p.Palette.Render()
	}

	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    ` + headMeta + `
</head>
<body>
    ` + headerBar + `
    <div class="goui-layout">
        <aside class="goui-sidebar">` + p.Nav.Render() + `</aside>
        <main class="goui-main">` + bodyHTML + `</main>
    </div>
    ` + paletteHTML + `
</body>
</html>`
}
