package components

// ThemeScript applies the saved theme before the first paint (no flash) and
// exposes gouiToggleTheme() globally. Icons are driven by CSS, not JS.
const ThemeScript = `<script>
(function(){var t=localStorage.getItem('goui-theme');if(t)document.documentElement.setAttribute('data-theme',t);})();
function gouiToggleTheme(){var el=document.documentElement;var dark=el.getAttribute('data-theme')==='dark'||(el.getAttribute('data-theme')!=='light'&&window.matchMedia('(prefers-color-scheme:dark)').matches);var next=dark?'light':'dark';el.setAttribute('data-theme',next);localStorage.setItem('goui-theme',next);}
function gouiToggleGroup(btn){var items=btn.nextElementSibling;var open=items.classList.toggle('open');btn.setAttribute('aria-expanded',open);}
function gouiToggleSidebar(){var sb=document.querySelector('.goui-sidebar');if(sb)sb.classList.toggle('mobile-open');}
window.goui = {
	action: async function(id, params = {}) {
		const url = new URL('/api/goui/action', window.location.origin);
		url.searchParams.append('id', id);
		for (let key in params) { 
			if (Array.isArray(params[key])) {
				params[key].forEach(v => url.searchParams.append(key, v));
			} else {
				url.searchParams.append(key, params[key]); 
			}
		}
		const r = await fetch(url);
		if (!r.ok) throw new Error('Falha na ação: ' + id);
		return await r.text();
	},
	selectOption: function(id, val, label) {
		const input = document.getElementById(id);
		if (input) {
			input.value = val;
			input.dispatchEvent(new Event('change', { bubbles: true }));
		}
		const btn = document.getElementById('btn-' + id);
		if (btn) {
			const txt = btn.querySelector('.selected-text');
			if (txt) txt.innerText = label;
		}
		const menu = document.getElementById('menu-' + id);
		if (menu) menu.style.display = 'none';
	}
};
window.gouiSelectOption = window.goui.selectOption;
function gouiToast(msg, type='info'){
	var cont = document.getElementById('goui-toast-container');
	if(!cont){
		cont = document.createElement('div');
		cont.id = 'goui-toast-container';
		cont.className = 'goui-toast-container';
		document.body.appendChild(cont);
	}
	var t = document.createElement('div');
	t.className = 'goui-toast goui-toast-' + type;
	t.innerHTML = '<span>' + msg + '</span>';
	cont.appendChild(t);
	setTimeout(() => { t.classList.add('show'); }, 10);
	setTimeout(() => { t.classList.remove('show'); setTimeout(() => t.remove(), 500); }, 4000);
}
window.onerror = function(msg, url, line){ gouiToast('Erro: ' + msg + ' (Linha ' + line + ')', 'error'); return false; };
window.onunhandledrejection = function(e){ gouiToast('Promessa Rejeitada: ' + e.reason, 'error'); };
</script>`

// Theme is the goui built-in CSS injected automatically via Headbar.
// Light: purple. Dark: black/gray — switches via OS setting or the toggle button.
const Theme = `
<style>
/* ============================================================
   goui Design System
   Light: Purple  |  Dark: Black & Gray
   Controlled by OS (prefers-color-scheme) or data-theme attr
   ============================================================ */

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --goui-bg: #f8fafc;
  --goui-surface: #ffffff;
  --goui-surface-2: #f1f5f9;
  --goui-border: #e2e8f0;
  --goui-border-code-block: #272727;
  --goui-text: #0f172a;
  --goui-text-muted: #64748b;
  --goui-primary: #7c3aed;
  --goui-primary-hover: #6d28d9;
  --goui-primary-dim: rgba(124, 58, 237, 0.1);
  --goui-radius: 8px;
}

/* --- Custom Scrollbar --- */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: var(--goui-border);
  border-radius: 10px;
  border: 2px solid var(--goui-surface);
}

::-webkit-scrollbar-thumb:hover {
  background: var(--goui-text-muted);
}

/* Firefox */
* {
  scrollbar-width: thin;
  scrollbar-color: var(--goui-border) transparent;
}

body {
  --goui-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --goui-font: 'Inter', -apple-system, system-ui, sans-serif;
  background: var(--goui-bg);
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: 16px;
  line-height: 1.6;
  transition: background .2s, color .2s;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

/* Dark via OS preference (unless user manually chose light) */
@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    --goui-primary:       #a78bfa;
    --goui-primary-hover: #8b5cf6;
    --goui-primary-dim:   rgba(167,139,250,.12);
    --goui-bg:            #0a0a0a;
    --goui-surface:       #141414;
    --goui-surface-2:     #1a1a1a;
    --goui-text:          #e5e5e5;
    --goui-text-muted:    #9ca3af;
    --goui-border:        #272727;
    --goui-shadow:        0 1px 4px rgba(0,0,0,.4);
  }
}

/* Dark forced by toggle button */
:root[data-theme="dark"] {
  --goui-primary:       #a78bfa;
  --goui-primary-hover: #8b5cf6;
  --goui-primary-dim:   rgba(167,139,250,.12);
  --goui-bg:            #0a0a0a;
  --goui-surface:       #141414;
  --goui-surface-2:     #1a1a1a;
  --goui-text:          #e5e5e5;
  --goui-text-muted:    #9ca3af;
  --goui-border:        #272727;
  --goui-shadow:        0 1px 4px rgba(0,0,0,.4);
}

body {
  background: var(--goui-bg);
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: 16px;
  line-height: 1.6;
  transition: background .2s, color .2s;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

/* --- Top Header (Headbar) --- */
.goui-header {
  flex-shrink: 0;
  background: var(--goui-primary);
  color: #fff;
  height: 52px;
  display: flex;
  align-items: center;
  padding: 0 1.5rem;
  box-shadow: var(--goui-shadow);
  position: sticky;
  top: 0;
  z-index: 200;
}

:root[data-theme="dark"] .goui-header {
  background: #111111;
  border-bottom: 1px solid var(--goui-border);
}

@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) .goui-header {
    background: #111111;
    border-bottom: 1px solid var(--goui-border);
  }
}

.goui-header-title {
  font-size: .95rem;
  font-weight: 600;
  letter-spacing: -.01em;
}

/* --- Layout (below header) --- */
.goui-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* --- Sidebar --- */
.goui-sidebar {
  width: 220px;
  flex-shrink: 0;
  background: var(--goui-surface);
  border-right: 1px solid var(--goui-border);
  display: flex;
  flex-direction: column;
  padding: 1.5rem 1rem;
  height: 100%;
  overflow-y: auto;
}

.goui-sidebar-logo {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--goui-primary);
  padding: 0 .5rem .5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid var(--goui-border);
  letter-spacing: -.02em;
}

.goui-sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: .15rem;
  flex: 1;
}

.goui-sidebar-link {
  color: var(--goui-text-muted);
  text-decoration: none;
  font-size: .875rem;
  font-weight: 500;
  padding: .5rem .75rem;
  border-radius: 6px;
  transition: color .15s, background .15s;
  display: block;
}

.goui-sidebar-link:hover {
  color: var(--goui-text);
  background: var(--goui-surface-2);
  text-decoration: none;
}

/* --- NavGroup (collapsible) --- */
.goui-nav-group { display: flex; flex-direction: column; }

.goui-nav-group-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: .5rem .75rem;
  background: none;
  border: none;
  border-radius: 6px;
  color: var(--goui-text-muted);
  font-family: var(--goui-font);
  font-size: .875rem;
  font-weight: 500;
  cursor: pointer;
  transition: color .15s, background .15s;
  text-align: left;
}

.goui-nav-group-btn:hover {
  color: var(--goui-text);
  background: var(--goui-surface-2);
}

.goui-nav-group-arrow {
  flex-shrink: 0;
  transition: transform .2s ease;
}

.goui-nav-group-btn[aria-expanded="true"] .goui-nav-group-arrow {
  transform: rotate(90deg);
}

.goui-nav-group-items {
  overflow: hidden;
  max-height: 0;
  transition: max-height .22s ease;
}

.goui-nav-group-items.open {
  max-height: 600px;
}

.goui-nav-group-inner {
  display: flex;
  flex-direction: column;
  gap: .05rem;
  padding: .2rem 0 .2rem .5rem;
}

.goui-sidebar-link-child {
  font-size: .825rem;
  padding: .4rem .75rem;
  border-left: 2px solid var(--goui-border);
  border-radius: 0 6px 6px 0;
}

.goui-sidebar-link-child:hover {
  border-left-color: var(--goui-primary);
}

/* ── Nested NavGroup (NavGroup inside NavGroup) ── */
.goui-nav-group-inner > .goui-nav-group > .goui-nav-group-btn {
  font-size: .8rem;
  padding: .38rem .75rem .38rem 1.1rem;
  color: var(--goui-text-muted);
}

.goui-nav-group-inner > .goui-nav-group > .goui-nav-group-items
  .goui-sidebar-link-child {
  padding-left: 2rem;
  font-size: .8rem;
}

.goui-sidebar-footer {
  padding-top: 1rem;
  border-top: 1px solid var(--goui-border);
}

/* --- Theme toggle (sidebar) --- */
.goui-code-block {
    background: #1e1e1e;
    color: #d4d4d4;
    padding: 1.5rem;
    border-radius: 8px;
    font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-size: 15px;
    line-height: 1.6;
    overflow-x: auto;
    border: 1px solid #333;
    margin: 0;
}

.goui-code-block code {
    font-family: inherit;
    font-size: inherit;
}

.goui-theme-toggle {
  background: none;
  border: none;
  border-radius: 6px;
  color: var(--goui-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: .5rem;
  width: 100%;
  padding: .5rem .75rem;
  font-size: .8rem;
  font-weight: 500;
  font-family: var(--goui-font);
  transition: color .15s, background .15s;
}

.goui-theme-toggle:hover {
  color: var(--goui-text);
  background: var(--goui-surface-2);
}

.goui-theme-toggle svg { display: block; flex-shrink: 0; }

/* Sun icon = shown in dark mode; Moon icon = shown in light mode */
.goui-icon-sun          { display: none;  }
.goui-icon-moon         { display: block; }
.goui-theme-label-sun   { display: none;  }
.goui-theme-label-moon  { display: block; }

:root[data-theme="dark"] .goui-icon-sun         { display: block; }
:root[data-theme="dark"] .goui-icon-moon        { display: none;  }
:root[data-theme="dark"] .goui-theme-label-sun  { display: block; }
:root[data-theme="dark"] .goui-theme-label-moon { display: none;  }

@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) .goui-icon-sun         { display: block; }
  :root:not([data-theme="light"]) .goui-icon-moon        { display: none;  }
  :root:not([data-theme="light"]) .goui-theme-label-sun  { display: block; }
  :root:not([data-theme="light"]) .goui-theme-label-moon { display: none;  }
}

/* --- Main content --- */
.goui-main {
  flex: 1;
  min-width: 0;
  padding: 3rem 4rem;
  overflow-y: auto;
}

/* --- Typography --- */
h1, h2, h3, h4, h5, h6 {
  color: var(--goui-text);
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -.02em;
}

h1 { font-size: 2.25rem; }
h2 { font-size: 1.75rem; }
h3 { font-size: 1.35rem; }
h4 { font-size: 1.1rem;  }

p { color: var(--goui-text-muted); line-height: 1.7; }

a { color: var(--goui-primary); text-decoration: none; }
a:hover { text-decoration: underline; }

/* --- Card --- */
.goui-card {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  padding: 2.5rem;
  box-shadow: 0 4px 12px rgba(0,0,0,.05);
}

/* --- Section / Hero --- */
.goui-section {
  background: var(--goui-surface-2);
  border-radius: var(--goui-radius);
  padding: 5rem 3rem;
  text-align: center;
  margin-bottom: 2rem;
}

/* --- Buttons --- */
.goui-btn {
  display: inline-flex;
  align-items: center;
  gap: .45rem;
  padding: .55rem 1.25rem;
  font-size: .875rem;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  text-decoration: none !important;
  transition: background .15s, transform .1s, box-shadow .15s;
}

.goui-btn:active { transform: translateY(1px); }
.goui-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

.goui-btn-primary {
  background: var(--goui-primary);
  color: #fff !important;
  box-shadow: 0 2px 8px var(--goui-primary-dim);
}

.goui-btn-primary:hover { background: var(--goui-primary-hover); }

.goui-btn-danger {
  background: #ef4444;
  color: #fff;
  border: 1px solid #dc2626;
}
.goui-btn-danger:hover { background: #dc2626; }

.goui-btn-ghost {
  background: var(--goui-primary-dim);
  color: var(--goui-primary) !important;
}

.goui-btn-ghost:hover {
  background: var(--goui-primary);
  color: #fff !important;
}

.goui-btn-block {
  display: flex;
  width: 100%;
  justify-content: center;
}

/* --- Badge Variants --- */
.goui-badge {
  display: inline-block;
  font-size: .75rem;
  font-weight: 600;
  padding: .2rem .65rem;
  border-radius: 999px;
}

.goui-badge-default { background: var(--goui-surface-2); color: var(--goui-text); }
.goui-badge-success { background: #dcfce7; color: #166534; }
.goui-badge-error   { background: #fee2e2; color: #991b1b; }
.goui-badge-warning { background: #fef3c7; color: #92400e; }
.goui-badge-info    { background: #dbeafe; color: #1e40af; }

/* Dark mode overrides for badges */
:root[data-theme="dark"] .goui-badge-success { background: #064e3b; color: #34d399; }
:root[data-theme="dark"] .goui-badge-error   { background: #7f1d1d; color: #f87171; }
:root[data-theme="dark"] .goui-badge-warning { background: #78350f; color: #fbbf24; }
:root[data-theme="dark"] .goui-badge-info    { background: #1e3a8a; color: #60a5fa; }

/* --- Utilities --- */
.goui-text-center    { text-align: center; }
.goui-text-muted     { color: var(--goui-text-muted) !important; }
.goui-mt-1 { margin-top: .5rem;  }
.goui-mt-2 { margin-top: 1rem;   }
.goui-mt-3 { margin-top: 1.5rem; }
.goui-mt-4 { margin-top: 2rem;   }
.goui-flex           { display: flex; }
.goui-gap-1          { gap: .5rem; }
.goui-gap-2          { gap: 1rem;  }
.goui-items-center   { align-items: center; }
.goui-justify-center { justify-content: center; }

/* --- Dropdown / Select --- */
/* --- Dropdown / Select --- */
.goui-dropdown {
  position: relative;
  display: block; /* Mudar para block ajuda na estabilidade do grid */
  width: 100%;
  z-index: 10;
}

.goui-dropdown-arrow {
  position: absolute;
  right: .85rem;
  top: 50%;
  margin-top: -8px; /* Usar margem em vez de transform:translateY evita bugs de renderização em alguns browsers */
  height: 16px;
  pointer-events: none;
  color: var(--goui-text-muted);
  display: flex;
  align-items: center;
  transition: transform .2s ease;
}

.goui-select {
  appearance: none;
  -webkit-appearance: none;
  width: 100%;
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: .95rem;
  padding: .7rem 2.5rem .7rem .9rem;
  cursor: pointer;
  outline: none;
  transition: all .2s ease;
  box-shadow: inset 0 1px 2px rgba(0,0,0,0.05);
}

.goui-select:focus {
  border-color: var(--goui-primary);
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

.goui-select:hover {
  border-color: var(--goui-primary);
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
}

.goui-select:hover .goui-dropdown-arrow {
  color: var(--goui-primary);
  transform: translateY(2px);
}

/* --- Custom Dropdown Menu --- */
.goui-dropdown-menu {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  padding: 0.5rem;
  list-style: none;
  max-height: 250px;
  overflow-y: auto;
}

.goui-dropdown-item {
  padding: 0.6rem 0.8rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--goui-text);
  transition: all 0.2s;
}

.goui-dropdown-item:hover {
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
}

.goui-dropdown-item.active {
  background: var(--goui-primary);
  color: #fff;
}

/* --- Custom Calendar --- */
.goui-calendar-popup {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
  padding: 1rem;
  width: 280px;
}

.goui-calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  font-weight: 600;
  gap: 5px;
}

.goui-calendar-header select {
  appearance: none;
  background: var(--goui-bg);
  border: 1px solid var(--goui-border);
  border-radius: 6px;
  padding: 4px 28px 4px 12px;
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--goui-text);
  cursor: pointer;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='14' height='14' viewBox='0 0 24 24' fill='none' stroke='%23888888' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
  transition: all 0.2s ease;
}

.goui-calendar-header select:hover {
  border-color: var(--goui-primary);
  background-color: var(--goui-primary-dim);
}

/* Reduzir fonte e compactar dropdowns dentro do cabeçalho do calendário */
.goui-calendar-header .goui-select {
  font-size: 0.8rem;
  padding: 4px 25px 4px 10px;
}

.goui-calendar-header .goui-dropdown-item {
  font-size: 0.75rem;
  padding: 0.4rem 0.6rem;
}

.goui-calendar-header button {
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
  border: none;
  border-radius: 4px;
  padding: 2px 8px;
  cursor: pointer;
  font-weight: bold;
}

.goui-calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.goui-calendar-day-head {
  text-align: center;
  font-size: 0.7rem;
  font-weight: bold;
  color: var(--goui-text-muted);
  padding-bottom: 5px;
  text-transform: uppercase;
}

.goui-calendar-day {
  text-align: center;
  padding: 0.5rem 0;
  cursor: pointer;
  border-radius: 4px;
  font-size: 0.8rem;
  transition: all 0.2s;
  background: rgba(0,0,0,0.02);
}

.goui-calendar-day:hover {
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
}

.goui-calendar-day.selected {
  background: var(--goui-primary) !important;
  color: #fff !important;
  font-weight: bold;
}

.goui-calendar-day.in-range {
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
  border-radius: 0;
}

.goui-calendar-day.today {
  border: 1px solid var(--goui-primary);
  color: var(--goui-primary);
  font-weight: bold;
}

.goui-calendar-day.outside {
  color: var(--goui-text-muted);
  opacity: 0.3;
}

.goui-calendar-footer {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--goui-border);
  font-size: 0.75rem;
  text-align: center;
  color: var(--goui-text-muted);
}

/* --- Toast Notifications --- */
.goui-toast-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.goui-toast {
  min-width: 280px;
  max-width: 400px;
  padding: 1rem 1.25rem;
  background: var(--goui-surface);
  color: var(--goui-text);
  border-radius: var(--goui-radius);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  border-left: 4px solid var(--goui-primary);
  display: flex;
  align-items: center;
  transform: translateX(120%);
  opacity: 0;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  backdrop-filter: blur(10px);
}

.goui-toast.show {
  transform: translateX(0);
  opacity: 1;
}

.goui-toast-success { border-left-color: #10b981; }
.goui-toast-error   { border-left-color: #ef4444; background: rgba(239, 68, 68, 0.05); }
.goui-toast-warning { border-left-color: #f59e0b; }
.goui-toast-info    { border-left-color: #3b82f6; }

/* --- Input with Icons --- */
.goui-input-icon-wrapper {
  position: relative;
  width: 100%;
  display: flex;
  align-items: center;
}

.goui-input-icon {
  position: absolute;
  right: 0.9rem;
  color: var(--goui-text-muted);
  pointer-events: none;
  display: flex;
  align-items: center;
  z-index: 2;
}

.goui-calendar-input {
  padding-right: 2.5rem !important; /* Abre espaço para o ícone */
}

/* Multiselect — modern card style */
.goui-dropdown-multi .goui-select {
  padding: 8px;
  min-height: 160px;
  border-radius: var(--goui-radius);
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  box-shadow: var(--goui-shadow);
  display: block;
}

.goui-select[multiple] option {
  padding: 10px 14px;
  margin: 4px 0;
  border-radius: 6px;
  color: var(--goui-text);
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.goui-select[multiple] option:hover {
  background: var(--goui-surface-2);
}

.goui-select[multiple] option:checked {
  background: var(--goui-primary) !important;
  color: #fff !important;
  box-shadow: 0 2px 4px rgba(124, 58, 237, 0.3);
}

/* Custom Scrollbar for Multiselect */
.goui-select[multiple]::-webkit-scrollbar {
  width: 6px;
}
.goui-select[multiple]::-webkit-scrollbar-track {
  background: transparent;
}
.goui-select[multiple]::-webkit-scrollbar-thumb {
  background: var(--goui-border);
  border-radius: 10px;
}
.goui-select[multiple]::-webkit-scrollbar-thumb:hover {
  background: var(--goui-text-muted);
}

/* --- Icon --- */
.goui-icon {
  display: inline-block;
  width: 24px;
  height: 24px;
  vertical-align: middle;
  fill: currentColor;
  stroke: currentColor;
}

/* --- Form inputs --- */
.goui-label {
  display: block;
  font-size: .72rem;
  font-weight: 600;
  color: var(--goui-text-muted);
  text-transform: uppercase;
  letter-spacing: .06em;
  margin-bottom: .35rem;
}

.goui-input {
  width: 100%;
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: .95rem;
  padding: .7rem .9rem;
  outline: none;
  transition: all .2s ease;
  box-shadow: inset 0 1px 2px rgba(0,0,0,0.05);
}

.goui-input:focus {
  border-color: var(--goui-primary);
  background: var(--goui-surface);
  box-shadow: 0 0 0 3px var(--goui-primary-dim), inset 0 1px 2px rgba(0,0,0,0.02);
}

.goui-input::placeholder { color: var(--goui-text-muted); opacity: .5; }

/* --- Playground --- */
.goui-playground {
  display: grid;
  grid-template-columns: 280px 1fr;
  grid-template-areas: 
    "controls preview"
    "controls code";
  gap: 1.5rem;
  align-items: start;
}

.goui-pg-controls { 
  grid-area: controls;
  display: flex;
  flex-direction: column;
  gap: .9rem;
}

.goui-pg-preview { 
  grid-area: preview;
  display: flex;
  flex-direction: column;
  gap: .5rem;
  margin-bottom: 1rem;
}

.goui-code-section {
  grid-area: code;
}

.goui-pg-preview-label {
  font-size: .72rem;
  font-weight: 600;
  color: var(--goui-text-muted);
  text-transform: uppercase;
  letter-spacing: .06em;
}

.goui-pg-preview-box {
  background: var(--goui-bg);
  border: 1px dashed var(--goui-border);
  border-radius: var(--goui-radius);
  padding: 2.5rem 1.5rem;
  min-height: 160px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* --- Code block --- */
.goui-code-section {
  width: 100%;
  margin: 0;
}
/*
.goui-code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: .45rem .9rem;
  background: var(--goui-surface-2);
  border: 1px solid var(--goui-border);
  border-bottom: none;
  border-radius: var(--goui-radius) var(--goui-radius) 0 0;
}
*/
.goui-code-header-label {
  font-size: .72rem;
  font-weight: 600;
  color: var(--goui-text-muted);
  text-transform: uppercase;
  letter-spacing: .06em;
}

.goui-code-block {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border-code-block);
  border-radius: 0 0 var(--goui-radius) var(--goui-radius);
  padding: 1.1rem 1.4rem;
  font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Menlo', monospace;
  font-size: .82rem;
  line-height: 1.75;
  overflow-x: auto;
  margin: 0;
  white-space: pre;
  color: var(--goui-text);
}

/* Go syntax highlight tokens */
.gc-pkg   { color: #7c3aed; font-weight: 500; }
.gc-fn    { color: #2563eb; }
.gc-str   { color: #059669; }
.gc-key   { color: #d97706; }
.gc-brace { color: var(--goui-text-muted); }

:root[data-theme="dark"] .gc-pkg  { color: #a78bfa; }
:root[data-theme="dark"] .gc-fn   { color: #60a5fa; }
:root[data-theme="dark"] .gc-str  { color: #34d399; }
:root[data-theme="dark"] .gc-key  { color: #fbbf24; }

@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) .gc-pkg  { color: #a78bfa; }
  :root:not([data-theme="light"]) .gc-fn   { color: #60a5fa; }
  :root:not([data-theme="light"]) .gc-str  { color: #34d399; }
  :root:not([data-theme="light"]) .gc-key  { color: #fbbf24; }
}

/* --- Classes reference table --- */
.goui-classes-table {
  width: 100%;
  border-collapse: collapse;
  font-size: .875rem;
}

.goui-classes-table thead {
  background: var(--goui-surface-2);
}

.goui-classes-table th {
  text-align: left;
  padding: .6rem 1rem;
  font-size: .72rem;
  font-weight: 600;
  color: var(--goui-text-muted);
  text-transform: uppercase;
  letter-spacing: .06em;
  border-bottom: 1px solid var(--goui-border);
}

.goui-classes-table td {
  padding: .75rem 1rem;
  border-bottom: 1px solid var(--goui-border);
  color: var(--goui-text);
  vertical-align: middle;
}

.goui-classes-table tr:last-child td { border-bottom: none; }

.goui-classes-table tr:hover td { background: var(--goui-surface-2); }

.goui-class-tag {
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: .8rem;
  padding: .15rem .5rem;
  border-radius: 4px;
  white-space: nowrap;
}

/* --- Toggle Switch --- */
.goui-toggle {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  cursor: pointer;
  flex-shrink: 0;
}

.goui-toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.goui-toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--goui-border);
  transition: .2s;
  border-radius: 24px;
}

.goui-toggle-slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: #fff;
  transition: .2s;
  border-radius: 50%;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

input:checked + .goui-toggle-slider {
  background-color: var(--goui-primary);
}

input:focus + .goui-toggle-slider {
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

input:checked + .goui-toggle-slider:before {
  transform: translateX(20px);
}

/* --- Slider --- */
.goui-slider-container {
  width: 100%;
  max-width: 300px;
  padding: 10px 0;
}

.goui-slider {
  -webkit-appearance: none;
  width: 100%;
  height: 6px;
  background: var(--goui-border);
  outline: none;
  border-radius: 3px;
  transition: background .2s;
  cursor: pointer;
}

.goui-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  background: var(--goui-primary);
  cursor: pointer;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  transition: transform .1s;
}

.goui-slider::-webkit-slider-thumb:hover {
  transform: scale(1.1);
}

.goui-slider::-moz-range-thumb {
  width: 18px;
  height: 18px;
  background: var(--goui-primary);
  cursor: pointer;
  border-radius: 50%;
  border: 2px solid #fff;
}

/* --- File Uploader / Dropzone --- */
.goui-uploader {
  border: 2px dashed var(--goui-border);
  border-radius: var(--goui-radius);
  padding: 2rem;
  text-align: center;
  background: var(--goui-surface);
  transition: all 0.2s;
  cursor: pointer;
  position: relative;
}

.goui-uploader:hover, .goui-uploader.dragover {
  border-color: var(--goui-primary);
  background: var(--goui-primary-dim);
}

.goui-uploader-icon {
  color: var(--goui-primary);
  margin-bottom: 1rem;
}

.goui-uploader-text {
  font-size: 0.9rem;
  color: var(--goui-text-muted);
}

.goui-uploader-input {
  position: absolute;
  top: 0; left: 0; width: 100%; height: 100%;
  opacity: 0;
  cursor: pointer;
}

/* --- Centered layout (LayoutCentered) --- */
.goui-body-centered {
  max-width: 900px;
  margin-left: auto;
  margin-right: auto;
  width: 100%;
}

/* --- Narrow layout (LayoutNarrow) — Streamlit-style organic feel for large screens --- */
.goui-body-narrow {
  max-width: 730px;
  margin-left: auto;
  margin-right: auto;
  width: 100%;
}

@media (min-width: 1400px) {
  .goui-body-narrow  { max-width: 760px; }
  .goui-body-centered { max-width: 960px; }
}

@media (min-width: 1800px) {
  .goui-body-narrow  { max-width: 1000px; }
  .goui-body-centered { max-width: 1200px; }
}

/* --- Goui Utilities (Our own Bootstrap) --- */
.goui-flex { display: flex !important; }
.goui-flex-col { display: flex !important; flex-direction: column !important; }
.goui-flex-row { display: flex !important; flex-direction: row !important; }
.goui-items-start { align-items: flex-start !important; }
.goui-items-center { align-items: center !important; }
.goui-items-end { align-items: flex-end !important; }
.goui-justify-between { justify-content: space-between !important; }

.goui-gap-10 { gap: 10px !important; }
.goui-gap-15 { gap: 15px !important; }
.goui-gap-20 { gap: 20px !important; }

.goui-mt-10 { margin-top: 10px !important; }
.goui-mt-20 { margin-top: 20px !important; }
.goui-mt-30 { margin-top: 30px !important; }
.goui-mb-10 { margin-bottom: 10px !important; }
.goui-mb-20 { margin-bottom: 20px !important; }

.goui-p-10 { padding: 10px !important; }
.goui-p-20 { padding: 20px !important; }

.goui-fw-bold { font-weight: bold !important; }
.goui-w-full { width: 100% !important; }
.goui-max-w-300 { max-width: 300px !important; }

.goui-card {
  background: var(--goui-surface-2);
  border-radius: 10px;
  padding: 20px;
  border: 1px solid var(--goui-border);
}

.goui-text-primary { color: var(--goui-primary) !important; }

/* --- Date Range --- */
.goui-date-range {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  padding: 8px 12px;
  border-radius: var(--goui-radius);
  display: inline-flex !important;
  align-items: center;
  gap: 12px;
}

.goui-date-range .goui-input {
  border: none !important;
  padding: 4px !important;
  background: transparent !important;
  box-shadow: none !important;
  width: auto !important;
}

.goui-date-range-separator {
  color: var(--goui-text-muted);
  font-weight: bold;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 1px;
}
/* --- DataTable --- */
.goui-table-container {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  overflow: hidden;
  position: relative;
  margin-top: 15px;
}

.goui-table-actions {
  display: flex;
  justify-content: flex-end;
  padding: 10px;
  background: var(--goui-surface-2);
  border-bottom: 1px solid var(--goui-border);
}

.goui-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.goui-table th {
  background: #f3f0ff;
  color: #553c9a;
  padding: 12px 15px;
  text-align: left;
  font-weight: bold;
  border-bottom: 2px solid #e9d8fd;
}

.goui-table td {
  padding: 10px 15px;
  border-bottom: 1px solid var(--goui-border);
}

.goui-table tr:last-child td {
  border-bottom: none;
}

.goui-table tr:hover {
  background: var(--goui-primary-dim);
}

.goui-table td[contenteditable="true"]:focus {
  outline: 2px solid var(--goui-primary);
  background: #fff;
}

/* --- Command Palette --- */
.goui-command-palette-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.5); backdrop-filter: blur(4px);
  z-index: 9999; display: flex; justify-content: center; padding-top: 15vh;
}
.goui-command-palette-modal {
  background: #1e1e1e; width: 100%; max-width: 600px; max-height: 400px;
  border-radius: 12px; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.5);
  border: 1px solid #333; overflow: hidden; display: flex; flex-direction: column;
}
.goui-command-palette-input {
  width: 100%; padding: 1.25rem; background: transparent; border: none;
  border-bottom: 1px solid #333; color: #fff; font-size: 1.1rem; outline: none;
}
.goui-command-palette-results { overflow-y: auto; padding: 0.5rem; }
.goui-command-palette-item {
  padding: 0.75rem 1rem; border-radius: 6px; cursor: pointer; color: #aaa;
}
.goui-command-palette-item.active { background: var(--goui-primary); color: #fff; }

/* --- Code Snippets (Forced Dark) --- */
.goui-code-section { 
  margin-top: 1.5rem; 
  border-radius: 8px; 
  overflow: hidden; 
  border: 1px solid #333; 
  background: #1e1e1e !important;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}
.goui-code-header { display: flex; align-items: center; justify-content: space-between; padding: 0.6rem 1.25rem; background: #252526; border-bottom: 1px solid #333; }
.goui-code-header-label { font-size: .72rem; font-weight: 600; color: #999; text-transform: uppercase; }
.goui-code-block { 
  background: #1e1e1e !important; 
  color: #d4d4d4 !important; 
  padding: 1.5rem; 
  margin: 0; 
  text-align: left !important; 
  overflow-x: auto;
  max-height: 400px;
  overflow-y: auto;
}
.goui-code-block code { background: transparent !important; color: inherit !important; white-space: pre; display: block; text-align: left; }

/* --- Section Groups --- */
.goui-section-group {
    margin-top: 3.5rem;
    padding: 1.5rem;
    border: 1px solid var(--goui-border);
    border-radius: var(--goui-radius);
    background: var(--goui-bg);
}
.goui-section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 2px solid var(--goui-primary-dim);
    padding-bottom: 0.5rem;
    margin-bottom: 1.5rem;
}
.goui-section-chevron {
    color: var(--goui-primary);
    font-size: 0.8rem;
    transition: transform 0.2s ease;
}
.goui-section-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

/* --- Badge variants --- */
.goui-badge-default { background: var(--goui-primary-dim); color: var(--goui-primary); }
.goui-badge-success { background: rgba(16,185,129,.12); color: #059669; }
.goui-badge-error   { background: rgba(239,68,68,.12);  color: #dc2626; }
.goui-badge-warning { background: rgba(245,158,11,.12); color: #d97706; }
.goui-badge-info    { background: rgba(59,130,246,.12); color: #2563eb; }

/* Extra color utilities */
.goui-text-error   { color: #dc2626 !important; }
.goui-text-success { color: #059669 !important; }
.goui-text-warning { color: #d97706 !important; }
.goui-text-info    { color: #2563eb !important; }

/* Extra spacing utilities */
.goui-mb-30 { margin-bottom: 30px !important; }
.goui-mb-40 { margin-bottom: 40px !important; }
.goui-mt-40 { margin-top: 40px !important; }

/* Grid utilities */
.goui-grid-2 { 
  display: grid !important; 
  grid-template-columns: 1fr 1fr !important; 
  gap: 20px !important;
  align-items: start;
}

.goui-grid-2 > * { min-width: 0; } /* Evita que o código estoure a coluna */

@media (max-width: 900px) {
  .goui-grid-2 { grid-template-columns: 1fr !important; }
}

.goui-grid-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 1.5rem; }
@media (max-width: 768px) {
  .goui-grid-2, .goui-grid-3 { grid-template-columns: 1fr; }
}

/* Fix playground code section grid placement */
.goui-playground .goui-code-section { grid-area: code; margin-top: 0; }

/* --- Tabs / Pill Menu --- */
.goui-tabs {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 2rem;
  padding: 4px;
}

.goui-tab {
  padding: 8px 20px;
  border-radius: 99px;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--goui-text-muted);
  transition: all 0.2s ease;
}

.goui-tab:hover {
  color: var(--goui-text);
  background: var(--goui-surface-2);
}

.goui-tab.active {
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
}


/* --- Layout Shell (Standard Page) --- */
.goui-layout {
  display: flex;
  min-height: 100vh;
  width: 100%;
}

.goui-sidebar {
  width: 260px;
  min-width: 260px;
  height: 97vh;
  position: sticky;
  top: 0;
  background: var(--goui-surface);
  border-right: 1px solid var(--goui-border);
  display: flex;
  flex-direction: column;
}

.goui-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: var(--goui-surface-2);
  justify-content: flex-start; /* Centraliza verticalmente */
  align-items: center;    /* Centraliza horizontalmente */
}

/* --- Centered layout (LayoutCentered) --- */
.goui-body-centered {
  max-width: 900px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.goui-header {
  height: 64px;
  min-height: 64px;
  display: flex;
  align-items: center;
  padding: 0 2rem;
  background: var(--goui-primary);
  color: white;
}

.goui-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

.goui-header-title {
  font-weight: 600;
  font-size: 1.1rem;
}

.goui-sidebar-title {
  padding: 1.5rem 1.5rem 1rem 1.5rem;
  font-weight: 700;
  font-size: 1.5rem;
  color: var(--goui-primary);
}

.goui-sidebar-nav {
  flex: 1;
  padding: 1rem 0;
  overflow-y: auto;
}

.goui-sidebar-footer {
  padding: 1rem;
  border-top: 1px solid var(--goui-border);
}

/* --- Grid System --- */
.goui-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 20px;
  width: 100%;
}
.goui-col-1 { grid-column: span 1; }
.goui-col-2 { grid-column: span 2; }
.goui-col-3 { grid-column: span 3; }
.goui-col-4 { grid-column: span 4; }
.goui-col-5 { grid-column: span 5; }
.goui-col-6 { grid-column: span 6; }
.goui-col-7 { grid-column: span 7; }
.goui-col-8 { grid-column: span 8; }
.goui-col-9 { grid-column: span 9; }
.goui-col-10 { grid-column: span 10; }
.goui-col-11 { grid-column: span 11; }
.goui-col-12 { grid-column: span 12; }

@media (max-width: 768px) {
  .goui-grid {
    grid-template-columns: repeat(1, 1fr);
    gap: 15px;
  }
  [class*="goui-col-"] {
    grid-column: span 1 !important;
  }
}

/* --- Responsive Mobile Rules --- */
@media (max-width: 768px) {
  .goui-sidebar {
    display: none; /* Hide sidebar on small screens for now */
  }
  .goui-layout {
    flex-direction: column;
  }
  .goui-main {
    padding: 1.25rem !important;
  }
  .goui-body-narrow, .goui-body-centered {
    max-width: 100% !important;
    padding: 0 10px;
  }
  .goui-pg-split {
    flex-direction: column !important;
  }
  .goui-pg-code, .goui-pg-preview {
    width: 100% !important;
    max-width: 100% !important;
  }
  .goui-header {
    padding: 0 1rem;
  }
  .goui-flex {
    flex-direction: column !important;
  }
  .goui-flex-row {
    flex-direction: column !important; /* Stack rows on mobile */
  }
  .goui-gap-20 {
    gap: 10px !important;
  }
  .goui-sidebar.mobile-open {
    display: flex !important;
    position: fixed;
    top: 52px;
    left: 0;
    width: 80%;
    max-width: 300px;
    height: calc(100vh - 52px);
    z-index: 1000;
    box-shadow: 20px 0 50px rgba(0,0,0,0.5);
    animation: slideIn 0.3s ease-out;
  }
  @keyframes slideIn {
    from { transform: translateX(-100%); }
    to { transform: translateX(0); }
  }
  .goui-mobile-menu-btn {
    display: flex !important;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    color: #fff;
    margin-right: 1rem;
    cursor: pointer;
  }
}

.goui-mobile-menu-btn { display: none; }

/* ─── Checkbox ─────────────────────────────────────────────────────────── */
.goui-checkbox-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
  font-size: .9rem;
  color: var(--goui-text);
}

.goui-checkbox-input {
  position: absolute;
  opacity: 0;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0,0,0,0);
}

.goui-checkbox-box {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  border: 2px solid var(--goui-border);
  border-radius: 4px;
  background: var(--goui-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all .15s ease;
}

.goui-checkbox-label:hover .goui-checkbox-box {
  border-color: var(--goui-primary);
}

.goui-checkbox-input:checked + .goui-checkbox-box {
  background: var(--goui-primary);
  border-color: var(--goui-primary);
}

.goui-checkbox-input:checked + .goui-checkbox-box::after {
  content: '';
  display: block;
  width: 5px;
  height: 9px;
  border: 2px solid #fff;
  border-top: none;
  border-left: none;
  transform: rotate(42deg) translateY(-1px);
}

.goui-checkbox-input:focus-visible + .goui-checkbox-box {
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

.goui-checkbox-input:disabled + .goui-checkbox-box {
  opacity: .4;
  cursor: not-allowed;
}

.goui-checkbox-group {
  display: flex;
  flex-direction: column;
  gap: .35rem;
}

.goui-checkbox-group-items {
  display: flex;
  flex-direction: column;
  gap: .5rem;
  margin-top: .4rem;
}

/* ─── Radio ─────────────────────────────────────────────────────────────── */
.goui-radio-group {
  display: flex;
  flex-direction: column;
  gap: .35rem;
}

.goui-radio-group-items {
  display: flex;
  flex-direction: column;
  gap: .5rem;
  margin-top: .4rem;
}

.goui-radio-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
  font-size: .9rem;
  color: var(--goui-text);
}

.goui-radio-input {
  position: absolute;
  opacity: 0;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0,0,0,0);
}

.goui-radio-circle {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  border: 2px solid var(--goui-border);
  border-radius: 50%;
  background: var(--goui-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all .15s ease;
}

.goui-radio-label:hover .goui-radio-circle {
  border-color: var(--goui-primary);
}

.goui-radio-input:checked + .goui-radio-circle {
  border-color: var(--goui-primary);
}

.goui-radio-input:checked + .goui-radio-circle::after {
  content: '';
  display: block;
  width: 8px;
  height: 8px;
  background: var(--goui-primary);
  border-radius: 50%;
}

.goui-radio-input:focus-visible + .goui-radio-circle {
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

.goui-radio-input:disabled + .goui-radio-circle {
  opacity: .4;
  cursor: not-allowed;
}

/* ─── Textarea ──────────────────────────────────────────────────────────── */
.goui-textarea {
  resize: none;
  overflow: hidden;
  min-height: 80px;
  line-height: 1.6;
}

/* ─── TagInput ──────────────────────────────────────────────────────────── */
.goui-taginput {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  padding: .5rem .75rem;
  min-height: 44px;
  transition: border-color .2s;
  cursor: text;
}

.goui-taginput:focus-within {
  border-color: var(--goui-primary);
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

.goui-taginput-chips {
  display: contents; /* chips flow inline with the text field */
}

.goui-taginput-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
  font-size: .8rem;
  font-weight: 500;
  padding: 2px 8px;
  border-radius: 999px;
  white-space: nowrap;
}

.goui-taginput-remove {
  background: none;
  border: none;
  color: var(--goui-primary);
  cursor: pointer;
  font-size: 1rem;
  line-height: 1;
  padding: 0 1px;
  display: flex;
  align-items: center;
  opacity: .7;
  transition: opacity .15s;
}

.goui-taginput-remove:hover { opacity: 1; }

.goui-taginput-field {
  border: none;
  outline: none;
  background: transparent;
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: .9rem;
  flex: 1;
  min-width: 100px;
}

.goui-taginput-field::placeholder { color: var(--goui-text-muted); opacity: .5; }

/* ─── ColorPicker ───────────────────────────────────────────────────────── */
.goui-colorpicker {
  display: inline-flex;
  align-items: center;
  gap: .75rem;
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: var(--goui-radius);
  padding: .45rem .85rem;
  cursor: pointer;
  transition: border-color .2s;
}

.goui-colorpicker:hover { border-color: var(--goui-primary); }

.goui-colorpicker-swatch {
  width: 26px;
  height: 26px;
  border-radius: 5px;
  border: 2px solid rgba(0,0,0,.1);
  flex-shrink: 0;
  cursor: pointer;
  transition: transform .15s;
}

.goui-colorpicker-swatch:hover { transform: scale(1.1); }

.goui-colorpicker-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
  pointer-events: none;
}

.goui-colorpicker-hex {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: .82rem;
  color: var(--goui-text-muted);
  letter-spacing: .04em;
}

/* ─── Form ──────────────────────────────────────────────────────────────── */
.goui-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* ─── Field error message ───────────────────────────────────────────────── */
.goui-field-error {
  display: block;
  font-size: .78rem;
  color: #dc2626;
  margin-top: .25rem;
}

/* Highlight invalid fields after user interaction */
.goui-input:user-invalid,
.goui-textarea:user-invalid {
  border-color: #dc2626;
  box-shadow: 0 0 0 3px rgba(220, 38, 38, .12);
}
</style>`
