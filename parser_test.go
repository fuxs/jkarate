package jkarate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EmptyDocuments(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().ReadString("")
	elem, err := tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsNull())

	tokenizer.ReadString(" ")
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsNull())

	tokenizer.ReadString("[]")
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal(0, elem.Size())

	tokenizer.ReadString(" [ ] ")
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal(0, elem.Size())

	tokenizer.ReadString(" { } ")
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())
	must.Equal(0, elem.Size())
}
