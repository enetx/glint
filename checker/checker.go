package checker

import "go/ast"

type Rule func(*Context, *ast.File)

type Checker struct {
	Rules []Rule
}

func New(rules ...Rule) *Checker {
	return &Checker{Rules: rules}
}

func (c *Checker) Run(ctx *Context, file *ast.File) {
	for _, rule := range c.Rules {
		rule(ctx, file)
	}
}
