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
// Or statments create new productions for convenience/consistency.
func ReadNonTerminals(fname string) map[string][]string {
	f, err := os.Open(fname)
	Check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	var m map[string][]string
	m = make(map[string][]string)

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
func ReadTerminals(fname string) map[string]*regexp.Regexp {
	f, err := os.Open(fname)
	Check(err)
	defer f.Close()

	r := bufio.NewReader(f)
	m := make(map[string]*regexp.Regexp)

	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		entry := strings.Split(string(l), "->")

		if len(entry) == 2 {

			token := strings.Trim(entry[0], " ")
			production := strings.Trim(entry[1], " ")

			if verbose {
				fmt.Println(token, production)
			}

			rex, err := regexp.Compile("(?i)" + production)
			if err == nil {
				m[token] = rex
			} else {
				fmt.Println("Internal error! Production failed to compile to regex!")
			}
		}
	}

	return m
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
func NullableList(nonterms map[string][]string, terms map[string]*regexp.Regexp) map[string]string {

	nullable := make(map[string]string)
	// using a map because go lacks a set builtin.
	for k := range terms {
		// load up all of the terminals into nullable set.
		nullable[k] = "t"
	}
	for {
		stable := true
		for k, v := range nonterms {
			_, present := nullable[k]
			if !present {
				loopBreak := false
				for i := 0; i < len(v); i++ {
					p := v[i]
					// my nonterms are stored in a map[string][]string
					// each string in the []string slice is a production.
					ps := strings.Split(p, " ")
					// so split that string to get an iterable slice.
					fmt.Println("checking production", k, "->", p)
					for j := 0; j < len(ps); j++ {
						pj := ps[j]
						// pj is the current symbol of the current production
						_, present = nullable[pj]
						if !present {
							// current symbol not in nullable set.
							fmt.Println(pj, "is not in nullable, breaking...")
							loopBreak = true
							break
						}
					}
					if !loopBreak {
						fmt.Println("no loop break for", k)
						_, present = nullable[k]
						if !present {
							stable = false
							nullable[k] = "nt"
							fmt.Println("Adding", k, "to nullable set")
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
