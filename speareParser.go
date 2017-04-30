package main

import "fmt"

type stack struct {
	items []string
}

func peekStack(s stack) string {
	if len(s.items) >= 1 {
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
	newItems = append(newItems, i)
	for _, s := range s.items {
		newItems = append(newItems, s)
	}
	s.items = newItems
	return s
}

func ParseInput(input []Token) bool {
	var parseStack stack
	parseStack = pushStack(parseStack, StartState)
	var tokens []string
	for _, t := range input {
		tokens = append(tokens, t.sym)
	}
	for {
		currToken := tokens[0]
		lookupCode := TableLookUp(peekStack(parseStack), currToken)
		fmt.Println("currSymbol:", peekStack(parseStack), "currToken:", currToken, "got lookupCode", lookupCode)
		if lookupCode == 0 {
			//pop stack, consume token, keep going
			parseStack = popStack(parseStack)
			tokens = tokens[1:]
			if len(tokens) == 0 {
				//All tokens consumed!
				break
			}
		} else if lookupCode == 1 {
			//input not in language.
			return false
		} else if lookupCode == 2 {
			fmt.Println(ParseTable[peekStack(parseStack)][currToken])
			return false
		} else if lookupCode == 3 {
			//Pop s from stack.
			popSym := peekStack(parseStack)
			parseStack = popStack(parseStack)
			//For each x in table[s][t]: push x (push in reverse)
			tableEntry := ParseTable[popSym][currToken]
			for i := len(tableEntry) - 1; i >= 0; i-- {
				parseStack = pushStack(parseStack, tableEntry[i])
				fmt.Print("pushing ", tableEntry[i], "... ")
			}
			fmt.Println()

		}
		fmt.Println("end loop " + peekStack(parseStack))
	}

	return true
}
