// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: 
//		jnml, labs.nic.cz
//		Miek Gieben, SIDN, miek@miek.nl

/*

Package lexer provides generating actionless scanners (lexeme recognizers) at run time.

Scanners are defined by regular expressions
and/or lexical grammars, mapping between those definitions, token numeric identifiers and
an optional set of starting id sets, providing simmilar functionality as switching start states in *nix LEX.
The generated FSMs are Unicode rune based and all unicode.Categories and unicode.Scripts are supported by the
regexp syntax using the \p{name} construct.

For additional information please see the README file.

TODO(jnml) complete package docs

*/
package lexer


import (
	"bytes"
	"ebnf"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"os"
	"regexp"
	"strconv"
	"unicode"
)


type Lexer struct {
	nfa    Nfa
	starts []*NfaState
	accept *NfaState
}


// StartSetID is a type of a lexer start set identificator.
// It is used by Begin and PushState.
type StartSetID int


//TODO:full docs
func CompileLexer(starts [][]int, tokdefs map[string]int, grammar, start string) (lexer *Lexer, err os.Error) {
	lexer = &Lexer{}

	defer func() {
		if e := recover(); e != nil {
			lexer = nil
			err = e.(os.Error)
		}
	}()

	var prodnames string
	res, xref := map[int]string{}, map[int]string{}

	for tokdef, id := range tokdefs {
		if _, ok := res[id]; ok {
			panic(fmt.Errorf("duplicate id %d for token %q", id, tokdef))
		}

		xref[id] = fmt.Sprintf("id-%d", id)
		if re, ok := isRE(tokdef); ok {
			res[id] = re
			continue
		}

		if grammar == "" || !isIdent(tokdef) {
			res[id] = regexp.QuoteMeta(tokdef)
			continue
		}

		if prodnames != "" {
			prodnames += " | "
		}
		prodnames += tokdef
		res[id] = ""
	}

	if prodnames != "" {
		var g ebnf.Grammar
		ebnfSrc := grammar + fmt.Sprintf("\n%s = %s .", start, prodnames)
		fset := token.NewFileSet()
		fset.AddFile(start, fset.Base(), len(ebnfSrc))
		if g, err = ebnf.Parse(fset, start, []byte(ebnfSrc)); err != nil {
			panic(err)
		}

		if err = ebnf.Verify(fset, g, start); err != nil {
			panic(err)
		}

		grammarREs := map[*ebnf.Production]string{}
		for tokdef, id := range tokdefs {
			if isIdent(tokdef) {
				res[id], xref[id] = ebnf2RE(g, tokdef, grammarREs), tokdef
			}
		}
	}

	if starts == nil { // create the default, all inclusive start set
		starts = [][]int{{}}
		for id := range res {
			starts[0] = append(starts[0], id)
		}
	}

	lexer.accept = lexer.nfa.NewState()
	lexer.starts = make([]*NfaState, len(starts))
	for i, set := range starts {
		state := lexer.nfa.NewState()
		lexer.starts[i] = state
		for _, id := range set {
			var in, out *NfaState
			re, ok := res[int(id)]
			if !ok {
				panic(fmt.Errorf("unknown token id %d in set %d", id, i))
			}

			if in, out, err = lexer.nfa.ParseRE(fmt.Sprintf("%s-%s", start, xref[int(id)]), re); err != nil {
				panic(err)
			}

			state.AddNonConsuming(&EpsilonEdge{int(id), in})
			out.AddNonConsuming(&EpsilonEdge{0, lexer.accept})
		}
	}

	lexer.nfa.reduce()
	return
}


// MustCompileLexer is like CompileLexer but panics if the definitions cannot be compiled.
// It simplifies safe initialization of global variables holding compiled Lexers. 
func MustCompileLexer(starts [][]int, tokdefs map[string]int, grammar, start string) (lexer *Lexer) {
	var err os.Error
	if lexer, err = CompileLexer(starts, tokdefs, grammar, start); err != nil {
		if list, ok := err.(scanner.ErrorList); ok {
			scanner.PrintError(os.Stderr, list)
		}
		panic(err)
	}
	return
}


func (lx *Lexer) String() (s string) {
	s = lx.nfa.String()
	for i, set := range lx.starts {
		s += fmt.Sprintf("\nstart set %d = {", i)
		for _, edge := range set.NonConsuming {
			s += " " + strconv.Itoa(int(edge.Target().Index))
		}
		s += " }"
	}
	s += "\naccept: " + strconv.Itoa(int(lx.accept.Index))
	return
}


func identFirst(rune int) bool {
	return unicode.IsLetter(rune) || rune == '_'
}


func identNext(rune int) bool {
	return identFirst(rune) || unicode.IsDigit(rune)
}


func isIdent(s string) bool {
	for i, rune := range s {
		if i == 0 && !identFirst(rune) {
			return false
		}

		if !identNext(rune) {
			return false
		}
	}
	return true
}


// isRE checks if a string starts and ends in '/'. If so, return the string w/o the leading and trailing '/' and true.
// Otherwise return the original string and false.
func isRE(s string) (string, bool) {
	if n := len(s); n > 2 && s[0] == '/' && s[n-1] == '/' {
		return s[1 : n-1], true
	}
	return s, false
}


var pipe = map[bool]string{false: "", true: "|"}


func ebnf2RE(g ebnf.Grammar, name string, res map[*ebnf.Production]string) (re string) {
	p := g[name]
	if r, ok := res[p]; ok {
		return r
	}

	buf := bytes.NewBuffer(nil)
	var compile func(string, interface{}, string)

	compile = func(pre string, item interface{}, post string) {
		buf.WriteString(pre)
		switch x := item.(type) {
		default:
			panic(fmt.Errorf("unexpected type %T", x))
		case ebnf.Alternative:
			for i, item := range x {
				compile(pipe[i > 0], item, "")
			}
		case *ebnf.Group:
			compile("(", x.Body, ")")
		case *ebnf.Name:
			buf.WriteString("(" + ebnf2RE(g, x.String, res) + ")")
		case *ebnf.Option:
			compile("(", x.Body, ")?")
		case *ebnf.Range:
			buf.WriteString(fmt.Sprintf("[%s-%s]", regexp.QuoteMeta(x.Begin.String), regexp.QuoteMeta(x.End.String)))
		case *ebnf.Repetition:
			compile("(", x.Body, ")*")
		case ebnf.Sequence:
			for _, item := range x {
				compile("", item, "")
			}
		case *ebnf.Token:
			if s, ok := isRE(x.String); ok {
				buf.WriteString(s)
			} else {
				buf.WriteString(regexp.QuoteMeta(s))
			}
		}
		buf.WriteString(post)
	}

	compile("", p.Expr, "")
	re = buf.String()
	res[p] = re
	return
}


// Scanner returns a new Scanner which can run the Lexer FSM. A Scanner is not safe for concurent access
// but many Scanners can safely share the same Lexer.
//
// The RuneReader can be nil. Then an EOFReader is supplied and
// the real RuneReader(s) can be Included anytime afterwards.
func (lx *Lexer) Scanner(fname string, r io.RuneReader) *Scanner {
	return newScanner(lx, NewScannerSource(fname, r))
}
