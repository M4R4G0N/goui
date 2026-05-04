package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Error: Project name required. Usage: goui new <project_name>")
			return
		}
		projectName := os.Args[2]
		ScaffoldProject(projectName)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("goui - The Go UI Framework")
	fmt.Println("Usage:")
	fmt.Println("  goui new <project_name>   Create a new goui project structure")
}
