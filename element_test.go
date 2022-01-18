package jkarate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Element(t *testing.T) {
	must := require.New(t)
	tokenizer := NewTokenizer().ReadString(`{"the":"quick", "brown":["fox", "jumps", {"over":"the"}, ["lazy","dog"]]}`)
	elem, err := tokenizer.Parse()

	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())

	_, err = elem.StringE("tha")
	must.Error(err)

	var str string
	str, err = elem.StringE("", "the")
	must.Equal("quick", str)
	must.NoError(err)

	_, err = elem.AtE(0)
	must.Error(err)

	arr := elem.Get("brown")
	must.True(arr.IsArray())
	must.Equal(4, arr.Size())

	var aux *Element
	aux, err = arr.AtE(3)
	must.NoError(err)
	must.True(aux.IsArray())

	_, err = arr.AtE(4)
	must.Error(err)

	_, err = arr.GetE("4")
	must.Error(err)

	_, err = elem.GetE("the", "4")
	must.Error(err)

	_, err = elem.Get("the").GetE("4")
	must.Error(err)
}
