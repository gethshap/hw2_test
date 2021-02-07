package term

import (
	"errors"
)

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <term>     ::= ATOM | NUM | VAR | <compound>
// <compound> ::= <functor> LPAR <args> RPAR
// <functor>  ::= ATOM
// <args>     ::= <term> | <term> COMMA <args>
//

// Parser is the interface for the term parser.
// Do not change the definition of this interface.
type Parser interface {
	Parse(string) (*Term, error)
}

type P struct {
	unique map[string]*Term
}

func wrapper_read() {

}

func checkEOF(l *lexer) (bool, error) {
	tok, err := l.next()
	if err == nil {
		if tok.typ == tokenEOF {
			return true, err
		} else {
			return false, err
		}

	}
	return false, err
}

func makeVariableTerm(tok *Token, dic *map[string]*Term) *Term {
	r := new(Term)
	r.Typ = TermVariable
	r.Literal = tok.literal
	if val, ok := (*dic)[r.Literal]; ok {
		return val
	} else {
		(*dic)[r.Literal] = r
		return r
	}
}

func makeNumberTerm(tok *Token, dic *map[string]*Term) *Term {
	r := new(Term)
	r.Typ = TermNumber
	r.Literal = tok.literal
	if val, ok := (*dic)[r.Literal]; ok {
		return val
	} else {
		(*dic)[r.Literal] = r
		return r
	}

}

func makeAtomTerm(tok *Token, dic *map[string]*Term) *Term {
	r := new(Term)
	r.Typ = TermAtom
	r.Literal = tok.literal
	if val, ok := (*dic)[r.Literal]; ok {
		return val
	} else {
		(*dic)[r.Literal] = r
		return r
	}

}

func getTermInArgs(lexer *lexer, dic *map[string]*Term) (error, bool, bool, *Term) {
	tok, err := lexer.next()
	if err == nil {
		if tok.typ == tokenVariable {
			if nextTok, _ := lexer.next(); nextTok.typ == tokenComma {
				return nil, true, false, makeVariableTerm(tok, dic)
			} else if nextTok.typ == tokenRpar {
				return nil, false, true, makeVariableTerm(tok, dic)
			} else {
				return ErrParser, false, false, nil
			}
		} else if tok.typ == tokenNumber {
			if nextTok, _ := lexer.next(); nextTok.typ == tokenComma {
				return nil, true, false, makeNumberTerm(tok, dic)
			} else if nextTok.typ == tokenRpar {
				return nil, false, true, makeNumberTerm(tok, dic)
			} else {
				return ErrParser, false, false, nil
			}
		} else if tok.typ == tokenAtom {
			if nextTok, _ := lexer.next(); nextTok.typ == tokenLpar {
				factor_term := new(Term)
				factor_term.Functor = nil
				factor_term.Typ = TermAtom
				factor_term.Args = nil
				factor_term.Literal = tok.literal

				r := new(Term)
				r.Typ = TermCompound
				r.Functor = factor_term
				for true {
					err, hasNext, isTermination, term := getTermInArgs(lexer, dic)
					if err == nil {
						if hasNext {
							r.Args = append(r.Args, term)
							continue
						} else if isTermination {
							r.Args = append(r.Args, term)
							if nextTok, _ := lexer.next(); nextTok.typ == tokenComma {
								return nil, true, false, r
							} else if nextTok.typ == tokenRpar {
								return nil, false, true, r
							} else {
								return ErrParser, false, false, nil
							}

						}
					}
				}
			} else if nextTok.typ == tokenComma {
				return nil, true, false, makeAtomTerm(tok, dic)
			} else if nextTok.typ == tokenRpar {
				return nil, false, true, makeAtomTerm(tok, dic)
			}

		}
	}
	return ErrParser, false, false, nil

}

func GetTerm(lexer *lexer, dic *map[string]*Term, inargstate int) (*Term, error) {
	tok, err := lexer.next()

	if inargstate == 0 {
		if err == nil {
			if tok.typ == tokenNumber {
				if val, ok := (*dic)[tok.literal]; ok {
					return val, nil
				}
				r := new(Term)
				r.Typ = TermNumber
				r.Literal = tok.literal
				r.Args = nil
				r.Functor = nil
				if check, _ := checkEOF(lexer); check {
					if val, ok := (*dic)[r.Literal]; ok {
						return val, nil
					} else {
						(*dic)[r.Literal] = r
						return r, nil
					}
				} else {
					return nil, ErrParser
				}
			} else if tok.typ == tokenVariable {
				if val, ok := (*dic)[tok.literal]; ok {
					return val, nil
				}
				r := new(Term)
				r.Typ = TermVariable
				r.Literal = tok.literal
				r.Args = nil
				r.Functor = nil
				if check, _ := checkEOF(lexer); check {
					if val, ok := (*dic)[r.Literal]; ok {
						return val, nil
					} else {
						(*dic)[r.Literal] = r
						return r, nil
					}
				} else {
					return nil, ErrParser
				}
				///  <term>     ::= ATOM | NUM | VAR | <compound>
				//         <compound> ::= <functor> ( <args> )
				//         <functor>  ::= ATOM
				//         <args>     ::= <term> | <term>, <args>
			} else if tok.typ == tokenAtom {
				if subtok, _ := lexer.next(); subtok.typ == tokenLpar {

					if val, ok := (*dic)[tok.literal]; ok {
						return val, nil
					}
					factor_term := new(Term)
					factor_term.Functor = nil
					factor_term.Typ = TermAtom
					factor_term.Args = nil
					factor_term.Literal = tok.literal

					r := new(Term)
					r.Typ = TermCompound
					r.Literal = ""
					r.Functor = factor_term
					for true {
						err, hasNext, Termination, term := getTermInArgs(lexer, dic)
						if err == nil {
							r.Args = append(r.Args, term)
						} else {
							return nil, ErrParser
						}
						if hasNext {
							continue
						}
						if Termination {
							if check, _ := checkEOF(lexer); check {
								return r, nil
							} else {
								return nil, ErrParser
							}

						}

					}
				} else if subtok.typ == tokenEOF {
					// single atom
					if val, ok := (*dic)[tok.literal]; ok {
						return val, nil
					}
					r := new(Term)
					r.Typ = TermAtom
					r.Literal = tok.literal
					r.Args = nil
					r.Functor = nil
					return r, nil
				}
			} else if tok.typ == tokenEOF {
				return nil, nil
			}
		}
		return nil, ErrParser

	} else if inargstate == 1 {
		if err == nil {
			if tok.typ == tokenNumber {
				if val, ok := (*dic)[tok.literal]; ok {
					return val, nil
				}
				r := new(Term)
				r.Typ = TermNumber
				r.Literal = tok.literal
				r.Args = nil
				r.Functor = nil
				if lexer.peek() == ',' || lexer.peek() == ')' {
					return r, nil
				} else {
					return nil, ErrParser
				}
			} else if tok.typ == tokenVariable {
				if val, ok := (*dic)[tok.literal]; ok {
					return val, nil
				}
				r := new(Term)
				r.Typ = TermVariable
				r.Literal = tok.literal
				r.Args = nil
				r.Functor = nil

				if lexer.peek() == ',' || lexer.peek() == ')' {
					return r, nil
				} else {
					return nil, ErrParser
				}

				///  <term>     ::= ATOM | NUM | VAR | <compound>
				//         <compound> ::= <functor> ( <args> )
				//         <functor>  ::= ATOM
				//         <args>     ::= <term> | <term>, <args>
			} else if tok.typ == tokenAtom {

				if subtok, _ := lexer.next(); subtok.typ == tokenLpar {
					if val, ok := (*dic)[tok.literal]; ok {
						return val, nil
					}
					factor_term := new(Term)
					factor_term.Functor = nil
					factor_term.Typ = TermAtom
					factor_term.Args = nil
					factor_term.Literal = tok.literal

					r := new(Term)
					r.Typ = TermCompound
					r.Literal = ""
					r.Functor = factor_term
					for true {
						subterm, err := GetTerm(lexer, dic, 1)
						if err == nil {
							r.Args = append(r.Args, subterm)
							argtok, err := lexer.next()
							if err == nil {
								if argtok.typ == tokenComma {
									continue
								} else if argtok.typ == tokenRpar {
									return r, nil
								} else {
									return nil, ErrParser
								}
							} else {
								return nil, ErrLexer
							}
						} else {
							return nil, ErrParser
						}
						break
					}
				} else if subtok.typ == tokenComma {
					if val, ok := (*dic)[tok.literal]; ok {
						return val, nil
					}
					r := new(Term)
					r.Typ = TermAtom
					r.Literal = tok.literal
					r.Args = nil
					r.Functor = nil
					return r, nil

				}
			}
		}
		return nil, ErrParser
	}
	return nil, ErrParser
}

func (p P) Parse(str string) (*Term, error) {
	lexer := NewLexer(str)
	dic := make(map[string]*Term)
	r, err := GetTerm(lexer, &dic, 0)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	var p P
	p.unique = make(map[string]*Term)
	return p
}
