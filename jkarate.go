package jkarate

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	invalidToken int = iota
	errorToken
	arrayLeftToken
	arrayRightToken
	boolToken
	doneToken
	colonToken
	commaToken
	nullToken
	numToken
	objectLeftToken
	objectRightToken
	stringToken
	whiteSpaceToken
)

const defaultBufferSize = 4096

var tokenType [256]int
var hexValid [256]bool
var hexRune [256]uint32

var RUE []byte = []byte{'r', 'u', 'e'}
var ALSE []byte = []byte{'a', 'l', 's', 'e'}
var ULL []byte = []byte{'u', 'l', 'l'}

func init() {
	tokenType['"'] = stringToken
	tokenType['-'] = numToken
	tokenType['0'] = numToken
	tokenType['1'] = numToken
	tokenType['2'] = numToken
	tokenType['3'] = numToken
	tokenType['4'] = numToken
	tokenType['5'] = numToken
	tokenType['6'] = numToken
	tokenType['7'] = numToken
	tokenType['8'] = numToken
	tokenType['9'] = numToken
	tokenType['f'] = boolToken
	tokenType['t'] = boolToken
	tokenType['n'] = nullToken
	tokenType['{'] = objectLeftToken
	tokenType['}'] = objectRightToken
	tokenType['['] = arrayLeftToken
	tokenType[']'] = arrayRightToken
	tokenType[':'] = colonToken
	tokenType[','] = commaToken
	tokenType[' '] = whiteSpaceToken
	tokenType['\t'] = whiteSpaceToken
	tokenType['\n'] = whiteSpaceToken
	tokenType['\r'] = whiteSpaceToken
	for i := '0'; i <= '9'; i++ {
		hexValid[i] = true
		hexRune[i] = uint32(i - '0')
	}
	for i := 'a'; i <= 'f'; i++ {
		hexValid[i] = true
		hexRune[i] = uint32(i - 'a' + 10)
	}
	for i := 'A'; i <= 'F'; i++ {
		hexValid[i] = true
		hexRune[i] = uint32(i - 'A' + 10)
	}
}

type Tokenizer struct {
	reader     io.Reader
	buffer     []byte
	strBuffer  []byte
	bufferSize int
	begin      int
	end        int
	number     int
	line       int
}

type Token struct {
	Type   int
	Str    string
	Err    error
	Number int
}

func (t *Token) Boolean() bool {
	return t.Number != 0
}

func NewTokenizer() *Tokenizer {
	return NewTokenizerSize(defaultBufferSize)
}

func NewTokenizerSize(s int) *Tokenizer {
	return &Tokenizer{
		buffer:     make([]byte, s),
		strBuffer:  make([]byte, defaultBufferSize*2),
		bufferSize: s,
	}
}

func (t *Tokenizer) ReadString(str string) *Tokenizer {
	return t.Read(strings.NewReader(str))
}

func (t *Tokenizer) ReadArray(arr []byte) *Tokenizer {
	return t.Read(bytes.NewReader(arr))
}

func (t *Tokenizer) Read(reader io.Reader) *Tokenizer {
	t.reader = reader
	t.begin = 0
	t.end = 0
	return t
}

func (t *Tokenizer) Next(token *Token) bool {
	var (
		i, j, n, tt int
		c           byte
		err         error
		current     []byte
	)
	t.number++
	token.Number = t.number
LOOP_BUFFER:
	if t.begin == t.end {
		// loop @buffer
	OUTER_BUFFER_LOOP:
		for {
			if n, err = t.reader.Read(t.buffer); n == 0 {
				if err != nil {
					if err == io.EOF {
						token.Type = doneToken
						token.Err = nil
						return true
					}
					token.Err = err
					return false
				}
			} else {
				for i = 0; i < n; i++ {
					c = t.buffer[i]
					if tt = tokenType[c]; tt != whiteSpaceToken {
						t.begin = i
						t.end = n
						break OUTER_BUFFER_LOOP
					}
					if c == '\n' {
						t.line++
					}
				}
			}
		}
	} else {
		// eat all whitespaces
		for i, n = t.begin, t.end; i < n; i++ {
			c = t.buffer[i]
			if tt = tokenType[c]; tt != whiteSpaceToken {
				break
			}
			if c == '\n' {
				t.line++
			}
		}
		t.begin = i
		if i == n {
			goto LOOP_BUFFER
		}
	}
	// i == t.begin
	// n == t.end
	switch tt {
	case arrayLeftToken:
		token.Type = arrayLeftToken
		token.Str = "["
		t.begin = i + 1
		t.end = n
		return true
	case arrayRightToken:
		token.Type = arrayRightToken
		token.Str = "]"
		t.begin = i + 1
		t.end = n
		return true
	case boolToken:
		if c == 't' {
			for {
				if i++; i == n {
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								token.Type = errorToken
								t.begin = 0
								t.end = 0
								return false
							}
						} else {
							t.begin = 0
							i = 0
							break
						}
					}
				}
				if t.buffer[i] != RUE[j] {
					token.Err = fmt.Errorf("invalid true value in line %v", t.line)
					token.Type = invalidToken
					t.begin = i
					t.end = n
					return false
				}
				j++
				if j == 3 {
					token.Type = boolToken
					token.Str = "true"
					token.Number = 1
					t.begin = i + 1
					t.end = n
					return true
				}
			}
		} else {
			for {
				if i++; i == n {
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								token.Type = errorToken
								t.begin = 0
								t.end = 0
								return false
							}
						} else {
							t.begin = 0
							i = 0
							break
						}
					}
				}
				if t.buffer[i] != ALSE[j] {
					token.Err = fmt.Errorf("invalid false value in line %v", t.line)
					token.Type = invalidToken
					t.begin = i
					t.end = n
					return false
				}
				j++
				if j == 4 {
					token.Str = "false"
					token.Number = 0
					token.Type = boolToken
					t.begin = i + 1
					t.end = n
					return true
				}
			}
		}
	case colonToken:
		token.Type = colonToken
		token.Str = ":"
		t.begin = i + 1
		t.end = n
		return true
	case commaToken:
		token.Type = commaToken
		token.Str = ","
		t.begin = i + 1
		t.end = n
		return true
	case nullToken:
		for {
			if i++; i == n {
				for {
					if n, err = t.reader.Read(t.buffer); n == 0 {
						if err != nil {
							token.Err = err
							token.Type = errorToken
							t.begin = 0
							t.end = 0
							return false
						}
					} else {
						t.begin = 0
						i = 0
						break
					}
				}
			}
			if t.buffer[i] != ULL[j] {
				token.Err = fmt.Errorf("invalid null value in line %v", t.line)
				token.Type = invalidToken
				t.begin = i
				t.end = n
				return false
			}
			j++
			if j == 3 {
				token.Str = "null"
				token.Number = 0
				token.Type = nullToken
				t.begin = i + 1
				t.end = n
				return true
			}
		}
	case numToken:
		current = t.strBuffer[:0]
		if c == '-' {
			if i++; i == n {
				current = append(current, '-')
				// read new buffer
				for {
					if n, err = t.reader.Read(t.buffer); n == 0 {
						if err != nil {
							token.Err = err
							return false
						}
					} else {
						t.begin = 0
						i = 0
						break
					}
				}
			}
			c = t.buffer[i]
			if tokenType[c] != numToken {
				token.Err = fmt.Errorf("invalid number in line %v", t.line)
				token.Type = invalidToken
				return false
			}
		}
		if c == '0' {
			// next char must be .
			if i++; i == n {
				current = append(current, t.buffer[t.begin:i]...)
				// read new buffer
				for {
					if n, err = t.reader.Read(t.buffer); n == 0 {
						if err != nil {
							token.Err = err
							return false
						}
					} else {
						t.begin = 0
						i = 0
						break
					}
				}
			}
			c = t.buffer[i]
			if c != '.' {
				token.Str = "0"
				token.Type = numToken
				t.begin = i
				t.end = n
				return true
			}
		} else {
			for tokenType[c] == numToken {
				if i++; i == n {
					current = append(current, t.buffer[t.begin:i]...)
					// read new buffer
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								return false
							}
						} else {
							t.begin = 0
							i = 0
							break
						}
					}
				}
				c = t.buffer[i]
			}
			if c != '.' {
				goto PARSE_EXPONENT
			}
		}
		// after the dot
		// next char must be number
		if i++; i == n {
			current = append(current, t.buffer[t.begin:i]...)
			// read new buffer
			for {
				if n, err = t.reader.Read(t.buffer); n == 0 {
					if err != nil {
						token.Err = err
						return false
					}
				} else {
					t.begin = 0
					i = 0
					break
				}
			}
		}
		c = t.buffer[i]
		if tokenType[c] != numToken {
			token.Err = fmt.Errorf("expecting number after . in line %v", t.line)
			token.Type = invalidToken
			return false
		}
		// scan numbers
		for tokenType[c] == numToken {
			if i++; i == n {
				current = append(current, t.buffer[t.begin:i]...)
				// read new buffer
				for {
					if n, err = t.reader.Read(t.buffer); n == 0 {
						if err != nil {
							token.Err = err
							return false
						}
					} else {
						t.begin = 0
						i = 0
						break
					}
				}
			}
			c = t.buffer[i]
		}
	PARSE_EXPONENT:
		if c == 'e' || c == 'E' {
			if i++; i == n {
				current = append(current, t.buffer[t.begin:i]...)
				// read new buffer
				for {
					if n, err = t.reader.Read(t.buffer); n == 0 {
						if err != nil {
							token.Err = err
							return false
						}
					} else {
						t.begin = 0
						i = 0
						break
					}
				}
			}
			c = t.buffer[i]
			// do we have a sign?
			if c == '-' || c == '+' {
				if i++; i == n {
					current = append(current, t.buffer[t.begin:i]...)
					// read new buffer
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								return false
							}
						} else {
							t.begin = 0
							i = 0
							break
						}
					}
				}
				c = t.buffer[i]
			}
			//
			if tokenType[c] != numToken {
				token.Err = fmt.Errorf("expecting number in line %v", t.line)
				token.Type = invalidToken
				return false
			}
			// scan numbers
			for tokenType[c] == numToken {
				if i++; i == n {
					current = append(current, t.buffer[t.begin:i]...)
					// read new buffer
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								return false
							}
						} else {
							t.begin = 0
							i = 0
							break
						}
					}
				}
				c = t.buffer[i]
			}
		}
		token.Str = string(append(current, t.buffer[t.begin:i]...))
		token.Type = numToken
		t.begin = i
		t.end = n
		return true
	case objectLeftToken:
		token.Type = objectLeftToken
		token.Str = "{"
		t.begin = i + 1
		t.end = n
		return true
	case objectRightToken:
		token.Type = objectRightToken
		token.Str = "}"
		t.begin = i + 1
		t.end = n
		return true
	case stringToken:
		// skip the first "
		if i++; i == n {
			// read new buffer
			for {
				if n, err = t.reader.Read(t.buffer); n == 0 {
					if err != nil {
						token.Err = err
						return false
					}
				} else {
					i = 0
					break
				}
			}
		}
		t.begin = i
		// search for the end
		for ; i < n; i++ {
			c = t.buffer[i]
			if c == '"' {
				token.Str = string(t.buffer[t.begin:i])
				token.Type = stringToken
				t.begin = i + 1
				t.end = n
				return true
			} else if c == '\\' {
				current = append(t.strBuffer[:0], t.buffer[t.begin:i]...)
				t.begin = i
				t.end = n
				return t.readRestOfString(token, current)
			}
		}
		current = append(t.strBuffer[:0], t.buffer[t.begin:i]...)
		for {
			// read new buffer
			for {
				if n, err = t.reader.Read(t.buffer); n == 0 {
					if err != nil {
						token.Err = err
						return false
					}
				} else {
					t.begin = 0
					i = 0
					break
				}
			}
			for ; i < n; i++ {
				c = t.buffer[i]
				if c == '"' {
					token.Str = string(append(current, t.buffer[t.begin:i]...))
					token.Type = stringToken
					t.begin = i + 1
					t.end = n
					return true
				} else if c == '\\' {
					current = append(current, t.buffer[t.begin:i]...)
					t.begin = i
					t.end = n
					return t.readRestOfString(token, current)
				}
			}
			current = append(current, t.buffer[t.begin:i]...)
		}

	default:
		token.Type = invalidToken
		token.Err = fmt.Errorf("invalid token in line %v", t.line)
		return false
	}
}

func (t *Tokenizer) readRestOfString(token *Token, current []byte) bool {
	var (
		err error
		c   byte
	)
	i, n := t.begin, t.end
	for {
		// scan rest of buffer
		for i < n {
		FAST_LOOP:
			c = t.buffer[i]
			if c == '"' {
				token.Str = string(append(current, t.buffer[t.begin:i]...))
				token.Type = stringToken
				t.begin = i + 1
				t.end = n
				return true
			} else if c == '\\' {
				// save current buffer
				current = append(current, t.buffer[t.begin:i]...)
				if i++; i == n {
					// read new buffer
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								return false
							}
						} else {
							t.begin = 0
							i = 0
							break
						}
					}
				}
			ESCAPED:
				switch t.buffer[i] {
				case '"':
					current = append(current, '"')
				case '\\':
					current = append(current, '\\')
				case '/':
					current = append(current, '/')
				case 'b':
					current = append(current, '\b')
				case 'f':
					current = append(current, '\f')
				case 'n':
					current = append(current, '\n')
				case 'r':
					current = append(current, '\r')
				case 't':
					current = append(current, '\t')
				case 'u':
					var (
						r uint32
					)
					for j := 0; j < 4; j++ {
						if i++; i == n {
							// read new buffer
							for {
								if n, err = t.reader.Read(t.buffer); n == 0 {
									if err != nil {
										token.Err = err
										return false
									}
								} else {
									t.begin = 0
									i = 0
									break
								}
							}
						}
						c = t.buffer[i]
						if hexValid[c] {
							r = r<<4 + hexRune[c]
						} else {
							//token.Err = fmt.Errorf("invalid unicode in line %v", t.line)
							current = append(current, 0xEF, 0xBF, 0xBD)
							t.begin = i
							goto FAST_LOOP
						}
					}
					// utf16 surrogate?
					if r >= 0xD800 && r <= 0xDFFF {
						if r >= 0xDC00 {
							// starts with low surrogate, that is not correct
							r = 0xFFFD // Unicode replacement char
						} else {
						GET_LOW_SURROGATE:
							// fetch the low surrogate, must start with '\'
							if i++; i == n {
								// read new buffer
								for {
									if n, err = t.reader.Read(t.buffer); n == 0 {
										if err != nil {
											token.Err = err
											return false
										}
									} else {
										t.begin = 0
										i = 0
										break
									}
								}
							}
							if t.buffer[i] != '\\' {
								// write unicode replacement char in utf8
								current = append(current, 0xEF, 0xBF, 0xBD)
								t.begin = i
								goto FAST_LOOP
							}
							// next char must be 'u'
							if i++; i == n {
								// read new buffer
								for {
									if n, err = t.reader.Read(t.buffer); n == 0 {
										if err != nil {
											token.Err = err
											return false
										}
									} else {
										t.begin = 0
										i = 0
										break
									}
								}
							}
							if t.buffer[i] != 'u' {
								// write unicode replacement char in utf8
								current = append(current, 0xEF, 0xBF, 0xBD)
								t.begin = i
								goto ESCAPED
							}
							// read the 4 hex values
							var r2 uint32
							for j := 0; j < 4; j++ {
								if i++; i == n {
									// read new buffer
									for {
										if n, err = t.reader.Read(t.buffer); n == 0 {
											if err != nil {
												token.Err = err
												return false
											}
										} else {
											t.begin = 0
											i = 0
											break
										}
									}
								}
								c = t.buffer[i]
								if hexValid[c] {
									r2 = r2<<4 + hexRune[c]
								} else {
									//token.Err = fmt.Errorf("invalid unicode in line %v", t.line)
									current = append(current, 0xEF, 0xBF, 0xBD)
									t.begin = i
									goto FAST_LOOP
								}
							}
							//
							if r2 >= 0xD800 && r2 <= 0xDFFF {
								if r2 >= 0xDC00 {
									// low surrogate
									r = (r&0b1111111111)<<10 | r2&0b1111111111 + 0x10000
								} else {
									// high surrogate
									current = append(current, 0xEF, 0xBF, 0xBD)
									r = r2
									goto GET_LOW_SURROGATE
								}
							} else {
								// second escaped value is not a surrogate
								current = append(current, 0xEF, 0xBF, 0xBD)
								r = r2
							}

						}
					}
					switch {
					// 7 bit 0xxx.xxxx
					case r < 128:
						// ascii
						current = append(current, byte(r))
					// 11 bit 110x.xxxx 10xx.xxxx
					case r < 2048:
						current = append(current, 0b11000000|byte(r>>6), 0b10000000|byte(r)&0b00111111)
					// 16 bit 1110.xxxx 10xx.xxxx 10xx.xxxx
					case r < 65536:
						current = append(current, 0b11100000|byte(r>>12), 0b10000000|byte(r>>6)&0b00111111, 0b10000000|byte(r)&0b00111111)
					// 21 bit 1111.0xxx 10xx.xxxx 10xx.xxxx 10xx.xxxx
					default:
						current = append(current, 0b11110000|byte(r>>18), 0b10000000|byte(r>>12)&0b00111111, 0b10000000|byte(r>>6)&0b00111111, 0b10000000|byte(r)&0b00111111)
					}
				default:
					// unknown escape sequence
					current = append(current, 0xEF, 0xBF, 0xBD)
				}
				if i++; i == n {
					// read new buffer
					for {
						if n, err = t.reader.Read(t.buffer); n == 0 {
							if err != nil {
								token.Err = err
								return false
							}
						} else {
							i = 0
							break
						}
					}
				}
				t.begin = i
				goto FAST_LOOP
			}
			i++
		}
		current = append(current, t.buffer[t.begin:i]...)
		for {
			if n, err = t.reader.Read(t.buffer); n == 0 {
				if err != nil {
					token.Err = err
					return false
				}
			} else {
				t.begin = 0
				i = 0
				break
			}
		}
	}
}
