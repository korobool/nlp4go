package grammar

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	standardNontermRe = regexp.MustCompile("([\\w/][\\w/^<>-]*)\\s*")
	arrowRe           = regexp.MustCompile("\\s*->\\s*")
	probabilityRe     = regexp.MustCompile("(\\[[\\d\\.]+\\])\\s*")
	terminalRe        = regexp.MustCompile("(\"[^\"]+\"|'[^']+')\\s*")
	disjunctionRe     = regexp.MustCompile("\\|\\s*")
)

type GrammarSymbol struct {
	Val         string
	NonTerminal bool
}

func (gs *GrammarSymbol) String() string {
	return gs.Val
}

func NonTerm(s string) *GrammarSymbol {
	return &GrammarSymbol{s, true}
}

func Term(s string) *GrammarSymbol {
	return &GrammarSymbol{s, false}
}

type Production struct {
	Lhs         *GrammarSymbol
	Rhs         []*GrammarSymbol
	Probability float64
}

type ProductionsList []Production

func (p Production) String() string {
	return fmt.Sprintf("%s: %s [%.2f]", p.Lhs, p.Rhs, p.Probability)
}

func (p Production) IsLexical() bool {
	for _, s := range p.Rhs {
		if !s.NonTerminal {
			return true
		}
	}
	return false
}

func nontermParser(str string, pos int) (*GrammarSymbol, int, error) {
	loc := standardNontermRe.FindStringIndex(str[pos:])
	if loc == nil {
		return nil, 0, errors.New(fmt.Sprintf("Expected a nonterminal, found: %s", str[pos:]))
	}
	return NonTerm(strings.TrimSpace(str[pos+loc[0] : pos+loc[1]])), pos + loc[1], nil
}

func readProductions(line string, probabilistic bool) (ProductionsList, error) {
	pos := 0
	// parse the left-hand side.
	lhs, pos, err := nontermParser(line, pos)
	if err != nil {
		return nil, err
	}

	// Skip over the arrow.
	loc := arrowRe.FindStringIndex(line[pos:])
	if loc == nil {
		return nil, errors.New("Expected an arrow")
	}
	pos = pos + loc[1]

	// Parse the right hand side.
	probabilities := []float64{0.0}
	rhsides := [][]*GrammarSymbol{[]*GrammarSymbol{}}
	for pos < len(line) {
		if probabilistic && line[pos] == '[' {
			loc = probabilityRe.FindStringIndex(line[pos:])
			// Probability
			var p float64
			s := strings.TrimSpace(line[pos+loc[0] : pos+loc[1]])
			p, err = strconv.ParseFloat(s[1:len(s)-1], 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid float value: %s", s))
			}
			if p > 1.0 {
				return nil, errors.New(fmt.Sprintf("Production probability %f, should not be greater than 1.0", p))
			}
			pos = pos + loc[1]
			probabilities[len(probabilities)-1] = p
		} else if line[pos] == '\'' || line[pos] == '"' {
			// String -- add terminal.
			loc = terminalRe.FindStringIndex(line[pos:])
			if loc == nil {
				return nil, errors.New("Unterminated string")
			}
			s := strings.TrimSpace(line[pos+loc[0] : pos+loc[1]])
			rhsides[len(rhsides)-1] = append(rhsides[len(rhsides)-1], Term(s[1:len(s)-1]))
			pos = pos + loc[1]
		} else if line[pos] == '|' {
			// Vertical bar -- start new rhside.
			loc = disjunctionRe.FindStringIndex(line[pos:])
			probabilities = append(probabilities, 0.0)
			rhsides = append(rhsides, []*GrammarSymbol{})
			pos = pos + loc[1]
		} else {
			// Anything else -- nonterminal.
			var nonterm *GrammarSymbol
			nonterm, pos, err = nontermParser(line, pos)
			if err != nil {
				return nil, err
			}
			rhsides[len(rhsides)-1] = append(rhsides[len(rhsides)-1], nonterm)
		}
	}

	productions := ProductionsList{}
	for i, rhs := range rhsides {
		var prob float64
		if len(probabilities) <= i {
			prob = 0.0
		} else {
			prob = probabilities[i]
		}
		productions = append(productions, Production{Lhs: lhs, Rhs: rhs, Probability: prob})
	}
	return productions, nil
}

func readGrammar(input string, probabilistic bool) (start *GrammarSymbol, productions ProductionsList, err error) {
	lines := strings.Split(input, "\n")
	continue_line := ""
	for linenum, line := range lines {
		line = continue_line + strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		if strings.HasSuffix(line, "\\") {
			continue_line = line[:len(line)-1] + " "
			continue
		}
		continue_line = ""
		if line[0] == '%' {
			// this is directive, must be START
			if !strings.HasPrefix(line, "%%start") {
				err = errors.New(fmt.Sprintf("Bad directive at line %d", linenum))
				break
			}
			start, _, err = nontermParser(line, 7)
			if err != nil {
				break
			}
		} else {
			// parse production
			var r_productions ProductionsList
			r_productions, err = readProductions(line, probabilistic)
			if err != nil {
				break
			}
			productions = append(productions, r_productions...)
		}
	}

	if err == nil && len(productions) == 0 {
		err = errors.New("No productions found")
	}
	if err != nil {
		return nil, nil, err
	}

	if start == nil {
		start = productions[0].Lhs
	}
	return start, productions, nil
}

type ContextFreeGrammar struct {
	start         *GrammarSymbol
	productions   ProductionsList
	categories    []*GrammarSymbol
	lhs_index     map[GrammarSymbol]ProductionsList
	rhs_index     map[GrammarSymbol]ProductionsList
	empty_index   map[GrammarSymbol]Production
	lexical_index map[GrammarSymbol]ProductionsList

	is_lexical            bool
	is_nonlexical         bool
	all_unary_are_lexical bool
	min_len               int
	max_len               int
}

func GetContextFreeGrammar(start *GrammarSymbol, productions ProductionsList) *ContextFreeGrammar {
	categories := []*GrammarSymbol{}
	lhs_index := make(map[GrammarSymbol]ProductionsList)
	rhs_index := make(map[GrammarSymbol]ProductionsList)
	empty_index := make(map[GrammarSymbol]Production)
	lexical_index := make(map[GrammarSymbol]ProductionsList)

	is_lexical := true
	is_nonlexical := true
	all_unary_are_lexical := true
	min_len := 10000
	max_len := 0

	for _, prod := range productions {
		categories = append(categories, prod.Lhs)

		// Left hand side.
		lhs := prod.Lhs
		prod_list, ok := lhs_index[*lhs]
		if ok {
			prod_list = append(prod_list, prod)
		} else {
			prod_list = ProductionsList{prod}
			lhs_index[*lhs] = prod_list
		}

		if len(prod.Rhs) > 0 {
			// First item in right hand side.
			rhs0 := prod.Rhs[0]
			rprod_list, ok := rhs_index[*rhs0]
			if ok {
				rprod_list = append(rprod_list, prod)
			} else {
				rprod_list = ProductionsList{prod}
				rhs_index[*rhs0] = rprod_list
			}
		} else {
			// The right hand side is empty.
			empty_index[*lhs] = prod
		}

		// Lexical tokens in the right hand side.
		for _, token := range prod.Rhs {
			if !token.NonTerminal {
				lprod_list, ok := lexical_index[*token]
				if ok {
					lprod_list = append(lprod_list, prod)
				} else {
					lprod_list = ProductionsList{prod}
					lexical_index[*token] = lprod_list
				}
			}
		}

		// check binary
		if !prod.IsLexical() {
			is_lexical = false
			if len(prod.Rhs) == 1 {
				all_unary_are_lexical = false
			}
		}
		if len(prod.Rhs) != 1 && prod.IsLexical() {
			is_nonlexical = false
		}
		if len(prod.Rhs) < min_len {
			min_len = len(prod.Rhs)
		}
		if len(prod.Rhs) > max_len {
			max_len = len(prod.Rhs)
		}
	}

	return &ContextFreeGrammar{
		start, productions, categories, lhs_index, rhs_index, empty_index, lexical_index,
		is_lexical, is_nonlexical, all_unary_are_lexical, min_len, max_len}
}

func GetContextFreeGrammarFromString(input string, probabilistic bool) (*ContextFreeGrammar, error) {
	start, productions, err := readGrammar(input, probabilistic)
	if err != nil {
		return nil, err
	}
	return GetContextFreeGrammar(start, productions), nil
}

func (cfg *ContextFreeGrammar) Start() *GrammarSymbol {
	return cfg.start
}

func (cfg *ContextFreeGrammar) CheckCoverage(tokens []*GrammarSymbol) (bool, *GrammarSymbol) {
	for _, token := range tokens {
		_, ok := cfg.lexical_index[*token]
		if !ok {
			return false, token
		}
	}
	return true, nil
}

func (cfg *ContextFreeGrammar) Productions(lhs *GrammarSymbol, rhs *GrammarSymbol, empty bool) ProductionsList {
	ret_val := ProductionsList{}

	// no constraints so return everything
	if lhs == nil && rhs == nil {
		if !empty {
			return cfg.productions
		} else {
			for _, val := range cfg.empty_index {
				ret_val = append(ret_val, val)
			}
			return ret_val
		}
	} else if lhs != nil && rhs == nil {
		if !empty {
			val, ok := cfg.lhs_index[*lhs]
			if !ok {
				val = ret_val
			}
			return val
		} else {
			val, ok := cfg.empty_index[*lhs]
			if ok {
				ret_val = append(ret_val, val)
			}
			return ret_val
		}
	} else if rhs != nil && lhs == nil {
		val, ok := cfg.rhs_index[*rhs]
		if ok {
			ret_val = append(ret_val, val...)
		}
		return ret_val

	} else {
		prods, ok := cfg.lhs_index[*lhs]
		if !ok {
			prods = ProductionsList{}
		}
		for _, prod := range prods {
			_, ok = cfg.rhs_index[*rhs]
			if ok {
				ret_val = append(ret_val, prod)
			}
		}
		return ret_val
	}
}
