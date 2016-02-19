# combo
Lexer combinator in Go(GOlang). Based on ideas from
[this](http://theorangeduck.com/page/you-could-have-invented-parser-combinators)

This library allows you do do modular hexing of text input by the means of
defining independent lexers then combining them to work together on a text
input.


# Installation

```bash
go get github.com/FarmRadioHangar/combo
```

# Example

For instance you want to lex the  string `SELECT * FROM ` for the sake of
brevity let us just stick with that string input

The example is taken from [Hand written parsers and Lexers in
Go](https://blog.gopheracademy.com/advent-2014/parsers-lexers/) but not as a
whole for just show casting how this library works.
