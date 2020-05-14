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
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"fork":   FORK,
	"true":   TRUE,
	"false":  FALSE,
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
	ASSIGN       = "="
	PLUS_ASSIGN  = "+="
	MINUS_ASSIGN = "-="
	MUL_ASSIGN   = "*="
	DIV_ASSIGN   = "/="
	PLUS         = "+"
	MINUS        = "-"
	MUL          = "*"
	DIV          = "/"
	LT           = "<"
	GT           = ">"
	BANG         = "!"
	EQUAL        = "=="
	NOT_EQUAL    = "!="

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	FORK     = "FORK"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)
