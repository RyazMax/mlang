package evaluator

import (
	"mlang/ast"
	"mlang/lexer"
	"mlang/object"
	"mlang/parser"
	"testing"
)

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func eval(stmts []ast.Statement) []object.Object {
	env := object.NewEnvironment()
	return EvalProgram(stmts, env)
}

func TestEvalExpression(t *testing.T) {
	input := `5 * 1 - 2 * (1 + 2)`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	evaluated := eval(program.Statements)

	if len(evaluated) != 1 {
		t.Fatalf("evaluated not has %d objects, got %d", 1, len(evaluated))
	}

	integer, ok := evaluated[0].(*object.Integer)

	if !ok {
		t.Fatalf("evaluated[0] should be object.Integer got %T", evaluated[0])
	}

	if integer.Value != -1 {
		t.Fatalf("integer.Value should be %d, got %d", -1, integer.Value)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (true) {
		1 + 1
		3
		} else {
			4
		}`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	evaluated := eval(program.Statements)

	if len(evaluated) != 1 {
		t.Fatalf("evaluated not has %d objects, got %d", 1, len(evaluated))
	}

	integer, ok := evaluated[0].(*object.Integer)

	if !ok {
		t.Fatalf("evaluated[0] should be object.Integer got %T", evaluated[0])
	}

	if integer.Value != 3 {
		t.Fatalf("integer.Value should be %d, got %d", 3, integer.Value)
	}
}

func TestFuncExpression(t *testing.T) {
	input := `
	f = func(x, y) {return x + y; x - y}
	-f(5, 4) `

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	evaluated := eval(program.Statements)

	if len(evaluated) != 2 {
		t.Fatalf("evaluated not has %d objects, got %d", 2, len(evaluated))
	}

	_, ok := evaluated[0].(*object.Function)
	if !ok {
		t.Fatalf("evaluated[0] should be object.Function, got %T", evaluated[0])
	}

	integer, ok := evaluated[1].(*object.Integer)
	if !ok {
		t.Fatalf("evaluated[1] should be object.Integer, got %T", evaluated[1])
	}

	if integer.Value != -9 {
		t.Fatalf("integer.Value should be %d, got %d", -9, integer.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		expr string
		res  bool
	}{
		{"true", true},
		{"false", false},
		{"true && false", false},
		{"true || false", true},
		{"true ^ false", true},
		{"false ^ false", false},
		{"5 < 1 || 4 > 2", true},
		{"true && (false || true)", true},
		{"5 > 5 || 5 == 5", true},
		{"!(5 == 5) && true", false},
		{"5 == null", false},
		{"false == null", false},
		{"5 || false", true},
	}

	for i, tt := range tests {
		l := lexer.New(tt.expr)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		evaluated := eval(program.Statements)

		if len(evaluated) != 1 {
			t.Fatalf("tests[%d]: evaluated should have %d objects, got %d", i, 1, len(evaluated))
		}

		boolean, ok := evaluated[0].(*object.Boolean)
		if !ok {
			t.Fatalf("tests[%d] evaluated[0] should be object.Boolean, got %T", i, evaluated[0])
		}

		if boolean.Value != tt.res {
			t.Fatalf("tests[%d] boolean.Value should be %v, got %v", i, tt.res, boolean.Value)
		}
	}
}

func TestBlockStatement(t *testing.T) {
	tests := []struct {
		expr string
		res  object.Object
	}{
		{"{1;2}", &object.Integer{Value: 2}},
		{"{}", NULL},
		{"{false}", FALSE},
	}

	for i, tt := range tests {
		l := lexer.New(tt.expr)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		evaluated := eval(program.Statements)

		if len(evaluated) != 1 {
			t.Fatalf("tests[%d]: evaluated should have %d objects, got %d", i, 1, len(evaluated))
		}

		if !isEqual(evaluated[0], tt.res) {
			t.Fatalf("tests[%d] result should be %v, got %v", i, tt.res, evaluated[0])
		}
	}
}

func TestWithError(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		{"5 / 0", "division by zero 5 / 0"},
		{"func(x,y) {x + y}(1)", "function expects 2 arguments, 1 was given"},
		{"func() {1} + 5", "type mismatch: FUNCTION + INTEGER"},
		{"b = 1 + a", "identifier not found: a"},
		{"1(5)", "not a function: INTEGER"},
		{"1 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"true + true", "unknown operator: BOOLEAN + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"{f = func(){f()};f()}", "max recursion level reached"},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		evaluated := eval(program.Statements)

		if len(evaluated) != 1 {
			t.Fatalf("tests[%d]: evaluated should have %d objects, got %d", i, 1, len(evaluated))
		}

		err, ok := evaluated[0].(*object.Error)
		if !ok {
			t.Fatalf("tests[%d] evaluated[0] should be error, got %T", i, evaluated[0])
		}

		if err.Message != tt.err {
			t.Fatalf("tests[%d] err.Message should be \"%s\", got \"%s\"", i, tt.err, err.Message)
		}
	}
}

/*func TestRecursion(t *testing.T) {
	input := `f = func(){f()}; f()`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	evaluated := eval(program.Statements)

	if len(evaluated) != 2 {
		t.Fatalf("evaluated should have %d objects, go")
	}
}*/

func isEqual(a object.Object, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *object.Integer:
		t, _ := b.(*object.Integer)
		return a.Value == t.Value
	default:
		return a == b
	}
}
