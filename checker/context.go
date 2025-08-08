package checker

import (
	"fmt"
	"go/token"
)

type Context struct {
	Fset *token.FileSet
	File string
}

func (ctx *Context) Reportf(pos token.Pos, format string, args ...any) {
	position := ctx.Fset.Position(pos)
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s: %s\n", position.String(), msg)
}
