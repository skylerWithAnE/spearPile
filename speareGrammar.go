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

// ReadTerminals Takes a filename, returns a map of terminals and their regexes.
func ReadTerminals(fname string) map[string]*regexp.Regexp {
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
