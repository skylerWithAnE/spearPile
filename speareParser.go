package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type treeNode struct {
	children []treeNode
	tok      Token
	sym, id  string
}

type stack struct {
	items []treeNode
}

var parseStack stack
var treeList []treeNode

func peekStack(s stack) string {
	if len(s.items) >= 1 {
		return s.items[0].sym
	}
	return "empty_stack"
}

func popStack(s stack) (stack, treeNode) {
	n := s.items[0]
	if len(s.items) == 0 {
		return s, n
	}
	s.items = s.items[1:]
	return s, n
}

func pushStack(s stack, i treeNode) stack {
	var newItems []treeNode
	// fmt.Print("original: ", s.items)
	newItems = append(newItems, i)
	for _, s := range s.items {
		newItems = append(newItems, s)
	}
	s.items = newItems
	// fmt.Println("new..: ", s.items)
	return s
}

func ParseInput(input []Token, fname string) bool {
	var ssToken Token
	ssToken.lex = "Start State"
	ssToken.sym = "S"
	ssToken.linenum = 0
	var ssNode treeNode
	ssNode.tok = ssToken
	ssNode.sym = "S"
	ssNode.id = "0"
	parseStack = pushStack(parseStack, ssNode)
	var tokens []string
	for _, t := range input {
		tokens = append(tokens, t.sym)
	}
	tokens = append(tokens, "$")
	nodeID := 1
	tokenIndex := 0
	for {
		currToken := tokens[0]
		lookupCode := TableLookUp(peekStack(parseStack), currToken)

		if Verbose {
			fmt.Println("currSymbol:", peekStack(parseStack), "currToken:", currToken, "got lookupCode", lookupCode)
		}
		// w.WriteString("currSymbol:" + peekStack(parseStack) + "currToken:" + currToken + "got lookupCode" + string(lookupCode))
		if lookupCode == 0 {
			//pop stack, consume token, keep going
			// fmt.Println("Pop!", peekStack(parseStack))
			var newNode treeNode
			parseStack, newNode = popStack(parseStack)
			treeList = append(treeList, newNode)
			// fmt.Print(newNode.sym, ": ")
			// for _, c := range newNode.children {
			// 	fmt.Print(c.sym, ", ")
			// }
			// fmt.Println()
			tokenIndex++
			tokens = tokens[1:]

			if len(tokens) == 0 && peekStack(parseStack) == "empty_stack" {
				//All tokens consumed!
				break
			}
		} else if lookupCode == 1 {

			//input not in language.
			if peekStack(parseStack) == lambda {
				parseStack, _ = popStack(parseStack)
			} else if peekStack(parseStack) == "empty_stack" && currToken == "$" {
				break
			} else {
				os.Mkdir("graphs", 0777)
				rex, err := regexp.Compile("t[\\d]+")
				Check(err)
				ftitle := rex.FindString(fname)

				f, err := os.Create("graphs/graph" + ftitle + ".gv")
				Check(err)
				w := bufio.NewWriter(f)
				w.WriteString("strict graph {\n")
				for _, n := range treeList {
					for _, c := range n.children {
						w.WriteString("\t" + n.sym + n.id + " -- " + c.sym + c.id + "\n")
					}
					w.WriteRune('\n')
				}
				w.WriteString("\n}")
				return false
			}
		} else if lookupCode == 2 {
			fmt.Println("parse failed.", ParseTable[peekStack(parseStack)][currToken])

			return false
		} else if lookupCode == 3 {
			//Pop s from stack.
			popSym := peekStack(parseStack)
			var newParent treeNode
			parseStack, newParent = popStack(parseStack)
			// For each x in table[s][t]: push x (push in reverse)
			tableEntry := ParseTable[popSym][currToken]
			for i := len(tableEntry) - 1; i >= 0; i-- {
				var pushNode treeNode
				pushNode.sym = tableEntry[i]
				if tokenIndex < len(input) {
					pushNode.tok = input[tokenIndex]
				}
				pushNode.id = strconv.Itoa(nodeID)
				nodeID++
				newParent.children = append(newParent.children, pushNode)
				parseStack = pushStack(parseStack, pushNode)
			}
			treeList = append(treeList, newParent)
		}
		// fmt.Println("end loop " + peekStack(parseStack))
	}

	os.Mkdir("graphs", 0777)
	rex, err := regexp.Compile("t[\\d]+")
	Check(err)
	ftitle := rex.FindString(fname)

	f, err := os.Create("graphs/graph" + ftitle + ".gv")
	Check(err)
	w := bufio.NewWriter(f)
	w.WriteString("strict graph {\n")
	for _, n := range treeList {
		for _, c := range n.children {
			w.WriteString("\t" + n.sym + n.id + " -- " + c.sym + c.id + "\n")
		}
		w.WriteRune('\n')
	}
	w.WriteString("\n}")

	w.Flush()
	f.Close()

	return true
}
