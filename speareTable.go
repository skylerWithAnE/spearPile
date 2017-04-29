package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
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

// Actual working function.
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

// BuildTable Build That Table!! MAKE SHAKESPEARE PARSE AGAIN!
func BuildTable() map[string]map[string][]string {
	m := make(map[string]map[string][]string)
	for ntk, pa := range NonTerminals {
		for _, production := range pa {
			splitProd := strings.Split(production, " ")
			m[ntk] = make(map[string][]string)
			for _, symbol := range FindFirst(splitProd, FollowMap[ntk]) {
				_, present := m[ntk][symbol]
				if present {
					fmt.Println("CONFLICT at ParseTable[", ntk, "][", symbol, "]")
				} else {
					m[ntk][symbol] = pa
				}

			}
			if AllNullable(splitProd) {
				for _, symbol := range FollowMap[ntk] {
					_, present := m[ntk][symbol]
					if present {
						fmt.Println("CONFLICT2 at ParseTable[", ntk, "][", symbol, "]")
					} else {
						m[ntk][symbol] = splitProd
					}
				}
			}
		}

	}
	return m
}

// WriteTableToFile create a csv of the parse table.
func WriteTableToFile() {
	f, err := os.Create("tableFile3.csv")
	Check(err)
	defer f.Close()
	w := csv.NewWriter(f)
	var tsl []string
	tsl = append(tsl, " ")
	tsl, _ = UnionSlices(tsl, TerminalSymbolList)
	w.Write(tsl)

	for ntk := range ParseTable {
		var writeSlice []string
		writeSlice = append(writeSlice, ntk)
		for _, tk := range TerminalSymbolList {
			for _, p := range ParseTable[ntk][tk] {
				writeSlice = append(writeSlice, p)
			}
		}
		for i := len(writeSlice) - 1; i < len(TerminalSymbolList)-len(writeSlice)-1; i++ {
			writeSlice = append(writeSlice, " ")
		}
		w.Write(writeSlice)
	}
}

func TableLookUp() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Row: ")
		row, err := r.ReadString('\n')
		Check(err)
		if row == "$" {
			break
		}
		fmt.Print("Col: ")
		col, err := r.ReadString('\n')
		Check(err)

		p, present := ParseTable[row][col]
		if present {
			fmt.Println(p)
		}
	}
}
