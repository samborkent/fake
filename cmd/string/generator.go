package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

const (
	tokenPackage = "package "

	tokenImport      = "import "
	tokenImportEntry = `"`
	tokenParamStart  = "("
	tokenParamEnd    = ")"

	tokenComment    = "//"
	tokenFunction   = "func "
	tokenScopeStart = "{"
	tokenScopeEnd   = "}"

	tokenInterface     = "interface"
	tokenInterfaceType = "|"
)

func main() {
	content, err := os.ReadFile("interface.go")
	if err != nil {
		log.Fatal("reading file: " + err.Error())
	}

	lines := strings.Split(string(content), "\n")
	lines = removePackage(lines)

	imports, lastImportIndex, err := getImports(lines)
	if err != nil {
		log.Fatal("getting imports: " + err.Error())
	}

	log.Printf("imports: %v\n", imports)

	lines = filterLines(lines[lastImportIndex+1:])

	for i, line := range lines {
		log.Printf("%d: %s", i, line)
	}
}

// func findInterfaces(lines []string) []string {
// }

func removePackage(lines []string) []string {
	index := 0

	for i, line := range lines {
		if strings.HasPrefix(line, tokenPackage) {
			index = i
			break
		}
	}

	// Skip all lines before package, and empty line after.
	return lines[index+2:]
}

func getImports(lines []string) ([]string, int, error) {
	firstImport := 0
	lastImport := 0

	for i, line := range lines {
		if strings.HasPrefix(line, tokenImport) {
			if !strings.Contains(line, tokenParamStart) {
				importStart := strings.Index(line, tokenImportEntry)
				if importStart == -1 {
					return nil, 0, errors.New("missing import start for single-lined import")
				}

				importEnd := strings.LastIndex(line, tokenImportEntry)
				if importEnd == -1 {
					return nil, 0, errors.New("missing import end for single-lined import")
				}

				return []string{line[importStart+1 : importEnd]}, i + 1, nil
			}

			firstImport = i + 1
		} else if line == tokenParamEnd {
			lastImport = i
			break
		}
	}

	imports := make([]string, 0, lastImport-firstImport)

	for _, line := range lines[firstImport:lastImport] {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			continue
		}

		imports = append(imports, strings.Trim(trimmedLine, tokenImportEntry))
	}

	return imports, lastImport + 1, nil
}

func filterLines(lines []string) []string {
	filteredLines := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], tokenComment) {
			continue
		} else if strings.HasPrefix(lines[i], tokenFunction) {
			for {
				if strings.Contains(lines[i], tokenScopeStart) {
					break
				}

				i++
			}

			i++
			numScopes := 1
			endScopes := 0

			for {
				if strings.Contains(lines[i], tokenScopeStart) {
					numScopes++
				} else if strings.Contains(lines[i], tokenScopeEnd) {
					endScopes++
				}

				i++

				if numScopes == endScopes {
					break
				}

			}
		} else if strings.Contains(lines[i], tokenInterface) {
			if strings.Contains(lines[i], tokenScopeEnd) {
				// Single-line interface
				if !strings.Contains(lines[i], tokenParamStart) && !strings.Contains(lines[i], tokenParamEnd) {
					continue
				}
			}
		}

		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		filteredLines = append(filteredLines, lines[i])
	}

	return filteredLines
}
