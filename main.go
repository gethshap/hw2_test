package main

import term "hw2_test/term"

func main() {
	lexer := term.NewLexer("0")
	dic := make(map[string]*term.Term, 100)
	r, err := term.GetTerm(lexer, &dic, 0)
	if err != nil {
		r.Args = r.Args
	}
}
