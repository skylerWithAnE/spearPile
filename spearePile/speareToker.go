package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var t map[string]*regexp.Regexp
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

func readTerminals(fname string) map[string]*regexp.Regexp {
	f, err := os.Open(fname)
	check(err)
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

			fmt.Println(token, production)

			rex, err := regexp.Compile(production)
			if err == nil {
				m[token] = rex
			} else {
				fmt.Println("Internal error! Production failed to compile to regex!")
			}
		}
	}

	return m
}

func printMap(m map[string]*regexp.Regexp) {
	for key, value := range m {
		fmt.Println("Token:\t\t", key, "\nProduction:\t", value, "\n")
	}
}

func tokenize(s string) (bool, string, token, int) {
	// fmt.Println(s)
	var rstr string
	var newtoken token
	linebreaks := 0
	pass := false
	for key, value := range t {
		result := value.FindStringIndex(s)
		if result != nil && result[0] == 0 {
			lexstr := s[result[0]:result[1]]
			rstr = s[result[1]:]
			newtoken.sym = key
			newtoken.lex = lexstr
			pass = true
			linebreaks = strings.Count(s[:result[1]], "\n")
			break
		}
	}
	return pass, rstr, newtoken, linebreaks
}

func main() {
	t = readTerminals("terminals.txt")
	// printMap(t)

	// input, err := os.Open("inputs/11-1.txt")
	data, err := ioutil.ReadFile("inputs/11-1.txt")
	check(err)
	input := string(data)
	// defer input.Close()
	// ir := bufio.NewReader(input)
	linenum := 1
	var tokens []token
	for {
		pass, leftover, newtoken, linebreaks := tokenize(input)
		if pass == false {
			fmt.Println("Error, could not tokenize input! At line", linenum) //, "\ninput:", input)
			return
		}
		newtoken.linenum = linenum
		fmt.Println("new token: SYM:", newtoken.sym, "LEX:", newtoken.lex, "LINE:", newtoken.linenum)
		tokens = append(tokens, newtoken)
		input = leftover
		linenum += linebreaks
	}
}
