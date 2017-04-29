package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// ReadNonTerminals pass a file name, get a map of nonterminals and their productions.
// Or statements create new productions for convenience/consistency.
func ReadNonTerminals(fname string) map[string][]string {
	f, err := os.Open(fname)
	Check(err)
	defer f.Close()

	r := bufio.NewReader(f)
	m := make(map[string][]string)

	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		entry := strings.Split(string(l), "->")

		if len(entry) == 2 {
			token := strings.Trim(entry[0], " ")
			production := strings.Trim(entry[1], " ")
			multiProduction := strings.Split(production, "|")
			if len(multiProduction) > 0 {
				for i := range multiProduction {
					m[token] = append(m[token], strings.Trim(multiProduction[i], " "))
				}
			} else {
				m[token] = append(m[token], production)
			}
		}
	}

	return m
}

// ReadTerminals Takes a filename, returns a map of terminals and their regexes.
func ReadTerminals(fname string) (map[string]*regexp.Regexp, []string) {
	f, err := os.Open(fname)
	Check(err)
	defer f.Close()

	r := bufio.NewReader(f)
	m := make(map[string]*regexp.Regexp)
	var symbolList []string

	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		entry := strings.Split(string(l), "->")

		if len(entry) == 2 {

			symbol := strings.Trim(entry[0], " ")
			production := strings.Trim(entry[1], " ")

			if Verbose {
				fmt.Println(symbol, production)
			}

			rex, err := regexp.Compile("(?i)" + production)
			if err == nil {
				m[symbol] = rex
				symbolList = append(symbolList, symbol)
			} else {
				fmt.Println("Internal error! Production failed to compile to regex!")
			}
		}
	}

	return m, symbolList
}

// TerminalList Generate a list of all terminals from the regex map.
func TerminalList(termMap map[string]*regexp.Regexp) []string {
	var l []string
	for key := range termMap {
		l = append(l, key)
	}
	return l
}

// NullableList generates list of nullables
func NullableList() map[string]string {

	nullable := make(map[string]string)
	// using a map because go lacks a set builtin.
	nullable[lambda] = lambda
	_, present := nullable[lambda]
	fmt.Println(present)
	for {
		stable := true
		for k, v := range NonTerminals {
			_, present := nullable[k]
			if !present {
				for i := 0; i < len(v); i++ {
					loopBreak := false
					p := v[i]
					// my nonterms are stored in a map[string][]string
					// each string in the []string slice is a production.
					ps := strings.Split(p, " ")
					// so split that string to get an iterable slice.
					// fmt.Println("checking production", k, "->", p)
					for j := 0; j < len(ps); j++ {
						pj := ps[j]
						// pj is the current symbol of the current production
						_, present = nullable[pj]
						// fmt.Println("is", pj, "present?", present)
						if !present {
							// current symbol not in nullable set.
							// fmt.Println(pj, "is not in nullable, breaking...")
							loopBreak = true
							break
						}
					}
					if !loopBreak {
						// fmt.Println("no loop break for", k)
						_, present = nullable[k]
						if !present {
							stable = false
							nullable[k] = "nt"
							// fmt.Println("Adding", k, "to nullable set")
						}

					}
				}
			}

		}
		if stable {
			break
		}
	}
	return nullable
}

// BuildFirstMap build the first map.
func BuildFirstMap() map[string][]string {

	first := make(map[string][]string)
	for k := range TerminalRegexMap {
		var set []string
		set = append(set, k)
		first[k] = set
	}
	stable := true
	for {
		stable = true
		for lhs, v := range NonTerminals {
			for _, production := range v {
				pa := strings.Split(production, " ")
				for _, symbol := range pa {
					targetFirst := first[symbol]
					for _, fSymbol := range targetFirst {
						if !StringInSlice(fSymbol, first[lhs]) {
							first[lhs] = append(first[lhs], fSymbol)
							stable = false
						}
					}
					if !IsNullable(symbol) {
						break
					}
				}
			}
		}
		if stable {
			break
		}
	}

	return first
}

// BuildFollowMap creates the follow map.
func BuildFollowMap() map[string][]string {
	follow := make(map[string][]string)
	follow[StartState] = append(follow[StartState], "$")
	for {
		//do-while follow is not stable
		unstable := false
		for N, v := range NonTerminals {
			//for all NonTerminals
			for _, production := range v {
				// fmt.Println("checking production", production, "in NonTerm", N)
				symbolArray := strings.Split(production, " ")
				for i, symbol := range symbolArray {
					// fmt.Println("N:", N, "symbol:", symbol)
					// for each symbol in the production.
					if IsNonTerminal(symbol) && i < len(symbolArray)-1 {
						// fmt.Println("symbol,", symbol, ", is nonterminal and the nextSymbol is:", nextSymbol)
						// fmt.Println(follow[nextSymbol])
						for j := i + 1; j < len(symbolArray); j++ {
							nextSymbol := symbolArray[j]
							follow[symbol], unstable = UnionSlices(follow[symbol], FirstMap[nextSymbol])
							if !IsNullable(nextSymbol) {
								// fmt.Println("not nullable next symbol..", nextSymbol)
								break
							}

						}
						follow[symbol], unstable = UnionSlices(follow[symbol], follow[N])

					}
				} //end symbol loop
			} //end production loop
		} //end nonterm loop
		if unstable == false {
			break
		}
	} //end infinite loop

	return follow
}
