package combo

import "testing"

const (
	ILLEGAL TokenType = iota
	EOF
	WS
	IDENT
	ASTERIC
	COMMA
	SELECT
	FROM
)

func printToken(typ TokenType) string {
	switch typ {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case IDENT:
		return "IDENT"
	case ASTERIC:
		return "ASTERIC"
	case COMMA:
		return "COMMA"
	case SELECT:
		return "SELECT"
	case FROM:
		return "FROM"
	}
	return ""
}

func TestLexer(t *testing.T) {
	c := &LexCombinator{}
	whiteSpace := c.ChainOr(RuneLex(' ', WS), RuneLex('\t', WS), RuneLex('\n', WS))
	s := StringLex("SELECT", SELECT)
	asterik := RuneLex('*', ASTERIC)
	from := StringLex("FROM", FROM)
	//comma := RuneLex(',', COMMA)
	c.And(s, whiteSpace, asterik, whiteSpace, from, whiteSpace)
	sample := "SELECT * FROM "
	tok, err := c.Lex([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	sampleTypes := []TokenType{
		SELECT, WS, ASTERIC, WS, FROM, WS,
	}
	for k, v := range tok {
		typ := sampleTypes[k]
		if typ != v.Type() {
			t.Errorf("expected %s got %s", printToken(typ), printToken(v.Type()))
		}
	}
}
