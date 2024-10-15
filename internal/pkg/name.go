package pkg

import (
	"go/ast"
)

func GetPackageName(file *ast.File) string {
	var packageName string

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.File:
			packageName = node.Name.String()
		}

		return true
	})

	return packageName
}
