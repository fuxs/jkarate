package jkarate

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_strings(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().ReadString("")
	must.NotNil(tokenizer)
	token := &Token{}
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	strs := []string{`""`, `"hello"`, ` "world"`, ` "quick" `, `"fox" `}
	exp := []string{"", "hello", "world", "quick", "fox"}
	for i, str := range strs {
		// empty string
		tokenizer.ReadString(str)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(exp[i], token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	tokenizer.ReadArray([]byte{})
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)
}

func Test_invalid(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizer().ReadString("")
	must.NotNil(tokenizer)
	token := &Token{}
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	strs := []string{"a", " a", "\na", ` a "fox"`}
	for _, str := range strs {
		// empty string
		tokenizer.ReadString(str)
		must.False(tokenizer.Next(token))
		must.Equal(invalidToken, token.Type)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

const wsBytes = " \n\r\t"

func randWS(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = wsBytes[rand.Int63()%int64(len(wsBytes))]
	}
	return string(b)
}

func Test_buffers(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(7)
	must.NotNil(tokenizer)
	token := &Token{}

	for i := 0; i < 128; i++ {
		str := randString(i)
		// empty string
		tokenizer.ReadString(`"` + str + `"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for i := 0; i < 128; i++ {
		str := randString(i)
		// empty string
		tokenizer.ReadString(`"` + str + "\"\n\"" + str + `"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for j := 0; j < 34; j++ {
		ws := randWS(j)
		for i := 0; i < 34; i++ {
			str := randString(i)

			// empty string
			tokenizer.ReadString(`"` + str + `"` + ws + `"` + str + `"`)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(str, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(str, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(doneToken, token.Type)

			tokenizer.ReadString(ws + `"` + str + `"` + ws + `"` + str + `"` + ws)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(str, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(str, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(doneToken, token.Type)
		}
	}
}

type errorReader struct {
	reader io.Reader
	size   int
	read   int
}

func newErrorReader(str string) *errorReader {
	return &errorReader{
		reader: strings.NewReader(str),
		size:   len(str),
	}
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	if n, err = e.reader.Read(p); err != nil {
		return n, err
	}
	e.read = e.read + n
	if e.read >= e.size {
		return 0, errors.New("test error")
	}
	return n, err
}

func Test_buffersErr(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}
	for i := 0; i < 256; i++ {
		str := randString(i)
		// empty string
		tokenizer.Read(newErrorReader(`"` + str + `"`))
		must.False(tokenizer.Next(token))
		must.Equal(token.Err.Error(), "test error")
	}

	tokenizer.Read(newErrorReader(`       ""`))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader(`    "\n "`))
	must.False(tokenizer.Next(token))
}

const escBytes = "nrtfb/\\\""

func randESC(n int) string {
	b := make([]byte, n*2)
	esc := true
	for i := range b {
		if esc {
			b[i] = '\\'
		} else {
			b[i] = escBytes[rand.Int63()%int64(len(escBytes))]
		}
		esc = !esc
	}
	return string(b)
}

func escResult(str string) string {
	b := make([]byte, len(str)/2)
	for i := range b {
		switch str[i*2+1] {
		case 'n':
			b[i] = '\n'
		case 'r':
			b[i] = '\r'
		case 't':
			b[i] = '\t'
		case 'f':
			b[i] = '\f'
		case 'b':
			b[i] = '\b'
		case '/':
			b[i] = '/'
		case '"':
			b[i] = '"'
		case '\\':
			b[i] = '\\'
		}
	}
	return string(b)
}

func Test_escaped(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	str := `"\n\r\f\b\t\\\"\/"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("\n\r\f\b\t\\\"/", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	for i := 0; i < 34; i++ {
		str := randString(i)
		// empty string
		tokenizer.Read(newErrorReader(`"` + str + `\n"`))
		must.False(tokenizer.Next(token))
		must.Equal(token.Err.Error(), "test error")
	}

	for i := 0; i < 128; i++ {
		str = randESC(i)
		tokenizer.ReadString(`"` + str + `"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(escResult(str), token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for i := 0; i < 128; i++ {
		str = randString(i)
		tokenizer.ReadString(`"` + str + `\n"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(str+"\n", token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for i := 0; i < 128; i++ {
		str = randString(i)
		tokenizer.ReadString(`"\b` + str + `"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal("\b"+str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for j := 0; j < 34; j++ {
		estr := randESC(j)
		restr := escResult(estr)
		for i := 0; i < 34; i++ {
			str = randString(i)
			tokenizer.ReadString(`"` + estr + str + `"`)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(restr+str, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(doneToken, token.Type)

			tokenizer.ReadString(`"` + str + estr + `"`)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(str+restr, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(doneToken, token.Type)

			tokenizer.ReadString(`"` + str + estr + str + `"`)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(str+restr+str, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(doneToken, token.Type)

			tokenizer.ReadString(`"` + estr + str + estr + `"`)
			must.True(tokenizer.Next(token))
			must.Equal(stringToken, token.Type)
			must.Equal(restr+str+restr, token.Str)
			must.True(tokenizer.Next(token))
			must.Equal(doneToken, token.Type)
		}
	}

	str = `"\x"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"hello world \x"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("hello world �", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)
}

func Test_unicode(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	str := `"\u0079"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("y", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\u00E4"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("ä", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\u20AC"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("€", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	for i := 0; i < 128; i++ {
		str = randString(i)
		tokenizer.ReadString(`"` + str + `\u20AC"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(str+"€", token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	str = `"\uD834\uDD1E"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("𝄞", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD853\uDF5C"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("𤽜", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	for i := 0; i < 128; i++ {
		str = randString(i)
		tokenizer.ReadString(`"` + str + `\uD853\uDF5C"`)
		must.True(tokenizer.Next(token))
		must.Equal(stringToken, token.Type)
		must.Equal(str+"𤽜", token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

}

func Test_unicodeErr(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	str := `"\uO079"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�O079", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\u000O"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�O", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uDD1E\uD834"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("��", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uDD1E\uD834\uDD1E"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�𝄞", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD834 futter"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("� futter", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD834\nbeilage"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�\nbeilage", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD834\uGD1G"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�GD1G", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD834\uDD1G"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�G", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD853\uD853\uDF5C"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("�𤽜", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"\uD853\uD853\u20AC"`
	tokenizer.ReadString(str)
	must.True(tokenizer.Next(token))
	must.Equal(stringToken, token.Type)
	must.Equal("��€", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	str = `"space\uD853\u20AC"`
	tokenizer.Read(newErrorReader(str))
	must.False(tokenizer.Next(token))

	str = `"flat\u20AC"`
	tokenizer.Read(newErrorReader(str))
	must.False(tokenizer.Next(token))

	str = `"_\uD853\u20AC"`
	//tokenizer = NewTokenizerSize(8)
	tokenizer.Read(newErrorReader(str))
	must.False(tokenizer.Next(token))

	str = `"\uD853\u20AC"`
	//tokenizer = NewTokenizerSize(8)
	tokenizer.Read(newErrorReader(str))
	must.False(tokenizer.Next(token))
}

const numBytes = "123456789"
const numBytesWZ = "0123456789"

func randNum(n int) string {
	b := make([]byte, n)
	for i := range b {
		if i == 0 {
			b[i] = numBytes[rand.Int63()%int64(len(numBytes))]
		} else {
			b[i] = numBytesWZ[rand.Int63()%int64(len(numBytesWZ))]
		}

	}
	return string(b)
}

func Test_number(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(7)
	must.NotNil(tokenizer)
	token := &Token{}

	tokenizer.ReadString("0 ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("0", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("  0  ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("0", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("   0   ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("0", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("    0    ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("0", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("      0 ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("0", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("       0 ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("0", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	for i := 1; i < 128; i++ {
		str := randNum(i) + "E+2"
		// empty string
		tokenizer.ReadString(str + " ")
		must.True(tokenizer.Next(token))
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for i := 1; i < 128; i++ {
		str := randNum(i) + ""
		// empty string
		tokenizer.ReadString(str + " ")
		must.True(tokenizer.Next(token))
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	for i := 1; i < 1000000; i = i + 91 {
		str := fmt.Sprintf("%v", i)
		tokenizer.ReadString(str + " ")
		must.True(tokenizer.Next(token))
		must.Equal(numToken, token.Type)
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)

		str = str + "." + str
		tokenizer.ReadString(str + " ")
		must.True(tokenizer.Next(token))
		must.Equal(numToken, token.Type)
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)

		str = str + "e" + fmt.Sprintf("%v", i)
		tokenizer.ReadString(str + " ")
		must.True(tokenizer.Next(token))
		must.Equal(numToken, token.Type)
		must.Equal(str, token.Str)
		must.True(tokenizer.Next(token))
		must.Equal(doneToken, token.Type)
	}

	tokenizer.ReadString("-1 -24 -356 ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("-1", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("-24", token.Str)
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("-356", token.Str)

	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)
}

func Test_numberErr(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(7)
	must.NotNil(tokenizer)
	token := &Token{}

	for i := 0; i <= 128; i++ {
		ws := randWS(i)
		tokenizer.ReadString(ws + "- ")
		must.False(tokenizer.Next(token))
		must.Equal(invalidToken, token.Type)
	}

	tokenizer.ReadString("1E+e ")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.ReadString("1. ")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)
	tokenizer.ReadString("1.o ")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.ReadString(" 1 1. ")
	must.True(tokenizer.Next(token))
	must.Equal(numToken, token.Type)
	must.Equal("1", token.Str)
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.Read(newErrorReader("      -1"))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("      0."))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("      1."))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("     2.4"))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("    3.45"))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("     2e1"))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("    5e-1"))
	must.False(tokenizer.Next(token))

	tokenizer.Read(newErrorReader("   7e-12"))
	must.False(tokenizer.Next(token))
}

func Test_bool(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	tokenizer.ReadString("true")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("    true")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("     true")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("    true    ")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("       true")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("        true      ")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("false")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("   false")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("    false")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("        false")
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString(" true false true false true false")

	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("truefalse")

	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.True(token.Boolean())
	must.True(tokenizer.Next(token))
	must.Equal(boolToken, token.Type)
	must.False(token.Boolean())
}

func Test_boolErr(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	tokenizer.ReadString("True")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.ReadString("truE")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.ReadString("False")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.ReadString("falsE")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.Read(newErrorReader("     true"))
	must.False(tokenizer.Next(token))
	must.Equal(errorToken, token.Type)

	tokenizer.Read(newErrorReader("     false"))
	must.False(tokenizer.Next(token))
	must.Equal(errorToken, token.Type)
}

func Test_null(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	tokenizer.ReadString("null")
	must.True(tokenizer.Next(token))
	must.Equal(nullToken, token.Type)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("    null")
	must.True(tokenizer.Next(token))
	must.Equal(nullToken, token.Type)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("     null")
	must.True(tokenizer.Next(token))
	must.Equal(nullToken, token.Type)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("    null    ")
	must.True(tokenizer.Next(token))
	must.Equal(nullToken, token.Type)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("       null")
	must.True(tokenizer.Next(token))
	must.Equal(nullToken, token.Type)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)

	tokenizer.ReadString("        null      ")
	must.True(tokenizer.Next(token))
	must.Equal(nullToken, token.Type)
	must.True(tokenizer.Next(token))
	must.Equal(doneToken, token.Type)
}

func Test_nullErr(t *testing.T) {
	must := require.New(t)
	// nothing
	tokenizer := NewTokenizerSize(8)
	must.NotNil(tokenizer)
	token := &Token{}

	tokenizer.ReadString("Null")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.ReadString("nulL")
	must.False(tokenizer.Next(token))
	must.Equal(invalidToken, token.Type)

	tokenizer.Read(newErrorReader("     null"))
	must.False(tokenizer.Next(token))
	must.Equal(errorToken, token.Type)
}
