package syntax

import (
	G "../grammar"
	"errors"
	"fmt"
)

// SyntaxTree item
type STItem interface {
	String() string
}

type SyntaxTree struct {
	label *G.GrammarSymbol
	items []STItem
}

func NewSyntaxTree(label *G.GrammarSymbol, items []STItem) *SyntaxTree {
	return &SyntaxTree{label, items}
}

func (t *SyntaxTree) Label() string {
	return t.label.Val
}

func (t *SyntaxTree) String() string {
	s_items := ""
	for _, item := range t.items {
		s_items = s_items + " " + item.String()
	}
	return fmt.Sprintf("(%s%s)", t.label.String(), s_items)
}

type ShiftReduceParser struct {
	grammar *G.ContextFreeGrammar
}

func (srp *ShiftReduceParser) matchRhs(rhs []*G.GrammarSymbol, rightmost_stack []STItem) bool {
	if len(rhs) != len(rightmost_stack) {
		return false
	}

	for i, sv := range rightmost_stack {
		switch val := sv.(type) {
		case *SyntaxTree:
			if !rhs[i].NonTerminal {
				return false
			}
			if val.Label() != rhs[i].Val {
				return false
			}
		case *G.GrammarSymbol:
			if rhs[i].NonTerminal {
				return false
			}
			if val.Val != rhs[i].Val {
				return false
			}
		}
	}
	return true
}

func (srp *ShiftReduceParser) reduce(stack []STItem) ([]STItem, *G.Production) {
	productions := srp.grammar.Productions(nil, nil, false)

	// Try each production, in order
	for _, production := range productions {
		rhslen := len(production.Rhs)
		if rhslen >= len(stack) {
			rhslen = len(stack)
		}

		// check if the RHS of a production matches the top of the stack
		if srp.matchRhs(production.Rhs, stack[len(stack)-rhslen:]) {
			// combine the tree to reflect the reduction
			tail := make([]STItem, rhslen)
			copy(tail, stack[len(stack)-rhslen:])
			tree := NewSyntaxTree(production.Lhs, tail)
			stack := append(stack[:len(stack)-rhslen], tree)
			return stack, &production
		}
	}

	return stack, nil
}

func trace_stack(stack []STItem, remaining_text []*G.GrammarSymbol) {
	ret := "  [ "
	for _, item := range stack {
		switch val := item.(type) {
		case *SyntaxTree:
			ret = ret + val.Label() + " "
		case *G.GrammarSymbol:
			ret = ret + val.Val + " "
		}
	}
	ret = ret + "* "
	for _, s := range remaining_text {
		ret = ret + s.Val + " "
	}
	ret = ret + "]"
	fmt.Println(ret)
}

func (srp *ShiftReduceParser) Parse(tokens []*G.GrammarSymbol, trace bool) (*SyntaxTree, error) {
	ok, first_invalid_token := srp.grammar.CheckCoverage(tokens)
	if !ok {
		return nil, errors.New(fmt.Sprintf(
			"Grammar does not cover some of the input words: %s",
			first_invalid_token))
	}

	// initialize the stack
	stack := []STItem{}
	var remaining_text []*G.GrammarSymbol
	remaining_text = tokens

	// iterate through the text, pushing the token onto
	// the stack, then reducing the stack
	for len(remaining_text) > 0 {
		stack = append(stack, remaining_text[0])
		remaining_text = remaining_text[1:]
		if trace {
			fmt.Println(fmt.Sprintf("Shift %s:", stack[len(stack)-1]))
			trace_stack(stack, remaining_text)
		}
		for {
			var production *G.Production
			stack, production = srp.reduce(stack)
			if trace && production != nil {
				rhs := ""
				for _, item := range production.Rhs {
					rhs = rhs + item.Val + " "
				}
				fmt.Println(fmt.Sprintf("Reduce %s <- %s:", production.Lhs.Val, rhs))
				trace_stack(stack, remaining_text)
			}
			if production == nil {
				break
			}
		}
	}

	// Did we reduce everything?
	if len(stack) == 1 {
		// Did we end up with the right category?
		if stack[0].(*SyntaxTree).Label() == srp.grammar.Start().Val {
			return stack[0].(*SyntaxTree), nil
		}
	}
	return nil, nil
}
