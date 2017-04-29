package main

import (
	"fmt"
	"regexp"
)

// StringInSlice Go has no generics.
func StringInSlice(s string, l []string) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == s {
			return true
		}
	}
	return false
}

// PrintMap util function to print maps.
func PrintMap(m map[string]string) {
	for key, value := range m {
		fmt.Println("Key:\t", key, "\nValue:\t", value)
	}
}

// PrintNullableMap util to cleanly print the nonterms of the nullable set
func PrintNullableMap() {
	for key, value := range Nullables {
		if value == "nt" {
			fmt.Println(key)
		}
	}
}

// PrintFirstMap util to cleanly print first map.
func PrintFirstMap() {
	printStringSliceMap(FirstMap, "First")
}

// PrintFollowMap util to cleanly print follow map.
func PrintFollowMap() {
	printStringSliceMap(FollowMap, "Follow")
}

func printStringSliceMap(m map[string][]string, sym string) {
	for key, value := range m {
		if IsNonTerminal(key) {
			fmt.Print(sym+"[", key, "] = {")
			for i := range value {

				fmt.Print(value[i])
				if i != len(value)-1 {
					fmt.Print(", ")
				}
			}

			fmt.Print("}\n")
		}
	}
}

// PrintRegexMap Util to print the regex map.
func PrintRegexMap(m map[string]*regexp.Regexp) {
	for key, value := range m {
		fmt.Println("Token:\t\t", key, "\nProduction:\t", value)
	}
}

// Check Lazy error handle.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// MyLog used to print if verbose is true
func MyLog(message string) {
	if Verbose {
		fmt.Println(message)
	}
}

// IsNullable returns true/false for symbol in nullable list.
func IsNullable(candidate string) bool {
	_, present := Nullables[candidate]
	return present
}

// IsNonTerminal returns true/false for symbol in nonterminal map
func IsNonTerminal(candidate string) bool {
	_, present := NonTerminals[candidate]
	return present
}

// UnionSlices treats slices as a set and unions two slices.
func UnionSlices(src []string, dst []string) ([]string, bool) {
	ret := src
	add := false
	for _, v := range dst {
		if StringInSlice(v, ret) == false {
			ret = append(ret, v)
			add = true
		}
	}
	return ret, add
}
