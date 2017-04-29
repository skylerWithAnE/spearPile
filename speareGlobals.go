package main

import "regexp"

var TerminalRegexMap map[string]*regexp.Regexp
var TerminalSymbolList []string
var NonTerminals map[string][]string
var NonTerminalSymbolList []string
var Nullables map[string]string
var FirstMap map[string][]string
var FollowMap map[string][]string
var ParseTable map[string]map[string][]string

var Verbose bool

const StartState string = "S"
