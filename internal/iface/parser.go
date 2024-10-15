package iface

import (
	"go/ast"
	"strings"
)

func GetInterfaces(file *ast.File) []Interface {
	interfaces := make([]Interface, 0)

	ast.Inspect(file, func(n ast.Node) bool {
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

			foundInterface := Interface{
				Name:    node.Name.String(),
				Methods: make([]Method, 0, len(interfaceType.Methods.List)),
			}

			for _, method := range interfaceType.Methods.List {
				if method == nil {
					continue
				}

				functionType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				var methodNames string

				for _, methodName := range method.Names {
					if methodName == nil || !methodName.IsExported() {
						continue
					}

					methodNames += methodName.String()
				}

				foundMethod := Method{
					Name: methodNames,
				}

				if functionType.Params != nil {
					foundMethod.Parameters = make([]Variable, 0, len(functionType.Params.List))

					for _, param := range functionType.Params.List {
						paramVars := make([]string, 0, len(param.Names))

						for _, paramName := range param.Names {
							if paramName == nil || paramName.Name == "" {
								continue
							}

							paramVars = append(paramVars, paramName.Name)
						}

						var paramsType string

						switch paramType := param.Type.(type) {
						case *ast.Ident:
							paramsType = paramType.String()
						case *ast.SelectorExpr:
							ident, ok := paramType.X.(*ast.Ident)
							if ok {
								paramsType = ident.String() + "." + paramType.Sel.String()
							} else {
								paramsType = paramType.Sel.String()
							}
						case *ast.StarExpr:
							switch starType := paramType.X.(type) {
							case *ast.Ident:
								paramsType = "*" + starType.String()
							case *ast.SelectorExpr:
								ident, ok := starType.X.(*ast.Ident)
								if ok {
									paramsType = "*" + ident.String() + "." + starType.Sel.String()
								} else {
									paramsType = "*" + starType.Sel.String()
								}
							}
						}

						if len(paramVars) == 0 {
							foundMethod.Parameters = append(foundMethod.Parameters, Variable{
								Name: deduceVarName(paramsType),
								Type: paramsType,
							})
						} else {
							for _, paramName := range paramVars {
								foundMethod.Parameters = append(foundMethod.Parameters, Variable{
									Name: paramName,
									Type: paramsType,
								})
							}
						}
					}
				}

				if functionType.Results != nil {
					foundMethod.Results = make([]Variable, 0, len(functionType.Results.List))

					for _, result := range functionType.Results.List {
						resultVars := make([]string, 0, len(result.Names))

						for _, resultName := range result.Names {
							if resultName == nil || resultName.Name == "" {
								continue
							}

							resultVars = append(resultVars, resultName.Name)
						}

						var resultsType string

						switch resultType := result.Type.(type) {
						case *ast.Ident:
							resultsType = resultType.String()
						case *ast.SelectorExpr:
							ident, ok := resultType.X.(*ast.Ident)
							if ok {
								resultsType = ident.String() + "." + resultType.Sel.String()
							} else {
								resultsType = resultType.Sel.String()
							}
						case *ast.StarExpr:
							switch starType := resultType.X.(type) {
							case *ast.Ident:
								resultsType = "*" + starType.String()
							case *ast.SelectorExpr:
								ident, ok := starType.X.(*ast.Ident)
								if ok {
									resultsType = "*" + ident.String() + "." + starType.Sel.String()
								} else {
									resultsType = "*" + starType.Sel.String()
								}
							}
						}

						if len(resultVars) == 0 {
							foundMethod.Results = append(foundMethod.Results, Variable{
								Name: deduceVarName(resultsType),
								Type: resultsType,
							})
						} else {
							for _, resultName := range resultVars {
								foundMethod.Results = append(foundMethod.Results, Variable{
									Name: resultName,
									Type: resultsType,
								})
							}
						}
					}
				}

				foundInterface.Methods = append(foundInterface.Methods, foundMethod)
			}

			interfaces = append(interfaces, foundInterface)
		}

		return true
	})

	return interfaces
}

func deduceVarName(paramType string) string {
	switch paramType {
	case "context.Context":
		return "ctx"
	case "error":
		return "err"
	case "http.Request", "*http.Request":
		return "req"
	case "http.Response", "*http.Response":
		return "resp"
	case "http.ResponseWriter":
		return "rw"
	case "slog.Logger", "*slog.Logger":
		return "log"
	case "testing.T", "*testing.T":
		return "t"
	default:
		return strings.ReplaceAll(paramType, ".", "")
	}
}
