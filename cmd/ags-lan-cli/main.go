package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gutierrezsuarez632-sys/ags-lan-runtime/internal/generator"
	"github.com/gutierrezsuarez632-sys/ags-lan-runtime/internal/runtime"
)

func main() {

	// 🔹 flags CLI
	filePath := flag.String("file", "", "Path to .ags script file")
	outputPath := flag.String("out", "./output", "Output directory")

	flag.Parse()

	// 🔹 validar input
	if *filePath == "" {
		fmt.Println("❌ missing --file flag")
		fmt.Println("Usage: ags-lan-cli --file project.ags --out ./output")
		os.Exit(1)
	}

	// 🔹 leer archivo
	content, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Println("❌ error reading file:", err)
		os.Exit(1)
	}

	// 🔹 ejecutar runtime (DSL)
	rt := runtime.NewRuntime()

	err = rt.Run(string(content))
	if err != nil {
		fmt.Println("❌ DSL execution error:", err)
		os.Exit(1)
	}

	project := rt.GetProject()

	if project == nil {
		fmt.Println("❌ no project defined in script")
		os.Exit(1)
	}

	// 🔹 generar estructura
	gen := generator.NewGenerator(*outputPath)

	err = gen.Generate(project)
	if err != nil {
		fmt.Println("❌ generation error:", err)
		os.Exit(1)
	}

	// 🔹 éxito
	fmt.Println("✅ Project generated successfully")
	fmt.Println("📁 Output:", *outputPath)

}
