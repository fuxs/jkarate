package jkarate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Element(t *testing.T) {
	must := require.New(t)
	tokenizer := NewTokenizer().ReadString(`{"the":"quick", "brown":["fox", "jumps", {"over":"the"}, ["lazy","dog"], true]}`)
	elem, err := tokenizer.Parse()

	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())

	must.False(elem.Has("tha"))
	_, err = elem.StringE("tha")
	must.Error(err)

	var str string
	str, err = elem.StringE("", "the")
	must.Equal("quick", str)
	must.NoError(err)

	must.Equal(5, elem.Get("the").Size())
	_, err = elem.AtE(0)
	must.Error(err)

	arr := elem.Get("brown")
	must.True(arr.IsArray())
	must.Equal(5, arr.Size())

	var aux *Element
	aux, err = arr.AtE(3)
	must.NoError(err)
	must.True(aux.IsArray())

	_, err = arr.AtE(5)
	must.Error(err)

	_, err = arr.GetE("5")
	must.Error(err)

	_, err = elem.GetE("the", "4")
	must.Error(err)

	_, err = elem.Get("the").GetE("4")
	must.Error(err)

	must.True(elem.Get("brown", "4").IsBoolean())
	must.Equal(0, elem.Get("brown", "4").Size())
	must.False(elem.Get("brown", "4").Has("value"))
}

func Test_ElementStack(t *testing.T) {
	must := require.New(t)
	tokenizer := NewTokenizer().ReadString(`{"the":"quick", "brown":["fox", "jumps", {"over":"the"}, ["lazy","dog"]]}`)
	elem, err := tokenizer.Parse()

	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())

	es := NewElementStackWithCapacity(16)
	must.Nil(es.Peek())
	es.Push(elem.Get("the"))
	es.Push(elem.Get("brown"))

	arr := es.Peek()
	must.True(arr.IsArray())
	must.Equal(4, arr.Size())
}
