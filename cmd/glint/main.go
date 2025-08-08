package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/enetx/glint/checker"
	"github.com/enetx/glint/internal/rules"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myrevive <file1.go> [file2.go...]")
		os.Exit(1)
	}

	fset := token.NewFileSet()

	ch := checker.New(
		rules.GPercentFormatRule,
		rules.AppendAliasRule,
	)

	for _, path := range os.Args[1:] {
		if filepath.Ext(path) != ".go" {
			continue
		}

		file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", path, err)
			continue
		}

		ctx := &checker.Context{
			Fset: fset,
			File: path,
		}

		ch.Run(ctx, file)
	}
}
