package evaluator

import (
	"fmt"
	"mlang/object"
	"time"
)

var builtins = map[string]*object.Builtin{
	"sum": &object.Builtin{
		Fn: sumFunc,
	},
	"print": &object.Builtin{
		Fn: printFunc,
	},
	"read": &object.Builtin{
		Fn: readFunc,
	},
	"bool": &object.Builtin{
		Fn: toBool,
	},
	"time": &object.Builtin{
		Fn: timeFunc,
	},
}

func timeFunc(args ...object.Object) object.Object {
	return &object.Integer{Value: time.Now().UnixNano()}
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

func toBool(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("bool expects only one argument, %d was given", len(args))
	}

	obj := args[0]
	switch obj := obj.(type) {
	case *object.Integer:
		return nativeBoolToBooleanObject(obj.Value != 0)
	default:
		return nativeBoolToBooleanObject(!(obj == NULL || obj == FALSE))
	}
}

func readFunc(args ...object.Object) object.Object {
	if len(args) != 0 {
		return newError("read expects no arguments")
	}

	var a int64
	_, err := fmt.Scanf("%d", &a)

	if err != nil {
		return NULL
	}

	return &object.Integer{Value: a}
}
