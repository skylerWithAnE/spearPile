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
func PrintNullableMap(m map[string]string) {
	for key, value := range m {
		if value == "nt" {
			fmt.Println(key)
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
	if verbose {
		fmt.Println(message)
	}
}
