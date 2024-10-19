package iface

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/samborkent/fake/internal/cases"
	"github.com/samborkent/fake/internal/gen"
)

const (
	appendFunc  = "append"
	contextType = "context.Context"
	lenFunc     = "len"
	nilVar      = "nil"
	receiverVar = "f"
)

func (i Interface) Compliance() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: "_",
					},
				},
				Type: &ast.Ident{
					Name: i.Name,
				},
				Values: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: &ast.Ident{
								Name: i.FakeName(),
							},
						},
					},
				},
			},
		},
	}
}

func (i Interface) Constructor() *ast.FuncDecl {
	constructorVars := make([]ast.Expr, len(i.Methods))

	for index, method := range i.Methods {
		constructorVars[index] = method.ConstructorVar(i.Name)
	}

	return &ast.FuncDecl{
		Name: &ast.Ident{Name: "NewFake" + i.Name},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{gen.TestingT()},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: i.FakeName()},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.Ident{Name: i.FakeName()},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   &ast.Ident{Name: "t"},
										Value: &ast.Ident{Name: "t"},
									},
									&ast.KeyValueExpr{
										Key: &ast.Ident{Name: "On"},
										Value: &ast.UnaryExpr{
											Op: token.AND,
											X: &ast.CompositeLit{
												Type: &ast.Ident{Name: i.ExpecterName()},
												Elts: constructorVars,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m Method) ConstructorVar(intfaceName string) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key: &ast.Ident{Name: cases.Camel(m.Name)},
		Value: &ast.CallExpr{
			Fun: &ast.Ident{Name: "make"},
			Args: []ast.Expr{
				&ast.ArrayType{
					Elt: &ast.StarExpr{
						X: &ast.Ident{Name: cases.Camel(intfaceName) + m.Name},
					},
				},
				&ast.BasicLit{
					Kind:  token.INT,
					Value: "0",
				},
			},
		},
	}
}

func (i Interface) FakeStruct() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: i.FakeName()},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							gen.TestingT(),
							{
								Names: []*ast.Ident{{Name: "On"}},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: i.ExpecterName()},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (i Interface) ExpecterStruct() *ast.GenDecl {
	expecterFields := make([]*ast.Field, 0, 2*len(i.Methods))

	for _, method := range i.Methods {
		expecterFields = append(expecterFields, method.ExpecterField(i.Name)...)
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: i.ExpecterName()},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: expecterFields,
					},
				},
			},
		},
	}
}

func (m Method) ExpecterField(intfaceName string) []*ast.Field {
	lcMethod := cases.Camel(m.Name)

	return []*ast.Field{
		{
			Names: []*ast.Ident{{Name: lcMethod}},
			Type: &ast.ArrayType{
				Elt: &ast.StarExpr{
					X: &ast.Ident{Name: m.ExpecterName(intfaceName)},
				},
			},
		},
		{
			Names: []*ast.Ident{{Name: m.CounterName()}},
			Type:  &ast.Ident{Name: strings.ToLower(token.INT.String())},
		},
	}
}

func (i Interface) NameConst() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.CONST,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{{Name: cases.Camel(i.Name) + "Name"}},
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: `"` + i.Name + `"`,
					},
				},
			},
		},
	}
}

func (m Method) Expecter(interfaceName string) []ast.Decl {
	params := make([]*ast.Field, 0, len(m.Parameters))
	expecterValues := make([]ast.Expr, 0, 2*len(m.Parameters))

	for _, param := range m.Parameters {
		if param.Type == contextType {
			continue
		}

		params = append(params, &ast.Field{
			Names: []*ast.Ident{{Name: cases.Camel(param.Name)}},
			Type:  &ast.Ident{Name: param.Type},
		})

		expecterValues = append(expecterValues, &ast.KeyValueExpr{
			Key:   &ast.Ident{Name: param.Name},
			Value: &ast.Ident{Name: param.Name},
		})
	}

	expecterParams := params

	if len(m.Results) > 0 {
		expecterParams = append(expecterParams, &ast.Field{
			Names: []*ast.Ident{{Name: "returns"}},
			Type: &ast.StarExpr{
				X: &ast.Ident{
					Name: m.ReturnName(interfaceName),
				},
			},
		})
	}

	expecterSelector := &ast.SelectorExpr{
		X:   &ast.Ident{Name: receiverVar},
		Sel: &ast.Ident{Name: cases.Camel(m.Name)},
	}

	return []ast.Decl{
		&ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: receiverVar}},
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: cases.Camel(interfaceName) + expecterSuffix},
						},
					},
				},
			},
			Name: &ast.Ident{Name: m.Name},
			Type: &ast.FuncType{
				Params: &ast.FieldList{List: params},
				Results: &ast.FieldList{List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: m.ExpecterName(interfaceName)},
						},
					},
				}},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Cond: &ast.BinaryExpr{
							Op: token.EQL,
							X:  &ast.Ident{Name: receiverVar},
							Y:  &ast.Ident{Name: nilVar},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.Ident{Name: nilVar},
									},
								},
							},
						},
					},
					&ast.AssignStmt{
						Tok: token.ASSIGN,
						Lhs: []ast.Expr{expecterSelector},
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{Name: appendFunc},
								Args: []ast.Expr{
									expecterSelector,
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: m.ExpecterName(interfaceName)},
											Elts: expecterValues,
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.IndexExpr{
								X: expecterSelector,
								Index: &ast.BinaryExpr{
									Op: token.SUB,
									X: &ast.CallExpr{
										Fun:  &ast.Ident{Name: lenFunc},
										Args: []ast.Expr{expecterSelector},
									},
									Y: &ast.BasicLit{
										Kind:  token.INT,
										Value: "1",
									},
								},
							},
						},
					},
				},
			},
		},
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: m.ExpecterName(interfaceName)},
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: expecterParams,
						},
					},
				},
			},
		},
	}
}
