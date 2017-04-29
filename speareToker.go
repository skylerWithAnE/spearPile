package main

import (
	"bufio"
	"os"
	//"regexp"
)

type Token struct {
	sym, lex string
	linenum  int
}

func Tokenize(s string) (bool, string, Token) {

	var newtoken Token
	rstr := s
	pass := false
	for _, key := range TerminalSymbolList {
		//for key, value := range TerminalRegexMap {	// Bug. Maps aren't ordered.
		value := TerminalRegexMap[key]
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

func WriteTokensToFile(tokens []Token) {
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
