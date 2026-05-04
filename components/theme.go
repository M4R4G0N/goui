package components

// ThemeScript applies the saved theme before the first paint (no flash) and
// exposes gouiToggleTheme() globally. Icons are driven by CSS, not JS.
const ThemeScript = `<script>
(function(){var t=localStorage.getItem('goui-theme');if(t)document.documentElement.setAttribute('data-theme',t);})();
function gouiToggleTheme(){var el=document.documentElement;var dark=el.getAttribute('data-theme')==='dark'||(el.getAttribute('data-theme')!=='light'&&window.matchMedia('(prefers-color-scheme:dark)').matches);var next=dark?'light':'dark';el.setAttribute('data-theme',next);localStorage.setItem('goui-theme',next);}
function gouiToggleGroup(btn){var items=btn.nextElementSibling;var open=items.classList.toggle('open');btn.setAttribute('aria-expanded',open);}
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
  --goui-primary:       #7c3aed;
  --goui-primary-hover: #6d28d9;
  --goui-primary-dim:   rgba(124,58,237,.12);
  --goui-bg:            #ffffff;
  --goui-surface:       #faf9ff;
  --goui-surface-2:     #f3f0ff;
  --goui-text:          #111827;
  --goui-text-muted:    #6b7280;
  --goui-border:        #e5e7eb;
  --goui-shadow:        0 1px 4px rgba(0,0,0,.08);
  --goui-radius:        10px;
  --goui-font:          -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
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
  padding: 2.5rem 2rem;
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
  padding: 2rem;
  box-shadow: var(--goui-shadow);
}

/* --- Section / Hero --- */
.goui-section {
  background: var(--goui-surface-2);
  border-radius: var(--goui-radius);
  padding: 3.5rem 2rem;
  text-align: center;
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

.goui-btn-primary {
  background: var(--goui-primary);
  color: #fff !important;
  box-shadow: 0 2px 8px var(--goui-primary-dim);
}

.goui-btn-primary:hover { background: var(--goui-primary-hover); }

.goui-btn-ghost {
  background: var(--goui-primary-dim);
  color: var(--goui-primary) !important;
}

.goui-btn-ghost:hover {
  background: var(--goui-primary);
  color: #fff !important;
}

/* --- Badge --- */
.goui-badge {
  display: inline-block;
  background: var(--goui-primary-dim);
  color: var(--goui-primary);
  font-size: .75rem;
  font-weight: 600;
  padding: .2rem .65rem;
  border-radius: 999px;
}

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
.goui-dropdown {
  position: relative;
  display: inline-block;
  width: 100%;
}

.goui-dropdown-arrow {
  position: absolute;
  right: .65rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: var(--goui-text-muted);
  display: flex;
  align-items: center;
}

.goui-select {
  appearance: none;
  -webkit-appearance: none;
  width: 100%;
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
  border-radius: 6px;
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: .875rem;
  padding: .55rem 2.25rem .55rem .75rem;
  cursor: pointer;
  outline: none;
  transition: border-color .15s, box-shadow .15s;
}

.goui-select:focus {
  border-color: var(--goui-primary);
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

.goui-select:hover {
  border-color: var(--goui-primary);
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
  background: var(--goui-bg);
  border: 1px solid var(--goui-border);
  border-radius: 6px;
  color: var(--goui-text);
  font-family: var(--goui-font);
  font-size: .875rem;
  padding: .5rem .75rem;
  outline: none;
  transition: border-color .15s, box-shadow .15s;
}

.goui-input:focus {
  border-color: var(--goui-primary);
  box-shadow: 0 0 0 3px var(--goui-primary-dim);
}

.goui-input::placeholder { color: var(--goui-text-muted); opacity: .6; }

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
  align-items: center;
  justify-content: center;
}

/* --- Code block --- */
.goui-code-section {
  margin-top: 1.25rem;
}

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

.goui-code-header-label {
  font-size: .72rem;
  font-weight: 600;
  color: var(--goui-text-muted);
  text-transform: uppercase;
  letter-spacing: .06em;
}

.goui-code-block {
  background: var(--goui-surface);
  border: 1px solid var(--goui-border);
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

/* --- Goui Utilities (Our own Bootstrap) --- */
.goui-flex { display: flex !important; }
.goui-flex-col { display: flex !important; flex-direction: column !important; }
.goui-flex-row { display: flex !important; flex-direction: row !important; }
.goui-items-center { align-items: center !important; }
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
  background: var(--goui-surface-2);
  padding: 12px 15px;
  text-align: left;
  font-weight: bold;
  border-bottom: 2px solid var(--goui-border);
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
.goui-code-section { margin-top: 1.5rem; border-radius: 8px; overflow: hidden; border: 1px solid #333; background: #1e1e1e !important; }
.goui-code-header { display: flex; align-items: center; justify-content: space-between; padding: 0.6rem 1.25rem; background: #252526; border-bottom: 1px solid #333; }
.goui-code-header-label { font-size: .72rem; font-weight: 600; color: #999; text-transform: uppercase; }
.goui-code-block { background: #1e1e1e !important; color: #d4d4d4 !important; padding: 1.5rem; margin: 0; text-align: left !important; }
.goui-code-block code { background: transparent !important; color: inherit !important; white-space: pre; display: block; text-align: left; }

</style>
`
