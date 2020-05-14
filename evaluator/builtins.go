package evaluator

import (
	"fmt"
	"mlang/object"
)

var builtins = map[string]*object.Builtin{
	"sum": &object.Builtin{
		Fn: sumFunc,
	},
	"print": &object.Builtin{
		Fn: printFunc,
	},
}

func printFunc(args ...object.Object) object.Object {
	for _, obj := range args {
		fmt.Printf("%+v ", obj.Inspect())
	}
	fmt.Println()
	return NULL
}

func sumFunc(args ...object.Object) object.Object {
	var acc int64
	for _, obj := range args {
		switch obj := obj.(type) {
		case *object.Integer:
			acc += obj.Value
		case *object.Boolean:
			if obj.Value {
				acc++
			}
		default:
			return newError("sum expects integers or boolean. got=%s", obj.Type())
		}
	}

	return &object.Integer{Value: acc}
}
