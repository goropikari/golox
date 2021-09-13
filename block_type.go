package golox

type BlockType int

const (
	NoneBlock BlockType = iota
	ForBlock
	IfBlock
	WhileBlock
)
