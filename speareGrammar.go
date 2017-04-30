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
			if StringInSlice(token, NonTerminalSymbolList) == false {
				NonTerminalSymbolList = append(NonTerminalSymbolList, token)
			}
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
				if StringInSlice(symbol, symbolList) == false {
					symbolList = append(symbolList, symbol)
				}
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
	// f, err := os.Create("tmp/followLog")
	// Check(err)
	// w := bufio.NewWriter(f)
	for {
		unstable := false
		// w.WriteString("\n***NEW LOOP***\n")
		//do-while follow is not stable

		//for N, v := range NonTerminals {
		for _, N := range NonTerminalSymbolList {
			unionSuccess := false
			// w.WriteString("\nNew NonTerm loop\n")
			if unstable {
				// w.WriteString("unstable\n")
			}
			v := NonTerminals[N]
			//for all NonTerminals
			for _, production := range v {
				// w.WriteString("checking production \"" + production + "\" in NonTerm " + N + "\n")
				symbolArray := strings.Split(production, " ")
				for i, symbol := range symbolArray {
					loopBreak := false
					// w.WriteString("N: " + N + "\nsymbol: " + symbol + "\n")
					// for each symbol in the production.
					if IsNonTerminal(symbol) {
						// w.WriteString(symbol + " is non terminal, iterate over remaining production\n")
						for j := i + 1; j < len(symbolArray); j++ {
							nextSymbol := symbolArray[j]
							follow[symbol], unionSuccess = UnionSlices(follow[symbol], FirstMap[nextSymbol])
							// w.WriteString("perform union on follow[symbol] U first[nextSymbol]\n")
							if !IsNullable(nextSymbol) {
								// w.WriteString("not nullable next symbol.." + nextSymbol + "\n")
								loopBreak = true
								break
							}

						}
						if !loopBreak {
							// w.WriteString("no loop break occurred, do follow[" + symbol + "]Ufollow[" + N + "]..\n")
							// old := follow[symbol]
							follow[symbol], unionSuccess = UnionSlices(follow[symbol], follow[N])
							if unstable {
								// fmt.Println("union occurred:", "follow[", symbol, "] is now ",
								// follow[symbol], "previously ", old, ".")
							}
						}

					} else {
						// w.WriteString("symbol is terminal... continue\n")
					}
				} //end production loop
			} //end productions with lhs N loop
			if unionSuccess {
				unstable = true
			}
		} //end nonterm loop
		if unstable == false {
			break
		}
		// w.Flush()
	} //end infinite loop

	// for _, ntk := range NonTerminalSymbolList {
	// 	// w.WriteString(ntk + " -> ")
	// 	for _, v := range follow[ntk] {
	// 		// w.WriteString(v + ", ")
	// 	}
	// 	w.WriteString("\n")
	//}
	// w.Flush()
	// f.Close()
	return follow
}
