package jkarate

import (
	"errors"
)

func (t *Tokenizer) Parse() (*Element, error) {
	var (
		current, child *Element
		name           string
	)
	token := &Token{}
	es := NewElementStack()
	if !t.Next(token) {
		return nil, token.Err
	}
	if token.Type == objectLeftToken {
		current = &Element{
			Type:   objectType,
			Object: make(map[string]*Element),
		}
		// check for empty object
		if !t.Next(token) {
			return nil, token.Err
		}
		if token.Type == objectRightToken {
			if !t.Next(token) {
				return nil, token.Err
			}
			if token.Type != doneToken {
				return nil, errors.New("unexpected token after top level object")
			}
			return current, nil
		}
		goto NAME_STATE_FIRST
	}
	if token.Type == arrayLeftToken {
		current = &Element{
			Type:  arrayType,
			Array: make([]*Element, 0, 8),
		}
		// check for empty array
		if !t.Next(token) {
			return nil, token.Err
		}
		if token.Type == arrayRightToken {
			if !t.Next(token) {
				return nil, token.Err
			}
			if token.Type != doneToken {
				return nil, errors.New("unexpected token after top level array")
			}
			return current, nil
		}
		goto ARRAY_STATE_FIRST
	}
	if token.Type == doneToken {
		current = &Element{
			Type: nullType,
		}
		return current, nil
	}
	return nil, errors.New("unexpected token, must be '[', '{' or empty file")

NAME_STATE:
	if !t.Next(token) {
		return nil, token.Err
	}
NAME_STATE_FIRST:
	if token.Type == stringToken {
		name = token.Str
		if !t.Next(token) {
			return nil, token.Err
		}
		if token.Type == colonToken {
			if !t.Next(token) {
				return nil, token.Err
			}
			switch token.Type {
			case stringToken:
				child = &Element{
					Type:  stringType,
					Value: token.Str,
				}
				goto NEXT_OBJECT_PAIR
			case numToken:
				child = &Element{
					Type:  numberType,
					Value: token.Str,
				}
				goto NEXT_OBJECT_PAIR
			case boolToken:
				child = &Element{
					Type:  boolType,
					Value: token.Str,
				}
				goto NEXT_OBJECT_PAIR
			case nullToken:
				child = &Element{
					Type:  nullType,
					Value: "",
				}
				goto NEXT_OBJECT_PAIR
			case objectLeftToken:
				es.Push(current)
				current = &Element{
					Type:   objectType,
					Object: make(map[string]*Element),
					Value:  name,
				}
				goto NAME_STATE
			case arrayLeftToken:
				es.Push(current)
				current = &Element{
					Type:  arrayType,
					Array: make([]*Element, 0, 8),
					Value: name,
				}
				goto ARRAY_STATE
			}
			return nil, errors.New("unexpected token")
		}
		return nil, errors.New("expected ':'")
	}
	return nil, errors.New("expected string")

NEXT_OBJECT_PAIR:
	current.Object[name] = child
	if !t.Next(token) {
		return nil, token.Err
	}
	if token.Type == commaToken {
		goto NAME_STATE
	}
	if token.Type == objectRightToken {
		child = current
		if current = es.Pop(); current == nil {
			if !t.Next(token) {
				return nil, token.Err
			}
			if token.Type != doneToken {
				return nil, errors.New("unexpected token after top level object")
			}
			return child, nil
		}
		if current.Type == objectType {
			current.Object[current.Value] = child
			goto NEXT_OBJECT_PAIR
		}
		// must be an array
		current.Array = append(current.Array, child)
		goto NEXT_ARRAY_ELEMENT
	}
	return nil, errors.New("unxpected token")

ARRAY_STATE:
	if !t.Next(token) {
		return nil, token.Err
	}
ARRAY_STATE_FIRST:
	switch token.Type {
	case stringToken:
		child = &Element{
			Type:  stringType,
			Value: token.Str,
		}
		goto NEXT_ARRAY_ELEMENT
	case numToken:
		child = &Element{
			Type:  numberType,
			Value: token.Str,
		}
		goto NEXT_ARRAY_ELEMENT
	case boolToken:
		child = &Element{
			Type:  boolType,
			Value: token.Str,
		}
		goto NEXT_ARRAY_ELEMENT
	case nullToken:
		child = &Element{
			Type: nullType,
		}
		goto NEXT_ARRAY_ELEMENT
	case objectLeftToken:
		es.Push(current)
		current = &Element{
			Type:   objectType,
			Object: make(map[string]*Element),
			Value:  name,
		}
		goto NAME_STATE
	case arrayLeftToken:
		es.Push(current)
		current = &Element{
			Type:  arrayType,
			Array: make([]*Element, 0, 8),
			Value: name,
		}
		goto ARRAY_STATE
	}
	return nil, errors.New("unexpected token")
NEXT_ARRAY_ELEMENT:
	current.Array = append(current.Array, child)
	if !t.Next(token) {
		return nil, token.Err
	}
	if token.Type == commaToken {
		goto ARRAY_STATE
	}
	if token.Type == arrayRightToken {
		child = current
		if current = es.Pop(); current == nil {
			if !t.Next(token) {
				return nil, token.Err
			}
			if token.Type != doneToken {
				return nil, errors.New("unexpected token after top level array")
			}
			return child, nil
		}
		if current.Type == objectType {
			current.Object[current.Value] = child
			goto NEXT_OBJECT_PAIR
		}
		// must be an arry
		current.Array = append(current.Array, child)
		goto NEXT_ARRAY_ELEMENT
	}
	return nil, errors.New("expected , or ]")
}
