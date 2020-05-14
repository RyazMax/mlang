package repl

import (
	"bufio"
	"io"
	"mlang/evaluator"
	"mlang/lexer"
	"mlang/object"
	"mlang/parser"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		io.WriteString(out, PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		errors := p.Errors()

		if len(errors) != 0 {
			printParserErrors(out, errors)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
