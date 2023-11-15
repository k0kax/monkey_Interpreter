package parser

import (
	"monkey_Interpreter/ast"
	"monkey_Interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
			let y = 10;
			let foobar = 838383;`
	l := lexer.New(input) //词词法分析器实例
	p := New(l)           //语法解析器

	program := p.ParseProgram() //解析程序，并将返回的抽象语法树（AST）存储在变量program中

	//非空检查
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	//程序语法长度检查
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements.got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string //预期标识符
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	//遍历
	//i是循环索引，表示当前迭代的索引位置，从0开始递增。tt 是一个临时变量，它代表 tests 切片中的每个元素。
	for i, tt := range tests {
		stmt := program.Statements[i] //语句结构体
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// 测试let语句的断言 期望字面量名称
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" { //字面量检查
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	//将s转换为*ast.LetStatement类型
	letStmt, ok := s.(*ast.LetStatement) //断言

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	//字面量值
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	//token的字面量
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}
