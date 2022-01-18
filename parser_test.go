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

func Test_ParseErrors(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().Read(newErrorReader("{}"))
	elem, err := tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer = NewTokenizerSize(4).Read(newErrorReader("   {}"))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)
	tokenizer.Read(newErrorReader("   []"))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader("  {} "))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader("  [] "))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString("{}{}")
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString("[][]")
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString("[1][]")
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`"broken"`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`{"broken"}`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`{"broken":}`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`[123 : 456]`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`[:]`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`{"broken":"world":}`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`{"broken":"world"}{}`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.ReadString(`{123:456}`)
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a":"b"}`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a" `))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a":"b",    `))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a" :   `))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a":"b"}   []`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a" : [ ]}`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`{"a" : { }}`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`[1,  2]`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`[1,2]   []`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)

	tokenizer.Read(newErrorReader(`[[] ]`))
	elem, err = tokenizer.Parse()
	must.Error(err)
	must.Nil(elem)
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

	tokenizer.ReadString(`{"a":[432]}`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())
	must.Equal("432", elem.Get("a").At(0).String())
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
	must.Equal("value4", elem.String("level1", "level2", "level3", "name4"))
	must.Equal("value5", elem.String("level1", "level2", "level3", "level4", "name5"))

	tokenizer.ReadString(`{
		"level1": {
			"level2": {
				"level3": {
					"level4": {
						"name5": "value5"
					},
					"name4": "value4"
				},
				"name3": "value3"
			},
			"name2": "value2"
		},
		"name1": "value1"
	}`)
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
	must.Equal("value4", elem.String("level1", "level2", "level3", "name4"))
	must.Equal("value5", elem.String("level1", "level2", "level3", "level4", "name5"))

	tokenizer.ReadString(`{
		"katze1":"hund1",
		"level1": {
			"katze2":"hund2",
			"level2": {
				"katze3":"hund3",
				"level3": {
					"katze4":"hund4",
					"level4": {
						"katze5":"hund5",
						"name5": "value5"
					},
					"name4": "value4"
				},
				"name3": "value3"
			},
			"name2": "value2"
		},
		"name1": "value1"
	}`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsObject())

	must.True(elem.Get("level1").IsObject())
	must.Equal(3, elem.Get("level1").Size())
	must.Equal("value1", elem.String("name1"))
	must.Equal("hund1", elem.String("katze1"))
	must.True(elem.Get("level1", "level2").IsObject())
	must.Equal(3, elem.Get("level1", "level2").Size())
	must.Equal("value2", elem.String("level1", "name2"))
	must.Equal("hund2", elem.String("level1", "katze2"))
	must.True(elem.Get("level1", "level2", "level3").IsObject())
	must.Equal(3, elem.Get("level1", "level2", "level3").Size())
	must.Equal("value3", elem.String("level1", "level2", "name3"))
	must.Equal("hund3", elem.String("level1", "level2", "katze3"))

	must.True(elem.Get("level1", "level2", "level3", "level4").IsObject())
	must.Equal(2, elem.Get("level1", "level2", "level3", "level4").Size())
	must.Equal("value4", elem.String("level1", "level2", "level3", "name4"))
	must.Equal("value5", elem.String("level1", "level2", "level3", "level4", "name5"))
	must.Equal("hund4", elem.String("level1", "level2", "level3", "katze4"))
	must.Equal("hund5", elem.String("level1", "level2", "level3", "level4", "katze5"))
}

func Test_SimpleArray(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().ReadString(`["hello","world"]`)
	elem, err := tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal("hello", elem.String("0"))
	must.Equal("world", elem.String("1"))
	must.Equal(2, elem.Size())

	tokenizer.ReadString(`[{"hello":"world"}]`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal("world", elem.At(0).String("hello"))
	must.Equal(1, elem.Size())

	tokenizer.ReadString(`["hello","world", {}, [], true, false, null, 123]`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal("hello", elem.String("0"))
	must.Equal("world", elem.String("1"))
	must.Equal("true", elem.String("4"))
	must.Equal("false", elem.String("5"))
	must.Equal("", elem.String("6"))
	must.Equal("123", elem.String("7"))
	must.Equal(8, elem.Size())

	tokenizer.ReadString(`[[[[]]]]`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal(1, elem.Size())
	must.True(elem.Get("0").IsArray())
	must.Equal(1, elem.Get("0").Size())
	must.True(elem.Get("0", "0").IsArray())
	must.Equal(1, elem.Get("0", "0").Size())
	must.True(elem.Get("0", "0", "0").IsArray())
	must.Equal(0, elem.Get("0", "0", "0").Size())

	tokenizer.ReadString(`[[[],[[],[]]],[[[],[[]]],[]]]`)
	elem, err = tokenizer.Parse()
	must.NoError(err)
	must.NotNil(elem)
	must.True(elem.IsArray())
	must.Equal(2, elem.Size())
	must.True(elem.Get("0").IsArray())
	must.Equal(2, elem.Get("0").Size())
	must.True(elem.Get("1").IsArray())
	must.Equal(2, elem.Get("1").Size())
	must.Equal(0, elem.Get("0", "0").Size())
	must.True(elem.Get("0", "1").IsArray())
	must.Equal(2, elem.Get("0", "1").Size())
	must.Equal(0, elem.Get("0", "1", "1").Size())

}
