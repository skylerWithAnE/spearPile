package main

import "fmt"

type stack struct {
	items []string
}

type treeNode struct {
	children []treeNode
	tok      Token
	sym      string
}

func peekStack(s stack) string {
	if len(s.items) >= 1 {
		// if s.items[0] == lambda {
		// 	return "lambda"
		// }
		return s.items[0]
	}
	return "empty_stack"
}

func peekNextStack(s stack) string {
	if len(s.items) > 1 {
		return s.items[1]
	}
	return "none"
}

func popStack(s stack) stack {
	if len(s.items) == 0 {
		return s
	}
	//ret := s.items[0]
	s.items = s.items[1:]
	return s
}

func pushStack(s stack, i string) stack {
	var newItems []string
	// fmt.Print("original: ", s.items)
	newItems = append(newItems, i)
	for _, s := range s.items {
		newItems = append(newItems, s)
	}
	s.items = newItems
	// fmt.Println("new..: ", s.items)
	return s
}

func ParseInput(input []Token) bool {
	var parseStack stack
	parseStack = pushStack(parseStack, StartState)
	var tokens []string
	for _, t := range input {
		tokens = append(tokens, t.sym)
	}
	tokens = append(tokens, "$")
	// f, err := os.Create("parseLog.txt")
	// Check(err)
	// w := bufio.NewWriter(f)

	for {
		currToken := tokens[0]
		lookupCode := TableLookUp(peekStack(parseStack), currToken)

		fmt.Println("currSymbol:", peekStack(parseStack), "currToken:", currToken, "got lookupCode", lookupCode)
		// w.WriteString("currSymbol:" + peekStack(parseStack) + "currToken:" + currToken + "got lookupCode" + string(lookupCode))
		if lookupCode == 0 {
			//pop stack, consume token, keep going
			// fmt.Println("Pop!", peekStack(parseStack))
			parseStack = popStack(parseStack)
			tokens = tokens[1:]

			if len(tokens) == 0 && peekStack(parseStack) == "empty_stack" {
				//All tokens consumed!
				break
			}
		} else if lookupCode == 1 {

			//input not in language.
			if peekStack(parseStack) == lambda {
				parseStack = popStack(parseStack)
			} else if peekStack(parseStack) == "empty_stack" && currToken == "$" {
				return true
			} else {
				return false
			}
		} else if lookupCode == 2 {
			fmt.Println("parse failed.", ParseTable[peekStack(parseStack)][currToken])

			return false
		} else if lookupCode == 3 {
			//Pop s from stack.
			popSym := peekStack(parseStack)
			parseStack = popStack(parseStack)
			// For each x in table[s][t]: push x (push in reverse)
			tableEntry := ParseTable[popSym][currToken]
			for i := len(tableEntry) - 1; i >= 0; i-- {
				parseStack = pushStack(parseStack, tableEntry[i])
				// fmt.Print("pushing ", tableEntry[i], "... ")
			}
			// fmt.Println()

		}
		// fmt.Println("end loop " + peekStack(parseStack))
	}

	// w.Flush()
	// f.Close()
	return true
}
