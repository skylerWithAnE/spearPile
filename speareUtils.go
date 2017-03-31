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
