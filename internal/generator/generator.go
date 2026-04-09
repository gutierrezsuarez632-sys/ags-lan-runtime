package generator

import (
	"fmt"

	"github.com/gutierrezsuarez632-sys/ags-lan-runtime/internal/runtime"
)

type Generator struct {
	basePath string
}

func NewGenerator(basePath string) *Generator {
	return &Generator{
		basePath: basePath,
	}
}

func (g *Generator) Generate(project *runtime.Project) error {

	if project == nil {
		return fmt.Errorf("project is nil")
	}

	root := join(g.basePath, project.Name)

	if err := createDir(root); err != nil {
		return err
	}

	for _, bc := range project.Contexts {
		if err := g.generateBoundedContext(root, bc); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateBoundedContext(root string, bc *runtime.BoundedContext) error {

	bcPath := join(root, bc.Name)

	if err := createDir(bcPath); err != nil {
		return err
	}

	// subcontextos
	for _, sub := range bc.Subcontexts {
		subPath := join(bcPath, sub)

		if err := createDir(subPath); err != nil {
			return err
		}
	}

	// hexagonal
	if bc.Hexagonal {
		if err := g.createHexagonalStructure(bcPath); err != nil {
			return err
		}
	}

	return nil
}

// 🔥 y este también debe existir
func (g *Generator) createHexagonalStructure(base string) error {

	dirs := []string{
		"application",
		"domain",
		"infrastructure",
	}

	for _, d := range dirs {
		if err := createDir(join(base, d)); err != nil {
			return err
		}
	}

	return nil
}
