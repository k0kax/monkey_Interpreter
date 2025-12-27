package token

// 词法单元类型
type TokenType string

// 词法单元
type Token struct {
	Type TokenType
	// 字面量
	Literal string
}

// 声明一些词法常量
const (
	// 特殊类型
	ILLEGAL = "ILLEGAL" // 未知字符
	EOF     = "EOF"     // 文件结尾

	// 标识符+字面量
	IDENT = "IDENT" // add, foobar, x, y
	INT   = "INT"   // 1343456

	// 运算符
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	// 分隔符
	COMMA     = ","
	SEMICOLON = ";"

	LT       = "<"
	GT       = ">"
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
	//*********
	IF     = "IF"
	ELIF   = "ELIF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	//比较运算符
	EQ     = "=="
	NOT_EQ = "!="

	//字符串
	STRING = "STRING"

	//ToDo
	//逻辑运算符
	AND = "&&" //和
	OR  = "||" //或

	//
	COLON = ":"
)

//！-/*5；
//			5 < 10 > 5;

// 关键词表 哈希表
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
	//*********
	"if":     IF,
	"else":   ELSE,
	"elif":   ELIF,
	"return": RETURN,

	"true":  TRUE,
	"false": FALSE,
}

// 通过关键词表判断给定的标识符是否是关键词
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
