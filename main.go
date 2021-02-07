package main

import term "hw2_test/term"


func main() {
	lexer := term.NewLexer("f(1)g")
	dic := new(map[string] *term.Term)
	r,err := term.GetTerm(lexer,dic,0)
	if err != nil{
		r.Args = r.Args
	}
}
