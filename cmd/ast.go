package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/samborkent/fake/internal/cases"
	"github.com/samborkent/fake/internal/iface"
	"github.com/samborkent/fake/internal/pkg"
)

const contextType = "context.Context"

func main() {
	fileSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileSet, "interface.go", nil, parser.Mode(0))
	if err != nil {
		log.Fatal("parsing file: " + err.Error())
	}

	packageName := pkg.GetPackageName(astFile)
	interfaces := iface.GetInterfaces(astFile)

	for _, iface := range interfaces {
		fmt.Printf("%s\n\n", iface.String())

		fileName := cases.Snake(iface.Name) + "_fake.go"
		fakeName := "fake" + iface.Name
		lcName := cases.Camel(iface.Name)
		expectName := lcName + "Expect"

		goFile, err := os.Create(fileName)
		if err != nil {
			log.Printf("error: creating file '%s': %s\n", fileName, err.Error())
			continue
		}

		data := new(bytes.Buffer)

		// Package
		_, _ = data.WriteString("package ")
		_, _ = data.WriteString(packageName)
		_, _ = data.WriteString("\n\n")

		// Imports
		_, _ = data.WriteString(`import (
			"context"
			"testing"
		)
		`)

		// Interface compliance
		_, _ = data.WriteString("var _ ")
		_, _ = data.WriteString(iface.Name)
		_, _ = data.WriteString(" = &")
		_, _ = data.WriteString(fakeName)
		_, _ = data.WriteString("{}\n\n")

		// Constructor
		_, _ = data.WriteString("func NewFake")
		_, _ = data.WriteString(iface.Name)
		_, _ = data.WriteString("(t *testing.T) *")
		_, _ = data.WriteString(fakeName)
		_, _ = data.WriteString(" {\n\treturn &")
		_, _ = data.WriteString(fakeName)
		_, _ = data.WriteString(`{
				t: t,
				On: &`)
		_, _ = data.WriteString(expectName)
		_, _ = data.WriteString("{\n")

		for _, method := range iface.Methods {
			_, _ = data.WriteString("\t\t\t")
			_, _ = data.WriteString(cases.Camel(method.Name))
			_, _ = data.WriteString(":\tmake([]*")
			_, _ = data.WriteString(lcName)
			_, _ = data.WriteString(method.Name)
			_, _ = data.WriteString(", 0),\n")
		}

		_, _ = data.WriteString(`		},
					}
				}
				
				`)

		// Fake struct
		_, _ = data.WriteString("type fake")
		_, _ = data.WriteString(iface.Name)
		_, _ = data.WriteString(` struct {
					t  *testing.T
					On *`)
		_, _ = data.WriteString(cases.Camel(iface.Name))
		_, _ = data.WriteString(`Expect
				}
		
				`)

		// Expect struct
		_, _ = data.WriteString("type ")
		_, _ = data.WriteString(cases.Camel(iface.Name))
		_, _ = data.WriteString("Expect struct {\n")

		for _, method := range iface.Methods {
			_, _ = data.WriteString("\t")
			_, _ = data.WriteString(cases.Camel(method.Name))
			_, _ = data.WriteString(" []*")
			_, _ = data.WriteString(cases.Camel(iface.Name))
			_, _ = data.WriteString(method.Name)
			_, _ = data.WriteString("\n")

			_, _ = data.WriteString("\t")
			_, _ = data.WriteString(cases.Camel(method.Name))
			_, _ = data.WriteString("Counter int\n\n")
		}

		_, _ = data.WriteString("}\n\n")

		// Interface name constant
		_, _ = data.WriteString("const ")
		_, _ = data.WriteString(cases.Camel(iface.Name))
		_, _ = data.WriteString(`Name = "`)
		_, _ = data.WriteString(iface.Name)
		_, _ = data.WriteString(`"`)
		_, _ = data.WriteString("\n\n")

		// Methods
		for _, method := range iface.Methods {
			lcMethod := cases.Camel(method.Name)

			_, _ = data.WriteString("func (f *")
			_, _ = data.WriteString(fakeName)
			_, _ = data.WriteString(") ")
			_, _ = data.WriteString(method.Name)
			_, _ = data.WriteString("(")

			for i, param := range method.Parameters {
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(" ")
				_, _ = data.WriteString(param.Type)

				if i != len(method.Parameters)-1 {
					_, _ = data.WriteString(", ")
				}
			}

			_, _ = data.WriteString(") ")

			if len(method.Results) > 0 {
				_, _ = data.WriteString("(")

				for i, result := range method.Results {
					_, _ = data.WriteString(result.Name)
					_, _ = data.WriteString(" ")
					_, _ = data.WriteString(result.Type)

					if i != len(method.Results)-1 {
						_, _ = data.WriteString(", ")
					}
				}

				_, _ = data.WriteString(") ")
			}

			_, _ = data.WriteString("{\n")

			//// Method body

			// Check fake initialization.
			_, _ = data.WriteString(`if f == nil || f.t == nil {
				panic(errFakeNotInitialized)
			}
			
			`)

			// Declare method name constant.
			_, _ = data.WriteString(`const methodName = "`)
			_, _ = data.WriteString(method.Name)
			_, _ = data.WriteString(`"
			
			`)

			// Check method initialization.
			_, _ = data.WriteString("if f.On == nil || f.On.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString(" == nil {\n\tf.t.Fatalf(errMethodNotInitialized, ")
			_, _ = data.WriteString(lcName)
			_, _ = data.WriteString("Name, methodName)\n}\n\n")

			// Get expectation index;
			_, _ = data.WriteString("\tindex := f.On.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString("Counter\n\n")

			// Check expectation index.
			_, _ = data.WriteString("if index+1 > len(f.On.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString(") {\n\tf.t.Fatalf(errExpectationsMissing, ")
			_, _ = data.WriteString(lcName)
			_, _ = data.WriteString("Name, methodName, index+1, len(f.On.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString("))\n}\n\n")

			// Check arguments.
			for _, param := range method.Parameters {
				if param.Type == contextType {
					_, _ = data.WriteString("if ")
					_, _ = data.WriteString(param.Name)
					_, _ = data.WriteString(" == nil {\n\tf.t.Fatalf(errContextNil, ")
					_, _ = data.WriteString(lcName)
					_, _ = data.WriteString("Name, methodName)\n}\n\n")

					_, _ = data.WriteString("if err := ")
					_, _ = data.WriteString(param.Name)
					_, _ = data.WriteString(".Err(); err != nil {\n\tf.t.Fatalf(errContextCancel, ")
					_, _ = data.WriteString(lcName)
					_, _ = data.WriteString("Name, methodName, err.Error())\n}\n\n")

					continue
				}

				_, _ = data.WriteString("if ")
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(" != f.On.")
				_, _ = data.WriteString(lcMethod)
				_, _ = data.WriteString("[index].")
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(" {\n\tf.t.Fatalf(errArgumentMismatch, ")
				_, _ = data.WriteString(lcName)
				_, _ = data.WriteString(`Name, methodName, "`)
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(`", `)
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(", f.On.")
				_, _ = data.WriteString(lcMethod)
				_, _ = data.WriteString("[index].")
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(")\n}\n\n")
			}

			// Increment expectation counter.
			_, _ = data.WriteString("f.On.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString("Counter++\n\n")

			// Return statement.
			if len(method.Results) > 0 {
				_, _ = data.WriteString("\treturn ")

				for i, result := range method.Results {
					_, _ = data.WriteString("f.On.")
					_, _ = data.WriteString(lcMethod)
					_, _ = data.WriteString("[index].returns.")
					_, _ = data.WriteString(cases.Camel(result.Name))

					if i != len(method.Results)-1 {
						_, _ = data.WriteString(", ")
					}
				}
			}

			_, _ = data.WriteString("\n}\n\n")

			//// Method expect.
			_, _ = data.WriteString("func (e *")
			_, _ = data.WriteString(lcName)
			_, _ = data.WriteString("Expect) ")
			_, _ = data.WriteString(method.Name)
			_, _ = data.WriteString("(")

			for i, param := range method.Parameters {
				if param.Type == contextType {
					continue
				}

				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(" ")
				_, _ = data.WriteString(param.Type)

				if i != len(method.Parameters)-1 {
					_, _ = data.WriteString(", ")
				}
			}

			expectName := lcName + method.Name

			_, _ = data.WriteString(") *")
			_, _ = data.WriteString(expectName)
			_, _ = data.WriteString(" {\n")
			_, _ = data.WriteString(`	if e == nil {
					return nil
				}
			
				e.`)
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString(" = append(e.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString(", &")
			_, _ = data.WriteString(expectName)
			_, _ = data.WriteString("{\n")

			for _, param := range method.Parameters {
				if param.Type == contextType {
					continue
				}

				_, _ = data.WriteString("\t\t")
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(": ")
				_, _ = data.WriteString(param.Name)
				_, _ = data.WriteString(",\n")
			}

			_, _ = data.WriteString("})\n\nreturn e.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString("[len(e.")
			_, _ = data.WriteString(lcMethod)
			_, _ = data.WriteString(")-1]\n}\n\n")
		}

		_, _ = data.WriteString("\n")

		formattedData, err := format.Source(data.Bytes())
		if err != nil {
			log.Printf("error: formatting file '%s': %s\n", fileName, err.Error())

			_, err = goFile.Write(data.Bytes())
			if err != nil {
				log.Printf("error: writing to file '%s': %s\n", fileName, err.Error())
			}
		} else {
			_, err = goFile.Write(formattedData)
			if err != nil {
				log.Printf("error: writing to file '%s': %s\n", fileName, err.Error())
			}
		}

		_ = goFile.Close()
	}
}
