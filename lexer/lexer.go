// lexer/lexer.go
package lexer

import "monkey_Interpreter/token"

type Lexer struct {
	input        string
	position     int  // 输入的字符串中的当前位置(指向当前字符)
	readPosition int  // 输入的字符串中的当前读取位置(指向当前字符串之后的一个字符(ch))
	ch           byte // 当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// 初始化 l.ch,l.position,l.readPosition
	l.readChar()
	return l
}

// 读取下一个字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
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
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}
