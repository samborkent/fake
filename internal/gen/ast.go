package gen

import "go/ast"

func TestingT() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{
			{
				Name: "t",
			},
		},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "testing",
				},
				Sel: &ast.Ident{
					Name: "T",
				},
			},
		},
	}
}
