// ast/ast.go
package ast

//Abstrcat Syntax Tree 抽象语法树
//普拉特语法分析器，自上而下
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
	Node            //继承node接口
	statementNode() //语句节点
}

// 接口3
// 表达式
type Expression interface {
	Node             //继承node接口
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

//对齐前面两个接口
func (ls *LetStatement) statementNode()       {}                          //语句节点
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal } //词法单元字面值

// 标识符
type Identifier struct {
	Token token.Token //token.IDENT词法单元
	Value string      //字面量值
}

// 语句节点
func (i *Identifier) expressionNode() {}

// 词法单元字面量
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
