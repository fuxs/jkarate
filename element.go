package jkarate

type TokenType int

const (
	invalidType TokenType = iota
	nullType
	arrayType
	boolType
	objectType
	numberType
	stringType
)

type Element struct {
	Type   TokenType
	Object map[string]*Element
	Array  []*Element
	Value  string
}

func (e *Element) IsNull() bool {
	return e.Type == nullType
}

func (e *Element) IsArray() bool {
	return e.Type == arrayType
}

func (e *Element) IsObject() bool {
	return e.Type == objectType
}

func (e *Element) Size() int {
	switch e.Type {
	case arrayType:
		return len(e.Array)
	case objectType:
		return len(e.Object)
	}
	return 0
}

type elementStack []*Element

func NewElementStack() elementStack {
	return make(elementStack, 0, 32)
}

func NewElementStackWithCapacity(capacity int) elementStack {
	return make(elementStack, 0, capacity)
}

func (es *elementStack) Push(el *Element) {
	*es = append(*es, el)
}

func (es *elementStack) Peek() *Element {
	l := len(*es)
	if l == 0 {
		return nil
	}
	return (*es)[l-1]
}

func (es *elementStack) Pop() *Element {
	l := len(*es)
	if l == 0 {
		return nil
	}
	result := (*es)[l-1]
	*es = (*es)[:l-1]
	return result
}
