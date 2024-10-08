package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fileSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileSet, "interface.go", nil, parser.Mode(0))
	if err != nil {
		log.Fatal("parsing file: " + err.Error())
	}

	ast.Inspect(astFile, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.TypeSpec:
			if node == nil || node.Name == nil || !node.Name.IsExported() {
				break
			}

			interfaceType, ok := node.Type.(*ast.InterfaceType)
			if !ok {
				break
			}

			if interfaceType == nil || interfaceType.Methods == nil {
				break
			}

			for _, method := range interfaceType.Methods.List {
				if method == nil {
					continue
				}

				var methodNames string

				for _, methodName := range method.Names {
					if methodName == nil || !methodName.IsExported() {
						continue
					}

					methodNames += methodName.String()
				}

				functionType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				var params string
				var results string

				if functionType.Params != nil {
					for _, param := range functionType.Params.List {
						for _, paramName := range param.Names {
							if paramName == nil {
								continue
							}

							if params == "" {
								params += paramName.String()
							} else {
								params += ", " + paramName.String()
							}
						}

						paramID, ok := param.Type.(*ast.Ident)
						if ok {
							fmt.Printf("%s\n", paramID.String())
						}
					}
				}

				if functionType.Results != nil {
					for _, result := range functionType.Results.List {
						for _, resultName := range result.Names {
							if resultName == nil {
								continue
							}

							if results == "" {
								results += resultName.String()
							} else {
								results += ", " + resultName.String()
							}
						}
					}
				}

				fmt.Printf("'%s': %s.%s(%s) (%s)\n", fileSet.Position(node.Pos()), node.Name, methodNames, params, results)
			}
		}

		return true
	})
}
