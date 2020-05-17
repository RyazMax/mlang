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

	in := os.Stdin
	out := os.Stdout
	interactive := true
	if len(os.Args) > 1 {
		in, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Can not open file %s", os.Args[1])
			return
		}
		if len(os.Args) > 2 {
			out, err = os.Create(os.Args[2])
			if err != nil {
				fmt.Printf("Can not open file %s", os.Args[1])
				return
			}
		}
		interactive = false
	}
	if interactive {
		fmt.Printf("Hello %s! This is the MLang programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
	}
	repl.Start(in, out, interactive)
}
