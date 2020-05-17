package token

type TokenType string

type TokenPosition struct {
	Row int
	Col int
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     TokenPosition
}

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"fork":   FORK,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, exist := keywords[ident]; exist {
		return tok
	}
	return IDENT
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Ids + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN           = "="
	PLUS             = "+"
	MINUS            = "-"
	MUL              = "*"
	DIV              = "/"
	LT               = "<"
	GT               = ">"
	BANG             = "!"
	EQUAL            = "=="
	NOT_EQUAL        = "!="
	DOUBLE_AMPERSAND = "&&"
	DOUBLE_PIPE      = "||"
	CARET            = "^"

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	FORK     = "FORK"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
)
