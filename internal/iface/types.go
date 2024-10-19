package iface

import (
	"go/ast"
	"strings"

	"github.com/samborkent/fake/internal/cases"
)

const (
	counterSuffix  = "Counter"
	expecterSuffix = "Expect"
	fakePrefix     = "fake"
	returnSuffix   = "Return"
)

type Interface struct {
	Name    string
	Methods []Method
	AST     *ast.InterfaceType
}

type Method struct {
	Name       string
	Parameters []Variable
	Results    []Variable
	AST        *ast.FuncType
}

func (m Method) CounterName() string {
	return cases.Camel(m.Name) + counterSuffix
}

func (m Method) ExpecterName(interfaceName string) string {
	return cases.Camel(interfaceName) + m.Name
}

func (m Method) ReturnName(interfaceName string) string {
	return m.ExpecterName(interfaceName) + returnSuffix
}

type Variable struct {
	Name string
	Type string
	AST  ast.Expr
}

func (i Interface) String() string {
	builder := new(strings.Builder)

	_, _ = builder.WriteString("type ")
	_, _ = builder.WriteString(i.Name)
	_, _ = builder.WriteString(" interface {\n")

	for _, method := range i.Methods {
		_, _ = builder.WriteString("\t")
		_, _ = builder.WriteString(method.Name)
		_, _ = builder.WriteString("(")

		for i, parameter := range method.Parameters {
			if i > 0 {
				_, _ = builder.WriteString(", ")
			}

			_, _ = builder.WriteString(parameter.Name)
			_, _ = builder.WriteString(" ")
			_, _ = builder.WriteString(parameter.Type)
		}

		_, _ = builder.WriteString(") ")

		if len(method.Results) > 0 {
			_, _ = builder.WriteString("(")

			for i, result := range method.Results {
				if i > 0 {
					_, _ = builder.WriteString(", ")
				}

				if result.Name == "" {
					_, _ = builder.WriteString(result.Type)
					continue
				}

				_, _ = builder.WriteString(result.Name)
				_, _ = builder.WriteString(" ")
				_, _ = builder.WriteString(result.Type)
			}

			_, _ = builder.WriteString(")")
		}

		_, _ = builder.WriteString("\n")
	}

	_, _ = builder.WriteString("}")

	return builder.String()
}

func (i Interface) ExpecterName() string {
	return cases.Camel(i.Name) + expecterSuffix
}

func (i Interface) FakeName() string {
	return fakePrefix + i.Name
}
