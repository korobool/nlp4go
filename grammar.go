package gonlp

// A context-free grammar.

type ContextFreeFrammar struct {
	start       NonTerminal
	productions []Production
}

func NewContextFreeFrammar() *ContextFreeFrammar {
	var cfg ContextFreeFrammar
	return &cfg
}

func (cfg *ContextFreeFrammar) FromString(input string, encoding string) {
	cfg.start, cfg.productions = readGrammar(input, StandardNonTermParser, encoding)
}

type Production struct{}

type NonTerminal struct{}

type NonTermParser struct{}

func (ntp *NonTermParser) StandardNonTermParser(str string, pos string) {}

// Reading Phrase Structure Grammars
func readGrammar(input string, ntp NonTermParser, probabilistic bool, encoding string) (NonTerminal, []Production) {
	start := NonTerminal{}
	production := []Production{}
	return start, production
}
