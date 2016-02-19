package combo

import (
	"bytes"
	"errors"
	"fmt"
)

type TokenType int

type Token interface {
	Position
	Value() string
	Type() TokenType
}

type Position interface {
	Left() int
	Right() int
}

type Lexer func(in *bytes.Reader) (Token, error)

type simpleToken struct {
	typ   TokenType
	val   string
	left  int
	right int
}

func (s *simpleToken) Value() string {
	return s.val
}

func (s *simpleToken) Type() TokenType {
	return s.typ
}

func (s *simpleToken) Left() int {
	return s.left
}

func (s *simpleToken) Right() int {
	return s.right
}

func errorMSG(msg string, pos int, ch, expect string) error {
	return fmt.Errorf(">> error>>%s at: %d, found %s expecting: %s ", msg, pos, ch, expect)
}

func newToken(typ TokenType, val string, left, right int) Token {
	return &simpleToken{
		typ:   typ,
		val:   val,
		left:  left,
		right: right,
	}
}

type LexCombinator struct {
	tokens []Token
	err    error
	lexers []Lexer
}

func (combo *LexCombinator) And(lexers ...Lexer) {
	combo.lexers = append(combo.lexers, combo.ChainAnd(lexers...))
}

func (combo *LexCombinator) ChainAnd(lexers ...Lexer) Lexer {
	return func(in *bytes.Reader) (Token, error) {
		if combo.err != nil {
			return nil, combo.err
		}
		if len(lexers) > 0 {
			var tokens []Token
			for _, v := range lexers {
				t, err := v(in)
				if err != nil {
					combo.err = err
					return nil, err
				}
				tokens = append(tokens, t)
			}
			combo.tokens = append(combo.tokens, tokens...)
			return tokens[len(tokens)-1], nil
		}
		return nil, errors.New("no lexers to combine")
	}

}

func (combo *LexCombinator) Or(lexers ...Lexer) {
	combo.lexers = append(combo.lexers, combo.ChainOr(lexers...))
}

func (combo *LexCombinator) ChainOr(lexers ...Lexer) Lexer {
	return func(in *bytes.Reader) (Token, error) {
		if combo.err != nil {
			return nil, combo.err
		}
		if len(lexers) > 0 {
			var currTok Token
			for _, v := range lexers {
				t, _ := v(in)
				if t != nil {
					currTok = t
					break
				}
			}
			if currTok != nil {
				return currTok, nil
			}
			return nil, errors.New("some shit with or combination")
		}
		return nil, errors.New("no lexers to combine")
	}
}

func (combo *LexCombinator) Lex(src []byte) ([]Token, error) {
	r := bytes.NewReader(src)
	for _, v := range combo.lexers {
		_, err := v(r)
		if err != nil {
			return nil, err
		}
	}
	return combo.tokens, combo.err
}

func StringLex(s string, typ TokenType) Lexer {
	return func(b *bytes.Reader) (Token, error) {
		left := int(b.Size()) - b.Len()
		size := 0
		for _, v := range s {
			ch, sz, err := b.ReadRune()
			if err != nil {
				b.Seek(int64(left), 0)
				return nil, err
			}
			size = size + sz
			if ch != v {
				b.Seek(int64(left), 0)
				return nil, errorMSG("unexpected token", left+size, string(ch), string(v))
			}
		}
		return newToken(typ, s, left, left+len(s)), nil
	}
}

func RuneLex(r rune, typ TokenType) Lexer {
	return func(b *bytes.Reader) (Token, error) {
		left := int(b.Size()) - b.Len()
		ch, size, err := b.ReadRune()
		if err != nil {
			return nil, err
		}
		if ch != r {
			b.UnreadRune()
			return nil, errorMSG("unexpected token", left+size, string(ch), string(r))
		}
		return newToken(typ, string(ch), left, left+size), nil
	}
}
