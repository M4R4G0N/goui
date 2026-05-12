# goUI

![Version](https://img.shields.io/badge/version-v0.2-blue.svg)
![Go](https://img.shields.io/badge/go-1.22%2B-%2300ADD8.svg?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**goUI** é uma biblioteca Go para construir interfaces web server-rendered usando um modelo de componentes — sem JavaScript obrigatório, sem build step, sem frameworks externos.

Escreva toda a sua UI em Go puro. Componentes, rotas, formulários e validação — tudo em um lugar só.

---

## Pré-requisitos

Go **1.22 ou superior**.

- **Windows**: [golang.org/dl](https://go.dev/dl/)
- **Linux / macOS**: [Guia oficial](https://go.dev/doc/install)

---

## Instalação

```bash
go mod init meuapp
go get github.com/M4R4G0N/goUI@v0.2
```

### Problemas de checksum ou 404

```bash
go clean -modcache
export GOPROXY=direct
export GONOSUMDB=github.com/M4R4G0N/goUI
go get github.com/M4R4G0N/goUI@v0.2
```

---

## Início Rápido

```go
package main

import (
    goui "github.com/M4R4G0N/goUI"
    "github.com/M4R4G0N/goUI/components"
)

func main() {
    app := goui.NewApp()

    app.RegisterRoute("/", components.NewPage(
        components.Headbar("Meu App"),
        components.Navbar("goUI",
            components.Link{Href: "/", Text: "Home"},
        ),
        components.Div(
            components.Text("Olá, goUI!", "h1"),
            components.Button("Clique aqui", components.Primary),
        ),
    ))

    app.Start("", "8080")
}
```

```bash
go run main.go
# → http://localhost:8080
```

---

## Reatividade Watch/Bind

Sincronize componentes entre si sem escrever JavaScript:

```go
func MyPage(title, path string) components.Component {
    input := components.Input(components.Placeholder("Digite algo..."))
    preview := components.Text(input, "h2") // sincroniza automaticamente

    return components.Div(input, preview)
}
```

---

## Componentes Disponíveis (v0.2)

### Texto & Conteúdo
| Componente | Descrição |
|------------|-----------|
| `Text` | Texto estático ou reativo com qualquer tag HTML |
| `Badge` | Pill de status com variantes de cor |
| `Icon` | Ícone Lucide renderizado como SVG inline |
| `Snippet` | Bloco de código com syntax highlight |
| `Textarea` | Campo de texto multilinha com auto-resize |
| `TagInput` | Input com chips/tags digitáveis |

### Botões & Controles
| Componente | Descrição |
|------------|-----------|
| `Button` | Primary, Secondary, Danger, Ghost |
| `Dropdown` | Select estilizado com seleção única ou múltipla |
| `Toggle` | Interruptor on/off |
| `Slider` | Range numérico com min/max/step |
| `ProgressBar` | Barra de progresso com `SetTotal` / `Add` e 5 variantes de cor |
| `Toast` | Notificações temporárias (success, error, warning, info) |
| `Checkbox` | Caixa de seleção individual |
| `CheckboxGroup` | Grupo com múltipla seleção |
| `RadioGroup` | Grupo de opções exclusivas |
| `ColorPicker` | Seletor de cor com swatch visual |

### Formulários
| Componente | Descrição |
|------------|-----------|
| `Input` | Campo de texto, email, senha, número, etc. |
| `Form` | Wrapper de formulário com proteção CSRF |
| `FormField` | Campo com label e texto de ajuda |
| `Validation` | Validação declarativa client-side |
| `ValidateForm` | Validação server-side por campo |
| `FieldError` | Mensagem de erro inline |

### Dados
| Componente | Descrição |
|------------|-----------|
| `Table` | Tabela estática com headers e linhas |
| `Calendar` | Seletor de data único |
| `CalendarRange` | Seletor de intervalo de datas |

### Estrutura
| Componente | Descrição |
|------------|-----------|
| `NewPage` | Documento HTML completo com sidebar e header |
| `Navbar` + `NavGroup` | Barra lateral com grupos colapsáveis e aninhados |
| `Headbar` | Barra superior com título |
| `Tabs` | Abas com ativação por hash de URL |
| `Card` | Container com estilo de card |
| `Section` | Seção com título e conteúdo agrupado |
| `CommandPalette` | Paleta de comandos com busca fuzzy (⌘K) |

### Upload
| Componente | Descrição |
|------------|-----------|
| `FileUploader` | Upload com drag & drop |
| `DownloadButton` | Botão que dispara download de arquivo |

---

## Formulário com CSRF

```go
// GET handler:
token := components.NewCSRFToken(w)

page := components.Form("/submit", "POST",
    components.CSRF(token),
    components.FormField("Nome",
        components.Input(
            components.Name("nome"),
            components.Validation{Required: true, MinLen: 2},
        ),
        "",
    ),
    components.Button("Enviar", components.Primary),
)

// POST handler:
if !components.ValidateCSRF(r) {
    http.Error(w, "Token inválido", 403)
    return
}

errs := components.ValidateForm(r, map[string]components.FieldRule{
    "nome": {Required: true, MinLen: 2},
})
```

---

## ProgressBar com SSE

Barra de progresso server-side com atualização em tempo real via **Server-Sent Events**.  
O frontend **nunca precisa fazer polling** — cada `bar.Add()` empurra o estado automaticamente para todos os clientes conectados.

```go
bar := components.ProgressBar(components.ProgressSuccess)
bar.SetTotal(100)

// Registra o gatilho: botão → HTTP → backend inicia goroutine
components.RegisterAction(btnID, func(r *http.Request) string {
    go func() {
        for i := 1; i <= bar.Total; i++ {
            bar.Add()                    // → SSE broadcast automático
            time.Sleep(25 * time.Millisecond)
        }
        // envia evento nomeado ao fim — frontend pode reagir
        components.SSEBroadcastEvent(bar.GetID(), "done", `{"done":true}`)
    }()
    return "started"
})
```

`bar.Render()` injeta um `<script>` com `EventSource` automaticamente — zero configuração manual no frontend.

### Como funciona

```
Botão click ──► goui.action(btnID)   [HTTP]   ──► backend inicia goroutine
                                                         │
                                                  bar.Add() loop
                                                         │ SSEBroadcast automático
                                     ◄───────── /api/goui/stream?id=X
EventSource (aberto pelo Render)                         │
  ├─ onmessage → atualiza fill + label no DOM            │
  └─ addEventListener("done") → dispara CustomEvent      ▼
                                                   fim do loop
```

### Regras do `Add`

`Add` aceita qualquer inteiro positivo. Valores `0` ou negativos causam **panic** imediato com mensagem clara.

```go
bar.Add()    // +1 (padrão)
bar.Add(10)  // +10
bar.Add(0)   // panic: value must be positive (> 0), got 0
bar.Add(-1)  // panic: value must be positive (> 0), got -1
```

### Variantes

`ProgressDefault` · `ProgressSuccess` · `ProgressError` · `ProgressWarning` · `ProgressInfo`

```go
components.ProgressBar(components.ProgressInfo, false) // sem label current/total/%
```

### SSE direta — sem ProgressBar

Para enviar qualquer dado via SSE a partir de qualquer componente:

```go
// broadcast simples
components.SSEBroadcast(id, `{"valor": 42}`)

// evento nomeado (o frontend ouve com addEventListener("meuEvento", ...))
components.SSEBroadcastEvent(id, "meuEvento", `{"ok": true}`)
```

---

## Estrutura do Projeto

```
goUI/
├── goui.go                 # App, NewApp, Start, RegisterRoute
├── router/
│   └── file_router.go      # RegisterPage, InjectRoutes
├── components/
│   ├── component.go        # Interface Component, AutoID, HTML, RegisterAction
│   ├── theme.go            # CSS global + ThemeScript (dark/light)
│   ├── text.go             # Text, Watch, Bind
│   ├── input.go            # Input, Validation
│   ├── button.go           # Button, Primary, Secondary, Danger, Ghost
│   ├── dropdown.go         # Dropdown, Option, Multi
│   ├── toggle.go           # Toggle
│   ├── slider.go           # Slider, Min, Max, Step
│   ├── progress_bar.go     # ProgressBar, SetTotal, Add, ProgressVariant
│   ├── checkbox.go         # Checkbox, CheckboxGroup, CheckboxItem
│   ├── radio.go            # RadioGroup, RadioItem
│   ├── textarea.go         # Textarea, Rows
│   ├── tag_input.go        # TagInput
│   ├── color_picker.go     # ColorPicker
│   ├── form.go             # Form, CSRF, Validation, ValidateForm, FieldError
│   ├── navbar.go           # Navbar, NavGroup, Link, NavItem
│   ├── tabs.go             # Tabs, TabItem (hash-based activation)
│   ├── badge.go            # Badge, SuccessBadge, ErrorBadge…
│   ├── icon.go             # Icon (Lucide)
│   ├── snippet.go          # Snippet (syntax highlight)
│   ├── table.go            # Table
│   ├── calendar.go         # Calendar, CalendarRange
│   ├── section.go          # Section
│   ├── card.go             # Card
│   ├── command_palette.go  # CommandPalette
│   ├── toast.go            # ToastContainer, ShowToast
│   ├── fileupload.go       # FileUploader
│   ├── download.go         # DownloadButton
│   └── script.go           # Watch, Bind, SyncText, SyncRange, SyncCSV…
└── examples/
    └── helloworld/         # Documentação interativa com playground
```

---

## File-Based Routing

Use `init()` para registrar páginas automaticamente:

```go
// pages/minha_pagina.go
package pages

import (
    "github.com/M4R4G0N/goUI/components"
    "github.com/M4R4G0N/goUI/router"
)

func init() {
    router.RegisterPage("/minha-rota", "Título — goUI", MinhaPagina)
}

func MinhaPagina(title, path string) components.Component {
    return components.Div(
        components.Text(title, "h1"),
    )
}
```

```go
// main.go
import _ "meuapp/pages" // dispara todos os init()

router.InjectRoutes(app, func(title, path string, body components.Component) components.Component {
    return components.NewPage(
        components.Headbar(title),
        components.Navbar("App", components.Link{Href: "/", Text: "Home"}),
        body,
    )
})
```

---

## Roadmap

Veja o roadmap completo com versões planejadas e checklist detalhado em [ROADMAP.md](./ROADMAP.md).

**Próximo: v0.3** — Grid/Columns, Modal, Drawer, Accordion, Breadcrumb, Pagination. SSE em desenvolvimento ativo.

---

## Licença

MIT — desenvolvido por [Marcelo Antonio Goncalves](https://github.com/M4R4G0N).
