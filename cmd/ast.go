package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"

	gofumptformat "mvdan.cc/gofumpt/format"

	"github.com/samborkent/fake/internal/cases"
	"github.com/samborkent/fake/internal/iface"
	"github.com/samborkent/fake/internal/pkg"
)

const (
	appendFunc      = "append"
	contextType     = "context.Context"
	fakePackageName = "fake"
	indexVar        = "index"
	lenFunc         = "len"
	nilVar          = "nil"
	panicFunc       = "panic"
	receiverVar     = "f"
)

var (
	fOn = &ast.SelectorExpr{
		X:   &ast.Ident{Name: receiverVar},
		Sel: &ast.Ident{Name: "On"},
	}

	fT = &ast.SelectorExpr{
		X:   &ast.Ident{Name: receiverVar},
		Sel: &ast.Ident{Name: "t"},
	}

	fTFatalf = &ast.SelectorExpr{
		X:   fT,
		Sel: &ast.Ident{Name: "Fatalf"},
	}
)

func main() {
	fileSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileSet, "gen/interface.go", nil, parser.Mode(0))
	if err != nil {
		log.Fatal("parsing file: " + err.Error())
	}

	packageName := pkg.Name(astFile)
	imports := pkg.Imports(astFile)
	interfaces := iface.GetInterfaces(astFile)

	for _, intface := range interfaces {
		fileName := cases.Snake(intface.Name) + "_fake.go"

		goFile, err := os.Create("gen/" + fileName)
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
			// Fake implementation
			params := make([]*ast.Field, 0, len(method.Parameters))
			paramChecks := make([]ast.Stmt, 0, len(method.Parameters))

			for _, param := range method.Parameters {
				params = append(params, &ast.Field{
					Names: []*ast.Ident{{Name: param.Name}},
					Type:  param.AST,
				})

				if param.Type == contextType {
					// Context checks
					paramChecks = append(paramChecks, &ast.IfStmt{
						Cond: &ast.BinaryExpr{
							Op: token.EQL,
							X:  &ast.Ident{Name: "ctx"},
							Y:  &ast.Ident{Name: "nil"},
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: fTFatalf,
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"fake: '%s.%s': %s"`,
										},
										&ast.Ident{Name: cases.Camel(intface.Name) + "Name"},
										&ast.Ident{Name: "methodName"},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "fake"},
													Sel: &ast.Ident{Name: "ErrContextNil"},
												},
												Sel: &ast.Ident{Name: "Error"},
											},
										},
									},
								},
							},
						}},
					}, &ast.IfStmt{
						Init: &ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{&ast.Ident{Name: "err"}},
							Rhs: []ast.Expr{&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "ctx"},
									Sel: &ast.Ident{Name: "Err"},
								},
							}},
						},
						Cond: &ast.BinaryExpr{
							Op: token.NEQ,
							X:  &ast.Ident{Name: "err"},
							Y:  &ast.Ident{Name: nilVar},
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: fTFatalf,
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"fake: '%s.%s': %s"`,
										},
										&ast.Ident{Name: cases.Camel(intface.Name) + "Name"},
										&ast.Ident{Name: "methodName"},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "err"},
												Sel: &ast.Ident{Name: "Error"},
											},
										},
									},
								},
							},
						}},
					})
				} else {
					// Parameter checks
					paramChecks = append(paramChecks, &ast.IfStmt{
						Cond: &ast.UnaryExpr{
							Op: token.NOT,
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "reflect"},
									Sel: &ast.Ident{Name: "DeepEqual"},
								},
								Args: []ast.Expr{
									&ast.Ident{Name: param.Name},
									&ast.SelectorExpr{
										X: &ast.IndexExpr{
											X: &ast.SelectorExpr{
												X:   fOn,
												Sel: &ast.Ident{Name: cases.Camel(method.Name)},
											},
											Index: &ast.Ident{Name: "index"},
										},
										Sel: &ast.Ident{Name: param.Name},
									},
								},
							},
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: fTFatalf,
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"fake: '%s.%s': %s: '%s': got '%+v', want '%+v'"`,
										},
										&ast.Ident{Name: cases.Camel(intface.Name) + "Name"},
										&ast.Ident{Name: "methodName"},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "fake"},
													Sel: &ast.Ident{Name: "ErrArgumentMismatch"},
												},
												Sel: &ast.Ident{Name: "Error"},
											},
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"` + param.Name + `"`,
										},
										&ast.Ident{Name: param.Name},
										&ast.SelectorExpr{
											X: &ast.IndexExpr{
												X: &ast.SelectorExpr{
													X:   fOn,
													Sel: &ast.Ident{Name: cases.Camel(method.Name)},
												},
												Index: &ast.Ident{Name: "index"},
											},
											Sel: &ast.Ident{Name: param.Name},
										},
									},
								},
							},
						}},
					})
				}
			}

			results := make([]*ast.Field, 0, len(method.Results))
			returnVars := make([]ast.Expr, 0, len(method.Results))

			for _, result := range method.Results {
				results = append(results, &ast.Field{
					Names: []*ast.Ident{{Name: result.Name}},
					Type:  result.AST,
				})

				returnVars = append(returnVars, &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X: &ast.IndexExpr{
							X: &ast.SelectorExpr{
								X:   fOn,
								Sel: &ast.Ident{Name: cases.Camel(method.Name)},
							},
							Index: &ast.Ident{Name: "index"},
						},
						Sel: &ast.Ident{Name: "returns"},
					},
					Sel: &ast.Ident{Name: cases.Camel(result.Name)},
				})
			}

			body := []ast.Stmt{
				// Receiver check
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.LOR,
						X: &ast.BinaryExpr{
							Op: token.EQL,
							X:  &ast.Ident{Name: receiverVar},
							Y:  &ast.Ident{Name: nilVar},
						},
						Y: &ast.BinaryExpr{
							Op: token.EQL,
							X:  fT,
							Y:  &ast.Ident{Name: nilVar},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.Ident{Name: panicFunc},
									Args: []ast.Expr{
										&ast.BinaryExpr{
											Op: token.ADD,
											X: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `"fake: "`,
											},
											Y: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   &ast.Ident{Name: fakePackageName},
														Sel: &ast.Ident{Name: "ErrFakeNotInitialized"},
													},
													Sel: &ast.Ident{Name: "Error"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				// Method name const
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.CONST,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{{Name: "methodName"}},
								Values: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"` + method.Name + `"`,
									},
								},
							},
						},
					},
				},
				// Expecter check
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.LOR,
						X: &ast.BinaryExpr{
							Op: token.EQL,
							X:  fOn,
							Y:  &ast.Ident{Name: nilVar},
						},
						Y: &ast.BinaryExpr{
							Op: token.EQL,
							X: &ast.SelectorExpr{
								X:   fOn,
								Sel: &ast.Ident{Name: cases.Camel(method.Name)},
							},
							Y: &ast.Ident{Name: nilVar},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: fTFatalf,
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"fake: '%s.%s': %s"`,
										},
										&ast.Ident{Name: cases.Camel(intface.Name) + "Name"},
										&ast.Ident{Name: "methodName"},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X:   &ast.Ident{Name: fakePackageName},
													Sel: &ast.Ident{Name: "ErrMethodNotInitialized"},
												},
												Sel: &ast.Ident{Name: "Error"},
											},
										},
									},
								},
							},
						},
					},
				},
				// Set index
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						&ast.Ident{Name: indexVar},
					},
					Rhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   fOn,
							Sel: &ast.Ident{Name: cases.Camel(method.Name) + "Counter"},
						},
					},
				},
				// Counter check
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.GTR,
						X: &ast.BinaryExpr{
							Op: token.ADD,
							X:  &ast.Ident{Name: indexVar},
							Y: &ast.BasicLit{
								Kind:  token.INT,
								Value: "1",
							},
						},
						Y: &ast.CallExpr{
							Fun: &ast.Ident{Name: lenFunc},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   fOn,
									Sel: &ast.Ident{Name: cases.Camel(method.Name)},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: fTFatalf,
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"fake: '%s.%s': %s: called '%d' time(s), '%d' expectation(s) registered"`,
										},
										&ast.Ident{Name: cases.Camel(intface.Name) + "Name"},
										&ast.Ident{Name: "methodName"},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "fake"},
													Sel: &ast.Ident{Name: "ErrExpectationsMissing"},
												},
												Sel: &ast.Ident{Name: "Error"},
											},
										},
										&ast.BinaryExpr{
											Op: token.ADD,
											X:  &ast.Ident{Name: indexVar},
											Y: &ast.BasicLit{
												Kind:  token.INT,
												Value: "1",
											},
										},
										&ast.CallExpr{
											Fun: &ast.Ident{Name: lenFunc},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   fOn,
													Sel: &ast.Ident{Name: cases.Camel(method.Name)},
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

			// Parameter checks
			body = append(body, paramChecks...)

			// Return check
			if len(method.Results) > 0 {
				body = append(body, &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						Op: token.EQL,
						X: &ast.SelectorExpr{
							X: &ast.IndexExpr{
								X: &ast.SelectorExpr{
									X:   fOn,
									Sel: &ast.Ident{Name: cases.Camel(method.Name)},
								},
								Index: &ast.Ident{Name: "index"},
							},
							Sel: &ast.Ident{Name: "returns"},
						},
						Y: &ast.Ident{Name: nilVar},
					},
					Body: &ast.BlockStmt{List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: fTFatalf,
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"fake: '%s.%s': %s: '%d'"`,
									},
									&ast.Ident{Name: cases.Camel(intface.Name) + "Name"},
									&ast.Ident{Name: "methodName"},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   &ast.Ident{Name: fakePackageName},
												Sel: &ast.Ident{Name: "ErrReturnMissing"},
											},
											Sel: &ast.Ident{Name: "Error"},
										},
									},
									&ast.BinaryExpr{
										Op: token.ADD,
										X:  &ast.Ident{Name: "index"},
										Y: &ast.BasicLit{
											Kind:  token.INT,
											Value: "1",
										},
									},
								},
							},
						},
					}},
				})
			}

			// Increment expecter counter
			body = append(body, &ast.IncDecStmt{
				Tok: token.INC,
				X: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   &ast.Ident{Name: receiverVar},
						Sel: &ast.Ident{Name: "On"},
					},
					Sel: &ast.Ident{Name: cases.Camel(method.Name) + "Counter"},
				},
			})

			// Return
			body = append(body, &ast.ReturnStmt{
				Results: returnVars,
			})

			declarations = append(declarations, &ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: receiverVar}},
							Type: &ast.StarExpr{
								X: &ast.Ident{Name: intface.FakeName()},
							},
						},
					},
				},
				Name: &ast.Ident{Name: method.Name},
				Type: &ast.FuncType{
					Params:  &ast.FieldList{List: params},
					Results: &ast.FieldList{List: results},
				},
				Body: &ast.BlockStmt{List: body},
			})

			// Expecter
			declarations = append(declarations, method.Expecter(intface.Name)...)

			// Result expecter
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

	for _, intface := range interfaces {
		fileName := "gen/" + cases.Snake(intface.Name) + "_fake.go"

		astFile, err := parser.ParseFile(fileSet, fileName, nil, parser.Mode(0))
		if err != nil {
			log.Fatal(fmt.Errorf("formatting file '%s': %w", fileName, err).Error())
		}

		gofumptformat.File(fileSet, astFile, gofumptformat.Options{
			LangVersion: "go1.24.0",
			ModulePath:  "github.com/samborkent/fake",
			ExtraRules:  true,
		})
	}
}
