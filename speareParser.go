package main

import "fmt"
import "strings"

const (
	initial   = iota
	scanned   = iota
	predicted = iota
	completed = iota
)

const lambda = "λ"

type item struct {
	col, row, dpos int
	how            [2]int
	lhs            string
	rhs            []string
}

func newItem(h0, h1, r, c, d int, lhs string, rhs []string) *item {
	i := new(item)
	i.how[0] = h0
	i.how[1] = h1
	i.row = r
	i.col = c
	i.lhs = lhs
	i.rhs = rhs
	i.dpos = d
	return i
}

func isComplete(itm *item) bool {
	return itm.how[0] == completed
}

// Globals

// var input string
var nt map[string][]string
var terminals []string
var s []*item
var inputLength int

// End Globals

// MakeTree Shut up linter.
func MakeTree(terms []string, tokens []token) {
	terminals = terms
	inputLength = len(tokens)

	// make list of terminals using the tokenized input.
	for i := 0; i < len(tokens); i++ {
		sym := tokens[i].sym
		if !StringInSlice(sym, terminals) {
			terminals = append(terminals, sym)
		}
	}

	// get a map of nonterminals
	nt = ReadNonTerminals("nonterminals.txt")
	s = append(s, newItem(initial, 0, 0, 0, 0, "S'", nt["S"]))
	earleyParse()

}

func earleyParse() {
	for i := 0; i < inputLength; i++ {
		if i >= len(s) {
			fmt.Println("Earley crash.")
			return
		}
		prods := s[i].rhs
		for j := 0; j < len(prods); j++ {
			if !isComplete(s[i]) {
				prod := strings.Split(prods[j], " ")[0]
				fmt.Println("next production:", prod)
				// check if next token is terminal
				if StringInSlice(prod, terminals) {
					fmt.Println("next production is terminal!")
					// Do scanner
				} else { // nonterminal

					nextNonTerm, valid := nt[prod]
					fmt.Println("nextNonTerm:", nextNonTerm)
					if valid {
						fmt.Println("next production is nonterminal:", nextNonTerm)
						predict(s[i], i)
						// Do predictor
					}
				}
			} else {
				// Do complete
			}
		}
	}
}

func predict(itm *item, state int) {
	nextSym := strings.Split(itm.rhs[itm.dpos], " ")[0]
	nextProd, valid := nt[nextSym]
	if valid {
		fmt.Println("next production..")
		for i := 0; i < len(nextProd); i++ {
			fmt.Println(nextProd[i])
			var newRHS []string
			newRHS = append(newRHS, nextSym)
			nextSyms := strings.Split(nextProd[i], " ")
			for j := 0; j < len(nextSyms); j++ {
				newRHS = append(newRHS, nextSyms[j])
			}
			newDpos := 0
			s = append(s, newItem(predicted, state, 0, 0, newDpos, itm.lhs, newRHS))
		}
	} else {
		fmt.Println("Internal error.. could not find nonterminal[", nextSym, "]")
	}

}

func scan() {

}

func complete() {

}
