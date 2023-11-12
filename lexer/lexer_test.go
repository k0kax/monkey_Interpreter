package lexer

import (
	"monkey_Interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `let five = 5;
			let ten =10; 
			let add = fn(x,y){
				x+y;
			};
			let result = add(five,ten);
			`
	tests := []struct {
		expectedType    token.TokenType //期待的词法单元类型
		expectedLiteral string          //期待的字面量
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	// 创建词法分析器实例
	l := New(input) //lexer

	//i：表示当前tests中的位置。
	//tt：表示当前迭代的expectedType（词法单元类型）和expectedLiteral（字面量）字段。
	for i, tt := range tests {

		// 调用词法分析器的NextToken方法获取下一个词法单元
		tok := l.NextToken() //token

		// 检查实际的词法单元类型是否与预期的类型一致
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got==%q", i, tt.expectedType, tok.Type)
		}
		// 检查实际的词法单元字面量是否与预期的字面量一致
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got==%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
