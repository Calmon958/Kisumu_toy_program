// Read Eval Print Loop(console or interactive mode)
// Reads inputs, sends it to the interpreter for Evaluation, Prints the result/output of the interpreter and starts again(Loop)
package repel

import (
	"bufio"
	"fmt"
	"io"
	Lex "token/lexer"
	Tok "token/token"
)

const PROMPT = ">> "

/* Read from input source until you encounter a newline
Take the read line and pass it to an instance of our lexer
Print all the tokens the lexer gives up to the end
*/

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := Lex.New(line)
		for tok := l.NextToken(); tok.Type != Tok.EOF; tok = l.NextToken(){
			fmt.Printf("%+v\n", tok)
		}
	}
}
