package evaluator

import (
	"fmt"
	"mlang/object"
	"time"
)

var builtins = map[string]*object.Builtin{
	"sum": {
		Fn: sumFunc,
	},
	"print": {
		Fn: printFunc,
	},
	"read": {
		Fn: readFunc,
	},
	"bool": {
		Fn: toBool,
	},
	"time": {
		Fn: timeFunc,
	},
}

func fornFunc(args ...object.Object) object.Object {
	if len(args) < 2 {
		return newError("not enouth arguments")
	}
	switch obj := args[0].(type) {
	case *object.Integer:
	default:
		return newError("forn expects integer as first argument. got=%s", obj.Type())
	}
	switch obj := args[1].(type) {
	case *object.Function:
	default:
		return newError("forn expects function as second argument. got=%s", obj.Type())
	}
	for i := int64(0); i < args[0].(*object.Integer).Value; i++ {
		fun := args[1].(*object.Function)
		fun.Env.Set("i", &object.Integer{Value: i})
		applyFunction(fun, args[2:])
	}
	return NULL
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
