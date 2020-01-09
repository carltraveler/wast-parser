package parser

import (
	"errors"
	"fmt"
	"github.com/ontio/wast-parser/lexer"
	"io"
)


type Parse interface {
	Parse(parser *ParserBuffer) error
}

type ParserBuffer struct {
	tokens []lexer.Token
	curr   int
}

func NewParserBuffer(input string) (*ParserBuffer, error) {
	lex := lexer.NewLexer(input)
	var tokens []lexer.Token
	for lex.Eof() == false {
		token , err := lex.Parse()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return &ParserBuffer{ tokens:tokens, curr:0 }, nil
}

func (self *ParserBuffer)Empty() bool {
	token := self.PeekToken()

	return token == nil || token.Type() == lexer.RParenType
}

func (self *ParserBuffer)Cursor() Cursor {
	return Cursor{
		parser:self,
		curr:self.curr,
	}
}

func (self *ParserBuffer)ExpectString() (string, error) {
	cursor := self.Cursor()
	str := cursor.String()
	if str == "" {
		return "", errors.New("expect string")
	}
	self.curr = cursor.curr

	return str, nil
}

func (self *ParserBuffer) ExpectKeywordMatch(expect string) error {
	kw, err := self.ExpectKeyword()
	if err != nil  {
		return err
	}
	if kw != expect {
		return fmt.Errorf("expect keyword: %s, got: %s", expect, kw)
	}

	return nil
}

func (self *ParserBuffer) ExpectLParen() error {
	cursor := self.Cursor()
	err := cursor.ExpectLparen()
	if err != nil {
		return err
	}

	self.curr = cursor.curr
	return nil
}

func (self *ParserBuffer) ExpectRParen() error {
	cursor := self.Cursor()
	err := cursor.ExpectRparen()
	if err != nil {
		return err
	}

	self.curr = cursor.curr
	return nil
}

func (self *ParserBuffer) ExpectKeyword() (string, error) {
	cursor := self.Cursor()
	kw := cursor.Keyword()
	if len(kw) == 0 {
		return "", errors.New("expect keyword")
	}

	self.curr = cursor.curr
	return kw, nil
}

func (self *ParserBuffer) ExpectInteger() (lexer.Integer, error) {
	cursor := self.Cursor()
	val, err := cursor.Integer()
	if err != nil {
		return val, err
	}

	self.curr = cursor.curr
	return val, nil
}

func (self *ParserBuffer)PeekToken() lexer.Token {
	if self.curr < len(self.tokens) {
		return self.tokens[self.curr]
	}

	return nil
}

func (self *ParserBuffer)Peek2Token() lexer.Token {
	if self.curr+1 < len(self.tokens) {
		return self.tokens[self.curr+1]
	}

	return nil
}


func (self *ParserBuffer)TryGetId() string {
	cursor := self.Cursor()
	id := cursor.Id()
	if id == "" {
		return ""
	}

	self.curr = cursor.curr
	return id
}

func (self *ParserBuffer)Parens(fn func (ps *ParserBuffer) error) error {
	err := self.ExpectLParen()
	if err != nil {
		return err
	}

	err = fn(self)
	if err != nil {
		return err
	}

	err = self.ExpectRParen()
	if err != nil {
		return err
	}

	return nil
}

func (self *ParserBuffer)Step(fn func (cursor *Cursor) error) error {
	cursor := self.Cursor()
	err := fn(&cursor)
	if err != nil {
		return err
	}

	self.curr = cursor.curr
	return nil
}

type Cursor struct {
	parser *ParserBuffer
	curr int
}

func (self *Cursor)Clone() *Cursor {
	return &Cursor{
		parser: self.parser,
		curr: self.curr,
	}
}

func (self *Cursor)readToken() lexer.Token {
	if self.curr >= len(self.parser.tokens) {
		return nil
	}

	token := self.parser.tokens[self.curr]
	self.curr += 1
	return token
}

func (self *Cursor)ExpectLparen() error {
	token := self.readToken()
	if token == nil {
		return errors.New("expect lparen, got eof")
	}
	if t, ok := token.(lexer.LParen); !ok {
		return fmt.Errorf("expect lparen, got %s", t)
	}

	return nil
}

func (self *Cursor)ExpectRparen() error {
	token := self.readToken()
	if token == nil {
		return errors.New("expect rparen, got eof")
	}
	if t, ok := token.(lexer.LParen); !ok {
		return fmt.Errorf("expect rparen, got %s", t)
	}

	return nil
}

func (self *Cursor)Keyword() string {
	if token := self.readToken(); token != nil {
		if t, ok := token.(lexer.Keyword); ok {
			return t.Val
		}
	}

	return ""
}

func (self *Cursor)Id() string {
	if token := self.readToken(); token != nil {
		if t, ok := token.(lexer.Identifier); ok {
			return string(t.Val[1:])
		}
	}

	return ""
}

func (self *Cursor)Integer() (val lexer.Integer, err error) {
	if token := self.readToken(); token != nil {
		if t, ok := token.(lexer.Integer); ok {
			return t, nil
		}
	}

	return val, errors.New("expect integer")
}

func (self *Cursor) String() string {
	if token := self.readToken(); token != nil {
		if t, ok := token.(lexer.String); ok {
			return string(t.Val)
		}
	}

	return ""
}






