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

func Test_SimpleObject(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().ReadString(`{"hello":"world"}`)
	elem, err := tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())
	must.True(elem.Has("hello"))
	must.Equal("world", elem.String("hello"))
	must.False(elem.Has("hallo"))
	must.Equal("", elem.String("hallo"))
	must.Equal(1, elem.Size())

	tokenizer.ReadString(`{"hello":"world", "dog":"food"}`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())
	must.Equal("world", elem.String("hello"))
	must.Equal("food", elem.String("dog"))
	must.Equal(2, elem.Size())

	tokenizer.ReadString(`{"string" : "example",
	"number": 123456789,
	"number2": 123.45678,
	"boolean": true,
	"Null": null,
	"Object": {},
	"Array": [] }`)

	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())
	must.Equal("example", elem.String("string"))
	must.Equal("123456789", elem.String("number"))
	must.Equal("123.45678", elem.String("number2"))
	must.Equal("true", elem.String("boolean"))
	must.Equal("", elem.String("Null"))
	must.True(elem.Get("Null").IsNull())
	must.True(elem.Get("Object").IsObject())
	must.True(elem.Get("Array").IsArray())

}

func Test_NestedObject(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().ReadString(`{"level1": {
		"level2": {
			"level3" : { "level4":{}}
		}
	}}`)
	elem, err := tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())

	must.True(elem.Get("level1").IsObject())
	must.Equal(1, elem.Get("level1").Size())
	must.True(elem.Get("level1", "level2").IsObject())
	must.Equal(1, elem.Get("level1", "level2").Size())
	must.True(elem.Get("level1", "level2", "level3").IsObject())
	must.Equal(1, elem.Get("level1", "level2", "level3").Size())
	must.True(elem.Get("level1", "level2", "level3", "level4").IsObject())
	must.Equal(0, elem.Get("level1", "level2", "level3", "level4").Size())

	tokenizer.ReadString(`{"name1":"value1", "level1": {
		"name2":"value2",
		"level2": {
			"name3":"value3",
			"level3" : { "name4":"value4", "level4":{"name5":"value5"}}
		}
	}}`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())

	must.True(elem.Get("level1").IsObject())
	must.Equal(2, elem.Get("level1").Size())
	must.Equal("value1", elem.String("name1"))
	must.True(elem.Get("level1", "level2").IsObject())
	must.Equal(2, elem.Get("level1", "level2").Size())
	must.Equal("value2", elem.String("level1", "name2"))
	must.True(elem.Get("level1", "level2", "level3").IsObject())
	must.Equal(2, elem.Get("level1", "level2", "level3").Size())
	must.Equal("value3", elem.String("level1", "level2", "name3"))

	must.True(elem.Get("level1", "level2", "level3", "level4").IsObject())
	must.Equal(1, elem.Get("level1", "level2", "level3", "level4").Size())
	must.Equal("value5", elem.String("level1", "level2", "level3", "level4", "name5"))
}
