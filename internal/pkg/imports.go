package pkg

import (
	"go/ast"
)

func Imports(file *ast.File) []ast.Spec {
	imports := make([]ast.Spec, 0)

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ImportSpec:
			if node != nil {
				imports = append(imports, ast.Spec(node))
			}
		}

		return true
	})

	return imports
}
