package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/samborkent/fake/internal/cases"
	"github.com/samborkent/fake/internal/iface"
	"github.com/samborkent/fake/internal/pkg"
)

const (
	appendFunc  = "append"
	contextType = "context.Context"
	lenFunc     = "len"
	nilVar      = "nil"
	receiverVar = "f"
)

func main() {
	fileSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileSet, "interface.go", nil, parser.Mode(0))
	if err != nil {
		log.Fatal("parsing file: " + err.Error())
	}

	packageName := pkg.Name(astFile)
	imports := pkg.Imports(astFile)
	interfaces := iface.GetInterfaces(astFile)

	for _, intface := range interfaces {
		fileName := cases.Snake(intface.Name) + "_fake.go"

		goFile, err := os.Create(fileName)
		if err != nil {
			log.Printf("error: creating file '%s': %s\n", fileName, err.Error())
			continue
		}

		importDecl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: imports,
		}

		declarations := []ast.Decl{
			importDecl,
			intface.Compliance(),
			intface.Constructor(),
			intface.FakeStruct(),
			intface.ExpecterStruct(),
			intface.NameConst(),
		}

		for _, method := range intface.Methods {
			// TODO: func

			declarations = append(declarations, method.Expecter(intface.Name)...)

			if len(method.Results) > 0 {
				results := make([]*ast.Field, 0, len(method.Results))
				returnerVars := make([]ast.Expr, 0, len(method.Results))

				for _, result := range method.Results {
					lcResult := cases.Camel(result.Name)

					results = append(results, &ast.Field{
						Names: []*ast.Ident{{Name: lcResult}},
						Type:  &ast.Ident{Name: result.Type},
					})

					returnerVars = append(returnerVars, &ast.KeyValueExpr{
						Key:   &ast.Ident{Name: lcResult},
						Value: &ast.Ident{Name: lcResult},
					})
				}

				declarations = append(declarations, &ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{Name: method.ReturnName(intface.Name)},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: results,
								},
							},
						},
					},
				})

				declarations = append(declarations, &ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{{Name: receiverVar}},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: method.ExpecterName(intface.Name)},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Return"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: results,
						},
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
										&ast.ReturnStmt{},
									},
								},
							},
							&ast.AssignStmt{
								Tok: token.ASSIGN,
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   &ast.Ident{Name: receiverVar},
										Sel: &ast.Ident{Name: "returns"},
									},
								},
								Rhs: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: method.ReturnName(intface.Name)},
											Elts: returnerVars,
										},
									},
								},
							},
						},
					},
				})
			}
		}

		genFile := &ast.File{
			Package: token.NoPos,
			Name: &ast.Ident{
				Name: packageName,
			},
			Decls: declarations,
		}

		if err := format.Node(goFile, fileSet, genFile); err != nil {
			log.Fatal("printing ast: " + err.Error())
		}
	}
}
