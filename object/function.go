package object

import (
	"bytes"
	"mlang/ast"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("<function (")
	for _, p := range f.Parameters {
		out.WriteString(p.Value)
		out.WriteString(" ")
	}
	out.WriteString(")>")
	return out.String()
}
