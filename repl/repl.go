package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"mlang/evaluator"
	"mlang/lexer"
	"mlang/object"
	"mlang/parser"
	"os"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer, interactive bool) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("Something went wrong: %v", x)
			if interactive {
				Start(in, out, interactive)
			}
		}
	}()
	if interactive {
		startShell(in, out)
	} else {
		startFile(in, out)
	}
}

func startFile(in io.Reader, out io.Writer) {
	text_bytes, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Println("Can not read input file")
		return
	}

	text := string(text_bytes)
	env := object.NewEnvironment()
	l := lexer.New(text)
	p := parser.New(l)

	program := p.ParseProgram()
	errors := p.Errors()

	if len(errors) != 0 {
		printParserErrors(out, errors)
		return
	}

	evaluated := evaluator.EvalProgram(program.Statements, env)
	if out != os.Stdout {
		for _, stmt := range evaluated {
			io.WriteString(out, stmt.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func startShell(in io.Reader, out io.Writer) {
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

		evaluated := evaluator.EvalProgram(program.Statements, env)
		if evaluated != nil {
			io.WriteString(out, evaluated[len(evaluated)-1].Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
