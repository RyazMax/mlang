package lexer

/*
	Осуществляет разбор входного текста на токены
	Лексер структура с методом NextToken() возвращающий следующий токен из входного текста
*/

import (
	"mlang/token"
	"unicode"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	lastToken    token.Token
}

func New(input string) *Lexer {
	l := &Lexer{input: input, lastToken: newToken(token.SEMICOLON, 0)}
	l.readChar()
	return l
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '\n':
		tok = newToken(token.SEMICOLON, ';')
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.DIV, l.ch)
	case '*':
		tok = newToken(token.MUL, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DOUBLE_AMPERSAND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DOUBLE_PIPE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '^':
		tok = newToken(token.CARET, l.ch)
	case 0:
		tok.Literal = "EOF"
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			l.lastToken = tok
			return tok
		} else if isDigit(l.ch) {
			var ok bool
			tok.Literal, ok = l.readNumber()
			if ok {
				tok.Type = token.INT
			} else {
				tok.Type = token.ILLEGAL
			}
			l.lastToken = tok
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	l.lastToken = tok
	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) && (l.ch != '\n' || l.lastToken.Type == token.SEMICOLON || l.lastToken.Type == token.LBRACE) {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	ok := true
	for isLetter(l.ch) {
		l.readChar()
		ok = false
	}

	return l.input[position:l.position], ok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
