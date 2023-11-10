package token

// 词法单元类型
type TokenType string

// 词法单元
type Token struct {
	Type TokenType
	// 字面量
	Literal string
}

const (
	// 特殊类型
	ILLEGAL = "ILLEGAL" // 未知字符
	EOF     = "EOF"     // 文件结尾

	// 标识符+字面量
	IDENT = "IDENT" // add, foobar, x, y
	INT   = "INT"   // 1343456

	// 运算符
	ASSIGN = "="
	PLUS   = "+"

	// 分隔符
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
