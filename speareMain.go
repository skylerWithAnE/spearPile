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
	TerminalRegexMap, TerminalSymbolList = ReadTerminals("tmp/toyterminals.txt")
	if Verbose {
		PrintRegexMap(TerminalRegexMap)
	}
	data, err := ioutil.ReadFile("tests/t1.txt")
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
	NonTerminals = ReadNonTerminals("tmp/toynonterminals.txt")
	Nullables = NullableList()
	PrintNullableMap()
	FirstMap = BuildFirstMap()
	PrintFirstMap()
	FollowMap = BuildFollowMap()
	PrintFollowMap()
	//	grammarOut()
	ParseTable = make(map[string]map[string][]string)
	for _, ntk := range NonTerminalSymbolList {
		ParseTable[ntk] = make(map[string][]string)
	}
	BuildTable()
	// TableLookUp()
	f, err := os.Create("parseTableExterior")
	Check(err)

	w := bufio.NewWriter(f)
	for _, ntk := range NonTerminalSymbolList {
		for _, tk := range TerminalSymbolList {
			out := "ParseTable[" + ntk + "][" + tk + "] :"
			for _, v := range ParseTable[ntk][tk] {
				out = out + " " + v
			}
			out = out + "\n"
			w.WriteString(out)
		}
	}
	w.Flush()
	f.Close()
	WriteTableToFile()
	// PrintTable()
	fmt.Println(":(", FollowMap["aof'"])
	ParseInput(tokens)
}
