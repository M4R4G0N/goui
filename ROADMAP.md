# goUI — Roadmap

Acompanhamento de tudo que já foi construído `[x]` e do que ainda está planejado `[ ]`.

> Status atual: **v0.2.1** — componentes de formulário completos, proteção CSRF, validação declarativa e documentação interativa com playground.
> Em desenvolvimento ativo: **ProgressBar** (entregue) · **SSE / canal de streaming** (em andamento).

---

## Versões Lançadas

### v0.1.1 — Base Funcional
- [x] Interface `Component`, `AutoID`, registro global
- [x] `Input`, `Dropdown`, `Toggle`, `Slider`, `Button`, `DownloadButton`
- [x] `Text`, `Icon`, `Badge`, `Snippet`, `Table`
- [x] `Navbar`, `NavGroup`, `Headbar`, `Tabs`, `Card`, `Section`
- [x] `Calendar`, `CalendarRange`, `ParseDate`
- [x] `FileUploader`, `CommandPalette`, `ToastContainer`
- [x] Reatividade `Watch`/`Bind` client-side
- [x] Tema dark/light com `ThemeScript` + `Theme` CSS
- [x] Sistema de layout `NewPage`, `LayoutNarrow`, `LayoutFull`, `LayoutCentered`

### v0.2.1 — Formulários Completos & Documentação Interativa ✅
- [x] `Checkbox` — caixa de seleção individual
- [x] `CheckboxGroup` — grupo com múltipla seleção
- [x] `RadioGroup` com `RadioItem` — seleção exclusiva
- [x] `Textarea` com auto-resize via JS
- [x] `TagInput` — chips digitáveis com hidden CSV input
- [x] `ColorPicker` — swatch clicável + hex display
- [x] `Form(action, method, opts...)` — wrapper de formulário
- [x] Proteção CSRF nativa via HMAC-SHA256 (cookie double-submit)
- [x] `Validation` struct — validação declarativa client-side (required, minLen, maxLen, minNum, maxNum, pattern)
- [x] `ValidateForm(r, rules)` — validação server-side com retorno de erros por campo
- [x] `FieldError(msg)` — renderiza mensagem de erro inline
- [x] NavGroup aninhado — `NavItem` interface para sub-grupos na sidebar
- [x] Hash-based tab navigation — `data-target` + `gouiActivateHash` no DOMContentLoaded
- [x] `hashchange` listener para navegação via hash sem reload
- [x] Playground master com todos os componentes
- [x] Documentação interativa individual por componente (20+ páginas)

---

## Versões Planejadas

### v0.3 — Layout e Navegação Avançados
- [ ] `Grid` e `Columns` — layout declarativo flexível
- [ ] `Modal` e `Drawer` — sobreposições com foco gerenciado
- [ ] `Accordion` — seções expansíveis
- [ ] `Breadcrumb`, `Pagination` e `Stepper`
- [ ] Navbar responsiva com colapso automático em mobile

### v0.4 — Tabela de Dados Poderosa
- [ ] Ordenação e filtragem de colunas
- [ ] Seleção de linhas com checkbox
- [ ] Paginação integrada e carregamento server-side
- [ ] Exportação para CSV direto no cliente

### v0.5 — Exibição e Feedback Visual
- [ ] `Alert`, `Tooltip`, `Avatar`
- [x] `ProgressBar` — barra de progresso com `SetTotal` / `Add`, 5 variantes de cor, label `current/total/%`, validação de valor positivo
- [ ] `Skeleton` (loading state)
- [ ] `Stat` — card de métrica com variação (KPI)
- [ ] `Timeline`
- [ ] Múltiplos temas pré-definidos e `CustomTheme`

### v0.6 — Mídia e Captura
- [ ] `CameraCapture` e `PhotoCapture`
- [ ] `VideoRecorder` — gravação via câmera com controles
- [ ] `AudioRecorder` — gravação de áudio com waveform visual
- [ ] `ScreenCapture` — captura de tela ou aba do navegador
- [ ] Upload automático de mídia gravada para o backend
- [ ] Seleção de dispositivo de entrada e formatos configuráveis

### v0.3.1 — SSE / Streaming Server-Side (em andamento)
- [ ] Canal SSE robusto — `SSEStream` com registry global thread-safe (`sync.RWMutex`)
- [ ] Endpoint `/api/goui/stream?id=X` com `http.Flusher` e heartbeat automático
- [ ] `ProgressBar.Add()` empurra estado via SSE automaticamente (sem chamada separada)
- [ ] Cleanup automático de clientes desconectados
- [ ] JS nativo via `EventSource` — zero dependências externas
- [ ] Graceful shutdown de streams ao encerrar a aplicação

### v0.7 — Reatividade Server-Side (avançado)
- [ ] `SyncVisible` e `SyncDisabled` — visibilidade e estado por valor
- [ ] WebSocket nativo integrado ao sistema de componentes
- [ ] `WatchURL` — deep link via parâmetros de URL
- [ ] Estado persistente de componente entre requests (sessão ou cookie)

### v0.8 — DX e Ferramentas
- [ ] CLI `goui dev` com hot reload automático
- [ ] Roteamento baseado em arquivos (file-based routing) via CLI
- [ ] Parâmetros de rota dinâmicos (`/user/:id`)
- [ ] Suporte a testes de componente com `RenderToString`
- [ ] Middleware de pipeline (auth, logging, recovery antes do render)

### v0.9 — Acessibilidade e Qualidade
- [ ] Atributos ARIA injetados automaticamente por componente
- [ ] Navegação por teclado em `Dropdown`, `CommandPalette`, `Modal`
- [ ] Suporte a `prefers-reduced-motion`
- [ ] Exemplos completos: CRUD, dashboard em tempo real, autenticação

### v1.0 — Estável para Produção
- [ ] API pública estável e versionada
- [ ] Documentação completa com playground online
- [ ] Suite de testes de integração cobrindo todos os componentes
- [ ] Performance benchmark e otimizações de render
- [ ] Suporte a múltiplas instâncias de `App` no mesmo processo

---

## Checklist Detalhado por Categoria

---

## Core

- [x] Interface `Component` com método `Render() string`
- [x] `AutoID()` — geração automática de IDs únicos por componente
- [x] `Register(id, comp)` — registro global de componentes para introspecção
- [x] `RegisterAction(id, handler)` — vínculo de handlers dinâmicos por ID
- [x] `GlobalActionHandler` — seletor central estilo Controller para todas as ações
- [x] Endpoint `/api/goui/action` — despachante automático de ações
- [x] `HTML` — wrapper para strings HTML brutas como `Component`
- [x] `SyncWith` — helper que implementa `SyncSource` a partir de um ID raw
- [ ] Estado persistente de componente entre requests (ex.: sessão ou cookie)
- [ ] Middleware de pipeline (ex.: auth, logging, recovery antes do render)
- [ ] Suporte a múltiplas instâncias de `App` no mesmo processo

---

## App & Roteamento

- [x] `NewApp()` — criação da aplicação
- [x] `RegisterRoute(path, component)` — registro de rotas manualmente
- [x] `RegisterHandler(path, handler)` — registro de endpoints HTTP customizados
- [x] `App.Start(ip, port)` — inicialização do servidor HTTP com log de URL
- [x] `router.RegisterPage(path, title, handler)` — registro via `init()` para file-based routing
- [x] `router.InjectRoutes(app, layoutFn)` — injeção de todas as rotas com layout compartilhado
- [ ] Roteamento baseado em arquivos gerado automaticamente por CLI
- [ ] Parâmetros de rota dinâmicos (ex.: `/user/:id`)
- [ ] Hot reload automático em modo de desenvolvimento
- [ ] Servidor TLS nativo (`App.StartTLS`)

---

## Layout de Página

- [x] `NewPage(head, nav, body)` — documento HTML completo com head, sidebar e main
- [x] `LayoutNarrow` — conteúdo centralizado em ~730px (padrão)
- [x] `LayoutFull` — conteúdo ocupa toda a largura disponível
- [x] `LayoutCentered` — conteúdo centralizado em ~900px
- [x] `Div(args...)` — container genérico com suporte a classes, atributos e filhos
- [x] `Section(title, args...)` — seção com título e conteúdo agrupado
- [x] `Card(args...)` — container com estilo de card
- [x] `FormField(label, input, helpText)` — campo de formulário com label e texto de ajuda
- [ ] `Grid(cols int, args...)` — layout em grid com número de colunas configurável
- [ ] `Columns(args...)` — layout horizontal proporcional tipo flexbox declarativo
- [ ] `Modal(trigger, content)` — janela modal com backdrop e foco gerenciado
- [ ] `Drawer(side, trigger, content)` — painel lateral deslizante (left/right)
- [ ] `Accordion(items...)` — seções expansíveis com animação

---

## Navegação

- [x] `Navbar(logo, items...)` — barra lateral de navegação com logo
- [x] `NavGroup(label, items...)` — grupo colapsável de links na sidebar
- [x] NavGroups aninhados — `NavItem` interface permite sub-grupos dentro de grupos
- [x] `Headbar(title)` — barra superior com título da página
- [x] `Tabs(items...)` — abas de navegação com conteúdo trocável
- [x] Hash-based tab activation — ativa aba via `window.location.hash` no DOMContentLoaded
- [x] `hashchange` listener — ativa aba ao navegar com hash sem reload de página
- [ ] `Breadcrumb(items...)` — trilha de navegação hierárquica
- [ ] `Pagination(total, page, pageSize)` — paginação com navegação entre páginas
- [ ] `Stepper(steps...)` — indicador de progresso em etapas (wizard)
- [ ] Navbar responsiva com colapso automático em telas pequenas

---

## Componentes de Texto e Exibição

- [x] `Text(content, opts...)` — texto com suporte a reatividade via `Watch`
- [x] `Badge(content, opts...)` — etiqueta com variantes de cor (SuccessBadge, ErrorBadge, WarningBadge, InfoBadge, DefaultBadge)
- [x] `Icon(name, opts...)` — ícone Lucide renderizado como SVG inline
- [x] `Snippet(title, code, opts...)` — bloco de código com highlight
- [x] `Table(headers, rows, opts...)` — tabela com dados estáticos
- [ ] `Alert(message, type)` — banner de alerta inline (info, warning, error, success)
- [ ] `Tooltip(trigger, text)` — texto flutuante ao passar o mouse
- [ ] `Avatar(src, fallback)` — imagem de perfil com fallback em iniciais
- [x] `ProgressBar(opts...)` — barra de progresso com `SetTotal(n)` / `Add(n...)`, 5 variantes de cor (Default, Success, Error, Warning, Info), label `current/total/%` configurável, validação de passo positivo com panic descritivo
- [ ] `Skeleton(opts...)` — placeholder de carregamento (loading state)
- [ ] `Timeline(events...)` — exibição cronológica de eventos
- [ ] `Stat(label, value, delta)` — card de métrica com variação (KPI)

---

## Formulários e Inputs

- [x] `Input(opts...)` — campo de texto, número, senha, email, etc.
- [x] `Dropdown(opts...)` — select com `Multi(true)` para multi-seleção
- [x] `Toggle(opts...)` — interruptor booleano estilizado
- [x] `Slider(opts...)` — controle deslizante numérico com min/max/step
- [x] `FileUploader(label, opts...)` — upload de arquivos com drag & drop
- [x] `Button(label, opts...)` — botão com variantes Primary/Secondary/Danger/Ghost
- [x] `DownloadButton(label, url, opts...)` — botão que dispara download de arquivo
- [x] `Checkbox(label, opts...)` — caixa de seleção individual
- [x] `CheckboxGroup(label, items...)` — grupo de checkboxes com múltipla seleção
- [x] `RadioGroup(label, items...)` — grupo de opções exclusivas (radio buttons)
- [x] `Textarea(opts...)` — campo de texto multilinha com auto-resize
- [x] `TagInput(opts...)` — input que aceita múltiplas tags digitadas
- [x] `ColorPicker(opts...)` — seletor de cor com preview de swatch
- [x] `Form(action, method, opts...)` — wrapper de formulário com proteção CSRF
- [x] Validação client-side declarativa — required, minLen, maxLen, minNum, maxNum, pattern, mensagens customizadas
- [x] Validação server-side — `ValidateForm(r, rules)` com retorno de erros por campo
- [x] `FieldError(msg)` — renderiza mensagem de erro inline ao lado do campo

---

## Segurança

- [x] `SetCSRFSecret(secret)` — configura a chave HMAC global
- [x] `NewCSRFToken(w)` — gera token assinado e armazena em cookie HttpOnly
- [x] `ValidateCSRF(r)` — valida token do formulário contra o cookie

---

## Calendário e Datas

- [x] `Calendar(opts...)` — seletor de data único
- [x] `CalendarRange(opts...)` / `DateRange(opts...)` — seletor de intervalo de datas
- [x] `ParseDate(r, name)` — parse de `DateValue` a partir de um request HTTP

---

## Comandos e Busca

- [x] `CommandPalette(cmds...)` — paleta de comandos com busca fuzzy (estilo `⌘K`)

---

## Toast / Notificações

- [x] `ToastContainer()` — elemento HTML necessário para exibir toasts
- [x] `ShowToast(message, type)` — exibe uma notificação temporária
- [x] Tipos: `ToastSuccess`, `ToastError`, `ToastWarning`, `ToastInfo`
- [x] `gouiToast(msg, type)` — helper JS global para disparar toasts a partir de onclick

---

## Reatividade Client-Side (Sem JS manual)

- [x] `Watch(source, action)` — vincula a saída de um componente a outro reativo
- [x] `Bind(template, sources)` — interpolação de template a partir de múltiplas fontes
- [x] `SyncText(srcID, targetID, fallback)` — reflete valor de input em texto em tempo real
- [x] `SyncRange(startID, endID, labelID, sep)` — exibe intervalo de dois inputs em um label
- [x] `SyncCSV(uploaderID, tableID)` — popula tabela a partir de CSV carregado no cliente
- [x] `SyncServer(uploaderID, tableID, apiPath)` — envia arquivo ao backend e preenche tabela
- [x] `SyncSubmit(sourceIDs, apiPath, targetID)` — coleta inputs e envia JSON a uma API
- [ ] `SyncVisible(srcID, targetID, condition)` — mostra/esconde componente com base em valor
- [ ] `SyncDisabled(srcID, targetID, condition)` — desabilita componente com base em valor
- [ ] `WatchURL` — atualiza parâmetros da URL ao mudar estado (deep link)
- [ ] Server-Sent Events (SSE) — canal robusto com `SSEStream`, endpoint `/api/goui/stream`, push automático em `Add()` (em andamento — v0.3.1)
- [ ] WebSocket nativo integrado ao sistema de componentes

### Ações do `Watch`

| Ação               | Efeito                                        |
|--------------------|-----------------------------------------------|
| `WatchText`        | Atualiza o texto visível                      |
| `WatchTag`         | Altera a tag HTML                             |
| `WatchColor`       | Altera a cor                                  |
| `WatchSize`        | Altera o tamanho                              |
| `WatchWeight`      | Altera o peso da fonte                        |
| `WatchClass`       | Adiciona/remove classes CSS                   |
| `WatchPlaceholder` | Atualiza o placeholder de um input            |
| `WatchValue`       | Atualiza o valor de um input                  |
| `WatchType`        | Altera o tipo do input (ex.: text → password) |

---

## Tema e Estilo

- [x] `ThemeScript` — script inline de alternância dark/light sem flash
- [x] `Theme` — CSS global da biblioteca (variáveis, reset, classes utilitárias)
- [x] Suporte nativo a dark mode via variáveis CSS e `data-theme`
- [x] Props utilitários: `Class`, `Style`, `ID`, `Attr` disponíveis em qualquer componente
- [x] CSS classes utilitárias — `goui-flex`, `goui-gap-*`, `goui-mb-*`, `goui-text-*`, etc.
- [ ] Múltiplos temas pré-definidos (ex.: `ThemeOcean`, `ThemeSlate`, `ThemeForest`)
- [ ] `CustomTheme(vars map[string]string)` — tema gerado dinamicamente via CSS variables
- [ ] Suporte a temas por página (não só global)
- [ ] Animações e transições declarativas via props

---

## Utilitários

- [x] `ParseStringAttr(v, class, id, attrs)` — parsing de opções variádicas
- [x] `props.go` — sistema de props compartilhado entre componentes (`Name`, `Value`, `Type`, `Placeholder`, `Min`, `Max`, `Step`, `Rows`, `Multi`, `Primary`, `Disabled`, etc.)
- [ ] `If(cond bool, comp Component)` — renderização condicional declarativa
- [ ] `Each(items, fn)` — helper de loop para renderizar listas de componentes
- [ ] `Lazy(loader func() Component)` — carregamento adiado de componente pesado

---

## Acessibilidade

- [ ] Atributos ARIA injetados automaticamente por componente (role, aria-label, aria-expanded)
- [ ] Navegação por teclado em `Dropdown`, `CommandPalette`, `Tabs` e `Modal`
- [ ] Foco gerenciado ao abrir/fechar `Modal` e `Drawer`
- [ ] Suporte a `prefers-reduced-motion` nas animações

---

## Mídia e Captura

- [ ] `CameraCapture(opts...)` — acesso à câmera do usuário com preview em tempo real
- [ ] `PhotoCapture(opts...)` — captura de foto via câmera com preview e botão de confirmar
- [ ] `VideoRecorder(opts...)` — gravação de vídeo pela câmera com controles
- [ ] `AudioRecorder(opts...)` — gravação de áudio pelo microfone com waveform visual
- [ ] `ScreenCapture(opts...)` — captura de tela ou aba do navegador (getDisplayMedia)
- [ ] Upload automático de mídia gravada para endpoint no backend
- [ ] Seleção de dispositivo de entrada (câmera/microfone quando há múltiplos)
- [ ] Suporte a formatos de saída configurável (webm, mp4, ogg, wav)

---

## DX / Ferramentas

- [x] `examples/helloworld` — documentação interativa com playground e 20+ páginas de componentes
- [ ] CLI `goui new <projeto>` — scaffold de projeto com estrutura padrão
- [ ] CLI `goui dev` — servidor de desenvolvimento com hot reload
- [ ] Suporte a testes de componente via `RenderToString` em testes Go nativos

---

## Exemplos

- [x] `examples/helloworld` — aplicação completa com todos os componentes documentados
- [ ] Exemplo de CRUD completo com tabela, modal e formulário
- [ ] Exemplo de dashboard com gráficos e métricas em tempo real (SSE)
- [ ] Exemplo de autenticação com sessão e rotas protegidas
