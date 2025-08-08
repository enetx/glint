package rules

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"

	"github.com/enetx/glint/checker"
)

var percentFormatPattern = regexp.MustCompile(`%(\[[^\]]+\])?[\+\-\#0 ]*(\d+|\*)?(\.(\d+|\*))?[bcdefgopqstvxXU%]`)

func GPercentFormatRule(ctx *checker.Context, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) < 1 {
			return true
		}

		var (
			funcName string
			isFromG  bool
		)

		switch fun := call.Fun.(type) {
		case *ast.SelectorExpr:
			if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "g" {
				funcName = fun.Sel.Name
				isFromG = true
			}
		case *ast.Ident:
			funcName = fun.Name
		}

		if !isTargetFunc(funcName) {
			return true
		}

		argOffset := 1
		if funcName == "Write" || funcName == "Writeln" {
			argOffset = 2
		}

		if len(call.Args) < argOffset {
			return true
		}

		lit, ok := call.Args[argOffset-1].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}

		text, err := strconv.Unquote(lit.Value)
		if err != nil {
			return true
		}

		placeholderCount := countPlaceholders(text)
		argCount := len(call.Args) - argOffset

		switch {
		case percentFormatPattern.MatchString(text):
			ctx.Reportf(
				lit.Pos(),
				"%s: `%%` formatting detected, use `{}` or `{name}` instead",
				fullFuncName(isFromG, funcName),
			)

		case placeholderCount == 0 && argCount > 0:
			ctx.Reportf(
				lit.Pos(),
				"%s: possible missing `{}` or `{name}` in format string",
				fullFuncName(isFromG, funcName),
			)

		case placeholderCount > 0 && placeholderCount < argCount:
			ctx.Reportf(lit.Pos(), "%s: not enough `{}` or `{name}` placeholders (%d placeholders, %d args)",
				fullFuncName(isFromG, funcName), placeholderCount, argCount)
		}

		return true
	})
}

func isTargetFunc(name string) bool {
	switch name {
	case "Print", "Println", "Eprint", "Eprintln", "Write", "Writeln", "Errorf", "Format":
		return true
	default:
		return false
	}
}

func fullFuncName(isFromG bool, name string) string {
	if isFromG {
		return "g." + name
	}

	return name
}

func countPlaceholders(s string) int {
	count := 0
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '{' {
			if i > 0 && runes[i-1] == '\\' {
				continue
			}

			braceCount := 1
			if i+1 < len(runes) && runes[i+1] == '{' {
				if i > 0 && runes[i-1] == '\\' {
					continue
				}
				braceCount = 2
				i++
			}

			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '}' {
					if braceCount == 2 && j+1 < len(runes) && runes[j+1] == '}' {
						count++
						i = j + 1
						break
					} else if braceCount == 1 {
						count++
						i = j
						break
					}
				}
			}
		}
	}

	return count
}
