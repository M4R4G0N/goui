package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	ColorBlue  = "\033[34m"
	ColorCyan  = "\033[36m"
	ColorGreen = "\033[32m"
	ColorBold  = "\033[1m"
	ColorReset = "\033[0m"
)

func ScaffoldProject(name string) {
	fmt.Printf("\n%s%sCreating goUI project:%s %s%s%s\n", ColorBold, ColorBlue, ColorReset, ColorCyan, name, ColorReset)

	dirs := []string{
		name,
		filepath.Join(name, "pages"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("%sError creating directory %s: %v%s\n", ColorBold, dir, err, ColorReset)
			return
		}
	}

	createFile(filepath.Join(name, "go.mod"), fmt.Sprintf(`module %s

go 1.22

require github.com/M4R4G0N/goUI v0.1.0
`, name))

	createFile(filepath.Join(name, "main.go"), `package main

import (
	_ "`+name+`/pages"

	"github.com/M4R4G0N/goUI"
	"github.com/M4R4G0N/goUI/components"
	"github.com/M4R4G0N/goUI/router"
)

func main() {
	app := goui.NewApp()

	router.InjectRoutes(app, func(title, path string, body components.Component) components.Component {
		return components.NewPage(
			components.Headbar(title),
			components.Navbar("goUI App",
				components.Link{Href: "/", Text: "Home"},
				components.Link{Href: "/sobre", Text: "Sobre"},
			),
			body,
			components.LayoutCentered,
		)
	})

	app.Start("", "8080")
}
`)

	createFile(filepath.Join(name, "pages", "index.go"), `package pages

import (
	"github.com/M4R4G0N/goUI/components"
	"github.com/M4R4G0N/goUI/router"
)

func init() {
	router.RegisterPage("/", "Home", Index)
}

func Index(title, path string) components.Component {
	return components.Div(
		components.Text("Bem-vindo ao goUI!", "h1"),
		components.Text("Sua primeira página, escrita 100% em Go.", "p"),
		components.Button("Ver documentação", components.Primary, "https://github.com/M4R4G0N/goUI"),
	)
}
`)

	createFile(filepath.Join(name, "pages", "sobre.go"), `package pages

import (
	"github.com/M4R4G0N/goUI/components"
	"github.com/M4R4G0N/goUI/router"
)

func init() {
	router.RegisterPage("/sobre", "Sobre", Sobre)
}

func Sobre(title, path string) components.Component {
	return components.Div(
		components.Text(title, "h1"),
		components.Text("Este app foi criado com goUI — interfaces web em Go puro, sem JavaScript.", "p"),
	)
}
`)

	fmt.Printf("\n%s%s✔ Project created successfully!%s\n", ColorBold, ColorGreen, ColorReset)

	fmt.Printf("\n%sNext steps:%s\n", ColorBold, ColorReset)
	fmt.Printf("  %s1.%s cd %s\n", ColorBlue, ColorReset, name)
	fmt.Printf("  %s2.%s go mod tidy\n", ColorBlue, ColorReset)
	fmt.Printf("  %s3.%s go run main.go\n", ColorBlue, ColorReset)

	fmt.Printf("\n%sAcesse http://localhost:8080 no seu navegador%s\n", ColorBold, ColorReset)
}

func createFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("%sError creating file %s: %v%s\n", ColorBold, path, err, ColorReset)
	}
}
