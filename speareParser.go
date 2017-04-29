package main

type stack struct {
	items []string
}

func peekStack(s stack) string {
	if len(s.items) > 1 {
		return s.items[1]
	}
	return "none"
}

func popStack(s stack) string {
	if len(s.items) == 0 {
		return "empty_stack"
	}
	ret := s.items[0]
	s.items = s.items[1:]
	return ret
}

func pushStack(s stack, i string) {
	var newItems []string
	newItems = append(newItems, i)
	for _, s := range s.items {
		newItems = append(newItems, s)
	}
	s.items = newItems
}

func ParseInput(input []Token) {
	var parseStack stack
	pushStack(parseStack, StartState)
}
