# Building a personal interpreter
### 28-29th August

Learning Lexer and how to implement varios things in a compiler, i.e, variables, function, constants.
Test for the first part of the creation of a Lexer.
First && Second commits contains the progress I've done so far.


### 30th August
Today I have worked on the extension of the lexer and token set to accomodate for various keywords and two-character tokens.
Understanding that everytime you change the input in the lexer testfile don't forget to right it in the struct and follow the procedure. Example for `let x == 5` will be;
```bash
// [...]
{token.LET, "let"},
{token.IDEN, "x"},
{token.EQ, "=="},
{token.INT, "5"},
// [...]
```
Also take key note of the changes you make on the lexer file. We have modified the switch case to account for the two-character tokens by creating a peekChar function which checks what the next character is without increamenting it.
Take note I said check not move.
we also add some mods on `func (l *lexer) NextToken() token.Token` within the switch case to accomodate the two-character tokens 

Created REPL(Read Eval Print Loop) function which essentially is a console to run our interpreter from. It reads input, sends it to the interpreter to be evaluated, prints the result and repeats the same process.


Created main.go file for running everything and tested the program so far. It's essentially a working progress but serves as a great start.
it Print the command enter broken down into `Type:` and `Literal:`; basically if we use the previous exampl our output will be

```bash
{Type:LET Literal:let}
{TYPE:IDENT Literal:x}
{TYPE:== Literal:==}
{TYPE:INT Literal:5}
```
Thats all for today.

### 2 Sep
Started on Parser....Was a bit more complex than I initially thought but it's the part of the program that has most of the work.
Essentially the parse takes source code and builds data structure analysing the input to check if it has a certain structure.
Did the first part for creating the AST with basic skeleton and the parse with basic parts. Run into an error to be dealt with when next endulging in the project. 