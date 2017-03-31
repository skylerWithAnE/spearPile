package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var t map[string]*regexp.Regexp

//var nt map[string][]string
var currState []string
var terminalCandidates []*regexp.Regexp
var verbose bool

type token struct {
	sym, lex string
	linenum  int
}

func tokenize(s string) (bool, string, token) {

	var newtoken token
	rstr := s
	pass := false
	for key, value := range t {
		result := value.FindStringIndex(s)
		if result != nil && result[0] == 0 {
			lexstr := s[result[0]:result[1]]
			rstr = s[result[1]:]
			newtoken.sym = key
			newtoken.lex = lexstr
			pass = true
			break
		}
	}
	return pass, rstr, newtoken
}

func writeTokensToFile(tokens []token) {
	f, err := os.Create("tokenFile")
	Check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	currLine := 1
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t.linenum > currLine {
			for ; currLine < t.linenum; currLine++ {
				w.WriteString("\n")
			}
		}
		w.WriteString(t.sym + "\t")
		w.Flush()
	}
}

func main() {
	verbose = false
	t = ReadTerminals("terminals.txt")
	if verbose {
		PrintRegexMap(t)
	}
	data, err := ioutil.ReadFile("inputs/11-1.txt")
	Check(err)
	input := string(data) + "$"
	linenum := 1
	var tokens []token
	for {
		pass, leftover, newtoken := tokenize(input)
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
			if verbose {
				fmt.Println("new token:\n\t\tSYM:", newtoken.sym, "\n\t\tLEX:", newtoken.lex, "\n\t\tLINE:", newtoken.linenum)
			}
			tokens = append(tokens, newtoken)
			input = leftover
		}
	}
	writeTokensToFile(tokens)
}
