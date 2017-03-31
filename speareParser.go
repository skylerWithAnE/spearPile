package main

const (
	initial   = iota
	scanned   = iota
	predicted = iota
	completed = iota
)

const lambda = "λ"

type item struct {
	col, row, dpos int
	how            [2]int
	lhs            string
	rhs            []string
}

func newItem(h0, h1, r, c int, lhs string, rhs []string) *item {
	i := new(item)
	i.how[0] = h0
	i.how[1] = h1
	i.row = r
	i.col = c
	i.lhs = lhs
	i.rhs = rhs
	i.dpos = 0
	return i
}

// Globals

var input string
var nt map[string][]string
var pt []string
var s []*item
var inputLength int

// End Globals

// MakeTree Shut up linter.
func MakeTree(newInput string, tokens []token) {
	input = newInput
	inputLength = len(tokens)

	// make list of terminals using the tokenized input.
	for i := 0; i < len(tokens); i++ {
		sym := tokens[i].sym
		if !StringInSlice(sym, pt) {
			pt = append(pt, sym)
		}
	}

	// get a map of nonterminals
	nt = ReadNonTerminals("nonterminals.txt")
	s = append(s, newItem(initial, 0, 0, 0, "S'", nt["S"]))
	earleyParse()

}

func makeTree(itm *item) {
	// start with final production.

	// if itm.how == initial {
	// 	node.sym = "S"
	// 	t = node
	// } else if itm.how == scanned {
	// 	tmp := makeTreeNode(itm)
	// 	tmp.thing = "shutuplinter"
	// 	//node.prependChild(tmp)
	// 	//call makeTree recursively
	// 	t = makeTree(*itm.prev, node)
	// } else if itm.how == predicted {
	// 	//return nothing. (stop using token type, probably.)
	// }

}

func earleyParse() {
	for i := 0; i < inputLength; i++ {
		// itm := s[i]
		// for j := 0; j < len(itm.rhs); j++ {
		// 	prod := itm.rhs[j]
		// 	s = append(s, newItem(prediction, 0, 0, itm.lhs))
		// }
	}
}

func predict() {

}

func scan() {

}

func complete() {

}
