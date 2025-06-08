// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

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

	"github.com/prometheus/prometheus/promql/parser"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loadAnalyzer)
}

var loadAnalyzer = &analysis.Analyzer{
	Name: "loadvet",
	Doc:  "reports ill-formatted prometheus load directives or PromQL expressions",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
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
			// only format expressions that are behind "query" or "load" keys
			parent, ok := stack[len(stack)-1].(*ast.KeyValueExpr)
			if !ok {
				return true
			}
			p, ok := parent.Key.(*ast.Ident)
			if !ok {
				return true
			}
			switch strings.ToLower(p.Name) {
			case "query", "load":
			default:
				return true
			}

			s, ok := n.(*ast.BasicLit)
			if !ok {
				return true
			}
			if s.Kind != token.STRING {
				return true
			}
			position := pass.Fset.Position(s.Pos())
			lineAtPosition, err := readLine(position.Filename, position.Line)
			if err != nil {
				return true
			}
			whiteSpace := leadingWhitespace(lineAtPosition)
			cont := s.Value[1 : len(s.Value)-1]
			// for consistency and ease of replacement, we replace the quotes with ` here
			quote := s.Value[0]
			switch {
			case looksLikeLoadStmt(cont):
				if formatted := formatLoadDirective(cont, whiteSpace); cont != formatted || quote != '`' {
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
										NewText: fmt.Appendf(nil, "%c%s%c", '`', formatted, '`'),
									},
								},
							},
						},
					})
					return false
				}
			case looksLikePromQL(cont):
				if formatted := formatPromQL(cont); cont != formatted || quote != '`' {
					pass.Report(analysis.Diagnostic{
						Pos:     s.Pos(),
						Message: "ill-formatted promql found",
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: fmt.Sprintf("Should replace '%s' with '%s'", cont, formatted),
								TextEdits: []analysis.TextEdit{
									{
										Pos:     s.Pos(),
										End:     s.End(),
										NewText: fmt.Appendf(nil, "%c%s%c", '`', formatted, '`'),
									},
								},
							},
						},
					})
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}

func looksLikeLoadStmt(cont string) bool {
	return strings.HasPrefix(strings.TrimSpace(cont), "load")
}
func looksLikePromQL(cont string) bool {
	_, err := parser.ParseExpr(cont)
	return err == nil
}

func formatLoadDirective(load, whiteSpace string) string {
	var res strings.Builder

	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(load)))
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

func formatPromQL(exprStr string) string {
	expr, _ := parser.ParseExpr(exprStr)
	pretty := expr.Pretty(0)
	if strings.Count(pretty, "\n") > 0 {
		return "\n" + pretty
	}
	return pretty
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
