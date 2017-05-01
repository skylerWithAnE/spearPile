package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func addTableEntry(x, y string, production []string) {

}

// FindFirstStrings takes two strings, splits extra, then returns FindFirst
func FindFirstStrings(production, extra string) []string {
	extraA := strings.Split(extra, " ")
	return FindFirstMixed(production, extraA)
}

// FindFirstMixed splits production, returns FindFirstSlices
func FindFirstMixed(production string, extra []string) []string {
	productionA := strings.Split(production, " ")
	return FindFirst(productionA, extra)
}

// FindFirst Actual working function.
func FindFirst(production, extra []string) []string {
	var result []string
	for _, symbol := range production {
		result, _ = UnionSlices(result, FirstMap[symbol])
		if IsNullable(symbol) == false {
			return result
		}
	}
	result, _ = UnionSlices(result, extra)
	return result
}

// AllNullable returns true or false if the production is all nullable.
func AllNullable(production []string) bool {
	for _, symbol := range production {
		if IsNullable(symbol) == false {
			return false
		}
	}
	return true
}

func compareConflict(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if b[i] != s {
			return false
		}
	}
	return true
}

// BuildTable Build That Table!! MAKE SHAKESPEARE PARSE AGAIN!
func BuildTable() /*map[string]map[string][]string*/ {
	f, err := os.Create("parseTableLog")
	Check(err)
	w := bufio.NewWriter(f)

	for _, ntk := range NonTerminalSymbolList {
		//iterate over the ntks in order
		productionArray := NonTerminals[ntk]
		//an array of string productions with lhs N
		for _, production := range productionArray {
			splitProd := strings.Split(production, " ")
			//split productions into a slice
			// ParseTable[ntk] = make(map[string][]string)
			// fmt.Println(pa)
			for _, symbol := range FindFirst(splitProd, FollowMap[ntk]) {
				_, present := ParseTable[ntk][symbol]
				if present {
					// fmt.Println("CONFLICT at ParseTable[", ntk, "][", symbol, "]")
					// if !compareConflict(ParseTable[ntk][symbol], splitProd) {
						w.WriteString("CONFLICT at ParseTable[" + ntk + "][" + symbol + "]=")
						for _, s := range ParseTable[ntk][symbol] {
							w.WriteString(s + " ")
						}
						w.WriteString("..incoming: ")
						for _, s := range splitProd {
							w.WriteString(s + " ")
						}
						w.WriteRune('\n')
					// }
				} else {
					ParseTable[ntk][symbol] = splitProd
				}

			}
			if AllNullable(splitProd) {
				for _, symbol := range FollowMap[ntk] {
					_, present := ParseTable[ntk][symbol]
					if present {
						// if !compareConflict(ParseTable[ntk][symbol], splitProd) {
							w.WriteString("CONFLICT(2) at ParseTable[" + ntk + "][" + symbol + "]=")
							for _, s := range ParseTable[ntk][symbol] {
								w.WriteString(s + " ")
							}
							w.WriteString("..incoming: ")
							for _, s := range splitProd {
								w.WriteString(s + " ")
							}
							w.WriteRune('\n')
						// }
					} else {
						ParseTable[ntk][symbol] = splitProd
					}
				}
			}
		}

	}
	w.Flush()
	f.Close()
	// return m
}

// WriteTableToFile create a csv of the parse table.
func WriteTableToFile() {
	fileTime := strings.Split(time.Now().String(), " ")[1]
	fileTime = strings.Split(fileTime, ".")[0]
	hms := strings.Split(fileTime, ":")
	fileTime = hms[0] + hms[1] + hms[2]
	f, err := os.Create("tables/tableFile" + fileTime + ".csv")
	if err != nil {
		fmt.Println("couldn't write table for some reason..")
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	var tsl []string
	tsl = append(tsl, " ")
	tsl, _ = UnionSlices(tsl, TerminalSymbolList)
	tsl = append(tsl, "$")
	w.Write(tsl)
	// fmt.Println(tsl)

	for _, ntk := range NonTerminalSymbolList {
		var writeSlice []string
		writeSlice = append(writeSlice, ntk)
		for _, tk := range append(TerminalSymbolList, "$") {
			writeString := " "
			for _, p := range ParseTable[ntk][tk] {
				//writeSlice = append(writeSlice, p)
				if p == lambda {
					p = "lambda"
				}
				writeString = writeString + "  " + p
				// fmt.Println("appending ", p)
			}
			writeSlice = append(writeSlice, writeString)
		}
		w.Write(writeSlice)
		// fmt.Println("wrote ", writeSlice)
		w.Flush()
	}

}

func TableLookUp(stackSymbol, tokenSymbol string) int {
	//Option 1: s is terminal, and s==t
	if !IsNonTerminal(stackSymbol) && stackSymbol == tokenSymbol {
		return 0
	}
	//Option 2: s is terminal and s!=t
	if !IsNonTerminal(stackSymbol) && stackSymbol != tokenSymbol {
		return 1
	}
	//Option 3: s is non terminal, but table[s][t] is empty.
	if IsNonTerminal(stackSymbol) && len(ParseTable[stackSymbol][tokenSymbol]) == 0 {
		return 2
	}
	//Option 4: s is non terminal and table[s][t] is nonempty
	if IsNonTerminal(stackSymbol) && len(ParseTable[stackSymbol][tokenSymbol]) > 0 {
		return 3
	}

	//Option 5: something terrible and wacky has occurred.
	return 4
}

func KBTableLookUp() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Row: ")
		row, err := r.ReadString('\n')
		Check(err)
		row = strings.TrimRight(row, "\r\n")
		if row == "$" {
			break
		}
		fmt.Print("Col: ")
		col, err := r.ReadString('\n')
		Check(err)
		col = strings.TrimRight(col, "\r\n")
		p, present := ParseTable[row][col]
		// fmt.Println("row:", row, "col:", col, ParseTable[row][col])
		if present {
			fmt.Println(p)
		}
	}
}

func PrintTable() {
	for ntk := range ParseTable {
		for _, tk := range TerminalSymbolList {
			fmt.Println("row", ntk, "col", tk, ParseTable[ntk][tk])
		}
	}
}
