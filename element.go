package jkarate

import (
	"fmt"
	"strconv"
)

type ElementType int

const (
	invalidType ElementType = iota
	nullType
	arrayType
	boolType
	objectType
	numberType
	stringType
)

var elmentStr []string = []string{
	"invalid",
	"null",
	"array",
	"bool",
	"object",
	"number",
	"string",
}

type Element struct {
	Type   ElementType
	Object map[string]*Element
	Array  []*Element
	Value  string
}

func (e *Element) String(path ...string) string {
	if result := e.Get(path...); result != nil {
		return result.Value
	}
	return ""
}

func (e *Element) StringE(path ...string) (string, error) {
	if result, err := e.GetE(path...); result != nil {
		return result.Value, err
	} else {
		return "", err
	}
}

func (e *Element) At(index int) *Element {
	result, _ := e.AtE(index)
	return result
}

func (e *Element) AtE(index int) (*Element, error) {
	if e.Type != arrayType {
		return nil, fmt.Errorf("element %s does not support indexed access", elmentStr[e.Type])
	}
	if index >= 0 && index < len(e.Array) {
		return e.Array[index], nil
	}
	return nil, fmt.Errorf("index %v is out of range [0:%v]", index, len(e.Array))
}

func (e *Element) Get(path ...string) *Element {
	result, _ := e.GetE(path...)
	return result
}

func (e *Element) GetE(path ...string) (*Element, error) {
	var (
		index int
		err   error
	)
	current := e
	for i, name := range path {
		if name == "" {
			continue
		}
		if current.Type == objectType {
			if current = current.Object[name]; current == nil {
				return nil, fmt.Errorf("could not find path element #%v with name %s", i, name)
			}
			continue
		}
		if current.Type == arrayType {
			if index, err = strconv.Atoi(name); err == nil {
				if index < len(current.Array) && index >= 0 {
					current = current.Array[index]
					continue
				}
				err = fmt.Errorf("path element #%v with index %s is out of range [0:%v]", i, name, len(current.Array))
			}
			return nil, err
		}
		if i == 0 {
			return nil, fmt.Errorf("current element does not support GetE")
		}
		return nil, fmt.Errorf("path element #%v of type %s does not support GetE", i-1, elmentStr[current.Type])
	}
	return current, nil
}

func (e *Element) Has(name string) bool {
	if e != nil && e.Object != nil {
		_, ok := e.Object[name]
		return ok
	}
	return false
}

func (e *Element) IsNull() bool {
	return e != nil && e.Type == nullType
}

func (e *Element) IsArray() bool {
	return e != nil && e.Type == arrayType
}

func (e *Element) IsObject() bool {
	return e != nil && e.Type == objectType
}

func (e *Element) Size() int {
	if e != nil {
		switch e.Type {
		case arrayType:
			return len(e.Array)
		case objectType:
			return len(e.Object)
		case stringType:
			return len(e.Value)
		}
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
