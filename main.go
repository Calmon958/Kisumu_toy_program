package main

import (
	"fmt"
	"os"
	"os/user"
	Repl "token/repl"
)

func main(){
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!\n----------Welcome to Kisumu Programming Language(KPL).----------", user.Username)
	fmt.Printf("Feel free to type in your commands :)\n")
	Repl.Start(os.Stdin, os.Stdout)
}