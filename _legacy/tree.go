package gonlp

type Leaves interface {
}

type Tree struct {
	Type   string
	Leaves []Leaves
}
