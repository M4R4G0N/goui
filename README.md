# goUI 🚀

![Version](https://img.shields.io/badge/version-v0.1.1-blue.svg)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)

**goUI** is a lightweight, high-performance Go library for building server-rendered web interfaces using a component-based model and file-based routing.

Inspired by simplicity, **goUI** allows Go developers to build rich, reactive UIs entirely in Go — no JavaScript frameworks, no complex build pipelines, just pure Go.

---

## 📋 Prerequisites

To use **goUI**, you need **Go 1.22 or later** installed on your system.

### Install Go
- **Windows**: Download and run the MSI installer from [golang.org/dl](https://go.dev/dl/).
- **Linux**: Follow the [official installation guide](https://go.dev/doc/install).
- **macOS**: Download the `.pkg` installer from [golang.org/dl](https://go.dev/dl/).

---

## 📦 Installation

```bash
go mod init $(basename "$PWD") 

go get github.com/M4R4G0N/goUI@v0.1.1
```

### 🔧 Troubleshooting: 404 / Checksum Errors
If you see a `404 Not Found` or `checksum verification` error during installation, run these commands to fetch directly from the source:

# try about checksum verification error in the terminal.
```bash
go clean -modcache
```

**Linux / macOS:**
```bash
export GOPROXY=direct
export GONOSUMDB=github.com/M4R4G0N/goUI
go get github.com/M4R4G0N/goUI@v0.1.1
```

**Windows (PowerShell):**
```powershell
$env:GOPROXY="direct"
$env:GONOSUMDB="github.com/M4R4G0N/goUI"
go get github.com/M4R4G0N/goUI@v0.1.1
```

---

## ⚡ Quick Start

### 1. Scaffold a new project
```bash
# Install the CLI
go install github.com/M4R4G0N/goUI/cmd/goUI@latest

# Create your project
goui new myapp
```

### 2. Run the application
```bash
cd myapp
go mod tidy
go run main.go
```
Visit `http://localhost:8080`.

---

## ✨ Core Features

- **🚀 SSR-First**: Ultra-fast server-side rendering.
- **⚡ Reactive Watch/Bind**: Component synchronization without manual JS.
- **🛠️ Premium Components**: Built-in `DataTable`, `CommandPalette` (Ctrl+K), `Snippet` (Syntax Highlighting), `Calendar`, and more.
- **📂 File-based Routing**: Each file in `pages/` automatically becomes a route via `init()`.

---

## 🧱 Example: Reactive UI
```go
func MyPage(title, path string) components.Component {
    input := components.Input("Type here...", components.ID("my-input"))
    text := components.Text("Hello", "h1", components.Watch(input, components.WatchText))

    return components.Div(input, text)
}
```

---

## 🛠️ Project Structure

```
goUI/
├── goui.go              # Main App and Server logic
├── router/              # Automatic file-based routing
├── components/          # The core component library
│   ├── command_palette.go
│   ├── snippet.go (Accordion + Syntax Highlight)
│   ├── table.go (DataTable)
│   ├── input.go (Reactive)
│   └── ...
└── cmd/goUI/            # CLI Tool
```

---

## License
MIT License. Developed by [Marcelo Antonio Goncalves](https://github.com/M4R4G0N).
