package gonlp

import (
	"fmt"
	// "reflect"
)

// Recursive Descent Parsing
// Shift-Reduce Parsing
// The Left-Corner Parser
// Well-Formed Substring Tables

type TaggedToken struct {
	Token string
	Tag   string
}

type SyntaxParser struct {
}

func NewSyntaxParser() *SyntaxParser {
	var sp SyntaxParser
	return &sp
}

func (sp *SyntaxParser) Parse(taggedTokens []TaggedToken) (Tree, error) {
	// "The dog saw a man in the park"

	var treeNP1 Tree
	treeNP1.Type = "NP"
	treeNP1.Leaves = append(treeNP1.Leaves, TaggedToken{"the", "DT"}, TaggedToken{"dog", "NN"})

	var treeNP2 Tree
	treeNP2.Type = "NP"
	treeNP2.Leaves = append(treeNP2.Leaves, TaggedToken{"a", "DT"}, TaggedToken{"man", "NN"})

	var treeNP3 Tree
	treeNP3.Type = "NP"
	treeNP3.Leaves = append(treeNP3.Leaves, TaggedToken{"the", "DT"}, TaggedToken{"park", "NN"})

	var treePP1 Tree
	treePP1.Type = "PP"
	treePP1.Leaves = append(treePP1.Leaves, TaggedToken{"in", "IN"}, treeNP3)

	var treeVP1 Tree
	treeVP1.Type = "VP"
	treeVP1.Leaves = append(treeVP1.Leaves, TaggedToken{"saw", "VDB"}, treeNP2, treePP1)

	var tree Tree
	tree.Type = "S"
	tree.Leaves = append(tree.Leaves, treeNP1, treeVP1)

	return tree, nil
}

func (sp *SyntaxParser) PrettyPrint(tree Tree, level int) {

	fmt.Println(getTabs(level), tree.Type)
	level += 1
	for _, value := range tree.Leaves {
		switch leaf := value.(type) {
		case TaggedToken:
			fmt.Println(getTabs(level), leaf)
		case Tree:
			sp.PrettyPrint(leaf, level)
		}
	}
}

func getTabs(level int) string {
	tabs := ""
	for i := 0; i < level; i++ {
		tabs += "    "
	}
	return tabs
}
