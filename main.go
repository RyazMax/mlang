package main

import (
	"fmt"
	"mlang/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the MLang programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	in := os.Stdin
	out := os.Stdout
	if len(os.Args) > 2 {
		in, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Can not open file %s", os.Args[1])
		}
		out, err = os.Create(os.Args[2])
		if err != nil {
			fmt.Printf("Can not open file %s", os.Args[1])
		}
	}
	repl.Start(in, out)
}
