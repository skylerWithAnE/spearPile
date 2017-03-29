package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var t map[string][]string
var nt map[string][]string
var currState []string
var terminalCandidates []*regexp.Regexp

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type token struct {
	sym, lex string
	linenum  int
}

func readTerminals(fname string) (map[string][]string, []*regexp.Regexp) {
	f, err := os.Open(fname)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	var m map[string][]string
	var re []*regexp.Regexp
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
			m[token] = append(m[token], production)
			rex, err := regexp.Compile(production)
			if err == nil {
				re = append(re, rex)
			}
		}
	}

	return m, re
}

func readNonTerminals(fname string) map[string][]string {
	f, err := os.Open(fname)
	check(err)
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

func printMap(m map[string][]string) {
	for key, value := range m {
		for _, production := range value {
			fmt.Println("Token:\t\t", key, "\nProduction:\t", production)
		}
	}
}

func expandToken(s string) (string, []string) {
	nontermvalue, nonterminal := nt[s]
	productions, terminal := t[s]
	var retval string
	var expandedToken []string
	if nonterminal {
		next := strings.Split(nontermvalue[0], " ")[0]
		retval, expandedToken = expandToken(next)
	} else if terminal {
		for i := range productions {
			retval = productions[i]
			expandedToken = strings.Split(s, " ")
		}
	}
	return retval, expandedToken
}

func currentTokenCandidates() {
	ntval, nonterm := nt[currState[0]]
	terminalCandidates = terminalCandidates[:0]
	if nonterm {
		for i := range ntval {
			tval, term := t[ntval[i]]
			if term {
				regex, err := regexp.Compile(tval[0])
				check(err)
				terminalCandidates = append(terminalCandidates, regex)
				fmt.Println("Regex: ", regex)
			}
		}

	}
}

func consumeTerminal(s string, terminal string) ([]string, error) {
	var retval []string
	_, exists := t[terminal]
	if !exists {
		return nil, errors.New("Internal error! Called consumeTerminal, provided terminal: \"" + terminal + "\" not found in map t.")
	}
	fmt.Println(s, terminal)
	nontermvalue, nonterminal := nt[s]
	if nonterminal {
		for i := range nontermvalue {
			regex, err := regexp.Compile(terminal)
			check(err)
			result := regex.FindStringIndex(nontermvalue[i])
			if len(result) > 0 {
				retval = append(retval, strings.Trim(nontermvalue[i][result[1]:], " "))
			} else {
				fmt.Println("This thing ", terminal, " doesn't match this thing: ", nontermvalue[i])
			}

		}
	} else {
		return nil, errors.New("Internal error! Called consumeTerminal(" + s + " ," + terminal + ") where s is a terminal!")
	}
	return retval, nil
}

func main() {
	var regexList []*regexp.Regexp
	t, regexList = readTerminals("terminals.txt")
	nt = readNonTerminals("nonterminals.txt")
	printMap(t)
	printMap(nt)

	fmt.Println("Start state: ", nt["S"])

	var tokens []token

	input, err := os.Open("inputs/11-1.txt")
	check(err)
	defer input.Close()
	ir := bufio.NewReader(input)

	var i = 0
	currState = strings.Split(nt["S"][0], " ")
	for {
		line, _, readErr := ir.ReadLine()
		if readErr == io.EOF {
			break
		}
		var lstr = string(line)

		// for {
		// 	if len(currState) > 0 {
		// 		regexstr, terminalKey := expandToken(currState[0])
		// 		regex, err := regexp.Compile(regexstr)
		// 		currentTokenCandidates()
		// 		check(err)
		// 		result := regex.FindStringIndex(lstr)
		// 		if result != nil {
		// 			fmt.Println("Start index: ", result[0], "\tEnd index: ", result[1])
		// 			fmt.Println(lstr)
		// 			lexstr := lstr[result[0]:result[1]]
		// 			lstr = lstr[result[1]:]
		// 			fmt.Println(lstr)
		// 			newState, err := consumeTerminal(currState[0], terminalKey[0])
		// 			check(err)
		// 			if len(newState) == 1 {
		// 				if newState[0] == "" {
		// 					currState = currState[1:len(currState)]
		// 				} else {
		// 					currState[0] = newState[0]
		// 				}

		// 			} else {
		// 				fmt.Println("Multiple outcomes...")
		// 			}
		// 			fmt.Println(currState)

		// 			var newtoken token
		// 			newtoken.sym = terminalKey[0]
		// 			newtoken.lex = lexstr
		// 			newtoken.linenum = i
		// 			tokens = append(tokens, newtoken)
		// 		}

		// 	}
		// 	if len(lstr) <= 0 {
		// 		break
		// 	}
		// }
	}
}
