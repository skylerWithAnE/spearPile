package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func grammarOut() {
	f, err := os.Create("grammarOut.txt")
	Check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	for ntk, ps := range NonTerminals {
		outNtk := strings.ToUpper(ntk)
		w.WriteString(outNtk + " -> ")
		for i, prod := range ps {
			if i > 0 {
				w.WriteString("\n| ")
			}
			splitProd := strings.Split(prod, " ")
			for _, s := range splitProd {
				if IsNonTerminal(s) {
					w.WriteString(strings.ToUpper(s))
				} else {
					w.WriteString(strings.ToLower(s))
				}
				w.WriteString(" ")
			}
		}
		w.WriteString(".\n")
	}
	for tk, re := range TerminalRegexMap {
		outTk := strings.ToLower(tk)
		w.WriteString(outTk + " -> ")
		w.WriteString(re.String() + ".\n")
	}
}

func main() {
	Verbose = false
	TerminalRegexMap, TerminalSymbolList = ReadTerminals("_terminals.txt")
	if Verbose {
		PrintRegexMap(TerminalRegexMap)
	}
	var fname string
	if len(os.Args) > 1 {
		fname = os.Args[1]
	} else {
		fname = "sdt/t20.txt"
	}
	// data, err := ioutil.ReadFile("tokenFile49")
	// Check(err)
	// input := string(data)
	// inputA := strings.Split(input, " ")
	// var tokens []Token
	// fmt.Println(len(inputA))
	// for _, t := range inputA {
	// 	var newToken Token
	// 	newToken.sym = t
	// 	tokens = append(tokens, newToken)
	// }
	data, err := ioutil.ReadFile(fname)
	Check(err)
	input := string(data) + "$"
	linenum := 1
	var tokens []Token
	for {
		pass, leftover, newtoken := Tokenize(input)
		if !pass {
			if len(input) > 1 {
				if input[0:1] == "\n" {
					linenum++
				}
				//BUG(?)
				input = input[1:]
			} else {
				if input == "$" {
					break
				}
				fmt.Println("Failed to tokenize!")
				return
			}

		} else {
			newtoken.linenum = linenum
			if Verbose {
				fmt.Println("new token:\n\t\tSYM:", newtoken.sym, "\n\t\tLEX:", newtoken.lex, "\n\t\tLINE:", newtoken.linenum)
			}
			tokens = append(tokens, newtoken)
			input = leftover
		}
	}
	WriteTokensToFile(tokens)
	NonTerminals = ReadNonTerminals("_nonterminals.txt")
	Nullables = NullableList()
	// PrintNullableMap()
	FirstMap = BuildFirstMap()
	// PrintFirstMap()
	FollowMap = BuildFollowMap()
	// PrintFollowMap()
	//	grammarOut()
	ParseTable = make(map[string]map[string][]string)
	for _, ntk := range NonTerminalSymbolList {
		ParseTable[ntk] = make(map[string][]string)
	}
	BuildTable()
	// WriteTableToFile()
	// PrintTable()
	success := ParseInput(tokens)
	if success {
		fmt.Println("Parsed.")
		os.Exit(0)
	} else {
		fmt.Println("Failed.")
		os.Exit(1)
	}
	// KBTableLookUp()
}
