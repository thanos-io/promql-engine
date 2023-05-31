package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loadAnalyzer)
}

var loadAnalyzer = &analysis.Analyzer{
	Name: "loadvet",
	Doc:  "reports ill-formatted prometheus load directives",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if !strings.HasSuffix(file.Name.String(), "_test") {
			continue
		}
		var stack []ast.Node
		ast.Inspect(file, func(n ast.Node) bool {
			defer func() {
				if n == nil {
					stack = stack[:len(stack)-1]
				} else {
					stack = append(stack, n)
				}
			}()
			if n == nil || len(stack) == 0 {
				return true
			}
			parent := stack[len(stack)-1]
			if _, ok := parent.(*ast.KeyValueExpr); !ok {
				return true
			}

			s, ok := n.(*ast.BasicLit)
			if !ok {
				return true
			}
			if s.Kind != token.STRING {
				return true
			}
			quote := s.Value[0]
			if quote != '`' {
				return true
			}
			cont := s.Value[1 : len(s.Value)-1]
			if !strings.HasPrefix(cont, "load") {
				return true
			}
			position := pass.Fset.Position(s.Pos())
			lineAtPosition, err := readLine(position.Filename, position.Line)
			if err != nil {
				return true
			}
			whiteSpace := leadingWhitespace(lineAtPosition)
			if formatted := formatLoadDirective(cont, whiteSpace); cont != formatted {
				pass.Report(analysis.Diagnostic{
					Pos:     s.Pos(),
					Message: "ill-formatted load directive found",
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: fmt.Sprintf("Should replace '%s' with '%s'", cont, formatted),
							TextEdits: []analysis.TextEdit{
								{
									Pos:     s.Pos(),
									End:     s.End(),
									NewText: []byte(fmt.Sprintf("%c%s%c", quote, formatted, quote)),
								},
							},
						},
					},
				})
			}

			return false
		})
	}
	return nil, nil
}

func formatLoadDirective(load, whiteSpace string) string {
	var res strings.Builder

	sc := bufio.NewScanner(strings.NewReader(load))
	for i := 0; sc.Scan(); i++ {
		if i == 0 {
			res.WriteString(sc.Text())
		} else {
			res.WriteString(whiteSpace)
			res.WriteString(strings.Repeat(" ", 4))
			res.WriteString(strings.TrimSpace(sc.Text()))
		}
		res.WriteByte('\n')
	}

	return strings.TrimSpace(res.String())
}

func readLine(filename string, line int) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	sc := bufio.NewScanner(bytes.NewReader(content))
	for i := 0; sc.Scan(); i++ {
		if i == line-1 {
			return sc.Text(), sc.Err()
		}
	}
	return "", io.EOF
}

func leadingWhitespace(line string) string {
	var res strings.Builder
	for _, r := range line {
		if unicode.IsSpace(r) {
			res.WriteRune(r)
		} else {
			break
		}
	}
	return res.String()
}
