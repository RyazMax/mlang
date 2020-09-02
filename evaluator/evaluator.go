package evaluator

/*
	Модуль осуществляющий выполнение абстрактного синтаксического дерева программы.
	Основные функции EvalProgram и Eval осуществляют исполнение программы или другого узла аст
*/

import (
	"fmt"
	"mlang/ast"
	"mlang/object"
)

const MAX_RECURSION_LEVEL = 90000

var lvl int

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func EvalProgram(stmts []ast.Statement, env *object.Environment) []object.Object {
	builtins["forn"] = &object.Builtin{Fn: fornFunc}
	var (
		results []object.Object
		result  object.Object
	)

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			results = append(results, result.Value)
			return results
		case *object.Error:
			results = append(results, result)
			return results
		}
		results = append(results, result)
	}

	return results
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	lvl += 1
	defer func() { lvl -= 1 }()
	if lvl > MAX_RECURSION_LEVEL {
		return newError("max recursion level reached")
	}
	switch node := node.(type) {
	case *ast.AssignStatement:
		return evalAssignStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Null:
		return NULL
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value != 0
	default:
		return !(obj == NULL || obj == FALSE)
	}
}

func evalAssignStatement(as *ast.AssignStatement, env *object.Environment) object.Object {
	val := Eval(as.Value, env)
	if isError(val) {
		return val
	}
	env.Set(as.Name.Value, val)
	return val
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(id.Value); ok {
		return val
	}

	if builtin, ok := builtins[id.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", id.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	result = NULL

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	return nativeBoolToBooleanObject(!isTruthy(right))
}

func isLogicOperator(operator string) bool {
	switch operator {
	case "||":
		fallthrough
	case "&&":
		fallthrough
	case "^":
		return true
	default:
		return false
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.NULL_OBJ && right.Type() == object.NULL_OBJ:
		return evalNullInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case isLogicOperator(operator):
		l, r := nativeBoolToBooleanObject(isTruthy(left)), nativeBoolToBooleanObject(isTruthy(right))
		return evalBooleanInfixExpression(operator, l, r)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalNullInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero %s %s %s", left.Inspect(), operator, right.Inspect())
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case "^":
		leftBool, rightBool := isTruthy(left), isTruthy(right)
		return nativeBoolToBooleanObject(leftBool != rightBool)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	case "&&":
		return nativeBoolToBooleanObject(left == TRUE && right == TRUE)
	case "||":
		return nativeBoolToBooleanObject(left == TRUE || right == TRUE)
	case "^":
		return nativeBoolToBooleanObject(left != right)
	case "<":
		return nativeBoolToBooleanObject(left == FALSE && right == TRUE)
	case ">":
		return nativeBoolToBooleanObject(left == TRUE && right == FALSE)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return newError("function expects %d arguments, %d was given",
				len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
