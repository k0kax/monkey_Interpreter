语法分析器parser是将前文生成的词法单元进一步转化为抽象语法树AST
### 一、抽象语法树AST

本次进实现let的语法分析
```shell
let five = 5;  
let ten =10;  
let add = fn(x,y){  
	x+y;  
};  
let result = add(five,ten);  
```
由上可知let主要用于将值绑定到给定的名称上，可以是方法也可以是变量
对let进行语法分析，也就是生成一个属于它的AST


#### 1.1三个接口
一个节点接口Node，返回字面量
```go
type Node interface {
	TokenLiteral() string
}
```
对应TokenLiteral()方法，用于返回字面值
```go
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } //词法单元字面值
```
语句接口Statement,包括节点和语句节点
```go
type Statement interface {
	Node
	statementNode() //语句节点
}
func (ls *LetStatement) statementNode()       {}  
```
上文涉及到一个标识符的结构体,包括对应的词法单元类型和值
```go
type Identifier struct {
	Token token.Token //token.IDENT词法单元
	Value string      //字面量值
}
```
表达式接口
```go
type Expression interface {
	Node
	expressionNode() //表达式节点
}

func (i *Identifier) expressionNode() {}
```
#### 1.2结构体
程序结构体，也就是根节点
```go
type Program struct {
	Statements []Statement //接口类型的切片
}
```
###### 1.3词法单元的字面量
```go
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
```
#### 1.4 let语句专有的结构体
此处的AST采取如下结构
![](https://cdn.jsdelivr.net/gh/k0kax/PicGo@main/image/202311142020165.png)
包括词法单元、标识符名称、值
```go
type LetStatement struct {
	Token token.Token // token.LET词法单元
	Name  *Identifier //保存绑定的标识符名称
	Value Expression  //产生值的表达式
}
```
总代码ast/ast.go
```go
// ast/ast.go
package ast

//Abstrcat Syntax Tree 抽象语法树

//语法分析器将文本或者词法单元形式的源码作为输入，产生一个表示该源码的数据结构。
import "monkey_Interpreter/token"

//三个接口

// 接口1
// 用于返回字面量

type Node interface {
	TokenLiteral() string
}

// 接口2
// 语句

type Statement interface {
	Node
	statementNode() //语句节点
}

// 接口3
// 表达式

type Expression interface {
	Node
	expressionNode() //表达式节点
}

// 程序 根节点

type Program struct {
	Statements []Statement //接口类型的切片
}

// Token字面量

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// 定义所需字段

type LetStatement struct {
	Token token.Token // token.LET词法单元
	Name  *Identifier //保存绑定的标识符名称
	Value Expression  //产生值的表达式
}

func (ls *LetStatement) statementNode()       {}                          //语句节点
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } //词法单元字面值

// 标识符

type Identifier struct {
	Token token.Token //token.IDENT词法单元
	Value string      //字面量值
}

// 表达式节点
func (i *Identifier) expressionNode() {}

// 词法单元字面量
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

```

### 二、语法分析器
##### 2.1语法分析器的结构
包括词法单元指针lexer，当前词法单元curToken，下一个词法单元peekToken，此处和[[1_1词法分析器]]的position/readPosition 类似
```go
// 语法分析器结构
type Parser struct {
	l *lexer.Lexer //指向词法分析器实例的指针

	curToken  token.Token //当前词法单元
	peekToken token.Token //当前词法单元的下一位
}
```
##### 2.2实例化语法分析器
需要先带入词法单元
然后设置curToken和peekToken，使得词法分析器不断执行
```go
// 实例化语法分析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l} //语法分析器实例

	//读取两个词法单元，以设置curToken和peekToken
	p.nextToken()
	p.nextToken()
	return p
}
```
### 2.3 让词法单元动起来
``` go
// 获取下一个词法单元 前移curToken和peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()//递归
}
```
### 2.4 解析

##### 解析程序
```go
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}              //
	program.Statements = []ast.Statement{} //

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken() //下移
	}
	return program
}
```
##### 解析总语句
以解析`let  x = 2;`为例 
注意由skipWhiteSpace()跳过空白
第一轮
```go
type Parser struct {
	l *lexer.Lexer //指向词法分析器实例的指针 

	curToken  token.Token //当前词法单元 let
	peekToken token.Token //当前词法单元的下一位 x
}
```

```go
type Token struct {  
Type TokenType //LET
// 字面量  
Literal string  //let
}
```
读到`LET`
对当前的词法单元的类型进行判断，如果是LET就进入专属的解析程序，返回一个let语句节点接口
```go
// 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}
```
##### 解析let语句
先把当前词法单元放进Let的词法结构体中
判断当前词法单元的下一个(x)
```
	curToken :let
	peekToken :x
```

```go
type Token struct {  
	Type TokenType //IDENT
	// 字面量  
	Literal string  //X
}
```
是否是标识符，是则进行移位操作，移位后如下
```
curToken:x
peekToken:=
```
不是则退出
```go
//字面量
	if !p.expectPeek(token.IDENT) {//移位
		return nil
	}
	//将当前词法单元作为标识符的 Token 字段，并将其字面值作为标识符的值赋给 stmt.Name
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
```
然后将相关的token/字面量塞到标识符结构体的token/值中
```go
type Identifier struct {  
	Token token.Token //token.IDENT  
	Value string //  x
}
```
此时Identifier就会变成
```go
Token : INENT//字面量
Value : x //值
```

进入下一轮
```
	curToken  ：x
	peekToken ：=
```
 ***
```go
	//移动到下一个词法单元
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
```
判断当前词法单元下一个(=)是否是赋值符号(=),是则移位
不是则退出

进入下一轮
```go
type Parser struct {
	l *lexer.Lexer //指向词法分析器实例的指针 

	curToken  token.Token //当前词法单元 2
	peekToken token.Token //当前词法单元的下一位 ;
}
```
***
```go
//移动到下一个词法单元  
if !p.curTokenIs(token.SEMICOLON) {  
	p.nextToken()  
}
```
判断当前是否是(；)，是则移动到下一个单元，不是则退出
```go
// 解析let语句 let x = 2;
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	/*
		curToken:LET
		peekToken:x
	*/
	//字面值 peekToken: x
	if !p.expectPeek(token.IDENT) {//移位 
		/*
			curToken:x
			peekToken:=
		*/
		return nil
	}
	//将当前词法单元作为标识符的 Token 字段，并将其字面值作为标识符的值赋给 stmt.Name
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// =
	// peekToken: =
	if !p.expectPeek(token.ASSIGN) {//移位
		/*
			curToken:=
			peekToken:2
		*/
		return nil
	}


	//移动到下一个词法单元
	// TODO: 跳过对表达式的处理，直到遇见分号
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()//移位
		
	}
	return stmt
}
```
##### 判断函数
```go
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
		return false
	}
}
```

#### 2.4总代码parse/parse.go
```go
package parser

//语法分析器
import (
	"monkey_Interpreter/ast"
	"monkey_Interpreter/lexer"
	"monkey_Interpreter/token"
)

// 语法分析器结构
type Parser struct {
	l *lexer.Lexer //指向词法分析器实例的指针

	curToken  token.Token //当前词法单元
	peekToken token.Token //当前词法单元的下一位
}

// 实例化语法分析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l} //语法分析器实例

	//读取两个词法单元，以设置curToken和peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// 获取下一个词法单元 前移curToken和peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 解析程序语法
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}              //
	program.Statements = []ast.Statement{} //

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
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
	case token.LET://识别出let
		return p.parseLetStatement()
	default:
		return nil
	}
}

// 解析let语句
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {//识别出 标识符,并进行移位跳转
		return nil
	}
	//将当前词法单元作为标识符的 Token 字段，并将其字面值作为标识符的值赋给 stmt.Name
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {//识别出 =
		return nil
	}

	//移动到下一个词法单元
	// TODO: 跳过对表达式的处理，直到遇见分号
	if !p.curTokenIs(token.SEMICOLON) {//识别出 ;
		p.nextToken()
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
		p.nextToken()//移动到下一个词法单元
		return true
	} else {
		return false
	}
}

```
### 三、测试函数



parser/parser_test.go

```go
package parser

import (
	"monkey_Interpreter/ast"
	"monkey_Interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	//输入
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
		expectedIdentifier string //预期的标识符
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	//遍历
	//i是循环索引，表示当前迭代的索引位置，从0开始递增。tt 是一个临时变量，它代表 tests 切片中的每个元素。
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// 测试let语句的断言 期望字面量名称
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" { //字面量检查
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())//报错
		return false
	}

	//将s转换为*ast.LetStatement类型
	letStmt, ok := s.(*ast.LetStatement) //断言
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	//字面值
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	//字面值
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

```
测试结果
![](https://cdn.jsdelivr.net/gh/k0kax/PicGo@main/image/202311142217388.png)