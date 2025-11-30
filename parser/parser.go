package parser

//语法分析器
import (
	"fmt"
	"monkey_Interpreter/ast"
	"monkey_Interpreter/lexer"
	"monkey_Interpreter/token"
)

// 语法分析器结构
type Parser struct {
	l *lexer.Lexer //指向词法分析器实例的指针

	curToken  token.Token //当前词法单元
	peekToken token.Token //当前词法单元的下一位

	errors []string //错误集合
}

// 实例化语法分析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors: []string{},
	} //语法分析器实例

	//读取两个词法单元，以设置curToken和peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// 获取下一个词法单元 前移curToken和peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken() //递归
}

// 入口点
// 解析程序AST语法
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}              //声明ast根节点program
	program.Statements = []ast.Statement{} //语句接口切片集

	for p.curToken.Type != token.EOF { //碰到词法法单元Token EOF文件结尾 表示已将遍历完终止
		stmt := p.parseStatement() //解析具体语法
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken() //下移
	}

	return program
}

// 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// 解析let语句 let x=5;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	//字面量
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	//将当前词法单元作为标识符的 Token 字段，并将其字面值作为标识符的值赋给 stmt.Name
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// =
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// ;
	//TODO: 跳过对表达式的处理，直到遇见分号
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken() //移动
	}
	return stmt
}

// 当前token判断
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 下一个token判断
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 用于判断下一个词法单元的类型是否与给定的类型匹配，并移动到下一个词法单元
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekErrors(t)
		return false
	}
}

// 错误检测
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be “%s”,got=%s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
