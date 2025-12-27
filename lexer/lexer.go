package lexer

/*
	词法分析器
*/
import (
	"monkey_Interpreter/token"
)

type Lexer struct {
	input        string // 输入
	position     int    // 输入的字符串中的当前位置 (指向当前字符)
	readPosition int    // 输入的字符串中的当前读取位置 (指向当前字符之后的一个字符(ch))
	ch           byte   // 当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// 初始化 l.ch,l.position,l.readPosition
	l.readChar()
	return l
}

// 读取input的下一个字符，并前移其在input中的位置
// 检查是否到到input的结尾
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { //下一个要读取的字符位置大于整个输入字符串的长度，已经读到输入的末尾
		l.ch = 0 // NUL的ASSII码(0)，表示尚未读取任何内容或文件结尾
	} else {
		//未到input结尾，将l.ch置为下一个字符
		l.ch = l.input[l.readPosition]
	}
	// 字符向前移
	l.position = l.readPosition
	l.readPosition += 1
}

// 创建词法单元的方法
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch), //字面值
	}
}

// 根据当前的ch创建词法单元，匹配对应的Type和字面量Literal
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() //跳过空格和一些

	switch l.ch {
	case '=':
		if l.peekChar() == '=' { //特殊处理 等于==
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
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
		if l.peekChar() == '=' { //特殊处理 不等于！=
			ch := l.ch //获取当前字符
			l.readChar()
			literal := string(ch) + string(l.ch) //！=
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
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

	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	//Todo
	//逻辑运算符
	case '&': //与运算
		if l.peekChar() == '&' {
			ch := l.ch //获取当前字符
			l.readChar()
			literal := string(ch) + string(l.ch) //&&
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			tok.Literal = ""
			tok.Type = token.EOF
		}
	case '|': //或运算
		if l.peekChar() == '|' {
			ch := l.ch //获取当前字符
			l.readChar()
			literal := string(ch) + string(l.ch) //||
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			tok.Literal = ""
			tok.Type = token.EOF
		}
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)

	//检查是否是标识符
	default:
		if isLetter(l.ch) { //判断是否是字母
			tok.Literal = l.readIdentifier()          //字面量indent
			tok.Type = token.LookupIdent(tok.Literal) //检查关键字
			return tok
		} else if isDigit(l.ch) { //检查是否是数字
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
		l.readChar() //读下一个字符 不断后移l.positionS
	}
	return l.input[position:l.position] //截取输出字符串 完整的标识符ident
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

// 多字符匹配 检测下一个字符
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) { //l.readPosition字段（此时字符的位置）大于输入的长度，代表已经读完
		return 0
	} else {
		return l.input[l.readPosition] //返回下一个字符的位置（string位置是从0开始的）
	}
}

// 读字符串
func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
