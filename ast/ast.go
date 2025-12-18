// ast/ast.go
package ast

//Abstrcat Syntax Tree 抽象语法树
//普拉特语法分析器，自上而下
//语法分析器将文本或者词法单元形式的源码作为输入，产生一个表示该源码的数据结构。
import (
	"bytes"
	"monkey_Interpreter/token"
	"strings"
)

//-----------------------------------------三个接口------------------------------------------

// 接口1 Node接口，连接成树，需要的节点
// 用于返回字面量
type Node interface {
	TokenLiteral() string //token字面量
	String() string       //token的string形式，用于调试
}

// 接口2
// 语句Statement
type Statement interface {
	Node            //继承node接口
	statementNode() //语句节点
}

// 接口3
// 表达式
type Expression interface {
	Node             //继承node接口
	expressionNode() //表达式节点
}

// ---------------------------------------------根节点Program----------------------------------------------------
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

// ------------------------------------------------标识符indent----------------------------------------------
type Identifier struct {
	Token token.Token //token.IDENT词法单元
	Value string      //字面量值
}

// 语句节点
func (i *Identifier) expressionNode() {}

// ----------------------------------------相关语句statement------------------------------------------------------
// ------------------------------------------let语句---------------------------------------------------------
type LetStatement struct {
	Token token.Token // token.LET词法单元
	Name  *Identifier //保存绑定的标识符名称
	Value Expression  //产生值的表达式
}

// 对齐前面两个接口
func (ls *LetStatement) statementNode()       {}                          //语句节点
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } //词法单元字面值

// 词法单元字面量
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// --------------------------------------------------return语句--------------------------------------------
type ReturnStatement struct {
	Token       token.Token //return词法单元
	ReturnValue Expression  //返回值
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// --------------------------------------------------表达式语句Expression--------------------------------------------
type ExpressionStatement struct {
	Token      token.Token //该表达式的第一个词法单元
	Expression Expression  //保存表达式
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// --------------------------------------------为所有ast添加String()方法-------------------------------------
// 主要用于调试时打印ast节点
func (p *Program) String() string {
	var out bytes.Buffer //创建缓冲区

	for _, s := range p.Statements {
		out.WriteString(s.String()) //将每条语句的String()方法返回值写入缓冲区
	}

	return out.String() //以字符串形式返回
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString("=")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier) String() string { return i.Value }

// -------------------------------------------int字面量-----------------------------------------
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// ------------------------------------------------前缀运算符的AST-----------------------------------------
// <前缀运算符><表达式>
type PrefixExpression struct {
	Token    token.Token
	Operator string     //可能是- ！
	Right    Expression //运算符右边的表达式
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// -------------------------------------------------------中缀运算符------------------------------
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()

}

// --------------------------------------------布尔型字面量-----------------------------------
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// ---------------------------------------------------if-else--------------------------------------------
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// 一些列词法单元的组合
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// -------------------------------------------函数字面量-----------------------------------
type FunctionLiteral struct {
	Token      token.Token //fn词法单元
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// -----------------------------------------------调用表达式-----------------------------
type CallExpression struct {
	Token     token.Token //(词法单元
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}
