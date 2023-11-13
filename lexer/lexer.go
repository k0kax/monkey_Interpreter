package lexer

/*
	词法分析器
*/
import (
	"monkey_Interpreter/token"
)

type Lexer struct {
	input        string //输入
	position     int    // 输入的字符串中的当前位置(指向当前字符)
	readPosition int    // 输入的字符串中的当前读取位置(指向当前字符串之后的一个字符(ch))
	ch           byte   // 当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// 初始化 l.ch,l.position,l.readPosition
	l.readChar()
	return l
}

// 读取下一个字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { //如果下一个字符的长度大于l.input长度，也就是下一个字符再input末尾
		l.ch = 0 // NUL的ASSII码(0)
	} else {
		// 读取
		l.ch = l.input[l.readPosition]
	}
	// 前移
	l.position = l.readPosition
	l.readPosition += 1
}

// 创建词法单元的方法
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// 根据当前的ch创建词法单元
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	//！-/*5；
	//			5 < 10 > 5;
	case '!':
		tok = newToken(token.BANG, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// ****************************************字母********************************
// 读一个标识符并前移词法分析器的位置，直到遇见非字母
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar() //读下一个字符
	}
	return l.input[position:l.position]
}

// 判断给定的参数是否为字母
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// ***********************************空白******************************************
// 跳过空白
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// *******************************************数字************************************
// 读一个标识符并前移词法分析器的位置，知道遇见非数字
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 判断是否是数字
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
